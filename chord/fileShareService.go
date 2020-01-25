package chord

import (
	. "github.com/dosarudaniel/CS438_Project/services/file_share_service"
)

// RPC implementation
func (chordNode *ChordNode) TransferFile(fileInfo *FileInfo, stream FileShareService_TransferFileServer) error {
	// Find the file with filename == fileInfo.Filename

	// Create an array of chunks named fileChunks
	fileChunks := make([][]byte, 5)
	for _, chunk := range fileChunks {
		fileChunk := FileChunk{Content: chunk}
		if err := stream.Send(&fileChunk); err != nil {
			return err
		}

	}
	return nil
}



// TODO: implement the client :  calls the TransferFile function on another node, waits for the chunks, reconstruct the file
//func printFeatures(client pb.RouteGuideClient, rect *pb.Rectangle) {
//	log.Printf("Looking for features within %v", rect)
//	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//	defer cancel()
//	stream, err := client.ListFeatures(ctx, rect)
//	if err != nil {
//		log.Fatalf("%v.ListFeatures(_) = _, %v", client, err)
//	}
//	for {
//		feature, err := stream.Recv()
//		if err == io.EOF {
//			break
//		}
//		if err != nil {
//			log.Fatalf("%v.ListFeatures(_) = _, %v", client, err)
//		}
//		log.Println(feature)
//	}
//}

