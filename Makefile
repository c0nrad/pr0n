run: build
	reset
	./pr0n --local

ai: build
	reset
	./pr0n --ai

build:
	go build

network:
	./pr0n --host c0nrad.io --port :1337

serve: build
	./pr0n --serve
