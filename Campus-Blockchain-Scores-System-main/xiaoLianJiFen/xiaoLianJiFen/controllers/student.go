package controllers

import (
	Models "xiaoLianJiFen/models"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	// "github.com/ethereum/go-ethereum/console"
)

type StudentController struct {
	beego.Controller
}

// 获取当前用户的角色状态
func (c *StudentController) isClubAdmin() bool {
	userID := c.GetSession("userId")
	if userID == nil {
		return false
	}

	o := orm.NewOrm()
	var user Models.Users
	err := o.QueryTable("users").Filter("username", userID).One(&user)
	if err != nil {
		beego.Error("查询用户信息失败：", err)
		return false
	}

	return user.Role_name == "社团管理员"
}

// ShowDashboard 显示学生端首页
func (c *StudentController) ShowDashboard() {
	beego.Info("进入学生端首页")
	c.Data["ActivePage"] = "home"
	c.Data["IsClubAdmin"] = c.isClubAdmin()
	c.TplName = "student_nav.html"
}

// ShowActivities 显示活动列表页面
func (c *StudentController) ShowActivities() {
	beego.Info("进入活动列表页面")
	c.Data["ActivePage"] = "activities"
	c.Data["IsClubAdmin"] = c.isClubAdmin()
	c.TplName = "student_activities.html"
}

// ShowShop 显示积分商城页面
func (c *StudentController) ShowShop() {
	beego.Info("进入积分商城页面")
	c.Data["ActivePage"] = "shop"
	c.Data["IsClubAdmin"] = c.isClubAdmin()
	c.TplName = "student_shop.html"
}

// ShowProfile 显示个人中心页面
func (c *StudentController) ShowProfile() {
	beego.Info("进入个人中心页面")
	c.Data["ActivePage"] = "profile"
	c.Data["IsClubAdmin"] = c.isClubAdmin()
	c.TplName = "student_profile.html"
}

// ShowClub 显示社团活动页面
func (c *StudentController) ShowClub() {
	beego.Info("进入社团活动页面")
	c.Data["ActivePage"] = "club"
	c.Data["IsClubAdmin"] = c.isClubAdmin()
	c.TplName = "student_club.html"
}

// ShowRules 显示积分规则页面
func (c *StudentController) ShowRules() {
	beego.Info("进入积分规则页面")
	c.Data["ActivePage"] = "rules"
	c.Data["IsClubAdmin"] = c.isClubAdmin()
	c.TplName = "student_rules.html"
}

// GetActivities 获取活动列表
func (c *StudentController) GetActivities() {
	beego.Info("开始获取活动列表")

	o := orm.NewOrm()
	var activities []Models.Activities

	// 查询所有活动
	_, err := o.QueryTable("activities").All(&activities)
	if err != nil {
		beego.Error("查询活动列表失败：", err)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "获取活动列表失败",
		}
		c.ServeJSON()
		return
	}

	beego.Info("成功获取活动列表，数量：", len(activities))
	c.Data["json"] = map[string]interface{}{
		"success":    true,
		"activities": activities,
	}
	c.ServeJSON()
}

// ApplyForAdmin 申请成为管理员
func (c *StudentController) ApplyForAdmin() {
	beego.Info("开始处理申请管理员请求")

	// 获取当前登录用户的ID
	userID := c.GetSession("userId")
	if userID == nil {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "请先登录",
		}
		c.ServeJSON()
		return
	}

	o := orm.NewOrm()
	var user Models.Users
	// 使用 username 查询用户
	err := o.QueryTable("users").Filter("username", userID).One(&user)
	if err != nil {
		beego.Error("查询用户信息失败：", err)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "用户不存在",
		}
		c.ServeJSON()
		return
	}

	// 设置申请状态为1
	user.IsAdminRequest = int8(1)
	_, err = o.Update(&user, "IsAdminRequest")
	if err != nil {
		beego.Error("更新申请状态失败：", err)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "申请失败，请稍后重试",
		}
		c.ServeJSON()
		return
	}

	c.Data["json"] = map[string]interface{}{
		"success": true,
		"message": "申请已提交，请等待审核",
	}
	c.ServeJSON()
}

// ShowStudentNav 显示学生导航页面
func (c *StudentController) ShowStudentNav() {
	// 检查用户是否登录
	userID := c.GetSession("userId")
	if userID == nil {
		c.Redirect("/", 302)
		return
	}

	// 检查用户角色
	o := orm.NewOrm()
	var user Models.Users
	err := o.QueryTable("users").Filter("username", userID).One(&user)
	if err != nil || (user.Role_name != "学生" && user.Role_name != "社团管理员") {
		c.Redirect("/", 302)
		return
	}

	// 设置导航栏当前页面和管理员状态
	c.Data["ActivePage"] = "home"
	c.Data["IsClubAdmin"] = user.Role_name == "社团管理员"
	// 渲染学生导航页面
	c.TplName = "student_nav.html"
}

// GetClubActivities 获取社团活动列表
func (c *StudentController) GetClubActivities() {
	beego.Info("开始获取社团活动列表")

	o := orm.NewOrm()
	var activities []Models.Activities

	// 查询所有活动
	_, err := o.QueryTable("activities").All(&activities)
	if err != nil {
		beego.Error("查询活动列表失败：", err)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "获取活动列表失败",
		}
		c.ServeJSON()
		return
	}

	beego.Info("成功获取活动列表，数量：", len(activities))
	c.Data["json"] = map[string]interface{}{
		"success": true,
		"data":    activities,
	}
	c.ServeJSON()
}

// SubmitActivity 提交新活动申请
func (c *StudentController) SubmitActivity() {
	beego.Info("开始处理活动申报请求")

	// 获取当前登录用户的ID
	userID := c.GetSession("userId")
	if userID == nil {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "请先登录",
		}
		c.ServeJSON()
		return
	}

	// 解析请求数据
	var activity Models.Activities
	if err := c.ParseForm(&activity); err != nil {
		beego.Error("解析表单数据失败：", err)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "提交数据格式错误",
		}
		c.ServeJSON()
		return
	}

	// 设置活动状态为待审核
	activity.Status = 0

	// 将活动保存到数据库
	o := orm.NewOrm()
	_, err := o.Insert(&activity)
	if err != nil {
		beego.Error("保存活动数据失败：", err)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "保存活动失败，请稍后重试",
		}
		c.ServeJSON()
		return
	}

	beego.Info("活动申报成功：", activity)
	c.Data["json"] = map[string]interface{}{
		"success": true,
		"message": "活动申报成功，等待审核",
		"data":    activity,
	}
	c.ServeJSON()
}

// RequestAdmin 处理学生申请成为管理员的请求
func (c *StudentController) RequestAdmin() {
	userID := c.GetSession("userId")
	if userID == nil {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "未登录",
		}
		c.ServeJSON()
		return
	}

	o := orm.NewOrm()
	var user Models.Users
	err := o.QueryTable("users").Filter("username", userID).One(&user)
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "用户不存在",
		}
		c.ServeJSON()
		return
	}

	if user.IsAdminRequest == 1 {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "已经申请过管理员",
		}
		c.ServeJSON()
		return
	}

	user.IsAdminRequest = 1
	_, err = o.Update(&user, "is_admin_request")
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "申请失败",
		}
		c.ServeJSON()
		return
	}

	c.Data["json"] = map[string]interface{}{
		"success": true,
		"message": "申请成功",
	}
	c.ServeJSON()
}
