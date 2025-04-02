package main

import (
	"strconv"
	_ "xiaoLianJiFen/models"
	_ "xiaoLianJiFen/routers"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func add1(a int) int {
	return a + 1
}

func init() {
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", "root:123456@tcp(127.0.0.1:3306)/xiaoLianJiFen?charset=utf8")
}

func main() {
	err := beego.AddFuncMap("ShowPrePage", HandlePerPage)
	if err != nil {

	}
	err = beego.AddFuncMap("ShowNextPage", HandleNextPage)
	if err != nil {
		return
	}
	beego.AddFuncMap("add1", add1)
	beego.Run()
}

func HandlePerPage(data int) string {

	pageIndex := data - 1

	pageIndex1 := strconv.Itoa(pageIndex)
	return pageIndex1
}

func HandleNextPage(data int) string {

	pageIndex := data + 1
	pageIndex1 := strconv.Itoa(pageIndex)
	return pageIndex1
}
