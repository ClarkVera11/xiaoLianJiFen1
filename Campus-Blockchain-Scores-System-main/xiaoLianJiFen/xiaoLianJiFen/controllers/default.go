package controllers

import (
	"time"
	Models "xiaoLianJiFen/models"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type MainController struct {
	beego.Controller
}

// 注册
func (c *MainController) ZhuCeGet() {
	c.TplName = "ZhuCe.html"
}

func (c *MainController) HuoQu() {
	o := orm.NewOrm()
	user := Models.Users{}

	userName := c.GetString("用户名")
	userPwd := c.GetString("密码")
	userEmail := c.GetString("邮箱")
	userPhone := c.GetString("手机号")

	// 判断数据是否合法，否则提示且循环
	if userName == "" || userPwd == "" || userEmail == "" || userPhone == "" {
		beego.Info("数据不能为空！！！")
		c.Redirect("/ZhuCe", 302)
		return
	}

	// 根据用户名判断身份
	var userRole string
	if len(userName) == 9 && userName[:4] == "2022" {
		userRole = "学生"
	} else if len(userName) >= 2 && userName[:2] == "JX" {
		userRole = "教师"
	} else {
		beego.Info("用户名格式错误！")
		c.Redirect("/ZhuCe", 302)
		return
	}

	// 1. 拿到数据
	err1 := o.Read(&user, "Username")
	// 2. 判断数据合法
	if err1 == nil {
		beego.Info("用户名已存在！")
		c.Redirect("/ZhuCe", 302)
		return
	}

	// 插入
	user.Username = userName
	user.Password = userPwd
	user.Email = userEmail
	user.Phone = userPhone
	user.Role_name = userRole
	user.Points = 0 // 初始化为0，后面通过UpdatePoints设置初始积分

	_, err := o.Insert(&user)
	if err != nil {
		beego.Info("插入失败！")
		c.Redirect("/ZhuCe", 302)
		return
	}

	// 设置初始积分为60
	user.Points = 60
	_, err = o.Update(&user, "Points")
	if err != nil {
		beego.Info("设置初始积分失败！")
		c.Redirect("/ZhuCe", 302)
		return
	}

	c.Redirect("/", 302)
}

// 登入
func (c *MainController) ShowLogin() {
	c.TplName = "DengRu.html"
}

// HandleLogin 处理登录请求
func (c *MainController) HandleLogin() {
	// 1. 获取表单数据
	userName := c.GetString("用户名")
	pwd := c.GetString("密码")

	// 2. 数据验证
	if userName == "" || pwd == "" {
		beego.Info("数据输入不合法")
		c.TplName = "DengRu.html"
		return
	}

	// 3. 查询数据库
	o := orm.NewOrm()
	var user Models.Users
	err := o.QueryTable("users").Filter("username", userName).Filter("password", pwd).One(&user)
	if err != nil {
		beego.Info("用户名或密码错误！！！")
		c.TplName = "DengRu.html"
		return
	}

	// 4. 登录成功，设置 Cookie 和 Session
	c.Ctx.SetCookie("Username", userName, time.Second*3600)
	c.SetSession("userId", userName)

	// 5. 根据角色跳转到不同页面
	beego.Info("当前用户角色：", user.Role_name)
	switch user.Role_name {
	case "学生", "社团管理员":
		beego.Info("学生/社团管理员登录成功，跳转到学生导航页面")
		c.Redirect("/student", 302)
	case "教师":
		beego.Info("教师登录成功，跳转到教师导航页面")
		c.Redirect("/teacher", 302)
	case "管理员":
		beego.Info("管理员登录成功，跳转到管理员端")
		c.Redirect("/admin", 302)
	default:
		beego.Info("未知角色，跳转到首页")
		beego.Info("用户角色：", user.Role_name)
		c.Redirect("/", 302)
	}
}

// 重置密码
func (c *MainController) ShowChongZhiMiMa() {
	c.TplName = "ChongZhiMiMa.html"
}

func (c *MainController) HandleChongZhiMiMa() {
	// 初始化ORM
	o := orm.NewOrm()
	user := Models.Users{}

	// 获取表单数据
	userName := c.GetString("学号")
	userEmail := c.GetString("邮箱")
	userPhone := c.GetString("手机号")

	// 记录日志
	beego.Info(userName, userEmail, userPhone)

	// 数据验证：检查必填项是否为空
	if userName == "" || userEmail == "" || userPhone == "" {
		beego.Info("数据不能为空！！！")
		c.Redirect("/ChongZhiMiMa", 302)
		return
	}
	user.Username = userName
	user.Email = userEmail
	user.Phone = userPhone

	// 先通过用户名查询用户数据
	err := o.Read(&user, "Username")
	if err != nil {
		beego.Info("用户名不存在！")
		c.Redirect("/ChongZhiMiMa", 302)
		return
	}
	// 验证邮箱和手机号是否匹配
	if user.Email != userEmail || user.Phone != userPhone {
		beego.Info("邮箱或手机号不匹配！！！")
		c.Redirect("/ChongZhiMiMa", 302)
		return
	}

	// 将密码重置为默认值
	user.Password = "123456"

	// 更新数据库中的用户信息
	_, err = o.Update(&user)
	if err != nil {
		beego.Info("密码重置失败", err)
		c.TplName = "ChongZhiMiMa.html"
		return
	} else {
		// 重置成功，跳转到首页
		beego.Info("密码重置成功！")
		c.Redirect("/", 302)
	}
}

// Logout 处理退出请求
func (c *MainController) Logout() {
	// 清除 Session
	c.DelSession("UserName")
	// 清除 Cookie
	c.Ctx.SetCookie("Username", "", -1) // 设置 Cookie 过期时间为 -1，表示删除

	// 重定向到登录页面
	c.Redirect("/", 302)
}

// UpdateStudentToAdmin 更新学生角色为管理员
func (c *MainController) UpdateStudentToAdmin() {
	studentID := c.GetString("studentID")

	o := orm.NewOrm()
	var user Models.Users
	err := o.QueryTable("users").Filter("username", studentID).One(&user)
	if err != nil {
		beego.Error("读取用户数据失败: ", err)
		c.Data["json"] = map[string]interface{}{"success": false, "message": "读取用户数据失败"}
		c.ServeJSON()
		return
	}

	if user.Role_name != "学生" {
		beego.Error("该用户不是学生，无法更新为管理员")
		c.Data["json"] = map[string]interface{}{"success": false, "message": "该用户不是学生，无法更新为管理员"}
		c.ServeJSON()
		return
	}

	// 更新学生角色为管理员
	user.Role_name = "管理员"
	user.IsAdminRequest = int8(0) // 清除申请状态
	_, err = o.Update(&user, "Role_name", "IsAdminRequest")
	if err != nil {
		beego.Error("更新学生角色失败: ", err)
		c.Data["json"] = map[string]interface{}{"success": false, "message": "设置管理员失败"}
		c.ServeJSON()
		return
	}

	// 返回成功结果
	beego.Info("成功将学生设置为管理员, 学生 ID: ", studentID)
	c.Data["json"] = map[string]interface{}{"success": true, "message": "已将该学生设置为管理员"}
	c.ServeJSON()
}
