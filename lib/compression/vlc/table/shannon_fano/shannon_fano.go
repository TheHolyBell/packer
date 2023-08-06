package shannon_fano

import (
	"fmt"
	"math"
	"packer/lib/compression/vlc/table"
	"sort"
	"strings"
)

type Generator struct {
}

type code struct {
	Char     rune
	Quantity int
	Bits     uint32
	Size     int
}

type charStat map[rune]int

type encodingTable map[rune]code

func NewGenerator() Generator {
	return Generator{}
}

func (g Generator) NewTable(text string) table.EncodingTable {
	stat := newCharStat(text)

	return build(stat).Export()
}

func (et encodingTable) Export() table.EncodingTable {
	res := table.EncodingTable{}

	for k, v := range et {
		byteStr := fmt.Sprintf("%b", v.Bits)

		if lenDiff := v.Size - len(byteStr); lenDiff > 0 {
			byteStr = strings.Repeat("0", lenDiff) + byteStr
		}

		res[k] = byteStr
	}

	return res
}

func build(stat charStat) encodingTable {
	codes := make([]code, 0, len(stat))

	for ch, qty := range stat {
		codes = append(codes, code{
			Char:     ch,
			Quantity: qty,
		})
	}

	sort.Slice(codes, func(i, j int) bool {
		if codes[i].Quantity != codes[j].Quantity {
			return codes[i].Quantity > codes[j].Quantity
		}

		return codes[i].Char < codes[j].Char
	})

	assignCodes(codes)

	res := encodingTable{}

	for _, c := range codes {
		res[c.Char] = c
	}

	return res
}

func assignCodes(codes []code) {
	if len(codes) < 2 {
		return
	}

	divider := bestDividerPosition(codes)
	for i := 0; i < len(codes); i++ {
		codes[i].Bits <<= 1
		codes[i].Size++

		if i >= divider {
			codes[i].Bits |= 1
		}
	}

	assignCodes(codes[:divider])
	assignCodes(codes[divider:])
}

func bestDividerPosition(codes []code) int {
	total := 0

	for _, code := range codes {
		total += code.Quantity
	}

	left := 0
	prevDiff := math.MaxInt
	bestPosition := 0

	for i := 0; i < len(codes)-1; i++ {
		left += codes[0].Quantity
		right := total - left

		diff := abs(right - left)
		if diff >= prevDiff {
			break
		}

		prevDiff = diff
		bestPosition = i + 1
	}

	return bestPosition
}

func newCharStat(str string) charStat {
	res := charStat{}

	for _, ch := range str {
		res[ch]++
	}

	return res
}

func abs(val int) int {
	if val < 0 {
		return -val
	}

	return val
}
