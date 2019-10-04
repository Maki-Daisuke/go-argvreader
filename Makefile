.PHONY: deps
deps:
	go get -d

.PHONY: devel-deps
devel-deps:
	go get github.com/mattn/go-forlines

.PHONY: test
test:
	which shelltest >/dev/null || ( echo "**** Install shelltest command to test! ****" && exit 1 )
	cd testdata  &&  sh -c 'shelltest *.test'
