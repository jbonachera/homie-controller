package github_releases

type firmware struct {
	id      string
	version string
	owner   string
	repo    string
}

func (f *firmware) Name() string {
	return f.id
}

func (f *firmware) Version() string {
	return f.version
}

func (f *firmware) Checksum() string {
	return ""
}

func (f *firmware) Payload() []byte {
	return []byte{}
}
