run:
	go run main.go

test:
	go test ./t

clean:
	rm -r virtual/*

.PHONY: run test clean
