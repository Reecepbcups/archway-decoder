BIN=arch-decode

install:
	go install .

build:
	go build .

test: install
	${BIN} tx decode-file ./test/test.json ./test/output.json