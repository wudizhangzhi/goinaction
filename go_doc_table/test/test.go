package main

import (
	"encoding/xml"
	"fmt"
	"strings"
)

const xmldata = `
<w:p>
	<w:pPr>
		<w:pStyle w:val="10"/>
		<w:tabs>
			<w:tab w:val="left" w:pos="1706"/>
		</w:tabs>
		<w:spacing w:before="198"/>
		<w:ind w:left="107"/>
		<w:rPr>
			<w:sz w:val="21"/>
			<w:lang w:eastAsia="en-US"/>
		</w:rPr>
	</w:pPr>
	<w:r>
		<w:rPr>
			<w:rFonts w:ascii="Times New Roman" w:hAnsi="Times New Roman" w:eastAsia="Times New Roman"/>
			<w:sz w:val="21"/>
			<w:lang w:eastAsia="en-US"/>
		</w:rPr>
		<w:t>□</w:t>
	</w:r>
	<w:r>
		<w:rPr>
			<w:sz w:val="21"/>
			<w:lang w:eastAsia="en-US"/>
		</w:rPr>
		<w:t>已签订</w:t>
	</w:r>
	<w:r>
		<w:rPr>
			<w:sz w:val="21"/>
			<w:lang w:eastAsia="en-US"/>
		</w:rPr>
		<w:tab/>
	</w:r>
	<w:r>
		<w:rPr>
			<w:rFonts w:ascii="Times New Roman" w:hAnsi="Times New Roman" w:eastAsia="Times New Roman"/>
			<w:sz w:val="21"/>
			<w:lang w:eastAsia="en-US"/>
		</w:rPr>
		<w:sym w:font="Wingdings 2" w:char="0052"/>
	</w:r>
	<w:r>
		<w:rPr>
			<w:sz w:val="21"/>
			<w:lang w:eastAsia="en-US"/>
		</w:rPr>
		<w:t>未签订</w:t>
	</w:r>
</w:p>
`

type Paragraph struct {
	Rs []R `xml:"r"`
}

type R struct {
	Text string
}

func (r *R) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	fmt.Println("开始", start)
	// var s string
	// if err := d.DecodeElement(&s, &start); err != nil {
	// 	return err
	// }
	r.Text = ""
	if start.Name.Local != "r" {
		return fmt.Errorf("Invalid start tag for Property")
	}
	for {
		tok, err := d.Token()
		fmt.Println("Token: ", tok)
		if tok == nil {
			break
		}

		if err != nil {
			return err
		}
		switch se := tok.(type) {
		case xml.StartElement:
			fmt.Println("开始")
			switch se.Name.Local {
			case "sym":
				//
				fmt.Printf("sympol: %+v\n", se)
				for _, attr := range se.Attr {
					if attr.Name.Local == "char" {
						r.Text += attr.Value
					}
				}
			case "t":
				// r.Text = ""
			}
		case xml.CharData:
			r.Text += strings.TrimSpace(string(se))
		}
	}
	return nil
}

func main() {
	var p Paragraph
	if err := xml.Unmarshal([]byte(xmldata), &p); err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", p)
}
