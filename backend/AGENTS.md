# Backend Agent Instructions

## Scope

These instructions apply to all files under `backend/`.

## Project Context

- This is the Go backend for the ChinaTravel app.
- The server entrypoint is `cmd/server/main.go` and listens on `http://localhost:8888`.
- APIs are mounted under `/api/v1` and are consumed by the Vite frontend proxy at `/api`.
- The codebase uses mostly standard-library `net/http`; avoid introducing a web framework unless explicitly requested.

## Architecture

- Keep feature code grouped under `services/<domain>/` using the existing layers:
  - `api/` for HTTP handlers and request/response handling.
  - `application/` for business logic and orchestration.
  - `domain/` for domain models and domain-level types.
  - `infrastructure/` for file repositories, caches, stats, and persistence adapters.
- Use `internal/` only for shared backend internals such as context keys, common models, DB helpers, or legacy code.
- Prefer updating the active `services/` implementation over the older `internal/handlers` and `internal/store` paths unless the task specifically targets legacy code.

## Coding Style

- Run `gofmt` on modified Go files.
- Keep handlers thin: parse inputs, call application services, encode responses.
- Keep business rules in `application/` packages rather than HTTP handlers.
- Return JSON responses with `Content-Type: application/json` for API responses.
- Preserve current error response style unless intentionally refactoring an endpoint.
- Do not add broad global state unless it matches the existing cache/repository patterns.

## Data And Persistence

- Local development data lives in `backend/data/` and `backend/cmd/server/data/`.
- Treat JSON files and SQLite files as development/demo storage.
- Do not hardcode absolute machine paths; use relative paths or existing repository helpers.
- Be careful when editing `data/*.json` because these files may contain local test/demo state.

## Auth And Security

- Protected endpoints should derive the current user from `contextkeys.UserID` set by auth middleware.
- User-specific actions such as bookings, notifications, history, and favorites must reject anonymous users.
- Do not log passwords, bearer tokens, reset tokens, or full authorization headers.
- The forgot-password flow currently returns a reset token for demo use; do not treat this as production-safe.

## API Compatibility

- Preserve existing `/api/v1` route names unless the user asks for a breaking change.
- Keep frontend compatibility in mind when changing response shapes; check `frontend/src/composables` and `frontend/src/views` for consumers.
- Respect the `Accept-Language` header behavior and default to `en` where existing code does so.

## Validation

- Use these commands from `backend/` when relevant:

```bash
go test ./...
go run ./cmd/server
```

- If changing API behavior, prefer targeted manual checks with `curl` against `http://localhost:8888/api/v1/...`.
- Do not commit generated binaries such as `server` or `server_old`.
