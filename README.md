# gowiki

A small wiki based on the [Writing Web Applications](https://go.dev/doc/articles/wiki/) tutorial, updated for modern Go.

## Prerequisites

- Go 1.22 or later

## Check your code

```sh
go fmt ./...
go vet ./...
```

Optional static analysis:

```sh
go install honnef.co/go/tools/cmd/staticcheck@latest
staticcheck ./...
```

## Run locally

```sh
go run .
```

Then open http://localhost:8080/view/FrontPage (or any page name made of letters and digits).

## Build and run in Docker

Build and start the app. The `wiki_data` volume persists wiki pages.

```sh
docker build -t gowiki .
docker run -d -p 8080:8080 --mount source=wiki_data,target=/home/appuser/wiki_data gowiki
```

## What changed from the original tutorial

- Go 1.22 path routing (`GET /view/{title}`) instead of regex-based handlers
- `embed` for HTML templates instead of reading files from disk at runtime
- `os.ReadFile` / `os.WriteFile` instead of deprecated `ioutil`
- HTML auto-escaping via `html/template` (plain `string` body fields)
- Path-safe file storage with `filepath.Join` and title validation
- HTTP method enforcement on save (`POST` only)
- Multi-stage Docker build with a pinned Go/Alpine version and non-root user
