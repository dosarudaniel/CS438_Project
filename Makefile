generate_from_chord_proto:
	protoc services/chord_service/* --plugin=grpc --go_out=.