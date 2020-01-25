generate_from_chord_proto:
	protoc services/chord_service/*.proto --go_out=plugins=grpc:.
	protoc services/file_share_service/*.proto --go_out=plugins=grpc:.
	protoc services/client_service/*.proto --go_out=plugins=grpc:.
