package config

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"slack/servertool/internal/resources"
)

var (
	ErrDuplicateResourceName = errors.New("duplicate Resource names found")
	ErrFetchingResourceName  = errors.New("could not fetch resource name")
	ErrFetchingResource      = errors.New("could not fetch resource")
)

func Parse(path string) (resources.Resources, resources.ResourceMap, error) {
	f, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, nil, err
	}

	var r resources.Resources

	err = json.Unmarshal(data, &r)
	if err != nil {
		return nil, nil, err
	}

	rm := make(resources.ResourceMap, len(r))

	for _, resource := range r {
		name := resource.GetName()
		if name == "" {
			return nil, nil, ErrFetchingResourceName
		}
		if _, ok := rm[name]; ok {
			// duplicate resource
			return nil, nil, ErrDuplicateResourceName
		}
		_resource := resource.GetResource()
		if _resource == nil {
			return nil, nil, ErrFetchingResource
		}
		rm[name] = _resource
	}
	return r, rm, nil
}
