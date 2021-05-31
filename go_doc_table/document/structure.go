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
	Paragraphs []Paragraph `xml:"p"`
}

type Paragraph struct {
	Texts []R `xml:"r"`
}

type R struct {
	Text string
}

// 内容解析
func (r *R) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	r.Text = ""
	if start.Name.Local != "r" {
		return fmt.Errorf("Invalid start tag for Property: %v", start)
	}
	for {
		tok, err := d.Token()
		if tok == nil {
			break
		}

		if err != nil {
			return err
		}
		switch se := tok.(type) {
		case xml.StartElement:
			// fmt.Println("开始")
			switch se.Name.Local {
			case "sym": // 符号
				// fmt.Printf("sympol: %+v\n", se)
				for _, attr := range se.Attr {
					if attr.Name.Local == "char" {
						// TODO 转化 ☑ ☒ ☐
						// r.Text += attr.Value
						switch {
						// checkedd F050, F051, F054, F053, F052
						case contains(wordUncheckSympolList, attr.Value):
							r.Text += sympolUncheck
						case contains(wordcheckedSympolList, attr.Value):
							r.Text += sympolChecked
							// case "0052":
							// 	r.Text += sympolChecked
							// // uncheck F02A 00A3
							// case "F02A":
							// 	r.Text += sympolUncheck
							// case "00A3":
							// 	r.Text += sympolUncheck
						}
					}
				}
			case "t": // 文字
				// r.Text = ""
			}
		case xml.CharData: // 文字内容
			r.Text += strings.TrimSpace(string(se))
		}
	}
	return nil
}

func (p Paragraph) String() string {
	r := ""
	for _, t := range p.Texts {
		r += t.Text
	}
	return r
}

func (c Col) String() string {
	r := ""
	for _, p := range c.Paragraphs {
		r += p.String()
	}
	return r
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
			// 第二第三都是value
			m[row.Cols[0].String()] = row.Cols[1].String() + row.Cols[2].String()
		case colNum == 4:
			// 2对
			m[row.Cols[0].String()] = row.Cols[1].String()
			m[row.Cols[2].String()] = row.Cols[3].String()
		}

	}
	return m
}

// table 转化为 处于好结果的map
func (t Table) ToParsedMap() map[string]string {
	m := t.ToMap()
	for k, v := range m {
		// val := m[key]
		// fmt.Println(val)
		v = parseParagraphVal(v)
		m[k] = v
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
