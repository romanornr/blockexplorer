package rebuilddb

import (
	"testing"
)

func TestBuild(t *testing.T) {
	BuildDatabaseBlocks()
}

func BenchmarkBuildDatabaseBlocks(b *testing.B) {
	BuildDatabaseBlocks()
}