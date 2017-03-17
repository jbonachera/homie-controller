package persistentCache

import (
	"bytes"
	"errors"
	"github.com/jbonachera/homie-controller/log"
	"io/ioutil"
	"os"
	"sync"
)

var lock sync.Mutex

func Start() {
	baseDir := "/var/cache/homie-controller"
	fi, err := os.Stat(baseDir)
	if err != nil {
		err = os.Mkdir(baseDir, 0650)
		if err != nil {
			panic("could not create cache directory: " + err.Error())
		}
	} else if mode := fi.Mode(); !mode.IsDir() {
		panic("could not start cache layer: " + baseDir + "is not a directory")
	}

}

func Get(key string) (*bytes.Buffer, error) {
	lock.Lock()
	defer lock.Unlock()
	dat, err := ioutil.ReadFile("/var/cache/homie-controller/" + key)
	if err != nil {
		return nil, errors.New("key " + key + " not found in cache")
	} else {
		return bytes.NewBuffer(dat), nil
	}
}

func Set(key string, buffer *bytes.Buffer) error {
	lock.Lock()
	defer lock.Unlock()
	err := ioutil.WriteFile("/var/cache/homie-controller/"+key, buffer.Bytes(), 0640)
	log.Debug("caching asset " + key)
	if err != nil {
		return err
	} else {
		return nil
	}
}
