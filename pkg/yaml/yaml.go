package yaml

import (
	"io/ioutil"

	"github.com/go-yaml/yaml"
)

// Marshal returns the YAML encoding of v.
func Marshal(v interface{}) ([]byte, error) {
	return yaml.Marshal(v)
}

// Unmarshal parses the YAML-encoded data and stores the result in the value pointed to by v.
func Unmarshal(data []byte, v interface{}) error {
	return yaml.Unmarshal(data, v)
}

// ReadFile reads yaml from file and unmarshals into result
func ReadFile(filename string, result interface{}) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, result)
}

// WriteFile marshals data and writes to yaml file
func WriteFile(filename string, data interface{}) error {
	bytes, err := yaml.Marshal(data)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, bytes, 0755)
}
