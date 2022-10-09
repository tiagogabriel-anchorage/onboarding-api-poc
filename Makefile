GOCMD=go
TARGET=onboarding_api

build:
	$(GOCMD) build -o bin/$(TARGET)

run: kill clear build
	./bin/$(TARGET) &

clear:
	@rm -fr bin/*

clear-all: kill clear

kill:
	@if pgrep $(TARGET) 1>/dev/null; then pkill $(TARGET); fi