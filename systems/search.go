package systems

import (
	"encoding/json"
	"log"
	"os"
)

// Define our nodes and tree root.
var root = makeTreeNode()

type treeNode struct {
	Children map[byte]treeNode
	Values   []int32
}

// BuildNameSearchTree reads the input file and builds a search tree with the name.
func BuildNameSearchTree() {
	filename := os.Args[1] // todo: Handle errors.

	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Failed to open file: %s", err)
	}

	decoder := json.NewDecoder(file)

	log.Println("Starting to build system tree.")

	// read open bracket
	_, err = decoder.Token()
	if err != nil {
		log.Fatal(err)
	}

	// Read the systems in the list.
	for decoder.More() {
		var system EDSMSystem
		err := decoder.Decode(&system)
		if err != nil {
			log.Fatal(err)
		}

		addSystem(system)
	}

	// read closing bracket
	_, err = decoder.Token()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("System tree done.")
}

// Helper functions
func addSystem(system EDSMSystem) {
	nameLength := len(system.Name)

	node := root
	for i := 0; i < nameLength; i++ {
		char := system.Name[i]

		if _, ok := node.Children[char]; !ok {
			node.Children[char] = makeTreeNode()
		}
		node = node.Children[char]
	}
	node.Values = append(node.Values, system.ID)
}

func makeTreeNode() treeNode {
	return treeNode{
		Children: make(map[byte]treeNode, 0),
		Values:   make([]int32, 0),
	}
}

// SearchTree searches through the generated tree.
func SearchTree(input string) []int32 {
	inputLength := len(input)
	result := make([]int32, 0)

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
