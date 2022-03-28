package resources

type Resource interface {
	GetName() string
	Updated() bool
	Reconcile(resourceMap ResourceMap) error
}

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

type Resources []*ResourceContainer

func (rc *ResourceContainer) GetName() string {
	rc.setResource()
	return rc.resource.GetName()
}

func (rc *ResourceContainer) GetResource() Resource {
	rc.setResource()
	return rc.resource
}

func (rc *ResourceContainer) Reconcile(resourceMap ResourceMap) error {
	// should probably call setResource here for safety
	return rc.resource.Reconcile(resourceMap)
}
