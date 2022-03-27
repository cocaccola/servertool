package resources

import (
	"bytes"
	"crypto/md5"
	"errors"
	"io"
	"os"
	"os/user"
	"strconv"
)

type FileResource struct {
	Path     string `json:"path"`
	User     string `json:"user"`
	Group    string `json:"group"`
	Mode     string `json:"mode"`
	Contents string `json:"contents"`
	updated  bool   `json:"-"`
}

func (fr *FileResource) GetName() string {
	return "file:" + fr.Path
}

func (fr *FileResource) Updated() bool {
	return fr.updated
}

func (fr *FileResource) Reconcile(_ ResourceMap) error {
	needsUpdate := false

	// check that the file has the correct contents
	f, err := os.OpenFile(fr.Path, os.O_RDONLY, 0)
	if errors.Is(err, os.ErrNotExist) {
		needsUpdate = true
	} else if err != nil {
		return err
	} else {
		actualHash := md5.New()
		if _, err := io.Copy(actualHash, f); err != nil {
			return err
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
		return err
	}

	// if needsUpdate is true, re-create the file according to the desired state
	if needsUpdate {
		f, err := os.OpenFile(fr.Path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.FileMode(mode))
		if err != nil {
			return err
		}

		_, err = f.WriteString(fr.Contents)
		if err != nil {
			return err
		}

		fr.updated = true

		f.Close()
	}

	// lazily enforce the desired owners and permissions
	f, err = os.OpenFile(fr.Path, os.O_RDWR, 0)
	if err != nil {
		return err
	}
	defer f.Close()

	u, err := user.Lookup(fr.User)
	if err != nil {
		return err
	}
	uid, err := strconv.Atoi(u.Uid)
	if err != nil {
		return err
	}

	g, err := user.LookupGroup(fr.Group)
	if err != nil {
		return err
	}
	gid, err := strconv.Atoi(g.Gid)
	if err != nil {
		return err
	}

	err = f.Chown(uid, gid)
	if err != nil {
		return err
	}

	err = f.Chmod(os.FileMode(mode))
	if err != nil {
		return err
	}

	return nil
}
