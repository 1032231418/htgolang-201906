package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/xlotz/gocmdb/server/models"
)

func main() {
	// 初始化命令行参数
	h := flag.Bool("h", false, "help")
	help := flag.Bool("help", false, "help")
	//init := flag.Bool("init", false, "init server")
	//syncdb := flag.Bool("syncdb", false, "sync db")
	//force := flag.Bool("force", false, "force sync db(drop table)")
	//verbose := flag.Bool("v", false, "verbose")

	flag.Usage = func() {
		fmt.Println("usage: rsync_cloud -h")
		flag.PrintDefaults()
	}
	// 解析命令行参数
	flag.Parse()

	if *h || *help {
		flag.Usage()
		os.Exit(0)
	}

	// 设置日志到文件
	beego.SetLogger("file", `{
		"filename" : "logs/rsync_cloud.log",
		"level" : 7}`,
	)
	//if !*verbose {
	//	//删除控制台日志
	//	beego.BeeLogger.DelLogger("console")
	//
	//} else {
	//	orm.Debug = true
	//}

	// 初始化orm
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", beego.AppConfig.String("dsn"))

	// 测试数据库连接是否正常
	if db, err := orm.GetDB(); err != nil || db.Ping() != nil {
		beego.Error("数据库连接错误")
		os.Exit(-1)
	}

	for now := range time.Tick(2* time.Second){

		fmt.Println(now)

		platforms, _, _ := models.DefaultCloudPlatformManager.Query("", 0,0)


		for _, platform := range platforms {
			if !platform.IsEnable(){
				continue
			}
			fmt.Println(platform)
		}
	}

}
