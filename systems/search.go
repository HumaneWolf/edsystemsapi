package systems

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
	SizeBytes int64
	Nodes     int64
	NodeSize  int
}

// GetIndexStats gets some basic stats about the index.
func GetIndexStats() IndexStats {
	totalSize := getTotalSize()
	return IndexStats{
		SizeBytes: totalSize,
		Nodes:     totalSize / nodeSize,
		NodeSize:  nodeSize,
	}
}
