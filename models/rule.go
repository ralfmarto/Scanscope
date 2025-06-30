package models

// Rule represents a scanning rule.
type Rule struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	Category       string   `json:"category"`
	FileExtensions []string `json:"file_extensions"`
	Regex          string   `json:"regex"`
	Prompt         string   `json:"prompt"`
	Children       []Rule   `json:"children"`
}
