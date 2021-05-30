package document

import (
	"encoding/xml"
	"fmt"
	"strings"
)

type Document struct {
	XMLName xml.Name `xml:"document"`
	Tables  []Table  `xml:"body>tbl"`
}

type Table struct {
	XMLName xml.Name `xml:"tbl"`
	Rows    []Row    `xml:"tr"`
}

type Row struct {
	xml.Name `xml:"tr"`
	Cols     []Col `xml:"tc"`
}

type Col struct {
	Paragraphs []string `xml:"p>r>t"`
	// InnerXML   struct {
	// 	InnerXML string `xml:",innerxml"`
	// } `xml:"p"`
}

// type Paragraph struct {

// }

func (c Col) String() string {
	return strings.Join(c.Paragraphs, "")
}

func (t Table) String() string {
	var s string
	for _, row := range t.Rows {
		cols_strings := []string{}
		for _, col := range row.Cols {
			col_str := col.String()
			cols_strings = append(cols_strings, col_str)
		}
		s += fmt.Sprintf("| %s |\n", strings.Join(cols_strings, " | "))
	}
	return s
}

// table转化为map
func (t Table) ToMap() map[string]string {
	m := make(map[string]string)
	for _, row := range t.Rows {
		colNum := len(row.Cols)
		switch {
		case colNum == 1:
			// 跳过
			continue
		case colNum == 2:
			// 第一个是key， 第二个是value
			m[row.Cols[0].String()] = row.Cols[1].String()
		case colNum == 3:
			//
		case colNum == 4:
			// 2对
			m[row.Cols[0].String()] = row.Cols[1].String()
			m[row.Cols[2].String()] = row.Cols[3].String()
		}

	}
	return m
}

// 占位数量
func SpaceCount(s string) int {
	bLen := len(s)                 // 字节数量
	sLen := len([]rune(s))         // 字数
	cAmount := (bLen - sLen) / 2   // 中文数量
	eAmount := (3*sLen - bLen) / 2 // 英文数量
	return eAmount + 2*cAmount
}
