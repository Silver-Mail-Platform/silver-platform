# Intake

Minimal Go HTTP service for receiving silver instance events.

## What it does

- Exposes `POST /v1/silver/events`
- Requires `Content-Type: application/json`
- Requires a `timestamp` field in RFC3339 format
- Returns `202 Accepted` when valid

No auth is implemented yet.

## Run instructions

Prerequisite: Go `1.22+`.

Run from repository root:

```bash
go run ./intake/cmd/intake
```

Run from the `intake` directory:

```bash
cd intake
go run ./cmd/intake
```

Run on a custom port:

```bash
PORT=9090 go run ./cmd/intake
```

Build and run a binary:

```bash
cd intake
go build -o bin/intake ./cmd/intake
./bin/intake
```

## Example request

```bash
curl -i \
  -X POST http://localhost:8080/v1/silver/events \
  -H 'Content-Type: application/json' \
  -d '{"timestamp":"2026-03-05T10:30:45Z","event":"scan_complete","instance_id":"silver-01"}'
```

## Response behavior

- `202` valid payload
- `400` malformed JSON or missing/invalid `timestamp`
- `405` method not allowed
- `415` unsupported media type
- `500` ingest failure
