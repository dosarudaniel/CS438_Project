generate_from_chord_proto:
	protoc services/chord_service/*.proto --plugin=grpc --go_out=.
