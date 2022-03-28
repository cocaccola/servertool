package resources

import (
	"os/exec"
)

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
	return "package:" + pr.Name
}

func (pr *PackageResource) Updated() bool {
	return pr.updated
}

func (pr *PackageResource) Reconcile(_ ResourceMap) error {
	var actualPackageState PackageState

	// determine package state
	cmd := exec.Command("/usr/bin/dpkg-query", "-W", pr.Name)
	if err := cmd.Run(); err != nil {
		if e, ok := err.(*exec.ExitError); ok && e.ExitCode() == 1 {
			// package was not found on the system
			actualPackageState = PackageStateAbsent
		} else {
			// the command ran into another kind of issue
			return err
		}
	} else {
		// package is installed
		actualPackageState = PackageStateInstalled
	}

	// the package is installed when we don't want it to be
	if pr.State != actualPackageState && actualPackageState == PackageStateInstalled {
		// update package database
		cmd = exec.Command("/usr/bin/apt-get", "-y", "-q", "update")
		cmd.Env = append(cmd.Env, "DEBIAN_FRONTEND=noninteractive")
		if err := cmd.Run(); err != nil {
			return err
		}

		// install the package
		cmd = exec.Command("/usr/bin/apt-get", "-y", "-q", "remove", pr.Name)
		cmd.Env = append(cmd.Env, "DEBIAN_FRONTEND=noninteractive")
		if err := cmd.Run(); err != nil {
			return err
		}
		pr.updated = true
	}

	// the package is not installed and we want it to be
	if pr.State != actualPackageState && actualPackageState == PackageStateAbsent {
		// update package database
		cmd = exec.Command("/usr/bin/apt-get", "-y", "-q", "update")
		cmd.Env = append(cmd.Env, "DEBIAN_FRONTEND=noninteractive")
		if err := cmd.Run(); err != nil {
			return err
		}

		// remove package
		cmd = exec.Command("/usr/bin/apt-get", "-y", "-q", "install", pr.Name)
		cmd.Env = append(cmd.Env, "DEBIAN_FRONTEND=noninteractive")
		if err := cmd.Run(); err != nil {
			return err
		}
		pr.updated = true
	}

	return nil
}
