.PHONY: build clean

build: bin/edgerunner

bin/edgerunner: $(shell find . -name '*.go')
	go build -o bin/edgerunner git.sr.ht/~arturovm/edgerunner

clean:
	rm -rf bin
