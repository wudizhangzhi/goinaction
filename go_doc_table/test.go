package main

import (
	// "encoding/xml"
	"encoding/xml"
	"fmt"
	"go-doc-table/document"

	"github.com/nguyenthenguyen/docx"
)

func main() {
	// xml.Unmarshal()
	filename := "test.docx"
	r, err := docx.ReadDocxFile(filename)
	if err != nil {
		panic(err)
	}
	docx1 := r.Editable()
	docXmlContent := docx1.GetContent()
	fmt.Println(docXmlContent)

	d := document.Document{}
	xml.Unmarshal([]byte(docXmlContent), &d)
	fmt.Printf("%+v\n", d)
	fmt.Println("")
	fmt.Printf("%s\n", d.Tables[0])
	fmt.Printf("cols: %d, rows: %d", len(d.Tables[0].Rows[0].Cols), len(d.Tables[0].Rows))
}
