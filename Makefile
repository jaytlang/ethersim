GOFLAGS = -race
FMTFLAGS = -s -w

TARGET = ethersim

.PHONY: all $(TARGET)

all: $(TARGET)

$(TARGET):
	go build $(GOFLAGS) -o $(TARGET)

run: $(TARGET)
	./$(TARGET)

.PHONY: clean format kill
clean:
	rm -rf $(TARGET)

kill:
	pkill $(TARGET)

format:
	gofmt $(FMTFLAGS) .
