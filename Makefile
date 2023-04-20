.PHONY: gomodtidy
gomodtidy:
	go mod tidy
	cd ./ops && go mod tidy

.PHONY: godoc
godoc:
	go install golang.org/x/tools/cmd/godoc@latest
	# http://localhost:6060/pkg/github.com/smartcontractkit/chainlink-relay/
	godoc -http=:6060
