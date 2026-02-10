SRC_DIRS = internal/* cmd/shift-weaver
SRC_FILES = $(foreach dir,$(SRC_DIRS),$(wildcard $(dir)/*.go))
BIN := shift-weaver
MAIN := ./cmd/shift-weaver/main.go

print:
	@echo $(SRC_FILES)

all: $(BIN)

$(BIN): $(SRC_FILES)
	go build -o $(BIN) $(MAIN)

clean:
	rm -f shift-weaver

update:
	go mod tidy
