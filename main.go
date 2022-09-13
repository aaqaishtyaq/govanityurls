package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
)

var templateText = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta name="go-import"
          content="{{ .Domain }}/{{ .Name }}
                   git {{ .URL }}" />
    <meta name="go-source"
          content="{{ .Domain }}/{{ .Name }}
                   {{ .URL }}
                   {{ .URL }}/tree/{{ .MainBranch }}{/dir}
                   {{ .URL }}/blob/{{ .MainBranch }}{/dir}/{file}#L{line}" />
    <meta http-equiv="refresh" content="0; url={{ .URL }}">
</head></html>
`

var indexFileContents = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta http-equiv="refresh" content="0; url=https://github.com/{{ .GithubOrgName }}">
</head></html>
`

type config map[string]configEntry

type configEntry struct {
	URL        string `json:"URL"`
	MainBranch string `json:"MainBranch"`
}

type entry struct {
	Name          string
	URL           string
	Domain        string
	GitHubOrgName string
	MainBranch    string
}

func main() {
	githubOrgName := os.Getenv("GITHUB_REPOSITORY_OWNER")
	if githubOrgName == "" {
		panic(fmt.Errorf("please provide a GitHub organization name via the GITHUB_REPOSITORY_OWNER environment variable"))
	}

	domainName := os.Getenv("DOMAIN_NAME")
	if domainName == "" {
		repository := os.Getenv("GITHUB_REPOSITORY")
		if repository == "" || !strings.HasPrefix(repository, fmt.Sprintf("%s/", githubOrgName)) {
			panic(fmt.Errorf("please provide a domain name via the DOMAIN_NAME environment variable or the repository via the GITHUB_REPOSITORY name"))
		}
		domainName = repository[len(githubOrgName)+1:]
	}

	ghPagesDir := "gh-pages"
	if err := os.MkdirAll(ghPagesDir, 0o755); err != nil {
		panic(fmt.Errorf("failed to create dir %s (%w)", ghPagesDir, err))
	}

	cnameFile := filepath.Join(ghPagesDir, "CNAME")
	if err := os.WriteFile(cnameFile, []byte(domainName), 0o644); err != nil {
		panic(fmt.Errorf("failed to write CNAME file %s (%w)", cnameFile, err))
	}

	data, err := os.ReadFile("packages.json")
	if err != nil {
		panic(fmt.Errorf("failed open %s (%w)", "packages.json", err))
	}
	cfg := &config{}
	if err := json.Unmarshal(data, cfg); err != nil {
		panic(fmt.Errorf("failed load %s (%w)", "packages.json", err))
	}

	indexFile := filepath.Join(ghPagesDir, "index.html")
	iTpl := template.Must(template.New("html").Parse(indexFileContents))
	indexFh, err := os.Create(indexFile)
	if err != nil {
		panic(fmt.Errorf("failed to open %s (%w)", indexFile, err))
	}

	anon := struct {
		GithubOrgName string
	}{
		GithubOrgName: githubOrgName,
	}
	if err := iTpl.Execute(indexFh, anon); err != nil {
		panic(fmt.Errorf("failed to write index file %s (%w)", indexFile, err))
	}
	if err := indexFh.Close(); err != nil {
		panic(fmt.Errorf("failed to close %s", indexFile))
	}

	// if err := os.WriteFile(indexFile, []byte(indexFileContents), 0644); err != nil {
	// 	panic(fmt.Errorf("failed to write index file %s (%w)", indexFile, err))
	// }

	tpl := template.Must(template.New("html").Parse(templateText))
	for name, cfgEntry := range *cfg {
		if cfgEntry.MainBranch == "" {
			cfgEntry.MainBranch = "main"
		}
		e := entry{
			Name:          name,
			URL:           cfgEntry.URL,
			Domain:        domainName,
			GitHubOrgName: githubOrgName,
			MainBranch:    cfgEntry.MainBranch,
		}

		dir := filepath.Join(ghPagesDir, name)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			panic(fmt.Errorf("failed to create dir %s (%w)", dir, err))
		}
		file := filepath.Join(ghPagesDir, name, "index.html")
		fh, err := os.Create(file)
		if err != nil {
			panic(fmt.Errorf("failed to open %s (%w)", file, err))
		}
		defer fh.Close()

		if err := tpl.Execute(fh, e); err != nil {
			panic(fmt.Errorf("failed to render template (%w)", err))
		}
	}
}
