package main

import (
	_ "xiaoLianJiFen/models"
	_ "xiaoLianJiFen/routers"
	"github.com/astaxie/beego"
	"strconv"
)

func main() {
	err := beego.AddFuncMap("ShowPrePage", HandlePerPage)
	if err != nil {

	}
	err = beego.AddFuncMap("ShowNextPage", HandleNextPage)
	if err != nil {
		return
	}
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
