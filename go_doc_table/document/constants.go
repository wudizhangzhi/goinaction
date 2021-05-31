package document

import "strings"

var (
	sympolUncheckList    []string = []string{"☐", "□"}
	sympolCheckedList    []string = []string{"☑"}
	sympolUncheckOneline string   = strings.Join(sympolUncheckList, "")
	sympolCheckedOneline string   = strings.Join(sympolCheckedList, "")
	sympolUncheck        string   = "☐"
	sympolChecked        string   = "☑"
)
