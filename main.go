package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"u3.com/u3query/hack"
	_ "u3.com/u3query/routers"
)


func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.Run()
}


//用于初始化解析数据文件，解析后原数据文件删了，然后使用B+树来保存数据
func init(){
	defer func() {
		if p := recover(); p != nil {
			logs.Error("main init Panic")
			logs.Error(p)
		}
	}()
	//生成测试数据
	hack.GenerateTestDataFile()
	//加载数据到B+树结构中,并分片保存
	hack.GenerateBtree()
}