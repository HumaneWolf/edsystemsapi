package systems

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
