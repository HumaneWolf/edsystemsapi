package systems

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"runtime"
	"strconv"
)

// SystemLine represents line in the systems csv file.
type systemLine struct {
	ID64 int64
	Name string

	// X float64
	// Y float64
	// Z float64
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
		system := systemLine{ID64: id, Name: record[nameIndex]}
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
func addSystem(system systemLine) {
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
