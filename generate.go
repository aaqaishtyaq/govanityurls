package main

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
)

type entry struct {
	Name          string
	URL           string
	Domain        string
	GitHubOrgName string
	MainBranch    string
}

type index struct {
	GithubOrgName string
}

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

func generateModulePages(cfg *config, domainName, ghOrgName, ghPagesDir string) error {
	t := template.Must(template.New("html").Parse(templateText))
	for name, cfgEntry := range *cfg {
		if cfgEntry.MainBranch == "" {
			cfgEntry.MainBranch = "main"
		}
		e := entry{
			Name:          name,
			URL:           cfgEntry.URL,
			Domain:        domainName,
			GitHubOrgName: ghOrgName,
			MainBranch:    cfgEntry.MainBranch,
		}

		dir := filepath.Join(ghPagesDir, name)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("failed to create dir %s (%w)", dir, err)
		}
		file := filepath.Join(ghPagesDir, name, "index.html")
		fh, err := os.Create(file)
		if err != nil {
			return fmt.Errorf("failed to open %s (%w)", file, err)
		}
		defer fh.Close()

		if err := t.Execute(fh, e); err != nil {
			return fmt.Errorf("failed to render template (%w)", err)
		}
	}

	return nil
}

func generateIndexPage(ghOrgName, ghPagesDir string) error {
	file := filepath.Join(ghPagesDir, "index.html")
	t := template.Must(template.New("html").Parse(indexFileContents))
	f, err := os.Create(file)
	if err != nil {
		return fmt.Errorf("failed to open %s (%w)", file, err)
	}

	if err := t.Execute(f, index{ghOrgName}); err != nil {
		return fmt.Errorf("failed to write index file %s (%w)", file, err)
	}

	if err := f.Close(); err != nil {
		return fmt.Errorf("failed to close %s", file)
	}

	return nil
}
