module github.com/dosarudaniel/CS438_Project

go 1.13

//    go mod init creates a new module, initializing the go.mod file that describes it.
//    go build, go test, and other package-building commands add new dependencies to go.mod as needed.
//    go list -m all prints the current module’s dependencies.
//    go get changes the required version of a dependency (or adds a new dependency).
//    go mod tidy removes unused dependencies.

require (
	github.com/deckarep/golang-set v1.7.1 // indirect
	github.com/gin-gonic/gin v1.5.0
	github.com/go-playground/universal-translator v0.17.0 // indirect
	github.com/golang/protobuf v1.3.3
	github.com/json-iterator/go v1.1.9 // indirect
	github.com/leodido/go-urn v1.2.0 // indirect
	github.com/masatana/go-textdistance v0.0.0-20191005053614-738b0edac985
	github.com/mattn/go-isatty v0.0.12 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/sirupsen/logrus v1.4.2
	golang.org/x/net v0.0.0-20200114155413-6afb5195e5aa // indirect
	golang.org/x/sys v0.1.0 // indirect
	google.golang.org/genproto v0.0.0-20200117163144-32f20d992d24 // indirect
	google.golang.org/grpc v1.26.0
	gopkg.in/go-playground/validator.v9 v9.31.0 // indirect
	gopkg.in/yaml.v2 v2.2.8 // indirect
)
