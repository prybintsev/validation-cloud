APP=validation-cloud
APP_EXECUTABLE="./out/$(APP)"

build:  $(shell find . -type f -name '*.go')
	mkdir -p out/
	go build -o $(APP_EXECUTABLE) ./cmd

run:
	make build
	$(APP_EXECUTABLE)