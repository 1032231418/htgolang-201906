package usersmem

import (

	"encoding/csv"

	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"crypto/md5"

	"github.com/howeyc/gopass"
)

const (
	MaxAuth  = 3
	//password = "1fdb7184e697ab9355a3f1438ddc6ef9"
    passfile = ".password"
    userfile = "users.csv"
)

type User struct {
	ID       int
	Name     string
	Birthday time.Time
	Tel      string
	Addr     string
	Desc     string
}

func (u User) String() string {
	return fmt.Sprintf("ID: %d\n名字:%s\n出生日期:%s\n联系方式:%s\n联系地址:%s\n备注:%s", u.ID, u.Name, u.Birthday.Format("2006-01-02"), u.Tel, u.Addr, u.Desc)
}




func checkFile(file string) bool {
    _, err := os.Stat(file)
    if err != nil {
        //fmt.Println("文件不存在")
        return false
    }
    return true
}

func setPass() bool {
	fmt.Print("请设置密码: ")

	bytes, _ := gopass.GetPasswd()

	pass, err := os.Create(passfile)

	if err == nil {

		_, err := pass.WriteString(fmt.Sprintf("%x", md5.Sum(bytes)))
		if err != nil {
			fmt.Println("密码设置成功.")
			return true

		}

	}
	return false

}


//var users map[int]User

//func init() {
//	users = make(map[int]User)
//}

type UserManager map[int]User
func NewUserManager() UserManager {
	return UserManager{}
}

func (m UserManager) load() {

	if file, err := os.Open(userfile);err == nil{
		defer file.Close()

		reader := csv.NewReader(file)
		for {
			line, err := reader.Read()
			if err != nil {
				if err != io.EOF {
					fmt.Println("[-]错误: ", err)
				}
				break
			}
			id,_ := strconv.Atoi(line[0])
			birthday, _ := time.Parse("2006-01-02", line[2])
			m[id] = User{
				ID: id,
				Name: line[1],
				Birthday: birthday,
				Tel: line[3],
				Addr: line[4],
				Desc: line[5],
			}
		}

	}else {
		if !os.IsNotExist(err){
			fmt.Println("[-]发生错误: ", err)
		}
	}

}

func (m UserManager) store() {
	// 重命名文件
	if _, err := os.Stat(userfile); err == nil {
		os.Rename(userfile, strconv.FormatInt(time.Now().Unix(), 10) + "." + userfile)
	}
	if names, err := filepath.Glob("*.users.csv"); err == nil {
		sort.Sort(sort.Reverse(sort.StringSlice(names)))
		fmt.Println(names)
		if len(names) > 3{
			for _, name := range names[3:] {
				os.Remove(name)
			}
		}
	}

	if file, err := os.Create(userfile); err == nil {
		defer file.Close()
		writer := csv.NewWriter(file)
		for _, user := range m {
			writer.Write([]string{
				strconv.Itoa(user.ID),
				user.Name,
				user.Birthday.Format("2006-01-02"),
				user.Tel,
				user.Addr,
				user.Desc,
			})
		}
		writer.Flush()
	}
}

func InputString(prompt string) string {
	var input string
	fmt.Print(prompt)
	fmt.Scan(&input)
	return strings.TrimSpace(input)
}

// 从命令行输入密码, 并进行验证
// 通过返回值告知验证成功还是失败
func Auth() bool {
	for i := 0; i < MaxAuth; i++ {

		// fmt.Scan(&input)

		if checkFile(passfile) {

			fmt.Print("请输入密码:")
			bytes, _ := gopass.GetPasswd()

			file, err := os.Open(passfile)

			if err == nil {

				defer file.Close()

				password, _ := ioutil.ReadAll(file)

				if string(password) == fmt.Sprintf("%x", md5.Sum(bytes)) {
					return true
				} else {
					fmt.Println("密码错误")
				}

			}

		}else {
			setPass()
		}


	}
	return false
}

func Query() {
	q := InputString("请输入查询内容:")
	list := make([]User, 0)
	fmt.Println("================================")

	users := NewUserManager()
	users.load()

	for _, v := range users {
		//name, birthday, tel, addr, desc
		if strings.Contains(v.Name, q) || strings.Contains(v.Tel, q) || strings.Contains(v.Addr, q) || strings.Contains(v.Desc, q) {
			list = append(list, v)
		}
	}
	if len(list) == 0 {
		fmt.Println("查询内容为空")
	} else {
		sortKey := InputString("请输入排序字段(id/name/tel/addr/desc):")
		sort.Slice(list, func(i, j int) bool {
			switch sortKey {
			case "id":
				return list[i].ID < list[j].ID
			case "name":
				return list[i].Name < list[j].Name
			case "tel":
				return list[i].Tel < list[j].Tel
			case "addr":
				return list[i].Addr < list[j].Addr
			case "desc":
				return list[i].Desc < list[j].Desc
			default:
				return list[i].ID < list[j].ID
			}
		})

		for _, user := range list {
			fmt.Println(user)
			fmt.Println("----------------------------")
		}
	}

	fmt.Println("================================")
}

func getId() int {
	var id int
	users := NewUserManager()
	users.load()

	for k := range users {
		if id < k {
			id = k
		}
	}
	return id + 1
}

func inputUser(id int) User {

	var user User
	user.ID = id
	user.Name = InputString("请输入名字:")
	birthday, _ := time.Parse("2006-01-02", InputString("请输入出生日期(2000-01-01):"))
	user.Birthday = birthday
	user.Tel = InputString("请输入联系方式:")
	user.Addr = InputString("请输入联系地址:")
	user.Desc = InputString("请输入备注:")
	return user
}

func validUser(user User) error {
	if user.Name == "" {
		return fmt.Errorf("输入的用户名为空")
	}

	users := NewUserManager()
	users.load()
	for _, tuser := range users {
		if user.Name == tuser.Name && user.ID != tuser.ID {
			return errors.New("输入的名字已经存在")
		}
	}
	return nil
}

func Add() {
	id := getId()
	user := inputUser(id)

	if err := validUser(user); err == nil {
		users := NewUserManager()
		users.load()
		users[id] = user
		users.store()

		fmt.Println("[+]添加成功")
	} else {
		fmt.Print("[-]添加失败:")
		fmt.Println(err)
	}
}

func Modify() {

	if id, err := strconv.Atoi(InputString("请输入修改用户ID:")); err == nil {

		users := NewUserManager()
		users.load()
		if user, ok := users[id]; ok {
			fmt.Println("将修改的用户信息:")
			fmt.Println(user)
			input := InputString("确定修改(Y/N)?")
			if input == "y" || input == "Y" {
				user := inputUser(id)

				if err := validUser(user); err == nil {
					users[id] = user
					users.store()

					fmt.Println("[+]修改成功")
				} else {
					fmt.Print("[-]修改失败:")
					fmt.Println(err)
				}
			}
		} else {
			fmt.Println("[-]用户ID不存在")
		}
	} else {
		fmt.Println("[-]输入ID不正确")
	}
}

func Del() {
	if id, err := strconv.Atoi(InputString("请输入删除用户ID:")); err == nil {
		users := NewUserManager()
		users.load()
		if user, ok := users[id]; ok {
			fmt.Println("将删除的用户信息:")
			fmt.Println(user)
			input := InputString("确定删除(Y/N)?")
			if input == "y" || input == "Y" {
				delete(users, id)

				users.store()
				fmt.Println("[+]删除成功")
			}
		} else {
			fmt.Println("[-]用户ID不存在")
		}
	} else {
		fmt.Println("[-]输入ID不正确")
	}
}

func ModifyPass() {
	password, err := ioutil.ReadFile(passfile)
	if err == nil {
		fmt.Print("请输入密码: ")
		bytes, _ := gopass.GetPasswd()
		if string(password) == fmt.Sprintf("%x", md5.Sum(bytes)){
			fmt.Print("请输入新的密码: ")
			bytes, _= gopass.GetPasswd()
			ioutil.WriteFile(passfile, []byte(fmt.Sprintf("%x", md5.Sum(bytes))), os.ModePerm)
			fmt.Println("[+]密码修改成功")
		}else {
			fmt.Println("[-]密码错误")
		}
	}
}