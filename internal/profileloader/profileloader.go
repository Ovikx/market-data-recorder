package profileloader

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/santhosh-tekuri/jsonschema"
)

type profile struct {
	Provider string   `json:"provider"`
	Symbols  []string `json:"symbols"`
	WSUrl    string   `json:"wsUrl"`
}

// Unmarshals the JSON content at the specified profile path into a profile struct.
func FromFile(profilePath string, schemaPath string) (*profile, error) {
	// Unmarshal the profile
	b, err := os.ReadFile(profilePath)
	if err != nil {
		return nil, err
	}

	var prof profile
	err = json.Unmarshal(b, &prof)
	if err != nil {
		return nil, err
	}

	// Validate the profile
	c, err := jsonschema.Compile(schemaPath)
	if err != nil {
		return nil, err
	}

	err = c.Validate(strings.NewReader(string(b)))
	if err != nil {
		return nil, err
	}

	return &prof, nil
}
