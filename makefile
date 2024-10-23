.PHONY: runtest
runtest:
	gotestsum --format=short -- -v ./... | grep -v "\[no test files\]"

