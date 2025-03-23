package controllers

import (
    "github.com/astaxie/beego"
)

type AdminController struct {
    beego.Controller
}

// ShowDashboard 显示管理员端首页
func (c *AdminController) ShowDashboard() {
    c.TplName = "admin.html" // 设置模板为 admin.html
}