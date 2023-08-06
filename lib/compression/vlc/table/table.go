package table

import "strings"

type Generator interface {
	NewTable(text string) EncodingTable
}

type EncodingTable map[rune]string

type decodingTree struct {
	Value string
	Zero  *decodingTree
	One   *decodingTree
}

func (et EncodingTable) Decode(str string) string {
	dt := et.decodingTree()

	return dt.Decode(str)
}

func (et EncodingTable) decodingTree() decodingTree {
	res := decodingTree{}

	for ch, code := range et {
		res.add(code, ch)
	}

	return res
}

func (dt *decodingTree) Decode(str string) string {
	var buf strings.Builder

	currNode := dt

	for _, ch := range str {
		if currNode.Value != "" {
			buf.WriteString(currNode.Value)
			currNode = dt
		}

		switch ch {
		case '0':
			currNode = currNode.Zero
			break
		case '1':
			currNode = currNode.One
			break
		}
	}

	if currNode.Value != "" {
		buf.WriteString(currNode.Value)
		currNode = dt
	}

	return buf.String()
}

func (dt *decodingTree) add(code string, value rune) {
	currentNode := dt

	for _, ch := range code {
		switch ch {
		case '0':
			if currentNode.Zero == nil {
				currentNode.Zero = &decodingTree{}
			}
			currentNode = currentNode.Zero
			break
		case '1':
			if currentNode.One == nil {
				currentNode.One = &decodingTree{}
			}
			currentNode = currentNode.One
			break
		}
	}

	currentNode.Value = string(value)
}
