package mygui

import (
	"fmt"
	"github.com/Xhofe/alist/model"
	"github.com/lxn/walk"
	log "github.com/sirupsen/logrus"
)

type MyAction string

const (
	DownLoad  MyAction = "DownLoad"
	Upload    MyAction = "Upload"
	Delete    MyAction = "Delete"
	Skip      MyAction = "Skip"
	CreateDir MyAction = "CreateDir"
	NoAction  MyAction = "NoAction"
	Error     MyAction = "Error"
)

var allActions = []string{"Skip", "DownLoad", "Upload", "Delete"}

type CallBackType string

const (
	NativeNoneFile CallBackType = "NativeNoneFile"
	SameSizeFile   CallBackType = "SameSizeFile"
	DiffSizeFile   CallBackType = "DiffSizeFile"
	BaiduNoneFile  CallBackType = "BaiduNoneFile"
)

type SyncFilePk struct {
	CallBackType   CallBackType `json:"call_back_type,omitempty"`
	NativeRootPath string       `json:"native_root_path,omitempty"`
	BaiduRootPath  string       `json:"baidu_root_path,omitempty"`
	File           model.File   `json:"File"`
	BaiduFile      model.File   `json:"baidu_file"`
	MyAction       MyAction     `json:"MyAction,omitempty"`
}

func (p SyncFilePk) doAction() bool {
	switch p.MyAction {
	case DownLoad:
		if p.CallBackType == NativeNoneFile {
			return syncOperator.CpDirOrFileFromBaiduWithIDMan(p.NativeRootPath, p.BaiduRootPath, p.BaiduFile)
		}
		break
	case Upload:
		if p.CallBackType == BaiduNoneFile {
			if p.File.IsDir() {
				//复制不存在的文件夹
				return syncOperator.uploadDirToBaidu(p.NativeRootPath, p.BaiduRootPath, p.File)
			} else {
				//复制不存在的文件
				return syncOperator.CpToBaiduFile(p.NativeRootPath, p.BaiduRootPath, p.File, false)
			}

		} else if p.CallBackType == DiffSizeFile {
			return syncOperator.CpToBaiduFile(p.NativeRootPath, p.BaiduRootPath, p.File, true)
		}
		break
	case Delete:
		if p.CallBackType == NativeNoneFile {
			return syncOperator.deleteBaiduFile(p.BaiduFile.Url)
		}
		break
	case CreateDir:
		if p.CallBackType == BaiduNoneFile {
			if p.File.IsDir() {
				//复制不存在的文件夹
				return syncOperator.createBaiduDir(p.BaiduRootPath + p.File.Url)
			}
		}
		break
	case Skip:
		log.Info("skip File:", p.String())
	default:
		break

	}
	return false

}
func (p SyncFilePk) String() string {
	if p.CallBackType == NativeNoneFile {
		return fmt.Sprintf("(%s)<<(%s)", "no support", p.BaiduFile.Url)
	} else if p.CallBackType == BaiduNoneFile {
		return fmt.Sprintf("(%s)>>(%s)", p.NativeRootPath+p.File.Url, p.BaiduRootPath+p.File.Url)
	} else if p.CallBackType == DiffSizeFile {
		return fmt.Sprintf("(%s)!=(%s)", p.NativeRootPath+p.File.Url, p.BaiduRootPath+p.File.Url)
	} else {
		return fmt.Sprintf("(%s)==(%s)", p.NativeRootPath+p.File.Url, p.BaiduRootPath+p.File.Url)
	}

}

type PkListModel struct {
	walk.ListModelBase
	items []SyncFilePk
}

func NewPkModel() *PkListModel {
	m := &PkListModel{items: make([]SyncFilePk, 0)}
	return m
}
func (m *PkListModel) ItemCount() int {
	return len(m.items)
}

func (m *PkListModel) Value(index int) interface{} {
	return m.items[index].String()
}
