package main

import (
	"fmt"

	"github.com/xlotz/gopkg"
)

// 用户管理系统界面
func initUser() {
	// 存储用户的信息
	users := make(map[string]map[string]string)
	id := 0

END:
	for {
		var op string
		fmt.Print(`请输入操作指令:
		1：add
		2: modify
		3: delete
		4: select
		5: exit
		     `)

		fmt.Scan(&op)
		switch op {
		case "1":
			id++
			gopkg.AddUser(id, users)
		case "2":
			gopkg.ModifyUser(users)
		case "3":
			gopkg.DeleteUser(users)
		case "4":
			gopkg.SelectUser(users)
		case "5":
			break END
		default:
			fmt.Println("指令错误！")
		}
	}
}

const (
	pass    = "123456"
	maxAuth = 3
)

func main() {

	if !gopkg.AuthPass(pass, maxAuth) {
		fmt.Printf("密码%d次错误， 程序退出\n", maxAuth)
		return
	}

	initUser()

}

/*
评分: 8
建议：
1. 内置密码使用md5形式，避免看到源码知晓输入密码
*/
