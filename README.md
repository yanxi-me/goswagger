## goswagger

One command to view beautiful documentation from a Swagger-compliant API.

#### Install:

```bash
git clone git@github.com:yanxi-me/goswagger.git && cd goswagger
go install cmd/goswagger.go
```

Make sure that GOBIN environment variable is set, and it is included in PATH environment variable.

#### Usage:

```bash
goswagger your/swgger/json/path [port]

# example
goswagger api-examples/pet-store.json
```
