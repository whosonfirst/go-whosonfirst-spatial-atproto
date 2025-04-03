GOMOD=$(shell test -f "go.work" && echo "readonly" || echo "vendor")
LDFLAGS=-s -w

CWD=$(shell pwd)

vuln:
	govulncheck ./...

# https://github.com/marcboeker/go-duckdb?tab=readme-ov-file#vendoring
# go install github.com/goware/modvendor@latest
modvendor:
	modvendor -copy="**/*.a **/*.h" -v

cli:
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/server cmd/server/main.go

debug:
	go run -mod $(GOMOD) cmd/server/main.go \
		-verbose \
		-spatial-database-uri 'duckdb://?uri=$(CWD)/fixtures/sf_county.parquet'
