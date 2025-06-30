package utils

import (
	"encoding/json"
	"io/ioutil"

	"github.com/example/scanscope/models"
)

// ParseRules reads JSON rules from path.
func ParseRules(path string) ([]models.Rule, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var rules []models.Rule
	if err := json.Unmarshal(data, &rules); err != nil {
		return nil, err
	}
	return rules, nil
}
