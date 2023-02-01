.DEFAULT_GOAL := help


.PHONE: golangci-lint_install
golangci-lint_install:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.50.1

.PHONE: golangci-lint_run
golangci-lint_run:
	golangci-lint run ./...

.PHONY: help
help:
	@echo "请输入特定的目标"