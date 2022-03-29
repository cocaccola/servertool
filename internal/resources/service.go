package resources

import (
	"fmt"
	"strings"

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
		return fmt.Errorf("could not establish connection to systemd via dbus: %w", err)
	}
	defer conn.Close()

	property, err := conn.GetUnitProperty(sr.Name+".service", "UnitFileState")
	if err != nil {
		return fmt.Errorf("could not get unit startup information: %w", err)
	}
	actualServiceOnStart := ServiceOnStart(strings.Trim(property.Value.String(), "\""))

	property, err = conn.GetUnitProperty(sr.Name+".service", "SubState")
	if err != nil {
		return fmt.Errorf("could not get unit status: %w", err)
	}
	actualServiceState := ServiceState(strings.Trim(property.Value.String(), "\""))
	if actualServiceState == "dead" {
		actualServiceState = ServiceStateStopped
	}

	// service is enabled but we don't want it to be
	if actualServiceOnStart != sr.OnStart && actualServiceOnStart == ServiceOnStartEnabled {
		_, err := conn.DisableUnitFiles([]string{sr.Name + ".service"}, false)
		if err != nil {
			return fmt.Errorf("could not disable unit %s: %w", sr.Name, err)
		}
		fmt.Printf("disabled service %s\n", sr.Name)
	}

	// service is disabled but we don't want it to be
	if actualServiceOnStart != sr.OnStart && actualServiceOnStart == ServiceOnStartDisabled {
		_, _, err := conn.EnableUnitFiles([]string{sr.Name + ".service"}, false, false)
		if err != nil {
			return fmt.Errorf("could not enable unit %s: %w", sr.Name, err)
		}
		fmt.Printf("enabled service %s\n", sr.Name)
	}

	// service is running but we don't want it to be
	if actualServiceState != sr.State && actualServiceState == ServiceStateRunning {
		// start service
		_, err := conn.StopUnit(sr.Name+".service", "replace", nil)
		if err != nil {
			return fmt.Errorf("could not stop unit %s: %w", sr.Name, err)
		}
		fmt.Printf("stopped service %s\n", sr.Name)
	}

	// service is stopped but we don't want it to be
	if actualServiceState != sr.State && actualServiceState == ServiceStateStopped {
		// stop service
		_, err := conn.StartUnit(sr.Name+".service", "replace", nil)
		if err != nil {
			return fmt.Errorf("could not start unit %s: %w", sr.Name, err)
		}
		serviceStarted = true
		fmt.Printf("started service %s\n", sr.Name)
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
		_, err := conn.RestartUnit(sr.Name+".service", "replace", nil)
		if err != nil {
			return fmt.Errorf("could not restart service %s: %w", sr.Name, err)
		}
		fmt.Printf("restarted service %s\n", sr.Name)
	}
	return nil
}
