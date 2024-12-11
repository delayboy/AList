package controllers

import (
	"bytes"
	"fmt"
	"github.com/Xhofe/alist/conf"
	"github.com/Xhofe/alist/drivers/base"
	"github.com/Xhofe/alist/mygui"
	"github.com/Xhofe/alist/server/common"
	"github.com/Xhofe/alist/utils"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"os/exec"
	"path/filepath"
	"strings"
)

/** shell脚本示例
my_array=("Apple" "Banana" "Orange" "Grapes" "Watermelon")

# 使用for循环按索引遍历数组并打印每个元素
for ((i = 0; i < ${#my_array[@]}; i++))
do
    echo "Index $i: ${my_array[i]}"
done
*/
// DownPost /*

func DownPost(c *gin.Context) {
	buffer := make([]byte, 100)
	var bigBuffer bytes.Buffer
	var keyMap = make(map[string]int)
	var body []string
	var links []*base.Link
	var shellScript string = "if echo \"准备下载脚本\"; then\n\n%s\n\necho \"一共${#my_filenames[@]}个文件，开始下载...\"\n# 使用for循环按索引遍历数组并打印每个元素\nfor ((i = 0; i < ${#my_dirs[@]}; i++))\ndo\n  if ! test -d \"${my_dirs[i]}\"; then\n    echo \"$i 创建文件夹: ${my_dirs[i]}\"\n    mkdir -p \"${my_dirs[i]}\"\n  else\n    echo \"$i 文件夹已经存在: ${my_dirs[i]}\"\n  fi\ndone\n\nfor ((i = 0; i < ${#my_filenames[@]}; i++))\ndo\n  if ! test -f \"${my_filenames[i]}\"; then\n    echo \"($((1+i))/${#my_filenames[@]}) 文件不存在,下载文件: ${my_filenames[i]}\"\n    wget  --user-agent=\"pan.baidu.com\" \"${my_links[i]}\" -O \"${my_filenames[i]}\"\n  else\n    echo \"($((1+i))/${#my_filenames[@]}) 文件已存在跳过下载: ${my_filenames[i]}\"\n  fi\ndone\nfi"
	var linkArrayStr string
	var fileNameArrayStr string
	n, _ := c.Request.Body.Read(buffer)
	for n > 0 {
		bigBuffer.Write(buffer[:n])
		n, _ = c.Request.Body.Read(buffer)

	}
	jsonStr := bigBuffer.String()
	err := utils.Json.UnmarshalFromString(jsonStr, &body)
	if err != nil {
		fmt.Println("解码失败", err.Error())
	}
	fmt.Println("读取到的数据：", jsonStr)
	rawPath := c.Param("path")
	for index, value := range body {
		tmpPath := utils.Join(rawPath, value)
		link := dealWithDownloadLink(c, tmpPath)
		if link != nil {
			fmt.Printf("%d>%s:%s", index, tmpPath, link.Url)

			stateStr := fmt.Sprintf("echo \"%s task:(%d/%d) filepath:%s\"", "begin", index+1, len(body), value)
			mygui.GuiSetText(stateStr)
			if link.FilePath != value {
				dirPath := strings.Replace(filepath.Dir(value), "\\", "/", -1)
				keyMap[dirPath] = 1
				link.FilePath = value
			}
			linkArrayStr += fmt.Sprintf(" \"%s\" \\\r\n", link.Url)
			fileNameArrayStr += fmt.Sprintf(" \"%s\" \\\r\n", "."+value)
			links = append(links, link)
		} else {
			break
		}
	}
	var shellArray string
	for key, _ := range keyMap {
		shellArray = shellArray + fmt.Sprintf(" \".%s\" ", key)
	}
	setArrayCmd := fmt.Sprintf("my_dirs=(%s)\r\nmy_links=(%s)\r\nmy_filenames=(%s)", shellArray, linkArrayStr, fileNameArrayStr)
	shellScript = fmt.Sprintf(shellScript, setArrayCmd)
	mygui.GuiSetText(shellScript)

}
func dealWithDownloadLink(c *gin.Context, rawPath string) *base.Link {
	rawPath = utils.ParsePath(rawPath)
	log.Debugf("down: %s", rawPath)
	account, path_, driver, err := common.ParsePath(rawPath)
	if err != nil {
		common.ErrorResp(c, err, 500)
		return nil
	}
	if driver.Config().OnlyProxy || account.Proxy || utils.IsContain(conf.DProxyTypes, utils.Ext(rawPath)) {
		Proxy(c)
		return nil
	}
	link, err := driver.Link(base.Args{Path: path_, IP: c.ClientIP()}, account)
	if err != nil {
		common.ErrorResp(c, err, 500)
		return nil
	}
	link.FilePath = utils.Base(path_)
	return link

}
func Down(c *gin.Context) {
	rawPath := c.Param("path")

	link := dealWithDownloadLink(c, rawPath)
	if link != nil {
		cmd := exec.Command("idman.exe", "/n", "/a", "/d", link.Url, "/p", "C:/", "/f", link.FilePath)
		if err := cmd.Run(); err != nil {
			log.Error(err)
		}
		c.Redirect(302, link.Url)
	}
	return
}
