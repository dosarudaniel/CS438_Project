package file_sharing

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
)

const MaxChunkSize int = 8192  // 1<<13

type Chunk []byte

func (chunk Chunk) isValidMetaFile() error {
	if len(chunk)%32 != 0 {
		return fmt.Errorf("invalid chunk length(not a multiple of 32 bytes)")
	}
	return chunk.isValidChunk()
}

func (chunk Chunk) isValidChunk() error {
	if len(chunk) > MaxChunkSize {
		return fmt.Errorf("invalid chunk length (more than MaxChunkSize=%v)", MaxChunkSize)
	}
	return nil
}

func NbChunk(size int) int {
	return int(math.Ceil(float64(size) / float64(MaxChunkSize)))
}

func fileToChunks(filename string) ([]Chunk, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	stats, statsErr := file.Stat()
	if statsErr != nil {
		return nil, statsErr
	}

	size := stats.Size()
	var chunks []Chunk
	if size == 0 {
		chunks = make([]Chunk, 1)
		chunks[0] = []byte{}
	} else {
		nbChunk := uint(math.Ceil(float64(size) / float64(MaxChunkSize)))
		chunks = make([]Chunk, nbChunk)

		reader := bufio.NewReader(file)
		for ; size > 0; size -= int64(MaxChunkSize) {
			chunkIndex := nbChunk - uint(math.Ceil(float64(size)/float64(MaxChunkSize)))

			if size >= int64(MaxChunkSize) {
				chunks[chunkIndex] = make(Chunk, MaxChunkSize)
				_, err = io.ReadFull(reader, chunks[chunkIndex])
				if err != nil {
					return nil, err
				}
			} else {
				chunks[chunkIndex] = make(Chunk, size)
				_, err = reader.Read(chunks[chunkIndex])
				if err != nil {
					return nil, err
				}
			}
		}
	}

	_ = file.Close()
	return chunks, nil
}
