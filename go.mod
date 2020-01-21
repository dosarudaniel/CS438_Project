module github.com/dosarudaniel/CS438_Project

go 1.13

//    go mod init creates a new module, initializing the go.mod file that describes it.
//    go build, go test, and other package-building commands add new dependencies to go.mod as needed.
//    go list -m all prints the current module’s dependencies.
//    go get changes the required version of a dependency (or adds a new dependency).
//    go mod tidy removes unused dependencies.

require github.com/golang/protobuf v1.3.2
