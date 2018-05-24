pkg_name=chainlink
pkg_origin=smartcontract
pkg_version=0.1.0
pkg_license=("MIT")
pkg_scaffolding=core/scaffolding-go
pkg_build_deps=(core/dep)

scaffolding_go_base_path=github.com/smartcontractkit

# Use dep to manage dependencies instead of `go get`
do_download() {
  build_line "Downloading Go build dependencies"
  pushd "$scaffolding_go_pkg_path" >/dev/null
  dep ensure
  popd >/dev/null
}

# Build with extra LDFLAGS
do_build() {
  export PATH="$PATH:$GOPATH/bin"
  pushd "$scaffolding_go_pkg_path" >/dev/null
	go build -ldflags "-X github.com/smartcontractkit/chainlink/store.Sha=`git rev-parse HEAD`" -o chainlink
  popd >/dev/null
}
