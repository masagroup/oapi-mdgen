# oapi-mdgen

`oapi-mdgen` is a command-line tool that converts OpenAPI 3 specifications into Markdown documentation. It generates a readable Markdown file detailing the endpoints, paths, operations, parameters, request bodies, responses, and schemas defined in your OpenAPI specification.

## Features

- Parses OpenAPI v3 specifications.
- Generates structured Markdown documentation grouping endpoints by tags.
- Includes details on requests, responses, and schemas.
- Built on top of [libopenapi](https://github.com/pb33f/libopenapi).

## Installation

Ensure you have Go installed on your system. You can build the tool by running:

```bash
make build
```

Alternatively, you can install or run it directly using `go`:

```bash
go install github.com/masagroup/oapi-mdgen
```

```bash
go run github.com/masagroup/oapi-mdgen -i input.yaml -o output.md
```

## Usage

```bash
oapi-mdgen --input <openapi.yaml> --output <documentation.md>
```

### Flags

- `-i`, `--input`: (Required) Path to the input OpenAPI specification file.
- `-o`, `--output`: (Required) Path to the output Markdown file.

### Example

```bash
oapi-mdgen -i openapi.yaml -o api-docs.md
```

## Development

The project includes a `Makefile` for common tasks:

- `make build`: Builds the project.
- `make test`: Runs the tests.
- `make lint`: Runs `golangci-lint` (requires Docker).
- `make coverage`: Generates test coverage data.
- `make coverage-html`: Generates an HTML report for test coverage.