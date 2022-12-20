package tableExample

import (
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"os"
)

type MyMainWindow struct {
	*walk.MainWindow
	edit *walk.TextEdit
	path string
}

func main() {
	mw := &MyMainWindow{}
	MW := MainWindow{AssignTo: &mw.MainWindow, Icon: "test.ico", Title: "文件选择对话框", MinSize: Size{150, 200}, Size: Size{300, 400}, Layout: VBox{}, Children: []Widget{TextEdit{AssignTo: &mw.edit}, PushButton{Text: "打开", OnClicked: mw.pbClicked}}}
	if _, err := MW.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
func (mw *MyMainWindow) pbClicked() {
	dlg := new(walk.FileDialog)
	dlg.FilePath = mw.path
	dlg.Title = "Select File"
	dlg.Filter = "Exe files (*.exe)|*.exe|All files (*.*)|*.*"
	if ok, err := dlg.ShowOpen(mw); err != nil {
		mw.edit.AppendText("Error : File Open\r\n")
		return
	} else if !ok {
		mw.edit.AppendText("Cancel\r\n")
		return
	}
	mw.path = dlg.FilePath
	s := fmt.Sprintf("Select : %s\r\n", mw.path)
	mw.edit.AppendText(s)
}
