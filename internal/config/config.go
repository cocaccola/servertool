package config

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"slack/servertool/internal/resources"
)

var ErrDuplicateResourceName = errors.New("duplicate Resource names found")

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
		if _, ok := rm[resource.GetName()]; ok {
			// duplicate resource
			return nil, nil, ErrDuplicateResourceName
		}
		_resource := resource.GetResource()
		rm[resource.GetName()] = &_resource
	}
	return r, rm, nil
}
