package document

import "strings"

var (
	sympolUncheckList     []string = []string{"☐", "□"}
	sympolCheckedList     []string = []string{"☑"}
	sympolUncheckOneline  string   = strings.Join(sympolUncheckList, "")
	sympolCheckedOneline  string   = strings.Join(sympolCheckedList, "")
	sympolUncheck         string   = "☐"
	sympolChecked         string   = "☑"
	wordUncheckSympolList []string = []string{"F02A", "00A3"}
	wordcheckedSympolList []string = []string{"0052", "F050", "F051", "F054", "F053", "F052"}
)
