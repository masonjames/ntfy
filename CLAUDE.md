# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

ntfy is an HTTP-based pub-sub notification service written in Go. It consists of:
- A server that handles message publishing, subscriptions, and delivery
- A CLI client for publishing and subscribing
- A React web app for browser-based access
- Mobile app support via Firebase Cloud Messaging (FCM)

## Build Commands

```bash
# Build everything (web app, docs, server/client) - slow
make build

# Build server+client for development
make cli-darwin-server    # macOS (current arch)
make cli-linux-server     # Linux (current arch)
make cli-linux-amd64      # Linux amd64 via GoReleaser

# Build client only (no server, no CGO)
make cli-client

# Web app
make web                  # Build web app (npm install + build)
make web-deps             # Install npm dependencies
make web-build            # Build only
make web-lint             # Run eslint
make web-fmt              # Run prettier

# Documentation (MkDocs)
make docs                 # Build docs
make docs-deps            # Install Python deps

# Run all checks (tests + linting)
make check
```

## Testing

```bash
# Run all tests
make test

# Run tests with verbose output
make testv

# Run tests with race detection
make race

# Run a single test
go test -v ./server -run TestServer_PublishAndPoll

# Run tests in a specific package
go test ./server/...
go test ./user/...
go test ./client/...
```

Test files use `stretchr/testify/require` for assertions. Test helpers are in `test/` package.

## Code Architecture

### Package Structure

- `main.go` - Entry point, builds CLI app via `cmd.New()`
- `cmd/` - CLI commands (serve, publish, subscribe, user, access, tier, token, webpush)
- `server/` - HTTP server, message handling, subscriptions, caches
- `client/` - Client library for publishing/subscribing
- `user/` - User management, authentication, SQLite-backed auth database
- `log/` - Structured logging
- `util/` - Utilities (rate limiting, time parsing, batching, gzip)
- `web/` - React frontend (Vite + MUI)

### Server Components

The `server.Server` struct orchestrates:
- HTTP/HTTPS/Unix socket listeners
- SMTP server for email-to-notification
- Firebase client for mobile push
- Web push for browser notifications
- Message cache (SQLite)
- File cache for attachments
- Visitor rate limiting
- User/auth management

Key files:
- `server/server.go` - Main server, request routing, pub/sub logic
- `server/message_cache.go` - SQLite message persistence
- `server/visitor.go` - Rate limiting per IP/user
- `server/server_firebase.go` - FCM integration
- `server/server_webpush.go` - Web push notifications

### Build Tags

- `noserver` - Build client-only (no server, no CGO required)
- `sqlite_omit_load_extension` - SQLite security hardening

### Configuration

Server config via flags, env vars, or YAML file (`/etc/ntfy/server.yml` by default).
All flags have env var equivalents prefixed with `NTFY_` (e.g., `NTFY_BASE_URL`).

### Web App

React app in `web/` using:
- Vite for bundling
- Material UI components
- Dexie.js for IndexedDB persistence
- i18next for internationalization

Built output goes to `server/site/` and is embedded in the Go binary via `//go:embed`.

## Linting

```bash
make fmt           # gofmt + prettier
make fmt-check     # Check formatting without changes
make vet           # go vet
make lint          # golint
make staticcheck   # staticcheck
make web-lint      # eslint for web app
```
