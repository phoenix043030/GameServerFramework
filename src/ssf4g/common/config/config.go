package config

import (
	"os"
	"path/filepath"
)

var (
	_ini_data *IniData
)

// Func - 初始化配置数据
func init() {
	ReloadConfig()
}

// Func - 重载配置数据
func ReloadConfig() {
	appPath, err := filepath.Abs(filepath.Dir(os.Args[0]))

	if err != nil {
		panic(err)
	}

	workPath, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	// conf.ini存在判断
	appConfigPath := filepath.Join(workPath, "conf.ini")

	_, err = os.Stat(appConfigPath)

	if err != nil && os.IsNotExist(err) == true {
		appConfigPath = filepath.Join(appPath, "conf.ini")

		_, err = os.Stat(appConfigPath)

		if err != nil && os.IsNotExist(err) == true {
			panic(err)
		}
	}

	iniMgr := &IniMgr{}

	_ini_data, err = iniMgr.Parse(appConfigPath)

	if err != nil {
		panic(err)
	}
}

// Func - 获取Ini配置数据
func GetIniData() *IniData {
	return _ini_data
}
