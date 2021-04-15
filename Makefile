GO = go
BUILD = build
RUN = run

TARGET = minmax

all: run
	
run: $(TARGET)
	./$(TARGET)

build:
	$(GO) $(BUILD) -o $(TARGET) $(TARGET).go

clean:
	rm -f *.o minmax