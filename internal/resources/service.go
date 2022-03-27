package resources

import (
	"fmt"

	systemd "github.com/coreos/go-systemd/dbus"
)

type ServiceState string
type ServiceOnStart string

const (
	ServiceStateRunning    ServiceState   = "running"
	ServiceStateStopped    ServiceState   = "stopped"
	ServiceOnStartEnabled  ServiceOnStart = "enabled"
	ServiceOnStartDisabled ServiceOnStart = "disabled"
)

type ServiceResource struct {
	Name      string         `json:"name"`
	State     ServiceState   `json:"state"`
	OnStart   ServiceOnStart `json:"onStart"`
	DependsOn []string       `json:"dependsOn,omitempty"`
	updated   bool           `json:"-"`
}

func (sr *ServiceResource) GetName() string {
	return "service:" + sr.Name
}

func (sr *ServiceResource) Updated() bool {
	return sr.updated
}

func (sr *ServiceResource) Reconcile(resourceMap ResourceMap) error {
	serviceStarted := false

	conn, err := systemd.New()
	if err != nil {
		return err
	}
	defer conn.Close()

	property, err := conn.GetUnitProperty(sr.Name, "UnitFileState")
	actualServiceOnStart := property.Value.String()

	property, err = conn.GetUnitProperty(sr.Name, "SubState")
	actualServiceState := property.Value.String()

	if actualServiceOnStart != string(sr.OnStart) && sr.OnStart == ServiceOnStartEnabled {
		_, _, err := conn.EnableUnitFiles([]string{sr.Name + ".service"}, false, false)
		if err != nil {
			return err
		}

	}
	if actualServiceOnStart != string(sr.OnStart) && sr.OnStart == ServiceOnStartDisabled {
		_, err := conn.DisableUnitFiles([]string{sr.Name + ".service"}, false)
		if err != nil {
			return err
		}
	}

	if actualServiceState != string(sr.State) && sr.State == ServiceStateRunning {
		// start service
		_, err := conn.StartUnit(sr.Name, "replace", nil)
		if err != nil {
			return err
		}
		serviceStarted = true
	}
	if actualServiceState != string(sr.State) && sr.State == ServiceStateStopped {
		// stop service
		_, err := conn.StopUnit(sr.Name, "replace", nil)
		if err != nil {
			return err
		}
	}

	// if the service was started we can end here
	// we do not need to check if the service's dependencies were modified
	if serviceStarted {
		return nil
	}

	needsRestart := false
	for _, resourceName := range sr.DependsOn {
		if r, ok := resourceMap[resourceName]; ok {
			if r.Updated() {
				needsRestart = true
			}
			continue
		}
		return fmt.Errorf("could not fetch resource %s from resource map", resourceName)
	}

	if needsRestart {
		// restart the service
		_, err := conn.RestartUnit(sr.Name, "replace", nil)
		if err != nil {
			return err
		}
	}
	return nil
}
