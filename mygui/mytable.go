package mygui

import (
	"github.com/lxn/walk"
	"github.com/lxn/walk/declarative"
	log "github.com/sirupsen/logrus"
	"os"
	"sort"
)

var tableView, actionTableView *walk.TableView
var titles = []declarative.TableViewColumn{{Title: "Index", Width: 20}, {Title: "Action", Width: 30}, {Title: "CallBackType"}, {Title: "Path", Width: 400}}
var selectedIndex []int
var combo *walk.ComboBox

type Condom struct {
	SyncFilePk `json:"sync_file_pk"`
	Index      int  `json:"index,omitempty"`
	Checked    bool `json:"Checked,omitempty"`
}

type CondomModel struct {
	walk.TableModelBase
	walk.SorterBase
	sortColumn int
	sortOrder  walk.SortOrder
	Items      []*Condom `json:"Items,omitempty"`
}

func (m *CondomModel) RowCount() int { return len(m.Items) }
func (m *CondomModel) Value(row, col int) interface{} {
	item := m.Items[row]
	switch col {
	case 0:
		return item.Index
	case 1:
		return item.MyAction
	case 2:
		return item.CallBackType
	case 3:
		return item.String()
	}
	return nil
}
func (m *CondomModel) Checked(row int) bool { return m.Items[row].Checked }
func (m *CondomModel) SetChecked(row int, checked bool) error {
	m.Items[row].Checked = checked
	return nil
}
func (m *CondomModel) Sort(col int, order walk.SortOrder) error {
	m.sortColumn, m.sortOrder = col, order
	sort.Stable(m)
	return m.SorterBase.Sort(col, order)
}
func (m *CondomModel) Len() int { return len(m.Items) }
func (m *CondomModel) Less(i, j int) bool {
	a, b := m.Items[i], m.Items[j]
	c := func(ls bool) bool {
		if m.sortOrder == walk.SortAscending {
			return ls
		}
		return !ls
	}
	switch m.sortColumn {
	case 0:
		return c(a.Index < b.Index)
	case 1:
		return c(a.MyAction < b.MyAction)
	case 2:
		return c(a.CallBackType < b.CallBackType)
	case 3:
		return c(a.String() < b.String())
	}
	return true
}
func (m *CondomModel) Swap(i, j int) { m.Items[i], m.Items[j] = m.Items[j], m.Items[i] }
func NewCondomModel() *CondomModel {
	m := new(CondomModel)
	m.Items = make([]*Condom, 0)
	return m
}
func addToTable(pk SyncFilePk) {
	pksTab, _ := tableView.Model().(*CondomModel)
	pksTab.Items = append(pksTab.Items, &Condom{pk, pksTab.Len(), false})

	err := tableView.SetModel(pksTab)
	if err != nil {
		log.Error(err)
	}
}
func addToActionTable() {
	pks, _ := tableView.Model().(*CondomModel)
	actionPks := NewCondomModel()
	selectedIndex = tableView.SelectedIndexes()
	for _, index := range selectedIndex {
		pks.Items[index].Checked = false
		actionPks.Items = append(actionPks.Items, pks.Items[index])
	}
	err := actionTableView.SetModel(actionPks)
	if err != nil {
		log.Error(err)
	}
}
func doActionTable() {
	pks, _ := tableView.Model().(*CondomModel)
	//actionPks, _ := actionTableView.Model().(*CondomModel)
	nowAction := MyAction(allActions[combo.CurrentIndex()])

	for _, index := range selectedIndex {
		pks.Items[index].Checked = true
		pks.Items[index].MyAction = nowAction
		if pks.Items[index].doAction() {
			pks.Items[index].CallBackType = SameSizeFile
		} else {
			pks.Items[index].MyAction = Error
		}

	}
	err := tableView.SetModel(pks)
	if err != nil {
		log.Error(err)
	}
}
func saveTableToFile() {
	dlg := new(walk.FileDialog)

	dlg.Title = "Select File"
	dlg.InitialDirPath = "./"
	dlg.Filter = "Json files (*.json)|*.json|All files (*.*)|*.*"
	if ok, err := dlg.ShowSave(mainWin); err != nil {
		log.Error(err)
		return
	} else if !ok {
		return
	}

	pks, _ := tableView.Model().(*CondomModel)
	listFilePath := dlg.FilePath
	var err error
	if PathExists(listFilePath) {
		listFile, err = os.OpenFile(listFilePath, os.O_RDWR, 0)
	} else {
		//新建文件
		listFile, _ = os.Create(listFilePath)

	}
	if err != nil {
		log.Fatal(err)
		return
	}
	listFile.WriteString(serializeJson(pks))
	listFile.Close()
	log.Infof("success save table file")

}
func readTableFromFile() *CondomModel {
	dlg := new(walk.FileDialog)
	pks := NewCondomModel()
	dlg.Title = "Select File"
	dlg.Filter = "Json files (*.json)|*.json|All files (*.*)|*.*"
	if ok, err := dlg.ShowOpen(mainWin); err != nil {
		log.Error(err)
		return pks
	} else if !ok {
		return pks
	}

	b, err := os.ReadFile(dlg.FilePath)
	if err != nil {
		log.Error(err)
		return pks
	}
	str := string(b)
	deSerializeJson(str, &pks)
	return pks

}
