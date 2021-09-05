LDFLAGS := -s -w

BINARY="lmap"
RELEASE_DIR="release"
VERSION=`git describe --tags $(git rev-list --tags --max-count=1 --branches master)`

os-archs=darwin:amd64 darwin:arm64 freebsd:386:softfloat freebsd:amd64 linux:386:softfloat linux:amd64 linux:arm linux:arm64 windows:386:softfloat windows:amd64 linux:mips64 linux:mips64le linux:mips:softfloat linux:mipsle:softfloat

packages=./cmd/lmap ./pkg/lmap

default:
	@go build ./cmd/lmap

version:
	@echo $(VERSION)

all:
	@$(foreach n, $(os-archs),\
		os=$(shell echo "$(n)" | cut -d : -f 1);\
		arch=$(shell echo "$(n)" | cut -d : -f 2);\
		GO386=$(shell echo "$(n)" | cut -d : -f 3);\
		target_suffix=$${os}_$${arch};\
		echo "Build $${os}-$${arch}...";\
		env CGO_ENABLED=0 GOOS=$${os} GOARCH=$${arch} GO386="$${GO386}" go build -trimpath -ldflags "$(LDFLAGS)" -o ./$(RELEASE_DIR)/$(BINARY)_$${target_suffix} ./cmd/lmap;\
		echo "Build $${os}-$${arch} done";\
	)
	@mv ./$(RELEASE_DIR)/$(BINARY)_windows_386 ./$(RELEASE_DIR)/$(BINARY)_windows_386.exe
	@mv ./$(RELEASE_DIR)/$(BINARY)_windows_amd64 ./$(RELEASE_DIR)/$(BINARY)_windows_amd64.exe

fmt:
	@gofmt -s -w ./

fmt-check:
	@diff=`gofmt -s -d ./`; \
	if [ -n "$$diff" ]; then \
		echo "Please run 'make fmt' and commit the result:"; \
		echo "$${diff}"; \
		exit 1; \
	fi;

install:
	@go mod tidy

test:
	@go test -cpu=1,2,4 -v -tags integration ./...

vet:
	@$(foreach n, $(packages),\
		go vet ${n} ;\
	)

clean:
	@if [ -f ${BINARY} ] ; then rm -r ${BINARY} ; fi
	@if [ -d ${RELEASE_DIR} ] ; then rm -r ${RELEASE_DIR} ; fi

.PHONY: default fmt fmt-check install test vet docker clean