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

// ShowBlockchainInfo displays the blockchain information page
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

// FormatTimestamp converts Unix timestamp to readable format
func (c *BlockchainController) FormatTimestamp(timestamp int64) string {
	if timestamp == 0 {
		return "未上链"
	}
	t := time.Unix(timestamp, 0)
	return t.Format("2006-01-02 15:04:05")
}
