# envcontract

`envcontract` is a Go library for detecting configuration drift between your
Go config structs and the environment variables your app actually receives.

It answers a simple question before your app ships or starts:

> Does this `.env` file or process environment still match what my Go config
> struct expects?

This helps catch missing, renamed, unused, or mistyped environment variables
before they turn into "works locally, breaks in production" incidents.

> Status: under active development. The core parser, source loaders, and
> validation engine exist, but the public API and CLI are not production-ready
> yet.

## Why This Exists

Environment variables are easy to add and easy to forget.

A developer or AI coding agent may add a new config field:

```go
type Config struct {
	OpenAIKey string `env:"OPENAI_API_KEY,required"`
	ModelName string `env:"MODEL_NAME,default=gpt-4.1-mini"`
	Port      int    `env:"PORT,required"`
}
```

But the matching `.env.example`, CI secrets, Docker Compose file, deployment
settings, or production environment may not get updated.

`envcontract` is built for that drift-detection lane. It treats your Go config
struct as the source of truth and checks environment data against it.

## Current Capabilities

- Parse Go structs with `env` tags into field contracts.
- Support required fields with `required`.
- Support defaults with `default=value`.
- Support nested structs.
- Support pointer fields.
- Load variables from `.env`-style files.
- Load variables from the system process environment.
- Use an in-memory mock source for tests.
- Report malformed `.env` lines as warnings.
- Validate missing variables, type mismatches, and unused variables internally.
- Export a struct's field contract as portable JSON (`ExportSchema`), and
  parse it back (`ParseSchema`).

Supported field kinds today:

- `string`
- `int`
- `int64`
- `float64`
- `bool`

## Example Contract

```go
type Config struct {
	Host  string  `env:"HOST,default=localhost"`
	Port  int     `env:"PORT,required"`
	Debug bool    `env:"DEBUG"`
	Rate  float64 `env:"RATE"`
}
```

From that struct, `envcontract` can derive the expected environment contract:

| Env key | Go field | Type | Required | Default |
|---|---|---|---|---|
| `HOST` | `Host` | `string` | no | `localhost` |
| `PORT` | `Port` | `int` | yes | |
| `DEBUG` | `Debug` | `bool` | no | |
| `RATE` | `Rate` | `float64` | no | |

## Example `.env`

```env
HOST=localhost
PORT=8080
DEBUG=true
RATE=3.14
```

If `PORT` is missing, the validation engine can report a missing required
variable. If `PORT=abc`, it can report a type mismatch because `abc` is not an
integer.

## Schema Export

Go structs only exist as real types inside a compiled Go program, so a
separate CLI binary has no way to inspect one directly. `ExportSchema` bridges
that gap: run it once from your own code to snapshot a struct's contract as
JSON, then anything else (including the upcoming CLI) can read that file
without needing your struct at all.

```go
data, err := envcontract.ExportSchema(&Config{})
if err != nil {
	return err
}

os.WriteFile("envcontract.schema.json", data, 0o644)
```

`ParseSchema` reverses the process, decoding that JSON back into
`[]envcontract.FieldContract`.

## Development

Run tests:

```bash
make test
```

Run lint:

```bash
make lint
```

If your local Go or linter cache is not writable, use temporary cache
directories:

```bash
GOCACHE=/tmp/envcontract-go-build go test ./... -race -cover
GOCACHE=/tmp/envcontract-go-build GOLANGCI_LINT_CACHE=/tmp/envcontract-golangci-lint-cache golangci-lint run ./...
```

## Roadmap

| Week | Focus | Status |
|---|---|---|
| 1 | Struct parser and `FieldContract` type | Done |
| 2 | Source adapters and env loading | Done |
| 3 | Validation engine | Internal implementation done |
| 4 | Public API and result type | Done |
| 5 | CLI | In progress |
| 6 | `.env.example` generator and schema export | Pending |
| 7 | Documentation, hardening, and `v0.1.0` | Pending |

## Planned Public API

The final package API is expected to look roughly like this:

```go
result, err := envcontract.Validate(&Config{},
	envcontract.WithFile(".env"),
	envcontract.WithUnusedCheck(),
)
if err != nil {
	return err
}

if result.HasErrors() {
	result.Print(os.Stdout)
	os.Exit(1)
}
```

This API is not finished yet. The current codebase is still focused on the
parser, source adapters, and validation engine.

## CLI

A CLI is under construction at `cmd/envcontract`. Build it locally with:

```bash
make build
./bin/envcontract version
```

Available today:

```bash
envcontract version
```

Planned:

```bash
envcontract check
envcontract init
```

## License

MIT
