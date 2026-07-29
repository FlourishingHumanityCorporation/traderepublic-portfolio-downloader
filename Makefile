.PHONY: all init generate lint test test-v1 test-v2 check architecture-check reset build-darwin-amd64 build-darwin-arm64 build-windows-386 build-windows-amd64 build-linux-386 build-linux-amd64 build-linux-arm64 test-policy

all: lint test

init:
	go run -mod=mod github.com/google/wire/cmd/wire ./...

generate:
	rm -f **/*_mock.go
	./scripts/generate-rest-client.sh
	go generate ./...
	go run ./cmd/example-generator

lint:
	golangci-lint run ./...

test: test-v1 test-v2

test-v1:
	go test -v ./...

test-v2:
	$(MAKE) -C v2 test

# Changed-file test policy gate. Runs in repo-native validation so marker,
# skip/xfail, and validation-lane policy drift cannot hide until a fleet scan.
test-policy:
	@repo=$$PWD; \
	files=$$(git -C "$$repo" ls-files -co --exclude-standard -- .); \
	if [ -z "$$files" ]; then \
		echo "test-policy: no repo files to scan"; \
	else \
		python3 "$$HOME/CodeProjects/.meta/scripts/test-policy-changed-files" --repo "$$repo" $$files; \
	fi
architecture-check:
	APPCHECK_HISTORY_ENABLED=false \
	APPCHECK_PROJECTS_JSON="$(CURDIR)/.appcheck/projects.json" \
	APPCHECK_CODEPROJECTS_ROOT="$(CURDIR)" \
	appcheck run traderepublic-portfolio-downloader --category architecture

check: architecture-check test-policy test

reset:
	rm -rf .session .refresh responses documents transactions.csv

build-darwin-amd64:
	GOOS=darwin GOARCH=amd64 go build -v -o /tmp/portfoliodownloader/public/portfoliodownloader-darwin-amd64 ./cmd/portfoliodownloader/public

build-darwin-arm64:
	GOOS=darwin GOARCH=arm64 go build -v -o /tmp/portfoliodownloader/public/portfoliodownloader-darwin-arm64 ./cmd/portfoliodownloader/public

build-windows-amd64:
	GOOS=windows GOARCH=amd64 go build -v -o /tmp/portfoliodownloader/public/portfoliodownloader-windows-amd64.exe ./cmd/portfoliodownloader/public

build-windows-arm64:
	GOOS=windows GOARCH=arm64 go build -v -o /tmp/portfoliodownloader/public/portfoliodownloader-windows-arm64.exe ./cmd/portfoliodownloader/public

build-linux-amd64:
	GOOS=linux GOARCH=amd64 go build -v -o /tmp/portfoliodownloader/public/portfoliodownloader-linux-amd64 ./cmd/portfoliodownloader/public

build-linux-arm64:
	GOOS=linux GOARCH=arm64 go build -v -o /tmp/portfoliodownloader/public/portfoliodownloader-linux-arm64 ./cmd/portfoliodownloader/public
