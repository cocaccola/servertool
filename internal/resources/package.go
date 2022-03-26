package resources

type PackageState string

const (
	PackageStateInstalled = "installed"
	PackageStateAbsent    = "absent"
)

type PackageResource struct {
	Name    string       `json:"name"`
	State   PackageState `json:"state"`
	updated bool         `json:"-"`
}

func (pr *PackageResource) GetName() string {
	return pr.Name
}

func (pr *PackageResource) Updated() bool {
	return pr.updated
}
