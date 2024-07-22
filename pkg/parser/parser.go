package parser

import (
	"io"
	"torrent/pkg/assert"
	"torrent/pkg/bencode"
)

type BencodeInfo struct {
	Pieces      string `bencode:"pieces"`
	PieceLength int    `bencode:"piece length"`
	Length      int    `bencode:"length"`
	Name        string `bencode:"name"`
}

type Torrent struct {
	Announce string      `bencode:"announce"`
	Info     BencodeInfo `bencode:"info"`
}

type TorrentFile struct {
	Announce    string
	InfoHash    [20]byte
	PieceHashes [][20]byte
	PieceLength int
	Length      int
	Name        string
}

func ParseTorrentFile(r io.Reader) Torrent {
	t := Torrent{}
	data, err := io.ReadAll(r)
	assert.NoError(err, "error while reading from reader")
	assert.NoError(bencode.Unmarshal(data, &t), "error while parsing torrent")
	return t
}
