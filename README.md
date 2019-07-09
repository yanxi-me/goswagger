## goswagger

One command to view Swagger API document, according to the Swagger json configuration.

#### Install:

```bash
git clone git@github.com:yanxi-me/goswagger.git && cd goswagger
go install cmd/goswagger.go
```

Make sure that GOBIN environment variable is set, and it is included in PATH environment variable.

#### Usage:

```bash
goswagger /swgger/json/path

# example
goswagger api-examples/pet-store.json
```
