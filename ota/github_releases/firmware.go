package github_releases

type firmware struct {
	id       string
	version  string
	checksum string
	repo     repoInfo
	payload  []byte
}

func (f *firmware) Name() string {
	return f.id
}

func (f *firmware) Version() string {
	return f.version
}

func (f *firmware) Checksum() string {
	return f.checksum
}

func (f *firmware) Payload() []byte {
	return f.payload
}
