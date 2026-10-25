.PHONY: test e2e lint lint-fix fmt

test:
	go test -race -shuffle=on ./...

## e2e: run against a real account; set HEVO_ACCESS_KEY, HEVO_SECRET_KEY and HEVO_BASE_URL
e2e:
	go test -tags e2e -count=1 ./tests/e2e/...

lint:
	golangci-lint run ./...

lint-fix:
	golangci-lint run --fix ./...

fmt:
	golangci-lint fmt ./...
