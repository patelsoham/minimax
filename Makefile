GO = go
BUILD = build
RUN = run

TARGET = minmax

all: build run
	
run: $(TARGET)
	./$(TARGET)

build:
	$(GO) $(BUILD) -o $(TARGET) *.go

clean:
	rm -f *.o minmax