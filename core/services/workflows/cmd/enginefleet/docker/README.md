# enginefleet Docker image

Packages the [`enginefleet`](../main.go) binary — which runs N instances of the
v2 workflow engine against a shared capabilities registry of in-process fakes
(HTTP action, consensus, chain write) — into a container, plus a script to run
it and measure the container's memory usage over the test.

## Build

```bash
./build.sh                 # builds enginefleet:latest
./build.sh -t my/tag:v1    # custom tag
```

`build.sh` compiles a Linux binary inside a `golang:<go.mod version>-bookworm`
container, reusing the host Go module cache read-only (the monorepo builds fully
offline, so no network or private-module auth is required), then copies the
binary into the slim runtime image. The WASM host uses `wasmtime` (CGO), so the
build runs with `CGO_ENABLED=1`. A named Docker volume caches Go build output
between runs.

## Run + measure memory

```bash
./run.sh -w /path/to/workflow.wasm -c /path/to/config.yaml -n 5 -t 120
```

Key flags (`-h` for all):

| Flag | Meaning |
|------|---------|
| `-w` | Path to the workflow WASM binary (required, **uncompressed** — see note) |
| `-c` | Path to the workflow config file |
| `-n` | Number of engine instances (default 1) |
| `-t` | Stop after N seconds (default: run until container exit / Ctrl-C) |
| `-i` | Sampling interval in seconds (default 1) |
| `-o` | CSV output path (default `./mem-<timestamp>.csv`) |
| `-m` | Container memory limit (e.g. `512m`, `2g`) |
| `-T` | Image tag to run (default `enginefleet:latest`) |
| `-d` | Enable engine debug logging |

The workflow binary and config are bind-mounted read-only into the container.
`run.sh` samples `docker stats` at the chosen interval into a CSV timeline, and on
exit prints the sampled peak/average plus — when available — the exact
kernel-tracked peak from the container cgroup's `memory.peak` (cgroup v2). The
cgroup peak is more accurate than sampling, which can miss spikes between
samples.

## Note: uncompressed WASM

The engine treats `-w` as uncompressed WASM. The example artifacts under
`../../cre/examples/v2/*/testdata/output.wasm.br` are brotli-compressed and must
be decompressed before use.
