
.PHONY: unit-tests

unit-tests:
	chmod +x ./test/run.sh; cd ./test; ./run.sh; cd ..
