package tableExample

import (
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"sort"
)

type Condom struct {
	Index   int
	Name    string
	Price   int
	checked bool
}
type CondomModel struct {
	walk.TableModelBase
	walk.SorterBase
	sortColumn int
	sortOrder  walk.SortOrder
	items      []*Condom
}

func (m *CondomModel) RowCount() int { return len(m.items) }
func (m *CondomModel) Value(row, col int) interface{} {
	item := m.items[row]
	switch col {
	case 0:
		return item.Index
	case 1:
		return item.Name
	case 2:
		return item.Price
	}
	return nil
}
func (m *CondomModel) Checked(row int) bool { return m.items[row].checked }
func (m *CondomModel) SetChecked(row int, checked bool) error {
	m.items[row].checked = checked
	return nil
}
func (m *CondomModel) Sort(col int, order walk.SortOrder) error {
	m.sortColumn, m.sortOrder = col, order
	sort.Stable(m)
	return m.SorterBase.Sort(col, order)
}
func (m *CondomModel) Len() int { return len(m.items) }
func (m *CondomModel) Less(i, j int) bool {
	a, b := m.items[i], m.items[j]
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
		return c(a.Name < b.Name)
	case 2:
		return c(a.Price < b.Price)
	}
	return true
}
func (m *CondomModel) Swap(i, j int) { m.items[i], m.items[j] = m.items[j], m.items[i] }
func NewCondomModel() *CondomModel {
	m := new(CondomModel)
	m.items = make([]*Condom, 3)
	m.items[0] = &Condom{Index: 0, Name: "杜蕾斯", Price: 20}
	m.items[1] = &Condom{Index: 1, Name: "杰士邦", Price: 18}
	m.items[2] = &Condom{Index: 2, Name: "冈本", Price: 19}
	return m
}

type CondomMainWindow struct {
	*walk.MainWindow
	model *CondomModel
	tv    *walk.TableView
}

func tableInit() {
	mw := &CondomMainWindow{model: NewCondomModel()}
	MainWindow{AssignTo: &mw.MainWindow, Title: "Condom展示",
		Size:   Size{800, 600},
		Layout: VBox{},
		Children: []Widget{
			Composite{
				Layout: HBox{MarginsZero: true},
				Children: []Widget{HSpacer{}, PushButton{Text: "Add", OnClicked: func() {
					mw.model.items = append(mw.model.items, &Condom{Index: mw.model.Len() + 1, Name: "第六感", Price: mw.model.Len() * 5})
					mw.model.PublishRowsReset()
					mw.tv.SetSelectedIndexes([]int{})
				}}, PushButton{Text: "Delete", OnClicked: func() {
					items := []*Condom{}
					remove := mw.tv.SelectedIndexes()
					for i, x := range mw.model.items {
						remove_ok := false
						for _, j := range remove {
							if i == j {
								remove_ok = true
							}
						}
						if !remove_ok {
							items = append(items, x)
						}
					}
					mw.model.items = items
					mw.model.PublishRowsReset()
					mw.tv.SetSelectedIndexes([]int{})
				}}, PushButton{Text: "ExecChecked", OnClicked: func() {
					for _, x := range mw.model.items {
						if x.checked {
							fmt.Printf("checked: %v\n", x)
						}
					}
					fmt.Println()
				}}, PushButton{Text: "AddPriceChecked", OnClicked: func() {
					for i, x := range mw.model.items {
						if x.checked {
							x.Price++
							mw.model.PublishRowChanged(i)
						}
					}
				}}}}, Composite{Layout: VBox{}, ContextMenuItems: []MenuItem{
				Action{Text: "I&nfo", OnTriggered: mw.tv_ItemActivated},
				Action{Text: "E&xit", OnTriggered: func() { mw.Close() }}},
				Children: []Widget{
					TableView{AssignTo: &mw.tv, CheckBoxes: true, ColumnsOrderable: true, MultiSelection: true, Columns: []TableViewColumn{{Title: "编号"}, {Title: "名称"}, {Title: "价格"}}, Model: mw.model, OnCurrentIndexChanged: func() {
						i := mw.tv.CurrentIndex()
						if 0 <= i {
							fmt.Printf("OnCurrentIndexChanged: %v\n", mw.model.items[i].Name)
						}
					}, OnItemActivated: mw.tv_ItemActivated}}},
		}}.Run()
}
func (mw *CondomMainWindow) tv_ItemActivated() {
	msg := ``
	for _, i := range mw.tv.SelectedIndexes() {
		msg = msg + "\n" + mw.model.items[i].Name
	}
	walk.MsgBox(mw, "title", msg, walk.MsgBoxIconInformation)
}
