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

// ShowBlockchainInfo displays the blockchain information page for students
func (c *BlockchainController) ShowBlockchainInfo() {
	// Check if user is logged in
	userID := c.GetSession("userId")
	if userID == nil {
		c.Redirect("/", 302)
		return
	}

	// Get current user's role
	o := orm.NewOrm()
	var currentUser Models.Users
	err := o.QueryTable("users").Filter("username", userID).One(&currentUser)
	if err != nil {
		c.Redirect("/", 302)
		return
	}

	// Set role for navigation
	c.Data["Role"] = currentUser.Role_name
	c.Data["ActivePage"] = "blockchain"
	c.Data["IsClubAdmin"] = currentUser.Role_name == "admin" || currentUser.Role_name == "club_admin"

	// Get all users with blockchain information
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
	// Check if user is logged in
	userID := c.GetSession("userId")
	beego.Info("session userId:", userID)
	if userID == nil {
		beego.Info("未登录，跳转首页")
		c.Redirect("/", 302)
		return
	}

	// Get current user's role
	o := orm.NewOrm()
	var currentUser Models.Users
	err := o.QueryTable("users").Filter("username", userID).One(&currentUser)
	beego.Info("查到用户:", currentUser, "err:", err)
	if err != nil {
		beego.Info("查不到用户，跳转首页")
		c.Redirect("/", 302)
		return
	}

	// Check if user is a teacher
	if currentUser.Role_name != "教师" && currentUser.Role_name != "teacher" {
		beego.Info("不是教师，跳转首页, 当前角色:", currentUser.Role_name)
		c.Redirect("/", 302)
		return
	}

	// Set role for navigation
	c.Data["Role"] = currentUser.Role_name
	c.Data["ActivePage"] = "blockchain"

	// Get search query
	searchQuery := c.GetString("search")

	// Get users with blockchain information
	var users []Models.Users
	qs := o.QueryTable("users")

	// Apply search filter if search query exists
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
