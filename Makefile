.PHONY: build
build: webui-build public-setup rice-setup go-build

.PHONY: webui-build
webui-build:
	pushd ./webui/ && npm run build && popd

.PHONY: public-setup
public-setup:
	rm -rf ./public
	mv ./webui/build ./public

.PHONY: rice-setup
rice-setup:
	if [ -x `which rice` ]; then \
		echo "found rice: `which rice`"; \
	else \
		echo "installing rice"; \
		go install github.com/GeertJohan/go.rice/rice@latest; \
	fi
	rice embed-go

.PHONY: go-build
go-build:
	$(eval MAIN_JS := $(shell basename `ls public/static/js/main.*.js`))
	$(eval MAIN_CSS := $(shell basename `ls public/static/css/main.*.css`))
	$(eval JS := $(shell echo $(MAIN_JS) | cut -d'.' -f2))
	$(eval CSS := $(shell echo $(MAIN_CSS) | cut -d'.' -f2))
	@echo "FIXME: Set Webui Version: main.$(JS).js main.$(CSS).css"
	go build -ldflags "-X 'github.com/pflow-dev/pflow-eth/config.cssBuild=$(CSS)' -X 'github.com/pflow-dev/pflow-eth/config.jsBuild=$(JS)' -s -w" -o pflow-eth main.go

.PHONY: archive
archive:
	git archive --format=zip --output=pflow-eth.$$(date -I).zip main

.PHONY: generate
generate:
	solc --overwrite --abi ./hardhat/contracts/MyStateMachine.sol -o build
	abigen --abi ./build/MyStateMachine.abi --pkg metamodel --out ./metamodel/metamodel.go 

.PHONY: restart-hardhat
restart-hardhat:
	cd docker && \
	echo 'restarting hardhat' && \
	docker-compose down hardhat && docker-compose up hardhat -d && sleep 1 && \
	cd - && \
	echo 'deploy contract' && \
	cd ./hardhat && npm run deploy && cd - && \
	echo 'reset_db' && \
	curl 'http://127.0.0.1:8083/v0/control?cmd=reset_db' && \
	curl 'http://127.0.0.1:8083/v0/control?cmd=init_block_numbers' && \
	curl 'http://127.0.0.1:8083/v0/control?cmd=sync' && \
	echo 'view model' && \
	curl 'http://127.0.0.1:8083/v0/model?addr=0x5FbDB2315678afecb367f032d93F642f64180aa3'