GO_BUILDOPT := -ldflags '-s -w'

gom:
	@if [`which go` = "go not found"]; then \
		echo "Please install Golang"; \
	else \
		go get github.com/mattn/gom; \
		gom install; \
	fi

link:
	mkdir -p $(GOPATH)/src/github.com/joniyjoniy
	ln -si $(CURDIR) $(GOPATH)/src/github.com/joniyjoniy/mouryou-dog

fmt:
	gom exec goimports -w *.go lib/*.go
	gom exec goimports -w sample-server/*.go

build: fmt
	gom build $(GO_BUILDOPT) -o bin/mouryou-dog main.go
	gom build $(GO_BUILDOPT) -o bin/sample-server sample-server/sample-server.go
