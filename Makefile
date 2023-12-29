TESTS?=$$(go list ./... | egrep -v "vendor|resources")

test:
	go test -count=1 -failfast -v $(TESTS)