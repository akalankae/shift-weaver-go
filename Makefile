all:
	go build ./cmd/shift-weaver

clean:
	rm -f shift-weaver

update:
	go mod tidy
