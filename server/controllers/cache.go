package controllers

import (
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Xhofe/alist/conf"
	"github.com/Xhofe/alist/drivers/base"
	"github.com/Xhofe/alist/model"
	"github.com/Xhofe/alist/server/common"
	"github.com/Xhofe/alist/utils"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var start = false
var pause = false

// 计数器
var wg sync.WaitGroup

// SyncOperator 如果类名首字母大写，表示其他包也能够访问
type SyncOperator struct {
	//如果说类的属性首字母大写, 表示该属性是对外能够访问的，否则的话只能够类的内部访问
	driversBaidu  base.Driver
	driversNative base.Driver
	accountBaidu  model.Account
	accountNative model.Account
	TotalFileNum  int64        `json:"TotalFileNum,omitempty"`
	FinishFileNum int64        `json:"FinishFileNum,omitempty"`
	NewFileNum    int          `json:"new_file_num,omitempty"`
	UpdateFileNum int          `json:"update_file_num,omitempty"`
	DeleteFileNum int          `json:"delete_file_num,omitempty"`
	TaskId        int          `json:"taskId" json:"task_id,omitempty"`
	Queue         []model.File `json:"queue" json:"queue,omitempty"`
}

func (operator *SyncOperator) deleteBaiduFile(filePath string) bool {
	retryNum := 0
	for {
		err := operator.driversBaidu.Delete(filePath, &operator.accountBaidu)
		if err != nil {
			log.Errorf("继续尝试(%d)文件删除失败(%s)", retryNum, filePath)
			retryNum++
			operator.pauseAndSaveListener(retryNum)
			time.Sleep(time.Second)
		} else {
			log.Infof("文件删除成功(%s)", filePath)
			operator.DeleteFileNum++
			return true
		}
	}

}
func (operator *SyncOperator) createBaiduDir(dirPath string) bool {
	retryNum := 0
	for {
		err := operator.driversBaidu.MakeDir(dirPath, &operator.accountBaidu)
		if err != nil {
			log.Errorf("继续尝试(%d)文件夹创建失败(%s)", retryNum, dirPath)
			retryNum++
			operator.pauseAndSaveListener(retryNum)
			time.Sleep(time.Second)
		} else {
			log.Infof("文件夹创建成功(%s)", dirPath)
			operator.UpdateFileNum++
			return true
		}
	}

}

func (operator *SyncOperator) CpToBaiduFile(nativePath string, baiduPath string, file model.File, update bool) bool {

	nativeFilePath := filepath.Join(operator.accountNative.RootFolder, nativePath, file.Url)
	aimPath := utils.Join(baiduPath, file.Url)
	if update {
		if !operator.deleteBaiduFile(aimPath) {
			log.Errorf("由于删除失败，导致上传到百度网盘失败!(%s)->(%s)", nativePath, aimPath)
			return false
		}
	}
	dir, name := filepath.Split(aimPath)
	retryNum := 0
	for {
		if open, err := os.Open(nativeFilePath); err == nil {
			fileStream := model.FileStream{
				File:       open,
				Size:       uint64(file.Size),
				ParentPath: dir,
				Name:       name,
				MIMEType:   "application/octet-stream",
			}

			err = operator.driversBaidu.Upload(&fileStream, &operator.accountBaidu)
			if err != nil {
				log.Errorf("继续尝试(%d)上传到百度网盘失败!(%s)->(%s)", retryNum, nativePath, aimPath)
				retryNum++
				operator.pauseAndSaveListener(retryNum)
				time.Sleep(time.Second)
			} else {
				log.Infof("已上传(%s)->(%s)", nativePath, aimPath)
				if update {
					operator.UpdateFileNum++
				} else {
					operator.NewFileNum++
				}
				return true
			}
			_ = open.Close()
		} else {
			log.Errorf("本地文件打开失败(%s)", nativePath)
			return false
		}
	}

}
func (operator *SyncOperator) GetBaiduFilesMap(dirPath string) (map[string]model.File, bool) {
	dict := make(map[string]model.File)
	files, err := operator.driversBaidu.Files(dirPath, &operator.accountBaidu)
	if err == nil {
		for _, file := range files {
			dict[file.Name] = file
		}
		return dict, true
	} else {

		return dict, false
	}
}

func (operator *SyncOperator) GetNativeFile(filePath string) *model.File {
	file, _ := operator.driversNative.File(filePath, &operator.accountNative)
	return file
}
func (operator *SyncOperator) GetNativeFilesMap(dirPath string) map[string]model.File {
	dict := make(map[string]model.File)
	files, _ := operator.driversNative.Files(dirPath, &operator.accountNative)
	for _, file := range files {
		dict[file.Name] = file
	}
	return dict
}
func (operator *SyncOperator) syncAll() {
	nativeRootPaths := getSyncPathSetting("nativeRootPaths", []string{"/Work/Notes", "/Computer/Program/Program_maker"})
	baiduRootPaths := getSyncPathSetting("baiduRootPaths", []string{"/apps/FileGee文件同步备份系统/Notes", "/apps/FileGee文件同步备份系统/Program/Program_maker"})
	count := len(nativeRootPaths)
	useSyncList := getSyncListSetting()
	for i := operator.TaskId; i < count; i++ {
		operator.TaskId = i
		operator.TotalFileNum = operator.statisticFilesNum(nativeRootPaths[i])
		if useSyncList {
			operator.listSync(nativeRootPaths[i], baiduRootPaths[i])
		} else {
			operator.sync(nativeRootPaths[i], baiduRootPaths[i])
		}
	}

}
func (operator *SyncOperator) statisticFilesNum(nativeRootPath string) int64 {
	file := operator.GetNativeFile(nativeRootPath)
	file.Url = "/"
	totalFileNum := int64(0)
	var Queue []model.File
	Queue = append(Queue, *file)
	for {
		if len(Queue) < 1 {
			break
		}
		nowFile := Queue[0]
		nowFileNativePath := utils.Join(nativeRootPath, nowFile.Url)
		nativeFilesMap := operator.GetNativeFilesMap(nowFileNativePath)
		for _, file := range nativeFilesMap {
			file.Url = utils.Join(nowFile.Url, file.Name)

			totalFileNum++
			if totalFileNum%5000 == 0 {
				log.Infof("%d>%s", totalFileNum, file.Url)
			}

			if file.IsDir() {
				Queue = append(Queue, file)

			}

		}
		Queue = Queue[1:] // 删除开头1个元素

	}
	log.Infof("统计完成，一共需要同步:%d个文件", totalFileNum)
	return totalFileNum
}
func (operator *SyncOperator) sync(nativeRootPath string, baiduRootPath string) {
	start = true
	file := operator.GetNativeFile(nativeRootPath)
	file.Url = "/"

	var queueTemp []model.File
	if operator.Queue == nil {
		operator.Queue = queueTemp
		operator.Queue = append(operator.Queue, *file)
	}

	for {
		if len(operator.Queue) < 1 {
			break
		}
		nowFile := operator.Queue[0]
		if !nowFile.IsDir() {
			log.Errorf("出现了非文件夹节点,请选择文件夹进行同步！！！")
			break
		}
		nowFileNativePath := utils.Join(nativeRootPath, nowFile.Url)
		nativeFilesMap := operator.GetNativeFilesMap(nowFileNativePath)
		nowFileAimPath := utils.Join(baiduRootPath, nowFile.Url)
		baiduFileMap, hasDir := operator.GetBaiduFilesMap(nowFileAimPath)
		operator.pauseAndSaveListener(0)

		if hasDir {
			for _, file := range nativeFilesMap {
				file.Url = utils.Join(nowFile.Url, file.Name)
				operator.FinishFileNum++
				if operator.FinishFileNum%50 == 0 {
					percent := 100.0
					percent = percent * float64(operator.FinishFileNum) / float64(operator.TotalFileNum)
					log.Infof("[%.2f%%](scan %d, new %d,update %d,delete %d)", percent, operator.FinishFileNum, operator.NewFileNum, operator.UpdateFileNum, operator.DeleteFileNum)
				}
				if baiduFile, hasFile := baiduFileMap[file.Name]; hasFile {
					operator.pauseAndSaveListener(0)
					if file.IsDir() {
						operator.Queue = append(operator.Queue, file)

					} else {
						//复制文件大小不一样的文件,使用go关键字异步执行
						if baiduFile.Size != file.Size {
							operator.CpToBaiduFile(nativeRootPath, baiduRootPath, file, true)
						} else {

							//log.Infof("skip(%s)", baiduFile.Url)

						}

					}
				} else {
					if file.IsDir() {
						//创建文件夹并入队列
						if operator.createBaiduDir(baiduRootPath + file.Url) {
							operator.Queue = append(operator.Queue, file)
						}

					} else {
						//复制不存在的文件
						operator.CpToBaiduFile(nativeRootPath, baiduRootPath, file, false)
					}

				}

			}
			baiduFileMap, _ := operator.GetBaiduFilesMap(nowFileAimPath) //重新获取百度网盘文件
			for _, baiduFile := range baiduFileMap {
				if _, hasFile := nativeFilesMap[baiduFile.Name]; !hasFile {
					operator.deleteBaiduFile(baiduFile.Url)
				}
			}
			operator.Queue = operator.Queue[1:] // 删除开头1个元素
		} else {
			log.Errorf("致命错误，百度网盘没有找到该文件夹(%s)", nowFileAimPath)
			return

		}

	}
	start = false
	log.Infof("所有文件夹，均已同步!!!")
	operator.Queue = nil
	err := os.Remove("./operator.gob")
	if err != nil {
		log.Errorf("删除operator缓存文件失败")
	}

}
func (operator *SyncOperator) listSync(nativeRootPath string, baiduRootPath string) {
	start = true
	file := operator.GetNativeFile(nativeRootPath)
	file.Url = "/"
	listFilePath := "./listSync.txt"
	var listFile *os.File
	var err error
	if PathExists(listFilePath) {
		listFile, err = os.OpenFile(listFilePath, os.O_APPEND, 0)
	} else {
		//新建文件
		listFile, _ = os.Create(listFilePath)

	}

	if err != nil {
		log.Fatal(err)
		return
	}
	defer listFile.Close()
	listFile.WriteString(fmt.Sprintf("记录文件夹:%s\n", nativeRootPath))
	var queueTemp []model.File
	if operator.Queue == nil {
		operator.Queue = queueTemp
		operator.Queue = append(operator.Queue, *file)
	}

	for {
		if len(operator.Queue) < 1 {
			break
		}
		nowFile := operator.Queue[0]
		if !nowFile.IsDir() {
			log.Errorf("出现了非文件夹节点,请选择文件夹进行同步！！！")
			break
		}
		nowFileNativePath := utils.Join(nativeRootPath, nowFile.Url)
		nativeFilesMap := operator.GetNativeFilesMap(nowFileNativePath)
		nowFileAimPath := utils.Join(baiduRootPath, nowFile.Url)
		baiduFileMap, hasDir := operator.GetBaiduFilesMap(nowFileAimPath)
		operator.pauseAndSaveListener(0)

		if hasDir {
			for _, file := range nativeFilesMap {
				file.Url = utils.Join(nowFile.Url, file.Name)
				operator.FinishFileNum++
				if operator.FinishFileNum%50 == 0 {
					percent := 100.0
					percent = percent * float64(operator.FinishFileNum) / float64(operator.TotalFileNum)
					log.Infof("[%.2f%%](scan %d, new %d,update %d,delete %d)", percent, operator.FinishFileNum, operator.NewFileNum, operator.UpdateFileNum, operator.DeleteFileNum)
				}
				if baiduFile, hasFile := baiduFileMap[file.Name]; hasFile {
					operator.pauseAndSaveListener(0)
					if file.IsDir() {
						operator.Queue = append(operator.Queue, file)

					} else {
						//复制文件大小不一样的文件,使用go关键字异步执行
						if baiduFile.Size != file.Size {
							operator.UpdateFileNum++
							listFile.WriteString(fmt.Sprintf("向网盘更新不同的文件(%s)->(%s)\n", nativeRootPath+file.Url, baiduRootPath+file.Url))
						}

					}
				} else {
					if file.IsDir() {
						//提示创建文件夹
						operator.NewFileNum++
						listFile.WriteString(fmt.Sprintf("在网盘创建文件夹(%s)\n", baiduRootPath+file.Url))

					} else {
						//复制不存在的文件
						operator.NewFileNum++
						listFile.WriteString(fmt.Sprintf("向网盘上传不存在的文件(%s)\n", nativeRootPath+file.Url))
					}

				}

			}
			for _, baiduFile := range baiduFileMap {
				if _, hasFile := nativeFilesMap[baiduFile.Name]; !hasFile {
					operator.DeleteFileNum++
					listFile.WriteString(fmt.Sprintf("删除网盘上存在的文件(%s)\n", baiduFile.Url))
				}
			}
			operator.Queue = operator.Queue[1:] // 删除开头1个元素
		} else {
			log.Errorf("致命错误，百度网盘没有找到该文件夹(%s)", nowFileAimPath)
			return

		}

	}
	start = false
	log.Infof("所有文件夹，均已扫描!!!")
	operator.Queue = nil
	err = os.Remove("./operator.gob")
	if err != nil {
		log.Errorf("删除operator缓存文件失败")
	}

}
func InitSyncOperator(accounts []model.Account) SyncOperator {
	var driversBaidu base.Driver
	var driversNative base.Driver
	var accountBaidu model.Account
	var accountNative model.Account
	for i := range accounts {
		account := accounts[i]
		driver, _ := base.GetDriver(account.Type)
		log.Infof("Account:%s  Type:%s", account.Name, driver.Config().Name)
		if driver.Config().Name == "Baidu.Disk" {
			log.Infof("发现百度网盘账号：%s", account.Name)
			accountBaidu = account
			driversBaidu = driver

		}
		if driver.Config().Name == "Native" {
			log.Infof("发现本地盘账号：%s", account.Name)
			driversNative = driver
			accountNative = account
		}
	}
	operator := SyncOperator{driversBaidu, driversNative, accountBaidu, accountNative, 0, 0, 0, 0, 0, 0, nil}
	lastOperator := getOperator()
	if lastOperator != nil {
		operator.TaskId = lastOperator.TaskId
		operator.Queue = lastOperator.Queue
		operator.TotalFileNum = lastOperator.TotalFileNum
		operator.FinishFileNum = lastOperator.FinishFileNum
		operator.NewFileNum = lastOperator.NewFileNum
		operator.UpdateFileNum = lastOperator.UpdateFileNum
		operator.DeleteFileNum = lastOperator.DeleteFileNum

	}
	return operator

}
func saveSyncOperator(operator SyncOperator) {
	file, err := os.Create("./operator.gob")
	if err != nil {
		fmt.Println("文件创建失败", err.Error())
		return
	}
	encoder := gob.NewEncoder(file)
	err = encoder.Encode(operator)
}
func getOperator() *SyncOperator {
	file, err := os.Open("./operator.gob")
	if err != nil {
		fmt.Println("文件打开失败", err.Error())
	}
	defer file.Close()
	decoder := gob.NewDecoder(file)
	var operator SyncOperator
	err = decoder.Decode(&operator)
	if err != nil {
		fmt.Println("解码失败", err.Error())
	} else {
		log.Infof("已找到找到之前缓存的文件队列，将从上一次运行结果开始")

		return &operator
	}
	return nil
}
func (operator *SyncOperator) pauseAndSaveListener(retryNum int) {
	operator.checkRetryAndPause(retryNum)
	if pause {

		saveSyncOperator(*operator)
		log.Infof("已暂停")
		wg.Add(1)
		go checkPause()
		wg.Wait()
		log.Infof("已启动")
	}
}
func (operator *SyncOperator) checkRetryAndPause(retryNum int) {
	if retryNum > 10 {
		log.Errorf("已尝试超过(%d)次，无法解决，请人工查看", retryNum)
		pause = true
	}
}
func checkPause() {
	for {
		if pause {
			time.Sleep(time.Second)
		} else {
			break
		}
	}
	wg.Done()
}
func serializeJson(g any) string {

	//把指针丢进去
	enc, _ := json.MarshalIndent(g, "", " ")
	//调用Encode进行序列化
	return string(enc)
}
func deSerializeJson(g string, e any) {

	//创建缓存
	buf := []byte(g)
	err := json.Unmarshal(buf, e)
	if err != nil {
		fmt.Println("解码失败", err.Error())
	}

}

// PathExists 判断所给路径文件/文件夹是否存在
func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	//isnotexist来判断，是不是不存在的错误
	if os.IsNotExist(err) { //如果返回的错误类型使用os.isNotExist()判断为true，说明文件或者文件夹不存在
		return false
	}
	return false
}
func getSyncListSetting() bool {
	setting, err := model.GetSettingByKey("ListSync")
	if err == nil {
		return setting.Value == "true"
	} else {

		newSetting := model.SettingItem{Key: "ListSync", Value: "true", Type: "bool", Description: "list fake operation to file instead of really do sync"}
		_ = model.SaveSetting(newSetting)
	}
	return true
}
func getSyncPathSetting(key string, defaultPath []string) []string {
	var rootPaths []string
	setting, err := model.GetSettingByKey(key)
	if err == nil {
		rootPaths = base.StrToList(setting.Value)
	} else {
		rootPaths = defaultPath
		j := base.ListToStr(rootPaths)
		log.Infof("正在缓存设置")
		newSetting := model.SettingItem{Key: key, Value: j, Type: "text"}
		_ = model.SaveSetting(newSetting)
	}
	return rootPaths
}
func ClearCache() error {
	err := conf.Cache.Clear(conf.Ctx)
	log.Info("cache has been cleared, we will redo file sync")
	return err
}
func SyncEntry(c *gin.Context) {
	err := ClearCache()
	base.InitFileHideList()
	if err != nil || pause {
		//common.ErrorResp(c, err, 500)
		common.ErrorResp(c, errors.New("暂停"), 500)
	} else {
		common.ErrorResp(c, errors.New("运行中"), 500)
	}

	if !start {
		accounts, _ := model.GetAccounts()
		syncOperator := InitSyncOperator(accounts)
		go syncOperator.syncAll()
	} else {
		pause = !pause
	}
}
