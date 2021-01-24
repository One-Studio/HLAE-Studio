package config

import (
	"HLAE-Studio/backend/tool"
	"bytes"
	"encoding/json"
	"os"
)

type CFG struct {
	VersionCode   string //大版本代号
	AppVersion    string //app程序版本
	HlaeVersion   string //hlae版本
	FFmpegVersion string //ffmpeg版本
	HlaeAPI       string //hlae 版本号+下载地址等 API
	HlaeCdnAPI    string //hlae CDN镜像 API
	FFmpegAPI     string //ffmpeg 版本号+下载地址等 API
	FFmpegCdnAPI  string //ffmpeg CDN镜像 API
	HlaePath      string //hlae安装位置
	Init          bool   //程序当前状态 是否已初始化
	Standalone    bool   //是否单独安装hlae+ffmpeg
	HlaeState     bool   //hlae安装状态
	FFmpegState   bool   //ffmpeg安装状态
}

//读设置
func ReadConfig(path string) (CFG, error) {
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

		//检查现有设置，做一定语法上的修正，处理版本更新带来的设置选项的变化
		CFGInst, err = checkConfig(CFGInst)
		if err != nil {
			return CFG{}, err
		}

		return CFGInst, nil
	} else {
		//设置文件不存在则初始化
		return defaultCFG, nil
	}
}

//写设置
func SaveConfig(cfg CFG, path string) error {
	//检查文件是否存在
	exist, err := tool.IsFileExisted(path)
	if err != nil {
		//路径等错误，返回err
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

	JsonData, err := Config2Json(cfg)
	err = tool.WriteFast(path, JsonData)
	if err != nil {
		return err
	}

	return nil
}

//设置转Json字符串
func Config2Json(cfg CFG) (string, error) {
	JsonData, err := json.Marshal(cfg) //第二个参数要地址传递
	if err != nil {
		return "", err
	}

	//json写入文件
	var str bytes.Buffer
	_ = json.Indent(&str, JsonData, "", "    ")
	return str.String(), nil
}

//检查设置，更新覆盖部分设置
func checkConfig(cfg CFG) (CFG, error) {
	cfg.VersionCode = defaultCFG.VersionCode
	cfg.AppVersion = defaultCFG.AppVersion
	cfg.HlaeAPI = defaultCFG.HlaeAPI
	cfg.HlaeCdnAPI = defaultCFG.HlaeCdnAPI
	cfg.FFmpegAPI = defaultCFG.FFmpegAPI
	cfg.FFmpegCdnAPI = defaultCFG.FFmpegCdnAPI

	return cfg, nil
}
