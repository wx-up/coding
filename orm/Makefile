.DEFAULT_GOAL := help

cmd := $(MAKECMDGOALS)

.PHONY: controller
controller:
ifeq ($(cmd),controller)
	@go run main.go make $(cmd) -h
else
	@go run main.go make $(cmd)
endif


.PHONY: help
help:
	@echo "请指定目标，可用目标：e2e、e2e_up、e2e_down、lint、fmt、tidy、ut、check"

.PHONY: e2e
e2e:
	sh ./script/integrate_test.sh

.PHONY: e2e_up
e2e_up:
	docker compose -f script/integrate_test_compose.yml up -d

.PHONY: e2e_down
e2e_down:
	docker compose -f script/integrate_test_compose.yml down

.PHONY: e2e_ps
e2e_ps:
	docker compose -f script/integrate_test_compose.yml ps

.PHONY: lint
lint:
	@golangci-lint run

.PHONY: fmt
fmt:
	# -w 覆盖原始内容
	# -l 列出与 goimports 格式不同的文件
	@goimports -l -w .

.PHONY: tidy
tidy:
	@go mod tidy -v

# 单元测试
.PHONY: ut
ut:
	# -race 竞争检测
	# -coverprofile 覆盖率
	# go tool cover -html=cover.out
	@go test -race -coverprofile=cover.out -v ./...


.PHONY: check
check:
	# $(MAKE) 执行 make 命令
	@$(MAKE) fmt
	@$(MAKE) tidy