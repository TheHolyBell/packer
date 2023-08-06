package lib

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

type encodingTable map[rune]string

type BinaryChunks []BinaryChunk

type BinaryChunk string

type HexChunk string
type HexChunks []HexChunk

const (
	chunkSize         = 8
	hexChunkSeparator = " "
)

func (bcs BinaryChunks) ToHex() HexChunks {
	res := make(HexChunks, 0, len(bcs))

	for _, chunk := range bcs {
		hexChunk := chunk.ToHex()

		res = append(res, hexChunk)
	}

	return res
}

func (bcs BinaryChunks) Join() string {
	var buf strings.Builder

	for _, bc := range bcs {
		buf.WriteString(string(bc))
	}

	return buf.String()
}

func NewHexChunks(str string) HexChunks {
	parts := strings.Split(str, hexChunkSeparator)

	res := make(HexChunks, 0, len(parts))

	for _, part := range parts {
		res = append(res, HexChunk(part))
	}

	return res
}

func (hcs HexChunks) ToBinary() BinaryChunks {
	res := make(BinaryChunks, 0, len(hcs))

	for _, chunk := range hcs {
		res = append(res, chunk.ToBinary())
	}

	return res
}

func (hc HexChunk) ToBinary() BinaryChunk {
	num, err := strconv.ParseUint(string(hc), 16, chunkSize)
	if err != nil {
		panic("can't parse hex chunk: " + err.Error())
	}

	res := fmt.Sprintf("%08b", num)

	return BinaryChunk(res)
}

func (hcs HexChunks) ToString() string {

	switch len(hcs) {
	case 0:
		return ""
	case 1:
		return string(hcs[0])
	}

	var buf strings.Builder

	buf.WriteString(string(hcs[0]))

	for _, hc := range hcs[1:] {
		buf.WriteString(hexChunkSeparator)
		buf.WriteString(string(hc))
	}

	return buf.String()
}

func (bc BinaryChunk) ToHex() HexChunk {
	num, err := strconv.ParseUint(string(bc), 2, chunkSize)
	if err != nil {
		panic("can't parse binary chunk: " + err.Error())
	}

	res := strings.ToUpper(fmt.Sprintf("%x", num))

	if len(res) == 1 {
		res = "0" + res
	}

	return HexChunk(res)
}

func splitByChunks(bStr string, chunkSize int) BinaryChunks {
	strLen := utf8.RuneCountInString(bStr)
	chunksCount := strLen / chunkSize

	if strLen/chunkSize != 0 {
		chunksCount++
	}

	res := make(BinaryChunks, 0, chunksCount)

	var buf strings.Builder

	for i, ch := range bStr {
		buf.WriteString(string(ch))

		if (i+1)%chunkSize == 0 {
			res = append(res, BinaryChunk(buf.String()))
			buf.Reset()
		}
	}

	if buf.Len() != 0 {
		lastChunk := buf.String()

		lastChunk += strings.Repeat("0", chunkSize-len(lastChunk))
		res = append(res, BinaryChunk(lastChunk))
	}

	return res
}
