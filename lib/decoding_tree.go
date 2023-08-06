package lib

import "strings"

type DecodingTree struct {
	Value string
	Zero  *DecodingTree
	One   *DecodingTree
}

func (et encodingTable) DecodingTree() DecodingTree {
	res := DecodingTree{}

	for ch, code := range et {
		res.Add(code, ch)
	}

	return res
}

func (dt *DecodingTree) Decode(str string) string {
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

func (dt *DecodingTree) Add(code string, value rune) {
	currentNode := dt

	for _, ch := range code {
		switch ch {
		case '0':
			if currentNode.Zero == nil {
				currentNode.Zero = &DecodingTree{}
			}
			currentNode = currentNode.Zero
			break
		case '1':
			if currentNode.One == nil {
				currentNode.One = &DecodingTree{}
			}
			currentNode = currentNode.One
			break
		}
	}

	currentNode.Value = string(value)
}
