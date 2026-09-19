# caching-proxy

A CLI caching proxy server. It forwards requests to an origin server, caches the
responses, and serves repeat requests from the cache instead of hitting the origin.

## Requirements

Go 1.26+

## Usage

```bash
caching-proxy --port <number> --origin <url>
```

| Flag            | Description                                                 |
| --------------- | ----------------------------------------------------------- |
| `--port`        | Port the proxy listens on (1–65535)                         |
| `--origin`      | Origin server URL to forward to (must be `http` or `https`) |
| `--clear-cache` | Clear the cache and exit                                    |

Requests to `http://localhost:<port>/path` are forwarded to `<origin>/path`.
Responses are returned with the origin's status code and headers, plus:

```
X-Cache: HIT     # served from cache
X-Cache: MISS    # fetched from the origin
```

## Example

```bash
make run PORT=3000 ORIGIN=http://dummyjson.com
```

```bash
curl -i "localhost:3000/products?limit=2"
```

Clear the cache:

```bash
go run . --clear-cache
```

## Development

```bash
make          # list targets
make build    # compile to ./caching-proxy
make run      # start the proxy
make fmt      # format
make lint     # vet + golangci-lint
make test     # run tests
make check    # format check + lint + test, before committing
```

`make fmt` uses [gofumpt](https://github.com/mvdan/gofumpt); `make lint` uses
[golangci-lint](https://golangci-lint.run).
