default:
    just -l

# Run this before every commit to ensure tests pass and all examples are regenerated
update-all: unit-tests example0 example1 example2 example3 example4 example5 readme-sync

unit-tests:
    cd test && ./run.sh

# Updates the README code blocks between sync markers from the example0 outputs
readme-sync: example0
    ./scripts/readme-sync.sh

example0:
    go build -o examples/00-docs/genconfig && cd examples/00-docs && go run generate.go

example1:
    go build -o examples/01-direct-call/genconfig && cd examples/01-direct-call && go run generate.go

example2:
    go build -o examples/02-executable-gogenerate/genconfig && cd examples/02-executable-gogenerate && go generate ./...

example3:
    cd examples/03-go-tool && go generate ./...

example4:
    cd examples/04-tools-dot-go && go generate ./...

example5:
    go build -o examples/05-executable-direct/genconfig && cd examples/05-executable-direct && ./run.sh
