package resources

import (
	"bytes"
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"os"
	"os/user"
	"strconv"
)

type FileEnsure string

const (
	FileEnsurePresent FileEnsure = "present"
	FileEnsureAbsent  FileEnsure = "absent"
)

type FileResource struct {
	Ensure   FileEnsure `json:"ensure"`
	Path     string     `json:"path"`
	User     string     `json:"user"`
	Group    string     `json:"group"`
	Mode     string     `json:"mode"`
	Contents string     `json:"contents"`
	updated  bool       `json:"-"`
}

func (fr *FileResource) GetName() string {
	return "file:" + fr.Path
}

func (fr *FileResource) Updated() bool {
	return fr.updated
}

func (fr *FileResource) Reconcile(_ ResourceMap) error {
	if fr.Ensure == FileEnsureAbsent {
		if err := os.Remove(fr.Path); err != nil {
			return err
		}
		return nil
	}

	needsUpdate := false

	// check that the file has the correct contents
	f, err := os.OpenFile(fr.Path, os.O_RDONLY, 0)
	if errors.Is(err, os.ErrNotExist) {
		needsUpdate = true
	} else if err != nil {
		return fmt.Errorf("could not open file %s: %w", fr.Path, err)
	} else {
		actualHash := md5.New()
		if _, err := io.Copy(actualHash, f); err != nil {
			return fmt.Errorf("could not process file %s: %w", fr.Path, err)
		}

		desiredHash := md5.New()
		io.WriteString(desiredHash, fr.Contents)

		if !bytes.Equal(actualHash.Sum(nil), desiredHash.Sum(nil)) {
			needsUpdate = true
		}
	}
	f.Close()

	mode, err := strconv.ParseUint(fr.Mode, 8, 32)
	if err != nil {
		return fmt.Errorf("could not parse desired file permissions for %s, %s: %w", fr.Path, fr.Mode, err)
	}

	// if needsUpdate is true, re-create the file according to the desired state
	if needsUpdate {
		f, err := os.OpenFile(fr.Path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.FileMode(mode))
		if err != nil {
			return fmt.Errorf("could not create file %s: %w", fr.Path, err)
		}

		_, err = f.WriteString(fr.Contents)
		if err != nil {
			return fmt.Errorf("could not write contents to file %s: %w", fr.Path, err)
		}

		fr.updated = true

		f.Close()
	}

	// lazily enforce the desired owners and permissions
	f, err = os.OpenFile(fr.Path, os.O_RDWR, 0)
	if err != nil {
		return fmt.Errorf("could not open file %s: %w", fr.Path, err)
	}
	defer f.Close()

	u, err := user.Lookup(fr.User)
	if err != nil {
		return fmt.Errorf("could not find user %s: %w", fr.User, err)
	}
	uid, err := strconv.Atoi(u.Uid)
	if err != nil {
		return fmt.Errorf("could not parse uid for user %s: %w", fr.User, err)
	}

	g, err := user.LookupGroup(fr.Group)
	if err != nil {
		return fmt.Errorf("could not find group %s: %w", fr.Group, err)
	}
	gid, err := strconv.Atoi(g.Gid)
	if err != nil {
		return fmt.Errorf("could not parse gid for group %s: %w", fr.Group, err)
	}

	err = f.Chown(uid, gid)
	if err != nil {
		return fmt.Errorf("could not modify ownership for file %s: %w", fr.Path, err)
	}

	err = f.Chmod(os.FileMode(mode))
	if err != nil {
		return fmt.Errorf("could not modify permissions for file %s: %w", fr.Path, err)
	}

	return nil
}
