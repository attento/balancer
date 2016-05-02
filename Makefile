.PHONY: tests clean deps compile

deps:
	go get -t ./...

tests:
	go test -v ./...

clean:
	rm -rf dist/

compile:
	mkdir -p dist/
	cd dist/; \
		env GOOS=darwin GOARCH=386 go build -o balancer_$(VERSION)_darwin_386 .. \
	    env GOOS=linux GOARCH=arm go build -o balancer_$(VERSION)_linux_arm .. \
		env GOOS=linux GOARCH=arm64 go build -o balancer_$(VERSION)_linux_arm64 .. \
		env GOOS=linux GOARCH=386 go build -o balancer_$(VERSION)_linux_386 ..

release: deps clean compile
build: deps tests
