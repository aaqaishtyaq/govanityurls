# GoVanityUrls

`govanityurls` generates HTML pages with meta tags for Go in order to
provide package URLs like `go.aaqa.dev/gmux`.
The pages are hosted on GitHub Pages and are autogenerated via GitHub Actions.

## Requirements

Add these secrets as part of your Github Actions Secrets

| variable                  | description     | example       |   |   |
|---------------------------|-----------------|---------------|---|---|
| `GITHUB_REPOSITORY_OWNER` | Github Username | `aaqaishtyaq` |   |   |
| `DOMAIN_NAME`             | Domain Name     | `go.aaqa.dev` |   |   |
|                           |                 |               |   |   |