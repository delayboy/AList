package mygui

import (
	"errors"
	"fmt"
	"github.com/Xhofe/alist/conf"
	"github.com/lxn/walk"
	"github.com/lxn/walk/declarative"
	log "github.com/sirupsen/logrus"
	"os/exec"
	"reflect"
)

var outTE *walk.TextEdit
var mainWin *walk.MainWindow

func toSyncFilePk(actual interface{}) ([]SyncFilePk, error) {
	var res []SyncFilePk
	value := reflect.ValueOf(actual)
	if value.Kind() != reflect.Slice && value.Kind() != reflect.Array {
		return nil, errors.New("parse error")
	}
	for i := 0; i < value.Len(); i++ {
		res = append(res, value.Index(i).Interface().(SyncFilePk))
	}
	return res, nil
}

func guiSetText(content string) {
	err := outTE.SetText(content)
	if err != nil {
		log.Error(err)
	}
}

// GuiInit 每一个源文件都可以包含一个init函数，该函数会在main函数执行前，被Go运行框架调用，也就是说init会在main函数前被调用
func GuiInit() {

	//加载自定义窗口
	fmt.Printf("load my gui window")
	var _, errWin = declarative.MainWindow{
		AssignTo: &mainWin,

		Title: "Benson Sync ToolBox",

		MinSize: declarative.Size(walk.Size{Width: 600, Height: 400}),

		Layout: declarative.VBox{},

		Children: []declarative.Widget{
			declarative.HSplitter{
				Children: []declarative.Widget{
					declarative.PushButton{

						Text: "Sync",

						OnClicked: func() {
							SyncEntry()
							//outTE.SetText(strings.ToUpper(inTE.Text()))

						},
					},
					declarative.PushButton{

						Text: "Show Settings",

						OnClicked: func() {
							url := fmt.Sprintf("http://127.0.0.1:%d/@manage/settings/0", conf.Conf.Port)
							cmd := exec.Command("cmd", "/c", "start", url)
							err := cmd.Start()
							if err != nil {
								log.Error(err)
							}
							//outTE.SetText(strings.ToUpper(inTE.Text()))

						},
					},
					declarative.PushButton{

						Text: "Open Table",

						OnClicked: func() {

							tableView.SetModel(readTableFromFile())
							//outTE.SetText(strings.ToUpper(inTE.Text()))

						},
					},
					declarative.PushButton{

						Text: "Save Table",
						OnClicked: func() {
							saveTableToFile()
							//outTE.SetText(strings.ToUpper(inTE.Text()))

						},
					},
				},
			},
			declarative.HSplitter{

				Children: []declarative.Widget{
					declarative.VSplitter{
						Children: []declarative.Widget{
							//declarative.ListBox{AssignTo: &actionListBox, MultiSelection: true, Model: NewPkModel()},
							declarative.TableView{AssignTo: &actionTableView, CheckBoxes: true, ColumnsOrderable: true, MultiSelection: true, Columns: titles, Model: NewCondomModel()},
							declarative.TextEdit{AssignTo: &outTE, ReadOnly: true},
							declarative.ComboBox{AssignTo: &combo, Model: allActions, CurrentIndex: 0},
							declarative.PushButton{
								Text: "do MyAction",
								OnClicked: func() {
									//outTE.SetText(strings.ToUpper(inTE.Text()))
									doActionTable()
								},
							},
						},
					},
					declarative.VSplitter{
						Children: []declarative.Widget{
							//declarative.ListBox{AssignTo: &listBox, MultiSelection: true, Model:  NewPkModel()},
							declarative.TableView{
								AssignTo:         &tableView,
								CheckBoxes:       true,
								ColumnsOrderable: true,
								MultiSelection:   true,
								Columns:          titles,
								Model:            NewCondomModel(),
							},
							declarative.PushButton{
								Text: "clear table",
								OnClicked: func() {
									tableView.SetModel(NewCondomModel())

								},
							},
							declarative.PushButton{
								Text: "add to MyAction List",
								OnClicked: func() {
									addToActionTable()
									//outTE.SetText(strings.ToUpper(inTE.Text()))

								},
							},
						},
					},
				},
			},
		},
	}.Run()

	if errWin != nil {
		fmt.Printf("Run err: %+v\n", errWin)
	}

}
