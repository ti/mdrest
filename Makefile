BUILD_DIR =  ./build/

.PHONY: build clean vendor

all: clean  build

build:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o $(BUILD_DIR)mdrest ./mdrest
	cp Dockerfile $(BUILD_DIR)
	cp mdrest/config.json $(BUILD_DIR)
clean:
	rm -rf $(BUILD_DIR)

