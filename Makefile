BUILD_DIR =  ./build/

.PHONY: build clean vendor

all: clean  build

build:
	go build -o $(BUILD_DIR)mdrest ./mdrest
	cp mdrest/config.json $(BUILD_DIR)
clean:
	rm -rf $(BUILD_DIR)

