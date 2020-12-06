package systems

import (
	"unicode"
)

func byteEqualsIgnoreCase(char1 byte, char2 byte) bool {
	return unicode.ToUpper(rune(char1)) == unicode.ToUpper(rune(char2))
}

func findCharacterNode(startOffset int64, character byte) *treeNode {
	offset := startOffset

	for {
		node := readNode(offset)

		if byteEqualsIgnoreCase(node.Character, character) {
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

	remaining := 40

	offset := int64(0)
	var node *treeNode
	var matchFullInput = false
	var correctlyCasedInput = ""
	for i := 0; i < len(input); i++ {
		char := input[i]
		node = findCharacterNode(offset, char)

		if i == (len(input) - 1) {
			matchFullInput = true
		}

		if node == nil {
			break
		} else if !matchFullInput && node.ChildOffset == -1 {
			correctlyCasedInput += string(node.Character)
			break
		} else {
			correctlyCasedInput += string(node.Character)
			offset = node.ChildOffset
		}
	}

	// Add exact match, if any
	if matchFullInput && node != nil && node.SystemCount != 0 {
		result = append(result, input)
	}
	remaining -= len(result)

	// Time to find systems which start with the given input, for autocomplete purposes. Right now we'll just return all of them, might want to set max limit.
	if node != nil && node.ChildOffset != -1 {
		result = append(result, returnChildrenNames(node.ChildOffset, correctlyCasedInput, remaining)...)
	}

	if len(result) > 40 {
		result = result[0:40]
	}

	return result
}

func returnChildrenNames(offset int64, name string, remaining int) []string {
	results := make([]string, 0)

	if remaining <= 0 {
		return results
	}
	// If we don't have any remaining slots to list, we don't need to read the node, so we do it after that check.
	node := readNode(offset)

	// If this node is a system, add it.
	if node.SystemCount != 0 {
		results = append(results, name+string(node.Character))
	}
	remaining -= len(results)

	if node.NextNodeOffset != -1 {
		r := returnChildrenNames(node.NextNodeOffset, name, remaining)
		results = append(results, r...)
	}
	remaining -= len(results)

	if node.ChildOffset != -1 {
		r := returnChildrenNames(node.ChildOffset, name+string(node.Character), remaining)
		results = append(results, r...)
	}
	remaining -= len(results)

	return results
}
