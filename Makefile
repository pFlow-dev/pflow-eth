# add a way to run foundry locally

FORGEPATH := $(HOME)/.foundry/bin/forge

forgeTest :
	echo "Running tests..."
	$(FORGEPATH) test


test : forgeTest
