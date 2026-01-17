APP_NAME := $(notdir $(CURDIR))

CSS_IN  := static/css/style.css
CSS_OUT := dist/static/css/tailwind.css

.PHONY: build clean

build:
	rm -rf dist/*

	pnpm exec tailwindcss -i $(CSS_IN) -o $(CSS_OUT) -m
	cp -r  static/font dist/static/font

	go build -o bin/$(APP_NAME) .
	./bin/$(APP_NAME)

clean:
	rm -rf bin
