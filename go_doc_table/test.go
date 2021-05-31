package main

import (
	// "encoding/xml"
	"encoding/json"
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
	fmt.Printf("tables: %d, cols: %d, rows: %d\n", len(d.Tables), len(d.Tables[0].Rows[0].Cols), len(d.Tables[0].Rows))
	// fmt.Printf("%+v\n", d.Tables[0].Rows[15].Cols[1].Paragraphs[0])
	// fmt.Printf("%+v\n", d.Tables[0].ToMap())
	j, _ := json.MarshalIndent(d.Tables[0].ToMap(), "", "  ")
	fmt.Println(string(j))
	j, _ = json.MarshalIndent(d.Tables[0].ToParsedMap(), "", "  ")
	fmt.Println(string(j))
}
