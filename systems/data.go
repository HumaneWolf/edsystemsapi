package systems

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"humanewolf.com/ed/systemapi/config"
)

var fileHandles = make(map[int]*os.File, 0)

const (
	characterOffset   = 0
	childOffsetOffset = characterOffset + 1
	nextOffsetOffset  = childOffsetOffset + 8
	systemCountOffset = nextOffsetOffset + 8

	nodeSize = systemCountOffset + 4
)

// Split the dataset into files and load the appropriate one based on the offset.
// offset / (100_000_000*21) = offst / (100 million * node size), adds up to close enough to the int32 limit.

type treeNode struct {
	Character      byte  // 1 byte
	ChildOffset    int64 // 8 bytes
	NextNodeOffset int64 // 8 bytes
	SystemCount    int32 // 4 bytes
}

func getFileAndOffset(offset int64) (*os.File, int64, int) {
	cfg := config.LoadConfig()

	fileMaxSize := cfg.FileStore.SystemsPerFile * nodeSize
	fileNumber := int(offset / int64(fileMaxSize))
	internalOffset := offset % int64(fileMaxSize)

	if file, exists := fileHandles[fileNumber]; exists {
		return file, internalOffset, fileNumber
	}

	fileName := fmt.Sprintf("%s/typeahead/index.%d.dat", cfg.FileStore.Path, fileNumber)

	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0775)
	if err != nil {
		log.Fatalf("Failed to open referenced index file: %s\n", err)
	}

	fileHandles[fileNumber] = file
	return file, internalOffset, fileNumber
}

func getTotalSize() int64 {
	cfg := config.LoadConfig()
	totalSize := int64(0)

	filepath.Walk(fmt.Sprintf("%s/typeahead/", cfg.FileStore.Path), func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// todo: Handle if someone for some reason has put unrelated files in this folder.
		totalSize += info.Size()
		return nil
	})

	return totalSize
}

func readNode(offset int64) treeNode {
	rawData := make([]byte, nodeSize)

	file, internalOffset, fileNumber := getFileAndOffset(offset)
	_, err := file.ReadAt(rawData, internalOffset)
	if err != nil {
		log.Fatalf("Failed to read node from index file (file number: %d, internal offset: %d): %s\n", fileNumber, internalOffset, err)
	}

	char := rawData[characterOffset:childOffsetOffset]
	childOffset, _ := binary.Varint(rawData[childOffsetOffset:nextOffsetOffset])
	nextOffset, _ := binary.Varint(rawData[nextOffsetOffset:systemCountOffset])
	systemCount, _ := binary.Varint(rawData[systemCountOffset:])

	return treeNode{
		Character:      char[0],
		ChildOffset:    childOffset,
		NextNodeOffset: nextOffset,
		SystemCount:    int32(systemCount),
	}
}

func updateNode(offset int64, node treeNode) {
	rawData := make([]byte, 0)

	rawData = append(rawData, node.Character)

	childOffsetBuffer := make([]byte, 8)
	binary.PutVarint(childOffsetBuffer, node.ChildOffset)
	rawData = append(rawData, childOffsetBuffer...)

	nextNodeBuffer := make([]byte, 8)
	binary.PutVarint(nextNodeBuffer, node.NextNodeOffset)
	rawData = append(rawData, nextNodeBuffer...)

	systemCountBuffer := make([]byte, 4)
	binary.PutVarint(systemCountBuffer, int64(node.SystemCount))
	rawData = append(rawData, systemCountBuffer...)

	for i := 0; i < len(rawData); i++ {
		file, internalOffset, fileNumber := getFileAndOffset(offset)
		_, err := file.WriteAt(rawData, internalOffset)
		if err != nil {
			log.Fatalf("Failed to update node in index file (file number: %d, internal offset: %d): %s\n", fileNumber, internalOffset, err)
		}
	}
}

func appendNode(node treeNode) int64 {
	offset := getTotalSize()
	updateNode(offset, node)
	return offset
}

// CloseFiles closes the handles of active files.
func CloseFiles() {
	for _, file := range fileHandles {
		file.Close()
	}
}
