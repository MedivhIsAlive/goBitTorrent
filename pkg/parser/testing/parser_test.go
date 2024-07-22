package testing

import (
	"fmt"
	"os"
	"testing"
	"torrent/pkg/assert"
	"torrent/pkg/bencode"
	"torrent/pkg/parser"
)

type Torrent struct {
	Announce string         `bencode:"announce"`
	Info     map[string]any `bencode:"info"`
}

func TestParser(t *testing.T) {
	wd, _ := os.Getwd()
	fmt.Println(wd)
	data, err := os.ReadFile("../../parser/testing/files/test_file.torrent")
	assert.NoError(err, "failed while reading file")
	torrent := Torrent{}
	assert.NoError(bencode.Unmarshal(data, &torrent), "failed while parsing torrent")
	fmt.Println(torrent.Announce)
}

func TestJopa(t *testing.T) {
	data, err := os.Open("../../parser/testing/files/test_file.torrent")
	assert.NoError(err, "")
	torrent := parser.ParseTorrentFile(data)
	fmt.Println(torrent.Announce)
}
