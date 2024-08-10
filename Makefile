run: main.go
	go run main.go

test: main_test.go
	go test

clean:
	rm -r fileserver/*

.PHONY: run test clean
