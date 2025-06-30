package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/example/scanscope/utils"
)

func main() {
	root, _ := os.Getwd()
	rulesPath := filepath.Join(root, "config", "rules.json")
	rules, err := utils.ParseRules(rulesPath)
	if err != nil {
		log.Fatalf("parse rules: %v", err)
	}

	cachePath := filepath.Join(root, "scanner-cache.json")
	cache, err := utils.LoadCache(cachePath)
	if err != nil {
		log.Fatalf("load cache: %v", err)
	}

	ai := utils.NewAI("gpt-3.5-turbo")

	findings, err := utils.Scan(root, rules, cache, ai)
	if err != nil {
		log.Fatalf("scan: %v", err)
	}

	if err := cache.Save(cachePath); err != nil {
		log.Printf("save cache error: %v", err)
	}

	if err := writeReports(findings); err != nil {
		log.Fatalf("write reports: %v", err)
	}

	fmt.Printf("%d findings saved to report.json and report.md\n", len(findings))
}

func writeReports(findings []utils.Finding) error {
	data, err := json.MarshalIndent(findings, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile("report.json", data, 0644); err != nil {
		return err
	}

	f, err := os.Create("report.md")
	if err != nil {
		return err
	}
	defer f.Close()

	for _, find := range findings {
		fmt.Fprintf(f, "### %s (%s)\n", find.RuleName, find.Category)
		fmt.Fprintf(f, "File: %s:%d\n\n", find.File, find.Line)
		fmt.Fprintf(f, "```\n%s\n```\n", find.Snippet)
		if find.AIValidated {
			fmt.Fprintln(f, "Risco validado pela IA.")
			fmt.Fprintln(f)
		}
	}
	return nil
}
