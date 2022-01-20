package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"
	"strconv"
)

var documents = "Documents\\tmp\\"
var path = ""
var fileName = "config.json"

var defaultAppName = "默认程序" //
var defaultAppPath = "自己去生成的文件配置文件下修改或者添加程序" //

//定义配置文件解析后的结构
type AppInfo struct {
	Name string
	Path string
}


func getHomeDir() string{
	user, err := user.Current()
	if err != nil {
		return ""
	}
	return user.HomeDir
}

func init(){
	homeDir := getHomeDir()
	if homeDir == ""{
		return
	}
	path = homeDir + "\\" + documents
}

func (m AppInfo) ToString() string {
	return "name:" + m.Name + ",path:" + m.Path
}

type AppConfig struct {
	Id      int
	AppInfo *AppInfo
}

func (c AppConfig) ToString() {
	log.Println("id:", c.Id, ",appInfo:", c.AppInfo.ToString())
}

func main() {
	if path == documents {
		return
	}
	fmt.Println("程序配置文件在 :",path)
	JsonParse := NewJsonStruct()
	vs := make([]*AppConfig, 0)
	//下面使用的是相对路径，config.json文件和main.go文件处于同一目录下
	JsonParse.Load(path + fileName, &vs)
	if len(vs) == 0 {
		vs = append(vs, &AppConfig{1, &AppInfo{defaultAppName,defaultAppPath}})
		bytes, err := json.Marshal(vs)
		if err != nil {
			log.Println(err.Error())
			return
		}
		JsonParse.CreateFile(path, fileName, bytes)
	}
	var exepath = ""
	if len(vs) == 1 {
		exepath = vs[0].AppInfo.Path
	} else {
		sint := SelectEXE(vs)
		if sint == -1 {
			return
		}
		exepath = vs[sint].AppInfo.Path
	}

	cmd := exec.Command("cmd.exe", "/c", "start "+exepath)
	if err := cmd.Run(); err != nil {
		fmt.Println("启动失败，原因百度:",err.Error())
	}

}

func SelectEXE(ac []*AppConfig) int {
	fmt.Println("\n请选择执行程序:")
	for i, v := range ac {
		fmt.Println(i, ":", v.AppInfo.Name)
	}
	fmt.Println("选择:")

	var qType string
	_, err := fmt.Scan(&qType)
	qTypeInt, err := strconv.Atoi(qType)
	if err != nil {
		log.Println("选择错误")
		return -1
	}

	return qTypeInt
}

type JsonStruct struct {
}

func NewJsonStruct() *JsonStruct {
	return &JsonStruct{}
}

func (jst *JsonStruct) Load(pathAll string, v interface{}) {
	if !PathExists(pathAll){
		return
	}
	//ReadFile函数会读取文件的全部内容，并将结果以[]byte类型返回
	byte, err := ioutil.ReadFile(pathAll)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	//读取的数据为json格式，需要进行解码
	err = json.Unmarshal(byte, v)
	if err != nil {
		return
	}
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func (jst *JsonStruct) CreateFile(path string, fileName string, bytes []byte) {
	if !PathExists(path) {
		err := os.Mkdir(path, 0777)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}

	f, err := os.Create(path + fileName)
	defer f.Close()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	_, err = f.Write(bytes)
	if err != nil {
		fmt.Println(err.Error())
	}
}
