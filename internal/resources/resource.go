package resources

type Resource interface {
	GetName() string
	Updated() bool
}

type ResourceMap map[string]*Resource

type ResourceContainer struct {
	File     *FileResource    `json:"file,omitempty"`
	Package  *PackageResource `json:"package,omitempty"`
	Service  *ServiceResource `json:"service,omitempty"`
	resource Resource         `json:"-"`
}

func (rc *ResourceContainer) setResource() {
	if rc.resource != nil {
		return
	}
	switch {
	case rc.File != nil:
		rc.resource = rc.File
	case rc.Package != nil:
		rc.resource = rc.Package
	case rc.Service != nil:
		rc.resource = rc.Service
	}
}

type Resources []ResourceContainer

func (rc *ResourceContainer) GetName() string {
	rc.setResource()
	return rc.resource.GetName()
}

func (rc *ResourceContainer) GetResource() Resource {
	switch {
	case rc.File != nil:
		return rc.File
	case rc.Package != nil:
		return rc.Package
	case rc.Service != nil:
		return rc.Service
	default:
		return nil
	}
}
