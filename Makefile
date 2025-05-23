
.PHONY: unit-tests example1 example2 example3 example4 example5

unit-tests:
	chmod +x ./test/run.sh; cd ./test; ./run.sh; cd ..

example0:
	go build -o examples/00-docs/genconfig; cd examples/00-docs; go run generate.go; cd ../..

example1:
	go build -o examples/01-direct-call/genconfig; cd examples/01-direct-call; go run generate.go; cd ../..

example2:
	cd examples/02-executable-gogenerate; go generate ./...; cd ../..

example3:
	cd examples/03-go-tool; go generate ./...; cd ../..

example4:
	cd examples/04-tools-dot-go; go generate ./...; cd ../..

example5:
	go build -o examples/05-executable-direct/genconfig; cd examples/05-executable-direct; chmod +x run.sh; ./run.sh; cd ../..

update-all: unit-tests example0 example1 example2 example3 example4 example5
