run: build
	reset
	./pr0n

ai: build
	reset
	./pr0n --ai

build:
	go build

serve: build
	./pr0n --serve
