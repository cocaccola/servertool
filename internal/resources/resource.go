package resources

type Resource interface {
	GetName() string
	Updated() bool
}

type ResourceMap map[string]*Resource

type ResourceContainer struct {
	File    *FileResource    `json:"file,omitempty"`
	Package *PackageResource `json:"package,omitempty"`
	Service *ServiceResource `json:"service,omitempty"`
}

type Resources []ResourceContainer

func (rc *ResourceContainer) GetName() string {
	switch {
	case rc.File != nil:
		return rc.File.GetName()
	case rc.Package != nil:
		return rc.Package.GetName()
	case rc.Service != nil:
		return rc.Service.GetName()
	default:
		return ""
	}
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
