.PHONY: build build-max clean

export GOOS := wasip1
export GOARCH := wasm
export CGO_ENABLED := 0

LDFLAGS := -buildid= -w -s
# 14 MB of random data pushes compressed binary to ~18 MB (under 20 MB limit)
BALLAST_MAX_KB := 14336

# Standard build — uses a minimal 1 KB ballast placeholder.
build: ballast.dat
	go build -o workflow.wasm -trimpath -ldflags="$(LDFLAGS)" .

# Max-binary build — generates a large ballast to push WASMCompressedBinarySizeLimit.
build-max:
	dd if=/dev/urandom of=ballast.dat bs=1024 count=$(BALLAST_MAX_KB)
	go build -o workflow.wasm -trimpath -ldflags="$(LDFLAGS)" .

ballast.dat:
	dd if=/dev/urandom of=$@ bs=1024 count=1

clean:
	rm -f workflow.wasm ballast.dat
