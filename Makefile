GO ?= go

.PHONY: fmt fmt-check vet test test-race verify sync-check

fmt:
	gofmt -w core service manifest extensionpoint host

fmt-check:
	test -z "$$(gofmt -l core service manifest extensionpoint host)"

vet:
	$(GO) vet ./...

test:
	$(GO) test ./...

test-race:
	$(GO) test -race ./...

verify: fmt-check vet test

sync-check:
	@test -f AGENTS.md
	@test -f .agents/rules/development-workflow.md
	@test -f .agents/rules/synchronization.md
	@grep -q "upstream repository" .agents/rules/synchronization.md
	@test -f host/README.md
	@test -f extensionpoint/README.md
	@test -z "$$(find . -maxdepth 1 -name '*.go' -print)"
