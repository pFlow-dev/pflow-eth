# add a way to run foundry locally
FORGE_BIN := $(HOME)/.foundry/bin/forge
NPX_BIN := $(HOME)/.nvm/versions/node/v20.5.1/bin/npx

forgeTest :
	echo "Running tests..."
	$(FORGE_BIN) test -vvvv --gas-report


test : forgeTest

flatten:
	$(FORGE_BIN) flatten ./src/TicTacToe.sol


remix:
	open https://remix.ethereum.org
	$(NPX_BIN) @remix-project/remixd -s ./src -u https://remix.ethereum.org
