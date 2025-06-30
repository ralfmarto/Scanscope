package utils

import (
	"io/fs"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/example/scanscope/models"
)

// Finding represents a scanning result.
type Finding struct {
	RuleID      string `json:"rule_id"`
	RuleName    string `json:"rule_name"`
	Category    string `json:"category"`
	File        string `json:"file"`
	Line        int    `json:"line"`
	Snippet     string `json:"snippet"`
	AIValidated bool   `json:"ai_validated"`
}

// Scan applies rules to files under rootDir.
func Scan(rootDir string, rules []models.Rule, cache *Cache, ai *AI) ([]Finding, error) {
	var findings []Finding
	err := filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return nil
		}
		for _, rule := range rules {
			if !matchExtension(path, rule.FileExtensions) {
				continue
			}
			hash, err := FileCategoryHash(path, rule.Category)
			if err != nil {
				continue
			}
			if cache.Has(hash) {
				continue
			}
			fileFindings := applyRule(string(data), rule, path, ai)
			if len(fileFindings) > 0 {
				findings = append(findings, fileFindings...)
			}
			cache.Add(hash)
		}
		return nil
	})
	return findings, err
}

func matchExtension(path string, exts []string) bool {
	if len(exts) == 0 {
		return true
	}
	for _, e := range exts {
		if strings.HasSuffix(path, e) {
			return true
		}
	}
	return false
}

func applyRule(content string, rule models.Rule, file string, ai *AI) []Finding {
	var results []Finding
	re, err := regexp.Compile(rule.Regex)
	if err != nil {
		return results
	}
	matches := re.FindAllStringIndex(content, -1)
	for _, m := range matches {
		snippet := content[m[0]:m[1]]
		line := strings.Count(content[:m[0]], "\n") + 1
		if len(rule.Children) > 0 {
			childResults := applyChildren(snippet, rule.Children, file, line, ai)
			results = append(results, childResults...)
		} else {
			validated := false
			if rule.Prompt != "" && ai != nil {
				ok, err := ai.Validate(rule.Prompt, snippet)
				if err == nil {
					validated = ok
				}
			}
			if rule.Prompt == "" || validated {
				results = append(results, Finding{
					RuleID:      rule.ID,
					RuleName:    rule.Name,
					Category:    rule.Category,
					File:        file,
					Line:        line,
					Snippet:     snippet,
					AIValidated: validated,
				})
			}
		}
	}
	return results
}

func applyChildren(content string, children []models.Rule, file string, line int, ai *AI) []Finding {
	var res []Finding
	for _, child := range children {
		re, err := regexp.Compile(child.Regex)
		if err != nil {
			continue
		}
		if !re.MatchString(content) {
			continue
		}
		if len(child.Children) > 0 {
			res = append(res, applyChildren(content, child.Children, file, line, ai)...)
		} else {
			validated := false
			if child.Prompt != "" && ai != nil {
				ok, err := ai.Validate(child.Prompt, content)
				if err == nil {
					validated = ok
				}
			}
			if child.Prompt == "" || validated {
				res = append(res, Finding{
					RuleID:      child.ID,
					RuleName:    child.Name,
					Category:    child.Category,
					File:        file,
					Line:        line,
					Snippet:     content,
					AIValidated: validated,
				})
			}
		}
	}
	return res
}
