# Preview Proxy

A lightweight reverse proxy server written in Go that enables dynamic subdomain-based routing for preview environments.

## Overview

Preview Proxy is a simple yet powerful reverse proxy that allows you to route traffic to different services based on subdomains. It's particularly useful for creating preview environments where each branch or pull request gets its own subdomain.

## Features

- Dynamic subdomain-based routing
- Configurable origin scheme and port
- Health check endpoint (`/proxy/healthz`)
- Built with Go for high performance
- Docker support
- Cross-platform support (Windows, macOS, Linux)

## Configuration

The proxy can be configured using environment variables:

- `PORT`: Port to listen on (default: 18080)
- `PROXY_DOMAIN`: Domain to listen on (required)
- `ORIGIN_PORT`: Port of the origin server (default: 443)
- `ORIGIN_BASE_DOMAIN`: Base domain for the origin server (required)
- `ORIGIN_SCHEME`: Scheme for the origin server (default: http)

## Installation

### Using Go

```bash
go install github.com/hope-inc/preview-proxy/cmd/preview-proxy
```

### Using Docker

```bash
docker pull ghcr.io/hope-inc/preview-proxy:latest
```

## Usage

### Running locally

```bash
# Build and run
make run

# Or run directly
preview-proxy
```

### Running with Docker

```bash
docker run -d \
  -p 18080:18080 \
  -e ORIGIN_BASE_DOMAIN=example.com \
  ghcr.io/hope-inc/preview-proxy:latest
```

## Development

### Building

```bash
make preview-proxy
```

### Testing

```bash
make test
```

### Creating a release

```bash
make packages
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
