.PHONY: run build dev clean

all: build run

run:
	cd build && ./cockpit

build:
	mkdir -p build
	cd web && npm run build
	cd server && go build -ldflags "-linkmode 'external' -extldflags '-static'" -tags netgo,osusergo -o ../build .

dev:
	rm server/cockpit.db* & cd web && npm run dev & cd server && go run .

clean:
	rm -r build
	rm -r server/build
	rm server/*.db server/*.db-wal server/*.db-shm
