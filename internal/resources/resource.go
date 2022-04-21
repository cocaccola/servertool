package resources

type Resource interface {
	GetName() string
	Updated() bool
	Reconcile(resourceMap ResourceMap) error
}

type Resources []*ResourceContainer

type ResourceMap map[string]Resource

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

func (rc *ResourceContainer) GetName() string {
	rc.setResource()
	return rc.resource.GetName()
}

func (rc *ResourceContainer) GetResource() Resource {
	rc.setResource()
	return rc.resource
}

func (rc *ResourceContainer) Reconcile(resourceMap ResourceMap) error {
	rc.setResource()
	return rc.resource.Reconcile(resourceMap)
}
