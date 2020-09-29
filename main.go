package main

import (
	"github.com/astaxie/beego"
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
//测试用数据命名为struct.binary的二进制数据文件
//测试用数据的KeySize\ValueSize的大小为uint的大小即64位，8byte。
func init(){
	//加载数据到B+树结构中

}