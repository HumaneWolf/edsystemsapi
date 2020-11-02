package systems

import "encoding/binary"

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
