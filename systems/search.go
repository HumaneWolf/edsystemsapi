package systems

import (
	"bytes"
	"encoding/csv"
	"io"
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
	}

	log.Printf("System tree done. %d systems added. %d skipped for missing ID64\n", counter, noID64)
}

// Helper functions
func addSystem(system SystemLine) {
	nameLength := len(system.Name)

	parent := root
	var nodeChar byte
	node := root
	for i := 0; i < nameLength; i++ {
		char := system.Name[i]

		if _, ok := node.Children[char]; !ok {
			node.Children[char] = makeTreeNode()
		}
		parent = node
		nodeChar = char
		node = node.Children[char]
	}

	node.Values = append(node.Values, system.ID64)
	parent.Children[nodeChar] = node
}

func makeTreeNode() treeNode {
	return treeNode{
		Children: make(map[byte]treeNode, 0),
		Values:   make([]int64, 0),
	}
}

// SearchTreeForIDs searches through the generated tree and returns a list of potential matches ID's..
func SearchTreeForIDs(input string) []int64 {
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
	result = append(result, node.Values...)

	// Time to find systems which start with the given input, for autocomplete purposes. Right now we'll just return all of them, might want to set max limit.
	for _, v := range node.Children {
		result = append(result, returnChildrenValues(v)...)
	}

	return result
}

func returnChildrenValues(node treeNode) []int64 {
	// This is currently depth-first, a width-first search might be better ofr our use case.
	results := make([]int64, 0)
	results = append(results, node.Values...)

	for _, v := range node.Children {
		results = append(results, returnChildrenValues(v)...)
	}
	return results
}

// SearchTreeForNames searches through the generated tree and returns a list of potential match names.
func SearchTreeForNames(input string) []string {
	inputLength := len(input)
	result := make([]string, 0)

	// Traverse the tree
	nameBuffer := bytes.NewBufferString(input)
	node := root
	for i := 0; i < inputLength; i++ {
		char := input[i]

		if val, ok := node.Children[char]; ok {
			node = val
		} else {
			return result
		}
	}

	// Add exact match, if any
	if len(node.Values) != 0 {
		result = append(result, input)
	}

	// Time to find systems which start with the given input, for autocomplete purposes. Right now we'll just return all of them, might want to set max limit.
	for k, v := range node.Children {
		tempNameBuffer := bytes.NewBuffer(nameBuffer.Bytes())
		tempNameBuffer.WriteByte(k)
		result = append(result, returnChildrenNames(v, *tempNameBuffer)...)
	}

	return result
}

func returnChildrenNames(node treeNode, nameBuffer bytes.Buffer) []string {
	// This is currently depth-first, a width-first search might be better for our use case.
	results := make([]string, 0)
	if len(node.Values) != 0 {
		results = append(results, nameBuffer.String())
	}

	for k, v := range node.Children {
		tempNameBuffer := bytes.NewBuffer(nameBuffer.Bytes())
		tempNameBuffer.WriteByte(k)
		results = append(results, returnChildrenNames(v, *tempNameBuffer)...)
	}
	return results
}
