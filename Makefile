.PHONY: go_of_proto client
client:
	cd client && go build

go_of_proto:
	./go_of_proto.sh # putting the content of go_of_proto.sh directly here messed with $GOPATH