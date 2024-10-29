.PHONY: install-bin uninstall-bin
BIN_DIR := $(CURDIR)/bin
$(BIN_DIR): ; mkdir -p $(BIN_DIR)
$(BIN_DIR)/%: | $(BIN_DIR) ; GOBIN=$(BIN_DIR) go install $(PACKAGE)
$(BIN_DIR)/mockgen:       PACKAGE=github.com/golang/mock/mockgen@v1.6.0
$(BIN_DIR)/golangci-lint: PACKAGE=github.com/golangci/golangci-lint/cmd/golangci-lint@v1.61.0
install-bin: | $(BIN_DIR)/mockgen $(BIN_DIR)/golangci-lint
uninstall-bin: ; rm -rf $(BIN_DIR)

.PHONY: build
BUILD_DIR := $(CURDIR)/build
build:
	go build -o $(BUILD_DIR)/server $(CURDIR)/cmd/server.go

.PHONY: lint
lint:
	golangci-lint run ./cmd/... ./internal/...

.PHONY: test
test:
	go test -v -cover -coverpkg ./internal/... -coverprofile=coverage.out ./internal/...

.PHONY: mock
GEN_MOCK_DIR := $(CURDIR)/generate/mock
mock:
	rm -rf $(GEN_MOCK_DIR)/**/*_mock.go

	$(BIN_DIR)/mockgen -package mock -source ./internal/redis/redis.go -destination $(GEN_MOCK_DIR)/redis/redis_mock.go
	$(BIN_DIR)/mockgen -package mock -source ./internal/repo/campaign.go -destination $(GEN_MOCK_DIR)/repo/campaign_mock.go
	$(BIN_DIR)/mockgen -package mock -source ./internal/repo/point_history.go -destination $(GEN_MOCK_DIR)/repo/point_history_mock.go
	$(BIN_DIR)/mockgen -package mock -source ./internal/repo/user_campaign.go -destination $(GEN_MOCK_DIR)/repo/user_campaign_mock.go
	$(BIN_DIR)/mockgen -package mock -source ./internal/repo/user.go -destination $(GEN_MOCK_DIR)/repo/user_mock.go