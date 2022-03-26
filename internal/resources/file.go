package resources

type FileResource struct {
	Path     string `json:"path"`
	User     string `json:"user"`
	Group    string `json:"group"`
	Mode     int    `json:"mode"`
	Contents string `json:"contents"`
	updated  bool   `json:"-"`
}

func (fr *FileResource) GetName() string {
	return fr.Path
}

func (fr *FileResource) Updated() bool {
	return fr.updated
}
