.PHONY: build build-max clean

export GOOS := wasip1
export GOARCH := wasm
export CGO_ENABLED := 0

LDFLAGS := -buildid= -w -s
# 14 MB of random data pushes compressed binary to ~18 MB (under 20 MB limit)
BALLAST_KB := 14336

# Standard build — exercises all call limits, config, and logging limits.
build:
	go build -o workflow.wasm -trimpath -ldflags="$(LDFLAGS)" .

# Max-binary build — also pushes WASMCompressedBinarySizeLimit (~18 MB compressed).
build-max: ballast.dat
	go build -tags maxbinary -o workflow.wasm -trimpath -ldflags="$(LDFLAGS)" .

ballast.dat:
	dd if=/dev/urandom of=$@ bs=1024 count=$(BALLAST_KB)

clean:
	rm -f workflow.wasm ballast.dat
