//go:generate protoc --go_out=. --go_opt=paths=source_relative events-metadata.proto
//go:generate protoc --go_out=. --go_opt=paths=source_relative events-workflow-started.proto
//go:generate protoc --go_out=. --go_opt=paths=source_relative events-workflow-finished.proto
//go:generate protoc --go_out=. --go_opt=paths=source_relative events-capability-started.proto
//go:generate protoc --go_out=. --go_opt=paths=source_relative events-capability-finished.proto
package pb
