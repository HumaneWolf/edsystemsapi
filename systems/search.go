package systems

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

// Define our nodes and tree root.
var root = makeTreeNode()

type treeNode struct {
	Children map[byte]treeNode
	Values   []int64
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

	counter := 0
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
			log.Fatalf("Failed to parse system ID: %s", err)
		}
		system := SystemLine{ID64: id, Name: record[nameIndex]}
		addSystem(system)

		counter++
		// log.Printf("Added system %s (%d)\n", system.Name, system.ID64)
	}

	dumpFile, _ := json.MarshalIndent(root, "", " ")
	_ = ioutil.WriteFile("dump.json", dumpFile, 0644)

	log.Printf("System tree done. %d systems added.\n", counter)
}

// Helper functions
func addSystem(system SystemLine) {
	nameLength := len(system.Name)

	node := root
	for i := 0; i < nameLength; i++ {
		char := system.Name[i]

		if _, ok := node.Children[char]; !ok {
			node.Children[char] = makeTreeNode()
		}
		node = node.Children[char]
	}

	node.Values = append(node.Values, system.ID64)
}

func makeTreeNode() treeNode {
	return treeNode{
		Children: make(map[byte]treeNode, 0),
		Values:   make([]int64, 0),
	}
}

// SearchTree searches through the generated tree.
func SearchTree(input string) []int64 {
	inputLength := len(input)
	result := make([]int64, 0)

	// Traverse the tree
	node := root
	for i := 0; i < inputLength; i++ {
		char := input[i]

		if val, ok := node.Children[char]; ok {
			node = val
		} else {
			return result
		}
	}

	// Add exact matches
	for i := 0; i < len(node.Values); i++ {
		result = append(result, node.Values[i])
	}

	return result
}
