build: web

actions = iconv -f ISO-8859-1 -t utf-8 | grep -v \' | grep -v '-' | iconv -f utf8 -t ascii//TRANSLIT//IGNORE | sed 's/[^a-zA-Z]//g' | tr '[:upper:]' '[:lower:]' | uniq | grep -E '^.{3}'
dictionaries:
	mkdir -p dictionaries
	curl -s https://web.archive.org/web/20090904013506id_/http://wortschatz.uni-leipzig.de/Papers/top10000en.txt | $(actions) > dictionaries/en.txt
	curl -s https://web.archive.org/web/20090909075401id_/http://wortschatz.uni-leipzig.de/Papers/top10000de.txt | $(actions) > dictionaries/de.txt
	curl -s https://web.archive.org/web/20090904105851id_/http://wortschatz.uni-leipzig.de/Papers/top10000fr.txt | $(actions) > dictionaries/fr.txt
	curl -s https://web.archive.org/web/20090904014314id_/http://wortschatz.uni-leipzig.de/Papers/top10000nl.txt | $(actions) > dictionaries/nl.txt

web: public
public: dictionaries
	mkdir -p public
	tinygo build -o public/web.wasm -target wasm web/web.go
	cp -rf dictionaries web/*.js web/*.html public/
	curl -s https://cdn.jsdelivr.net/npm/water.css@2/out/water.min.css > public/water.min.css

dist:
	goreleaser build --snapshot --clean

clean:
	rm -rf public/ dist/

lint:
	go vet $(shell go list ./... | grep -v /web)
	GOOS=js GOARCH=wasm go vet ./web
	go fix -diff $(shell go list ./... | grep -v /web)
	GOOS=js GOARCH=wasm go fix -diff ./web
	go tool staticcheck $(shell go list ./... | grep -v /web)
	# Run the host-built staticcheck binary directly: a GOOS/GOARCH prefix on
	# `go tool` would cross-compile staticcheck itself instead of the analysed packages.
	GOOS=js GOARCH=wasm $$(go tool -n staticcheck) ./web
	go mod tidy -diff

test:
	go test $(shell go list ./... | grep -v /web)
	GOOS=js GOARCH=wasm go test ./web
