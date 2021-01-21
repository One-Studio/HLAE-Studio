package config

import (
	"HLAE-Studio/backend/tool"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
)

type API struct {
	win   string
	mac   string
	linux string
}

type CFG struct {
	version       string
	srcPath       string
	dstPath       string
	param         string
	api           API
	ffmpegPath    string
	ffmpegVersion string
	ffmpegRegExp  string
	Init		  bool
}

func ReadConfig(cfg CFG) {
	fmt.Println("hello")
	//读取操作 TODO
	cfg = defaultCFG
	//程序设置初始化到前端
	
	fmt.Println(cfg)
}

func SaveConfig(cfg CFG) {
	fmt.Println(cfg)
	//保存操作 TODO

}

func readCFGs(path string) (CFG, error) {
	//检查文件是否存在
	exist, err := tool.IsFileExisted(path)
	if err != nil {
		return CFG{}, err
	} else if exist == true {
		//存在则读取文件
		content, err := tool.ReadAll(path)
		if err != nil {
			return CFG{}, err
		}

		//初始化实例并解析JSON
		var CFGInst CFG
		err = json.Unmarshal([]byte(content), &CFGInst) //第二个参数要地址传递
		if err != nil {
			return CFG{}, err
		}

		return CFGInst, nil
	} else {

		return CFG{}, nil
	}
}

func saveCFGs(cfg CFG,path string) error {
	//检查文件是否存在
	exist, err := tool.IsFileExisted(path)
	if err != nil {
		return err
	} else if exist == true {
		//存在则删除文件
		ok, err := tool.IsFileExisted(path)
		if err != nil {
			return err
		} else if ok == true {
			err := os.Remove(path)
			if err != nil {
				return err
			}
		}
	}

	JsonData, err := json.Marshal(cfg) //第二个参数要地址传递
	if err != nil {
		return err
	}

	//json.Indent(JsonData, )
	var str bytes.Buffer
	_ = json.Indent(&str, JsonData, "", "    ")
	//fmt.Println("formated: ", str.String())

	err = tool.WriteFast(path, str.String())
	if err != nil {
		return err
	}

	return nil
}