# envcontract

A Go library and CLI for environment config contracts — validate, generate, and enforce `.env` consistency from your config struct.

> ⚠️ This project is under active development. Not ready for production use yet.

## What it does

EnvContract validates that your environment variables match the contract defined in your Go config struct — catching missing variables, type mismatches, and config drift before your application starts.

## Planned features

- `envcontract check` — validate `.env` against your config struct
- `envcontract init` — generate `.env.example` directly from your struct
- `envcontract schema` — export a machine-readable schema of expected variables
- Build time enforcement via `go generate`
- Pluggable sources — `.env` files, system env, and more

## Development
```bash
make test    # run tests
make lint    # run linter
make build   # build the CLI
```

## Status

| Week | Focus | Status |
|------|-------|--------|
| 1 | Struct parser and FieldContract type | ✅ Done |
| 2 | Source adapters and env loading | 🔄 In progress |
| 3 | Validation engine | ⏳ Pending |
| 4 | Public API and Result type | ⏳ Pending |
| 5 | CLI | ⏳ Pending |
| 6 | Generator and schema export | ⏳ Pending |
| 7 | Hardening and v0.1.0 | ⏳ Pending |

## License

MIT