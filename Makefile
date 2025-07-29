.PHONY: ci
ci: vet test


.PHONY: dep
dep:
	go get -u ./... && go mod tidy


.PHONY: test
test:
	go test -v ./... -count=1 --race


.PHONY: test_coverage
test_coverage:
	go test ./... -coverprofile=coverage.out


.PHONY: vet
vet:
	go vet ./...


.PHONY: lint
lint:
	golangci-lint run --disable errcheck
