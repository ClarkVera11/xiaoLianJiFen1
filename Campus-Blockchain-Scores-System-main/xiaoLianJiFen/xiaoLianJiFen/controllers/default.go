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
	// 对密码进行加密
	hashedPwd, err := Models.HashPassword(userPwd)
	if err != nil {
		beego.Info("密码加密失败！")
		c.Redirect("/ZhuCe", 302)
		return
	}
	user.Password = hashedPwd
	user.Email = userEmail
	user.Phone = userPhone
	user.Role_name = userRole
	user.Points = 0 // 初始化为0，后面通过UpdatePoints设置初始积分

	_,err = o.Insert(&user)
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
	o := orm.NewOrm()

	// 获取最新的三个活动，按开始时间倒序
	var activities []Models.Activities
	_, err := o.QueryTable("activities").
		OrderBy("-start_time").
		Limit(3).
		All(&activities)
	if err != nil {
		beego.Error("获取活动失败：", err)
	}

	// 获取积分排名前三的学生，按积分倒序
	var topStudents []Models.Users
	_, err = o.QueryTable("users").
		Filter("role_name", "学生").
		OrderBy("-points").
		Limit(3).
		All(&topStudents)
	if err != nil {
		beego.Error("获取积分排名失败：", err)
	}

	// 传递数据给模板
	c.Data["Activities"] = activities
	c.Data["TopStudents"] = topStudents

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

	// 3. 查询数据库（只查用户名）
	o := orm.NewOrm()
	
	var user Models.Users
	err := o.QueryTable("users").Filter("username", userName).One(&user)
	if err != nil {
		beego.Info("用户名不存在")
		c.TplName = "DengRu.html"
		return
	}

	// 4. 比较加密密码
	if !Models.CheckPasswordHash(user.Password, pwd) {
		beego.Info("密码错误")
		c.TplName = "DengRu.html"
		return
	}

	// 5. 登录成功，设置 Cookie 和 Session
	c.Ctx.SetCookie("Username", userName, time.Second*3600)
	c.SetSession("userId", userName)

	// 6. 跳转页面
	beego.Info("当前用户角色：", user.Role_name)
	switch user.Role_name {
	case "学生", "社团管理员":
		c.Redirect("/student", 302)
	case "教师":
		c.Redirect("/teacher", 302)
	case "管理员":
		c.Redirect("/admin", 302)
	default:
		c.Redirect("/", 302)
	}
}


// 重置密码
func (c *MainController) ShowChongZhiMiMa() {
	c.TplName = "ChongZhiMiMa.html"
}

func (c *MainController) HandleChongZhiMiMa() {
	o := orm.NewOrm()
	user := Models.Users{}

	userName := c.GetString("学号")
	userEmail := c.GetString("邮箱")
	userPhone := c.GetString("手机号")

	if userName == "" || userEmail == "" || userPhone == "" {
		beego.Info("数据不能为空！！！")
		c.Redirect("/ChongZhiMiMa", 302)
		return
	}

	user.Username = userName
	err := o.Read(&user, "Username")
	if err != nil {
		beego.Info("用户名不存在！")
		c.Redirect("/ChongZhiMiMa", 302)
		return
	}

	if user.Email != userEmail || user.Phone != userPhone {
		beego.Info("邮箱或手机号不匹配！！！")
		c.Redirect("/ChongZhiMiMa", 302)
		return
	}

	// 对默认密码进行加密
	hashedPwd, err := Models.HashPassword("123456")
	if err != nil {
		beego.Info("密码加密失败", err)
		c.Redirect("/ChongZhiMiMa", 302)
		return
	}

	user.Password = hashedPwd
	_, err = o.Update(&user)
	if err != nil {
		beego.Info("密码重置失败", err)
		c.TplName = "ChongZhiMiMa.html"
		return
	}

	beego.Info("密码重置成功！")
	c.Redirect("/", 302)
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
