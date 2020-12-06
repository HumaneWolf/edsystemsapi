package systems

import (
	"log"
	"sync"
	"time"

	"humanewolf.com/ed/systemapi/config"
)

var (
	nodeCache     = make(map[int64]cachedNode)
	nodeCacheLock = &sync.RWMutex{}
)

type cachedNode struct {
	node         treeNode
	lastAccessed int64
}

func saveCachedNode(offset int64, node treeNode) {
	nodeCacheLock.RLock()
	defer nodeCacheLock.RUnlock()

	nodeCache[offset] = cachedNode{
		node:         node,
		lastAccessed: time.Now().Unix(),
	}
}

func getCachedNode(offset int64) *treeNode {
	nodeCacheLock.RLock()
	defer nodeCacheLock.RUnlock()

	if entry, exists := nodeCache[offset]; exists {
		entry.lastAccessed = time.Now().Unix() // Random overwrites here are fine, doesn't matter.
		return &entry.node
	}
	return nil
}

func cleanCache() int {
	nodeCacheLock.Lock()
	defer nodeCacheLock.Unlock()

	cfg := config.LoadConfig()

	deletedEntries := 0
	now := time.Now().Unix()
	for offset, entry := range nodeCache {
		if (now - entry.lastAccessed) >= int64(cfg.FileStore.MemoryCacheMaxAge) {
			delete(nodeCache, offset)
			deletedEntries++
		}
	}

	return deletedEntries
}

// StartCacheCleaner starts a ticker which will regularly clean the in memory node cache. Should be ran in own goroutine.
func StartCacheCleaner() {
	ticker := time.NewTicker(30 * time.Second)

	for {
		_ = <-ticker.C

		start := time.Now().UnixNano()
		deletedEntries := cleanCache()
		nanosSpent := time.Now().UnixNano() - start
		log.Printf("Cleaned node cache. (entries deleted=%d, time spent=%d ms).\n", deletedEntries, (nanosSpent / 1_000_000))
	}
}
