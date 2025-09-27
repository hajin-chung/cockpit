.PHONY: run build dev clean

all: build run

run:
	cd build && ./cockpit

build:
	mkdir -p build
	cd server && go build -o ../build
	cd web && npm run build

dev:
	rm server/cockpit.db* & cd web && npm run dev & cd server && go run .

clean:
	rm -r build
