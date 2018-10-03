package api

import (
	"fmt"
	"testing"
)

func TestSongList(t *testing.T) {
	list, err := GetSongList("99686749")
	fmt.Println(list, err)
}
