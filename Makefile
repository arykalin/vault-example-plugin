TEST?=$$(go list ./... |grep -v 'vendor')
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)

PLUGIN_NAME := example-vault-plugin
PLUGIN_DIR := bin
PLUGIN_PATH := $(PLUGIN_DIR)/$(PLUGIN_NAME)

VAULT_CONT := $(notdir $(shell pwd))_vault_1
DOCKER_CMD := docker exec -it $(VAULT_CONT)
VAULT_CMD := $(DOCKER_CMD) vault

MOUNT := $(PLUGIN_NAME)
SHA256 := $$(shasum -a 256 "$(PLUGIN_PATH)" | cut -d' ' -f1)
SHA256_DOCKER_CMD := sha256sum "/vault_plugin/$(PLUGIN_NAME)" | cut -d' ' -f1



### Exporting variables for demo and tests
.EXPORT_ALL_VARIABLES:
VAULT_ADDR = http://127.0.0.1:8200
#Must be set,otherwise cloud certificates will timeout
VAULT_CLIENT_TIMEOUT = 180s

fmt:
	gofmt -w $(GOFMT_FILES)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"


#Need to unset VAULT_TOKEN when running vault with dev parameter.
unset:
	unset VAULT_TOKEN

#Build
build:
	env GOOS=linux   GOARCH=amd64 go build -ldflags '-s -w' -o $(PLUGIN_DIR)/linux/$(PLUGIN_NAME) || exit 1
	env GOOS=linux   GOARCH=386   go build -ldflags '-s -w' -o $(PLUGIN_DIR)/linux86/$(PLUGIN_NAME) || exit 1
	env GOOS=darwin  GOARCH=amd64 go build -ldflags '-s -w' -o $(PLUGIN_DIR)/darwin/$(PLUGIN_NAME) || exit 1
	env GOOS=darwin  GOARCH=386   go build -ldflags '-s -w' -o $(PLUGIN_DIR)/darwin86/$(PLUGIN_NAME) || exit 1
	env GOOS=windows GOARCH=amd64 go build -ldflags '-s -w' -o $(PLUGIN_DIR)/windows/$(PLUGIN_NAME).exe || exit 1
	env GOOS=windows GOARCH=386   go build -ldflags '-s -w' -o $(PLUGIN_DIR)/windows86/$(PLUGIN_NAME).exe || exit 1
	chmod +x $(PLUGIN_DIR)/*

compress:
	mkdir -p $(DIST_DIR)
	rm -f $(DIST_DIR)/*
	zip -j "${CURRENT_DIR}/$(DIST_DIR)/${PLUGIN_NAME}_${VERSION}_linux.zip" "$(PLUGIN_DIR)/linux/$(PLUGIN_NAME)" || exit 1
	zip -j "${CURRENT_DIR}/$(DIST_DIR)/${PLUGIN_NAME}_${VERSION}_linux86.zip" "$(PLUGIN_DIR)/linux86/$(PLUGIN_NAME)" || exit 1
	zip -j "${CURRENT_DIR}/$(DIST_DIR)/${PLUGIN_NAME}_${VERSION}_darwin.zip" "$(PLUGIN_DIR)/darwin/$(PLUGIN_NAME)" || exit 1
	zip -j "${CURRENT_DIR}/$(DIST_DIR)/${PLUGIN_NAME}_${VERSION}_darwin86.zip" "$(PLUGIN_DIR)/darwin86/$(PLUGIN_NAME)" || exit 1
	zip -j "${CURRENT_DIR}/$(DIST_DIR)/${PLUGIN_NAME}_${VERSION}_windows.zip" "$(PLUGIN_DIR)/windows/$(PLUGIN_NAME).exe" || exit 1
	zip -j "${CURRENT_DIR}/$(DIST_DIR)/${PLUGIN_NAME}_${VERSION}_windows86.zip" "$(PLUGIN_DIR)/windows86/$(PLUGIN_NAME).exe" || exit 1

#Developement server tasks
dev_server: unset
	pkill vault || echo "Vault server is not running"
	vault server -log-level=debug -dev -config=scripts/config/vault/vault-config.hcl

dev: build_dev mount_dev

build_dev:
	go build -o $(PLUGIN_PATH) || exit 1
	chmod +x $(PLUGIN_PATH)

mount_dev: unset
	vault write sys/plugins/catalog/$(PLUGIN_NAME) sha_256="$(SHA256)" command="$(PLUGIN_NAME)"
	vault secrets disable $(MOUNT) || echo "Secrets already disabled"
	vault secrets enable -path=$(MOUNT) -plugin-name=$(PLUGIN_NAME) plugin

dev_test:
	vault write $(MOUNT)/user/user1 password=MySecret12
	vault write $(MOUNT)/user/user2 generate=true
	vault write $(MOUNT)/user/user3 generate=false || echo "Error"
	vault read $(MOUNT)/user/user1
	vault read $(MOUNT)/user/user2
	@echo "Listing users:"
	vault list $(MOUNT)/users

#Production server tasks
prod_server_up:
	docker-compose up -d
	@echo "Run: docker-compose logs"
	@echo "to see the logs"
	@echo "Run: docker exec -it $(VAULT_CONT) sh"
	@echo "to login into vault container"
	@echo "Waiting until server start"
	sleep 4


prod_server_init:
	$(VAULT_CMD) operator init -key-shares=1 -key-threshold=1
	@echo "To unseal the vault run:"
	@echo "$(VAULT_CMD) operator unseal UNSEAL_KEY"

prod_server_unseal:
	@echo Enter unseal key:
	$(VAULT_CMD) operator unseal

prod_server_auth:
	@echo Enter root token:
	$(VAULT_CMD) login

prod_server_down:
	docker-compose down --remove-orphans

prod_server_logs:
	docker-compose logs -f

prod_server_sh:
	$(DOCKER_CMD) sh

prod: prod_server_down prod_server_up prod_server_init prod_server_unseal prod_server_auth mount_prod
	@echo "Vault started. To run make command export VAULT_TOKEN variable and run make with -e flag, for example:"
	@echo "export VAULT_TOKEN=enter-root-token-here"
	@echo "make cloud -e"

mount_prod:
	$(eval SHA256 := $(shell echo $$($(DOCKER_CMD) $(SHA256_DOCKER_CMD))))
	echo $$SHA256
	$(VAULT_CMD) write sys/plugins/catalog/$(PLUGIN_NAME) sha_256="$$SHA256" command="$(PLUGIN_NAME)"
	$(VAULT_CMD) secrets disable $(MOUNT) || echo "Secrets already disabled"
	$(VAULT_CMD) secrets enable -path=$(MOUNT) -plugin-name=$(PLUGIN_NAME) plugin