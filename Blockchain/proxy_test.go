package Blockchain

import (
	"testing"
	"fmt"
)

func TestBlockListProxy_FindBlock(t *testing.T) {
	//var b []string
	//b[0] = "7539b2ae01fd492adcc16c2dd8747c1562a702f9057560fee9ca647b67b729e2"
	b2 := "ded7508b6b6452bfc99961366e3206a7a258cf897d3148b46e590bbf6f23f3d9"

	proxy := BlockListProxy{
		Database: &BlockList{},
	}

	fmt.Println(proxy.FindBlock(b2))

}