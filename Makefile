generate_from_chord_proto:
	protoc services/chord_service/*.proto --go_out=plugins=grpc:.
