package systems

import (
	"encoding/binary"
	"encoding/csv"
	"io"
	"log"
	"os"
	"runtime"
	"strconv"
)

var data = make([]byte, 0)

const (
	characterOffset   = 0
	childOffsetOffset = characterOffset + 1
	nextOffsetOffset  = childOffsetOffset + 8
	systemCountOffset = nextOffsetOffset + 8

	nodeSize = systemCountOffset + 4
)

// todo: For when there are too many systems for this method, split the dataset into files and load the appropriate one based on the offset.
// offset / (100_000_000*21) = offst / (100 million * node size), adds up to close enough to the int32 limit.
// When that time comes: Update readnode, updatenode, and appendnode to work with new structure.
// Just filling one file then the next should be good enough, we'll probably read one file a lot, then the next, etc.

type treeNode struct {
	Character      byte  // 1 byte
	ChildOffset    int64 // 8 bytes
	NextNodeOffset int64 // 8 bytes
	SystemCount    int32 // 4 bytes
}

func readNode(offset int64) treeNode {
	rawData := data[offset : offset+nodeSize]

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
		// todo: When we can deal with int64, we have to update this.
		data[int(offset)+i] = rawData[i]
	}
}

func appendNode(node treeNode) int64 {
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

	offset := len(data)
	data = append(data, rawData...)

	return int64(offset)
}

func findOrCreateCharacterNode(startOffset int64, character byte) (int64, treeNode) {
	offset := startOffset

	for {
		node := readNode(offset)

		if node.Character == character {
			return offset, node
		} else if node.NextNodeOffset == -1 {
			newNode := treeNode{
				Character:      character,
				ChildOffset:    -1, // No child
				NextNodeOffset: -1, // No next node
				SystemCount:    0,
			}
			node.NextNodeOffset = appendNode(newNode) // Add new node
			updateNode(offset, node)                  // Update current node
			return node.NextNodeOffset, newNode
		} else {
			offset = node.NextNodeOffset
		}
	}
}

func findCharacterNode(startOffset int64, character byte) *treeNode {
	offset := startOffset

	for {
		node := readNode(offset)

		if node.Character == character {
			return &node
		} else if node.NextNodeOffset == -1 {
			return nil
		} else {
			offset = node.NextNodeOffset
		}
	}
}

// BuildNameSearchTree reads the input file and builds a search tree with the name.
func BuildNameSearchTree() {
	filename := os.Args[1] // todo: Handle errors if it doesn't exist.

	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Failed to open file: %s", err)
	}

	reader := csv.NewReader(file)

	var idIndex int
	var nameIndex int

	header, err := reader.Read()
	if err != nil {
		log.Fatalf("Failed to read header record: %s", err)
	}
	for i := 0; i < len(header); i++ {
		switch header[i] {
		case "ed_system_address":
			idIndex = i
		case "name":
			nameIndex = i
		}
	}

	log.Printf("System list header read, name=%d, id=%d.\n", nameIndex, idIndex)

	// Add a start node to simplify the adding of systems.
	appendNode(treeNode{
		Character:      'a',
		ChildOffset:    -1, // No child
		NextNodeOffset: -1, // No next node
		SystemCount:    0,
	})

	counter := 0
	noID64 := 0
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Failed to read record: %s", err)
		}

		id, err := strconv.ParseInt(record[idIndex], 10, 64)
		if err != nil {
			// todo: Handle systems without an ID64?
			// log.Printf("Failed to parse system ID: %s\n", err)
			noID64++
			continue
		}
		system := SystemLine{ID64: id, Name: record[nameIndex]}
		addSystem(system)

		counter++

		if counter%1_000_000 == 0 {
			log.Printf("Tree progress: %d systems added.\n", counter)
		}
	}

	log.Printf("System tree done. %d systems added. %d skipped for missing ID64\n", counter, noID64)
	runtime.GC() // We can use quite a lot of memory during the tree creation, let's clean that up before we do serious work.
}

// Helper functions
func addSystem(system SystemLine) {
	offset := int64(0)
	childOffset := int64(0)
	var node treeNode
	for i := 0; i < len(system.Name); i++ {
		char := system.Name[i]

		if childOffset == -1 { // If there is no child
			newNode := treeNode{
				Character:      char,
				ChildOffset:    -1, // No child
				NextNodeOffset: -1, // No next node
				SystemCount:    0,
			}
			newOffset := appendNode(newNode)
			node.ChildOffset = newOffset
			updateNode(offset, node)
			offset, node = newOffset, newNode

		} else { // if there is a child.
			offset, node = findOrCreateCharacterNode(childOffset, char)
		}
		childOffset = node.ChildOffset
	}

	node.SystemCount++
	updateNode(offset, node)
}

// SearchTreeForNames searches through the generated tree and returns a list of potential match names.
func SearchTreeForNames(input string) []string {
	result := make([]string, 0)

	offset := int64(0)
	var node *treeNode
	var matchFullInput = false
	for i := 0; i < len(input); i++ {
		char := input[i]
		node = findCharacterNode(offset, char)

		if i == (len(input) - 1) {
			matchFullInput = true
		}

		if node == nil {
			break
		} else if !matchFullInput && node.ChildOffset == -1 {
			break
		} else {
			offset = node.ChildOffset
		}
	}

	// Add exact match, if any
	if matchFullInput && node != nil && node.SystemCount != 0 {
		result = append(result, input)
	}

	// Time to find systems which start with the given input, for autocomplete purposes. Right now we'll just return all of them, might want to set max limit.
	if node != nil && node.ChildOffset != -1 {
		result = append(result, returnChildrenNames(node.ChildOffset, input)...)
	}

	return result
}

func returnChildrenNames(offset int64, name string) []string {
	// This is currently depth-first, a width-first search might be better for our use case.
	results := make([]string, 0)
	node := readNode(offset)

	// If this node is a system, add it.
	if node.SystemCount != 0 {
		results = append(results, name+string(node.Character))
	}

	if node.NextNodeOffset != -1 {
		r := returnChildrenNames(node.NextNodeOffset, name)
		results = append(results, r...)
	}

	if node.ChildOffset != -1 {
		r := returnChildrenNames(node.ChildOffset, name+string(node.Character))
		results = append(results, r...)
	}

	return results
}

// IndexStats is a struct containing basic stats about the search engine.
type IndexStats struct {
	AllocatedCapacity int
	SizeBytes         int
	Nodes             int
	NodeSize          int
}

// GetIndexStats gets some basic stats about the index.
func GetIndexStats() IndexStats {
	return IndexStats{
		AllocatedCapacity: cap(data),
		SizeBytes:         len(data),
		Nodes:             len(data) / nodeSize,
		NodeSize:          nodeSize,
	}
}
