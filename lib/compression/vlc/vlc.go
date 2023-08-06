package vlc

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"log"
	"packer/lib/compression/vlc/table"
	"strconv"
	"strings"
	"unicode"
)

type EncoderDecoder struct {
	tableGenerator table.Generator
}

func New(tableGen table.Generator) EncoderDecoder {
	return EncoderDecoder{tableGenerator: tableGen}
}

func (ed EncoderDecoder) Encode(str string) []byte {
	tbl := ed.tableGenerator.NewTable(str)

	encoded := encodeBin(str, tbl)

	return buildEncodedFile(tbl, encoded)
}

func (ed EncoderDecoder) Decode(encodedData []byte) string {
	tbl, data := parseFile(encodedData)

	return tbl.Decode(data)
}

func parseFile(data []byte) (table.EncodingTable, string) {
	const (
		tableSizeBytesCount = 4
		dataSizeBytesCount  = 4
	)
	tableSizeBinary, data := data[:tableSizeBytesCount], data[tableSizeBytesCount:]
	dataSizeBinary, data := data[:dataSizeBytesCount], data[dataSizeBytesCount:]

	tableSize := binary.BigEndian.Uint32(tableSizeBinary)
	dataSize := binary.BigEndian.Uint32(dataSizeBinary)

	tblBinary, data := data[:tableSize], data[tableSize:]

	tbl := decodeTable(tblBinary)

	body := NewBinChunks(data).Join()

	return tbl, body[:dataSize]
}

func buildEncodedFile(tbl table.EncodingTable, data string) []byte {
	encodedTable := encodeTable(tbl)

	var buf bytes.Buffer

	buf.Write(encodeInt(len(encodedTable)))
	buf.Write(encodeInt(len(data)))
	buf.Write(encodedTable)
	buf.Write(splitByChunks(data, chunkSize).Bytes())

	return buf.Bytes()
}

func encodeInt(num int) []byte {
	res := make([]byte, 4)
	binary.BigEndian.PutUint32(res, uint32(num))

	return res
}

func encodeTable(tbl table.EncodingTable) []byte {
	var tableBuf bytes.Buffer

	if err := gob.NewEncoder(&tableBuf).Encode(tbl); err != nil {
		log.Fatal("can't serialize table: ", err)
	}

	return tableBuf.Bytes()
}

func decodeTable(tblBinary []byte) table.EncodingTable {
	var tbl table.EncodingTable

	r := bytes.NewReader(tblBinary)
	if err := gob.NewDecoder(r).Decode(&tbl); err != nil {
		log.Fatal("can't deserialize table: ", err)
	}

	return tbl
}

func encodeBin(str string, table table.EncodingTable) string {
	var buf strings.Builder

	for _, ch := range str {
		buf.WriteString(bin(ch, table))
	}

	return buf.String()
}

func bin(ch rune, table table.EncodingTable) string {
	res, ok := table[ch]

	if !ok {
		panic("unknown character: " + string(ch) + " value: " + strconv.Itoa(int(ch)))
	}

	return res
}

func exportText(str string) string {
	var buf strings.Builder

	for i := 0; i < len(str); i++ {
		if str[i] == '!' {
			buf.WriteRune(unicode.ToUpper(rune(str[i+1])))
			i++
		} else {
			buf.WriteRune(rune(str[i]))
		}
	}

	return buf.String()
}
