ctrl:
	go run main.go -controller

fs:
	go run main.go -flagserver

test:
	go test ./t

clean:
	rm -r virtual/*

.PHONY: ctrl fs test clean
