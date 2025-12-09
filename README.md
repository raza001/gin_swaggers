# gin_swagger

A small, reusable helper for protecting and mounting Swagger (or other API documentation UIs) in Gin applications.  
It provides a configurable middleware (IP allowlist, static header token, enable/disable) and a convenient `AttachSwagger` helper function.

## Table of Contents

- Introduction
- Features
- Installation
- Usage
  - Basic example
  - With functional options
  - Using with gin-swagger
- Configuration
- Behavior & HTTP responses
- Examples
- Troubleshooting
- Contributors
- License

## Introduction

`gin_swagger` is a utility for protecting Swagger UI routes in Gin-based services.

## Features

- Enable/disable swagger route
- IP allowlist
- Static header token auth
- Custom route path via options

## Installation

Place this file in your Go project or install via:

```
go get github.com/<you>/gin_swagger
```

## Usage

```go
gin_swagger.AttachSwagger(r, handler)
```

## Configuration

See provided Go source for full configuration details.

## Contributors

- You

## License

MIT License (replace or update as needed)
