package resources

type ServiceState string
type ServiceOnStart string

const (
	ServiceStateEnabled    ServiceState   = "running"
	ServiceStateStopped    ServiceState   = "stopped"
	ServiceOnStartEnabled  ServiceOnStart = "enabled"
	ServiceOnStartDisabled ServiceOnStart = "disabled"
)

type ServiceResource struct {
	Name    string         `json:"name"`
	State   ServiceState   `json:"state"`
	OnStart ServiceOnStart `json:"onStart"`
	updated bool           `json:"-"`
}

func (sr *ServiceResource) GetName() string {
	return sr.Name
}

func (sr *ServiceResource) Updated() bool {
	return sr.updated
}
