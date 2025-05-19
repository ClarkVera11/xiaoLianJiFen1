package controllers

import (
	"time"
	Models "xiaoLianJiFen/models"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type BlockchainController struct {
	beego.Controller
}

// ShowBlockchainInfo displays the blockchain information page for students and admins
func (c *BlockchainController) ShowBlockchainInfo() {
	beego.Info("进入区块链信息页面")

	// 检查用户是否登录
	userID := c.GetSession("userId")
	if userID == nil {
		c.Redirect("/", 302)
		return
	}

	// 获取当前用户信息
	o := orm.NewOrm()
	var currentUser Models.Users
	err := o.QueryTable("users").Filter("username", userID).One(&currentUser)
	if err != nil {
		c.Redirect("/", 302)
		return
	}

	// 判断角色
	role := currentUser.Role_name
	isAdmin := role == "admin" || role == "club_admin" || role == "社团管理员"
	isStudent := role == "student" || role == "学生"

	// 若不是社团管理员或学生，重定向
	if !isAdmin && !isStudent {
		c.Redirect("/", 302)
		return
	}

	// 若不是管理员或学生，重定向
	if !isAdmin && !isStudent {
		c.Redirect("/", 302)
		return
	}

	// 设置模板数据
	c.Data["Role"] = role
	c.Data["ActivePage"] = "blockchain"
	c.Data["IsClubAdmin"] = isAdmin
	c.Data["IsStudent"] = isStudent

	// 获取所有用户链上信息
	var users []Models.Users
	_, err = o.QueryTable("users").All(&users)
	if err != nil {
		beego.Error("获取用户列表失败：", err)
		c.Data["Error"] = "获取用户列表失败"
		c.TplName = "blockchain_info.html"
		return
	}

	c.Data["Users"] = users
	c.TplName = "blockchain_info.html"
}

// ShowTeacherBlockchainInfo displays the blockchain information page for teachers
func (c *BlockchainController) ShowTeacherBlockchainInfo() {
	beego.Info("进入教师区块链信息页面")

	// 检查是否登录
	userID := c.GetSession("userId")
	if userID == nil {
		beego.Info("未登录，跳转首页")
		c.Redirect("/", 302)
		return
	}

	// 获取当前用户信息
	o := orm.NewOrm()
	var currentUser Models.Users
	err := o.QueryTable("users").Filter("username", userID).One(&currentUser)
	if err != nil {
		beego.Info("查不到用户，跳转首页")
		c.Redirect("/", 302)
		return
	}

	// 判断是否为教师
	role := currentUser.Role_name
	if role != "teacher" && role != "教师" {
		beego.Info("不是教师，跳转首页, 当前角色:", role)
		c.Redirect("/", 302)
		return
	}

	// 设置模板数据
	c.Data["Role"] = role
	c.Data["ActivePage"] = "blockchain"

	// 获取搜索条件
	searchQuery := c.GetString("search")

	// 查询用户信息（支持搜索）
	var users []Models.Users
	qs := o.QueryTable("users")
	if searchQuery != "" {
		qs = qs.Filter("username__icontains", searchQuery)
	}
	_, err = qs.All(&users)
	if err != nil {
		beego.Error("获取用户列表失败：", err)
		c.Data["Error"] = "获取用户列表失败"
		c.TplName = "teacher_blockchain_info.html"
		return
	}

	c.Data["Users"] = users
	c.TplName = "teacher_blockchain_info.html"
}

// FormatTimestamp converts Unix timestamp to readable format
func (c *BlockchainController) FormatTimestamp(timestamp int64) string {
	if timestamp == 0 {
		return "未上链"
	}
	t := time.Unix(timestamp, 0)
	return t.Format("2006-01-02 15:04:05")
}
