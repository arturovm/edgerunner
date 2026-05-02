.PHONY: app clean

app: bin/edgerunner

bin/edgerunner: $(shell find . -path '**/*.go')
	go build -o bin/edgerunner git.sr.ht/~arturovm/edgerunner

clean:
	rm -rf bin
