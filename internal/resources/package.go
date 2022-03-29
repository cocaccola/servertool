package resources

import (
	"errors"
	"fmt"
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
		var exitErr *exec.ExitError
		if ok := errors.As(err, &exitErr); ok && exitErr.ExitCode() == 1 {
			// the package is not installed
			actualPackageState = PackageStateAbsent
		} else {
			// some other error happened
			return fmt.Errorf("could not query package database: %w", err)
		}
	}

	// the above command ran cleanly, so the package is installed
	if actualPackageState != PackageStateAbsent {
		actualPackageState = PackageStateInstalled
	}

	// the package is installed when we don't want it to be
	if pr.State != actualPackageState && actualPackageState == PackageStateInstalled {
		// update package database
		cmd = exec.Command("/usr/bin/apt-get", "-y", "-q", "update")
		cmd.Env = append(cmd.Env, "DEBIAN_FRONTEND=noninteractive")
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("could not update package database: %w", err)
		}

		// install the package
		cmd = exec.Command("/usr/bin/apt-get", "-y", "-q", "remove", pr.Name)
		cmd.Env = append(cmd.Env, "DEBIAN_FRONTEND=noninteractive")
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("could not remove package %s: %w", pr.Name, err)
		}
		pr.updated = true
	}

	// the package is not installed and we want it to be
	if pr.State != actualPackageState && actualPackageState == PackageStateAbsent {
		// update package database
		cmd = exec.Command("/usr/bin/apt-get", "-y", "-q", "update")
		cmd.Env = append(cmd.Env, "DEBIAN_FRONTEND=noninteractive")
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("could not update package database: %w", err)
		}

		// remove package
		cmd = exec.Command("/usr/bin/apt-get", "-y", "-q", "install", pr.Name)
		cmd.Env = append(cmd.Env, "DEBIAN_FRONTEND=noninteractive")
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("could not install package %s: %w", pr.Name, err)
		}
		pr.updated = true
	}

	return nil
}
