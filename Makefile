.PHONY: run build dev clean

all: build run

run:
	cd build && ./cockpit

build:
	mkdir -p build
	cd server && go build -o ../build
	cd web && npm run build

dev:
	cd web && npm run dev

clean:
	rm -r build
