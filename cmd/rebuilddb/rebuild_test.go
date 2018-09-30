package rebuilddb

import "testing"

func TestBuildDatabase(t *testing.T) {
	BuildDatabase()
}

func BenchmarkBuildDatabase(b *testing.B) {
	BuildDatabase()
}
