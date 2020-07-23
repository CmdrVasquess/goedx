.PHONY: build depgraph.svg

build:
	go generate ./...
	go build

deps: depgraph.svg

depgraph.svg:
	go mod graph | modgraphviz | dot -Tsvg -o $@
