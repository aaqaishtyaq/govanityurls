package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type config map[string]configEntry

type configEntry struct {
	URL        string `json:"URL"`
	MainBranch string `json:"MainBranch"`
}

func main() {
	githubOrgName, err := githubOrg()
	panicIfNil(err)

	domainName, err := domainName(githubOrgName)
	panicIfNil(err)

	ghPagesDir := "gh-pages"
	if err := os.MkdirAll(ghPagesDir, 0o755); err != nil {
		panic(fmt.Errorf("failed to create dir %s (%w)", ghPagesDir, err))
	}

	cnameFile := filepath.Join(ghPagesDir, "CNAME")
	if err := os.WriteFile(cnameFile, []byte(domainName), 0o644); err != nil {
		panic(fmt.Errorf("failed to write CNAME file %s (%w)", cnameFile, err))
	}

	err = generateIndexPage(githubOrgName, ghPagesDir)
	panicIfNil(err)

	data, err := os.ReadFile("packages.json")
	if err != nil {
		panic(fmt.Errorf("failed open %s (%w)", "packages.json", err))
	}

	cfg := &config{}
	if err := json.Unmarshal(data, cfg); err != nil {
		panic(fmt.Errorf("failed load %s (%w)", "packages.json", err))
	}

	err = generateModulePages(cfg, domainName, githubOrgName, ghPagesDir)
	panicIfNil(err)
}

func panicIfNil(err error) {
	if err != nil {
		panic(err)
	}
}

func githubOrg() (string, error) {
	githubOrgName := os.Getenv("GITHUB_REPOSITORY_OWNER")
	if githubOrgName == "" {
		return "", fmt.Errorf("please provide a GitHub organization name via the GITHUB_REPOSITORY_OWNER environment variable")
	}

	return githubOrgName, nil
}

func domainName(githubOrgName string) (string, error) {
	domainName := os.Getenv("DOMAIN_NAME")
	if domainName == "" {
		repository := os.Getenv("GITHUB_REPOSITORY")
		if repository == "" || !strings.HasPrefix(repository, fmt.Sprintf("%s/", githubOrgName)) {
			return "", fmt.Errorf("please provide a domain name via the DOMAIN_NAME environment variable or the repository via the GITHUB_REPOSITORY name")
		}
		domainName = repository[len(githubOrgName)+1:]
	}

	return domainName, nil
}
