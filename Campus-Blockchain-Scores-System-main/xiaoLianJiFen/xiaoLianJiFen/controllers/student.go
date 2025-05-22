package controllers

import (
	"fmt"
	"regexp"
	"strings"
	"xiaoLianJiFen/blockchain"
	Models "xiaoLianJiFen/models"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"

	// "github.com/ethereum/go-ethereum/console"
	"os"
	"time"

	// 由于报错信息显示导入包时出现问题，可能是依赖未正确安装或者版本不兼容。
	// 尝试先清理缓存并重新安装依赖。
	// 请在终端执行以下命令：
	// go clean -modcache
	// go mod tidy
	// 如果问题仍然存在，可能需要检查 go-ethereum 版本是否兼容。
	// 以下是暂时保留导入语句，等待依赖问题解决后的代码：
	"github.com/ethereum/go-ethereum/accounts/keystore"
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

    // 获取当前用户积分
    userID := c.GetSession("userId")
    if userID != nil {
        o := orm.NewOrm()
        var user Models.Users
        err := o.QueryTable("users").Filter("username", userID).One(&user)
        if err == nil {
            c.Data["UserPoints"] = user.Points
            c.Data["UserTitle"] = user.Title
			c.Data["ActivityCount"] = user.ActivityCount
        } else {
            beego.Error("查询用户信息失败:", err)
        }
    } else {
        beego.Error("未获取到用户ID")
    }

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

	// 新增：获取当前用户积分
	userID := c.GetSession("userId")
	if userID != nil {
		o := orm.NewOrm()
		var user Models.Users
		err := o.QueryTable("users").Filter("username", userID).One(&user)
		if err == nil {
			c.Data["UserPoints"] = user.Points
		} else {
			c.Data["UserPoints"] = 0
		}
	} else {
		c.Data["UserPoints"] = 0
	}

	c.TplName = "student_shop.html"
}

// ShowProfile 显示个人中心页面
func (c *StudentController) ShowProfile() {
	beego.Info("进入个人中心页面")
	c.Data["ActivePage"] = "profile"
	c.Data["IsClubAdmin"] = c.isClubAdmin()

	// 获取当前用户信息
	userID := c.GetSession("userId")
	if userID != nil {
		o := orm.NewOrm()
		var user Models.Users
		err := o.QueryTable("users").Filter("username", userID).One(&user)
		if err == nil {
			c.Data["User"] = user
		}
	}

	c.TplName = "student_profile.html"
}

// ShowClub 显示社团活动页面
func (c *StudentController) ShowClub() {
	beego.Info("进入社团活动页面")
	c.Data["ActivePage"] = "club"
	c.Data["IsClubAdmin"] = c.isClubAdmin()
	c.TplName = "student_club.html"
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

	// 获取当前用户信息
	o := orm.NewOrm()
	var user Models.Users
	err := o.QueryTable("users").Filter("username", userID).One(&user)
	if err != nil || (user.Role_name != "学生" && user.Role_name != "社团管理员") {
		c.Redirect("/", 302)
		return
	}

	// 设置基础数据
	c.Data["ActivePage"] = "home"
	c.Data["IsClubAdmin"] = user.Role_name == "社团管理员"
	c.Data["UserPoints"] = user.Points
	c.Data["UserTitle"] = user.Title
	c.Data["ActivityCount"] = user.ActivityCount
	c.Data["IsAdminRequest"] = user.IsAdminRequest == 1

	var allUsers []Models.Users
	_, err = o.QueryTable("users").Filter("role_name__in", "学生", "社团管理员").OrderBy("-points").All(&allUsers)
	if err == nil {
		total := len(allUsers)
		rank := 0
		lastPoints := -1
		userRank := 0

		for i, u := range allUsers {
			if u.Points != lastPoints {
				rank = i + 1
				lastPoints = u.Points
			}
			if u.Username == user.Username {
				userRank = rank
				break
			}
		}

		var percent float64
		if total > 1 {
			percent = float64(total - userRank) / float64(total - 1) * 100
		} else {
			percent = 100
		}

		c.Data["UserRank"] = userRank
		c.Data["TotalStudents"] = total
		c.Data["PercentSurpassed"] = fmt.Sprintf("%.0f", percent)

		if percent >= 90 {
			c.Data["Encouragement"] = fmt.Sprintf("你已超过 %.0f%% 的同学，太厉害了！", percent)
			c.Data["EncouragementColor"] = "#27ae60" // 绿色
		} else {
			c.Data["Encouragement"] = fmt.Sprintf("你已超过 %.0f%% 的同学，请继续努力！", percent)
			c.Data["EncouragementColor"] = "#e67e22" // 橙色
		}
	}
// 查询当前用户最新三条积分记录
	var pointsRecords []Models.PointsRecord
	_, err = o.QueryTable("points_record").
		Filter("user_id", user.Id).
		OrderBy("-created_at").
		Limit(3).
		All(&pointsRecords)
	if err != nil {
		pointsRecords = []Models.PointsRecord{}
	}

	c.Data["LatestPointsRecords"] = pointsRecords
	// 渲染页面
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

	// 获取前端传入的时间字段
	startTimeStr := c.GetString("StartTime")
	endTimeStr := c.GetString("EndTime")

	// 定义时间格式（匹配前端的格式）
	layout := "2006-01-02T15:04"

	// 解析开始时间和结束时间
	startTime, err := time.Parse(layout, startTimeStr)
	if err != nil {
		beego.Error("解析开始时间失败：", err)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "无效的开始时间格式",
		}
		c.ServeJSON()
		return
	}

	endTime, err := time.Parse(layout, endTimeStr)
	if err != nil {
		beego.Error("解析结束时间失败：", err)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "无效的结束时间格式",
		}
		c.ServeJSON()
		return
	}

	// 增加8小时
	eightHours := 8 * time.Hour
	startTime = startTime.Add(-eightHours)
	endTime = endTime.Add(-eightHours)

	// 将解析后的时间赋值给活动对象
	activity.StartTime = startTime
	activity.EndTime = endTime

	// 设置活动状态为待审核
	activity.Status = 0

	// 将活动保存到数据库
	o := orm.NewOrm()
	_, err = o.Insert(&activity)
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

// RegisterActivity 处理活动报名
func (c *StudentController) RegisterActivity() {
	beego.Info("开始处理活动报名请求")

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

	// 获取活动ID
	activityId, err := c.GetInt64("activityId")
	if err != nil {
		beego.Error("获取活动ID失败：", err)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "无效的活动ID",
		}
		c.ServeJSON()
		return
	}

	o := orm.NewOrm()
	// 开启事务
	err = o.Begin()
	if err != nil {
		beego.Error("开启事务失败：", err)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "系统错误，请稍后重试",
		}
		c.ServeJSON()
		return
	}

	// 查询活动信息（使用 FOR UPDATE 锁定记录）
	var activity Models.Activities
	err = o.Raw("SELECT * FROM activities WHERE id = ? FOR UPDATE", activityId).QueryRow(&activity)
	if err != nil {
		o.Rollback()
		beego.Error("查询活动信息失败：", err)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "活动不存在",
		}
		c.ServeJSON()
		return
	}

	// 检查活动状态
	if activity.Status != 1 {
		o.Rollback()
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "该活动不在可报名状态",
		}
		c.ServeJSON()
		return
	}

	// 查询用户信息
	var user Models.Users
	err = o.QueryTable("users").Filter("username", userID).One(&user)
	if err != nil {
		o.Rollback()
		beego.Error("查询用户信息失败：", err)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "用户信息获取失败",
		}
		c.ServeJSON()
		return
	}

	// 检查是否已报名
	exist := o.QueryTable("activity_registrations").
		Filter("activity_id", activityId).
		Filter("user_id", user.Id).
		Filter("status", 1).
		Exist()

	if exist {
		o.Rollback()
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "您已报名该活动",
		}
		c.ServeJSON()
		return
	}

	// 检查活动容量
	count, err := o.QueryTable("activity_registrations").
		Filter("activity_id", activityId).
		Filter("status", 1).
		Count()

	if err != nil {
		o.Rollback()
		beego.Error("查询报名人数失败：", err)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "系统错误，请稍后重试",
		}
		c.ServeJSON()
		return
	}

	if int(count) >= activity.Capacity {
		o.Rollback()
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "活动名额已满",
		}
		c.ServeJSON()
		return
	}

	// 创建报名记录
	registration := Models.ActivityRegistrations{
		ActivityId: activityId,
		UserId:     user.Id,
		Status:     1,
	}

	_, err = o.Insert(&registration)
	if err != nil {
		o.Rollback()
		beego.Error("创建报名记录失败：", err)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "报名失败，请稍后重试",
		}
		c.ServeJSON()
		return
	}

	// 更新活动剩余名额
	activity.Capacity = activity.Capacity - 1
	_, err = o.Update(&activity, "Capacity")
	if err != nil {
		o.Rollback()
		beego.Error("更新活动容量失败：", err)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "报名失败，请稍后重试",
		}
		c.ServeJSON()
		return
	}

	// 提交事务
	err = o.Commit()
	if err != nil {
		o.Rollback()
		beego.Error("提交事务失败：", err)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "报名失败，请稍后重试",
		}
		c.ServeJSON()
		return
	}

	c.Data["json"] = map[string]interface{}{
		"success": true,
		"message": "报名成功",
	}
	c.ServeJSON()
}

// GetRegistrationStatus 获取用户报名状态
func (c *StudentController) GetRegistrationStatus() {
	beego.Info("开始获取报名状态")

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

	// 获取活动ID
	activityId, err := c.GetInt64("activityId")
	if err != nil {
		beego.Error("获取活动ID失败：", err)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "无效的活动ID",
		}
		c.ServeJSON()
		return
	}

	o := orm.NewOrm()

	// 查询用户信息
	var user Models.Users
	err = o.QueryTable("users").Filter("username", userID).One(&user)
	if err != nil {
		beego.Error("查询用户信息失败：", err)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "用户信息获取失败",
		}
		c.ServeJSON()
		return
	}

	// 查询报名状态
	var registration Models.ActivityRegistrations
	err = o.QueryTable("activity_registrations").
		Filter("activity_id", activityId).
		Filter("user_id", user.Id).
		One(&registration)

	if err == orm.ErrNoRows {
		c.Data["json"] = map[string]interface{}{
			"success":      true,
			"isRegistered": false,
		}
	} else if err != nil {
		beego.Error("查询报名状态失败：", err)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "获取报名状态失败",
		}
	} else {
		c.Data["json"] = map[string]interface{}{
			"success":      true,
			"isRegistered": registration.Status == 1,
		}
	}
	c.ServeJSON()
}

// ShowMyActivities 显示我的活动页面
func (c *StudentController) ShowMyActivities() {
	beego.Info("进入我的活动页面")
	c.Data["ActivePage"] = "myactivities"
	c.Data["IsClubAdmin"] = c.isClubAdmin()
	c.TplName = "student_my_activities.html"
}

// GetMyActivities 获取我参加的活动列表
func (c *StudentController) GetMyActivities() {
	beego.Info("开始获取我的活动列表")

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

	// 查询用户信息
	var user Models.Users
	err := o.QueryTable("users").Filter("username", userID).One(&user)
	if err != nil {
		beego.Error("查询用户信息失败：", err)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "用户信息获取失败",
		}
		c.ServeJSON()
		return
	}

	// 使用 SQL 关联查询获取活动信息
	var activities []struct {
		RegistrationId     int64     `json:"registration_id"`
		ActivityId         int64     `json:"activity_id"`
		ActivityName       string    `json:"activity_name"`
		Description        string    `json:"description"`
		Location           string    `json:"location"`
		StartTime          time.Time `json:"start_time"`
		EndTime            time.Time `json:"end_time"`
		RegisteredAt       time.Time `json:"registered_at"`
		RegistrationStatus int       `json:"registration_status"`
		ActivityStatus     int       `json:"activity_status"`
		Score              int       `json:"score"`
	}

	sql := `
	SELECT 
    ar.id as registration_id,
    a.id as activity_id,
    a.name as activity_name,
    a.description,
    a.location,
    a.start_time,
    a.end_time,
    ar.created_at as registered_at,
    ar.status as registration_status,
    a.status as activity_status,
    a.points as score  -- 修改这里，将 a.score 改为 a.points
FROM activity_registrations ar
JOIN activities a ON ar.activity_id = a.id
WHERE ar.user_id = ? AND ar.status = 1
ORDER BY ar.created_at DESC
	`

	_, err = o.Raw(sql, user.Id).QueryRows(&activities)
	if err != nil {
		beego.Error("查询活动列表失败：", err)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "获取活动列表失败",
		}
		c.ServeJSON()
		return
	}

	c.Data["json"] = map[string]interface{}{
		"success":    true,
		"activities": activities,
	}
	c.ServeJSON()
}

// CancelRegistration 取消活动报名
func (c *StudentController) CancelRegistration() {
	beego.Info("开始处理取消报名请求")

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

	// 获取报名ID
	registrationId, err := c.GetInt64("registrationId")
	if err != nil {
		beego.Error("获取报名ID失败：", err)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "无效的报名ID",
		}
		c.ServeJSON()
		return
	}

	o := orm.NewOrm()

	// 查询用户信息
	var user Models.Users
	err = o.QueryTable("users").Filter("username", userID).One(&user)
	if err != nil {
		beego.Error("查询用户信息失败：", err)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "用户信息获取失败",
		}
		c.ServeJSON()
		return
	}

	// 查询报名记录
	var registration Models.ActivityRegistrations
	err = o.QueryTable("activity_registrations").Filter("id", registrationId).Filter("user_id", user.Id).One(&registration)
	if err != nil {
		beego.Error("查询报名记录失败：", err)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "报名记录不存在",
		}
		c.ServeJSON()
		return
	}

	// 查询活动信息
	var activity Models.Activities
	err = o.QueryTable("activities").Filter("id", registration.ActivityId).One(&activity)
	if err != nil {
		beego.Error("查询活动信息失败：", err)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "活动信息获取失败",
		}
		c.ServeJSON()
		return
	}

	// 检查活动开始时间是否在24小时内
	isWithin24Hours := time.Now().Add(24 * time.Hour).After(activity.StartTime)
	if isWithin24Hours {
		// 使用UpdatePoints函数扣除积分，这会自动更新区块链
		c.UpdatePoints(-5)

		// 扣除5积分，4月25日lr重新添加扣除积分功能
		user.Points -= 5
		_, err = o.Update(&user, "Points")
		if err != nil {
			beego.Error("更新用户积分失败：", err)
			c.Data["json"] = map[string]interface{}{
				"success": false,
				"message": "积分扣除失败",
			}
			c.ServeJSON()
			return
		}
	}

	// 更新报名状态为已取消(2)
	registration.Status = 2
	_, err = o.Update(&registration, "Status")
	if err != nil {
		beego.Error("更新报名状态失败：", err)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "取消报名失败，请稍后重试",
		}
		c.ServeJSON()
		return
	}

	// 创建积分记录
	if isWithin24Hours {
		pointsRecord := Models.PointsRecord{
			UserId:      user.Id,
			ActivityId:  registration.ActivityId,
			Points:      -5,
			Description: "恶意取消报名扣除积分",
			CreatedAt:   time.Now(),
		}
		_, err = o.Insert(&pointsRecord)
		if err != nil {
			beego.Error("创建积分记录失败：", err)
		}
	}

	c.Data["json"] = map[string]interface{}{
		"success":      true,
		"message":      "取消报名成功",
		"deductPoints": isWithin24Hours,
	}
	c.ServeJSON()
}

// ShowActivityRecords 显示扣除积分管理页面
func (c *StudentController) ShowActivityRecords() {
	beego.Info("进入扣除积分管理页面")

	// 检查用户是否是社团管理员
	userID := c.GetSession("userId")
	if userID == nil {
		c.Redirect("/", 302)
		return
	}

	o := orm.NewOrm()
	var user Models.Users
	err := o.QueryTable("users").Filter("username", userID).One(&user)
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "获取用户信息失败",
		}
		c.ServeJSON()
		return
	}

	if user.Role_name != "社团管理员" {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "您没有权限访问此页面",
		}
		c.ServeJSON()
		return
	}

	c.Data["ActivePage"] = "activity_records"
	c.Data["IsClubAdmin"] = true
	c.TplName = "student_activity_records.html"
}

// GetActivityRecords 获取扣除积分列表
func (c *StudentController) GetActivityRecords() {
	beego.Info("开始获取扣除积分列表")

	// 检查用户是否是社团管理员
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
	err := o.QueryTable("users").Filter("username", userID).One(&user)
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "获取用户信息失败",
		}
		c.ServeJSON()
		return
	}

	if user.Role_name != "社团管理员" {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "您没有权限访问此页面",
		}
		c.ServeJSON()
		return
	}

	page, _ := c.GetInt("page", 1)
	pageSize := 6
	offset := (page - 1) * pageSize

	// 定义结构体来接收查询结果
	type RecordResult struct {
		Id              int64     `json:"id"`
		ActivityId      int64     `json:"activity_id"`
		ActivityName    string    `json:"activity_name"`
		AttendanceCount int       `json:"attendance_count"`
		CreatedBy       int64     `json:"created_by"`
		CreatedAt       time.Time `json:"created_at"`
		StartTime       time.Time `json:"start_time"`
		EndTime         time.Time `json:"end_time"`
		Location        string    `json:"location"`
		Points          int       `json:"points"`
		AbsentStudents  string    `json:"absent_students"` // ✅ 加这一行
	}

	var records []RecordResult
	var total int64

	// 获取总记录数
	total, _ = o.QueryTable("activity_records").Count()

	// 构建查询
	sql := `
	SELECT 
		ar.id,
		ar.activity_id,
		a.name as activity_name,
		ar.attendance_count,
		ar.created_by,
		ar.created_at,
		a.start_time,
		a.end_time,
		a.location,
		a.points,
		ar.absent_students
	FROM activity_records ar
	LEFT JOIN activities a ON ar.activity_id = a.id
	ORDER BY ar.created_at DESC
	LIMIT ? OFFSET ?
`

	_, err = o.Raw(sql, pageSize, offset).QueryRows(&records)
	if err != nil {
		beego.Error("查询扣除积分记录失败：", err)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "获取扣除积分记录失败",
			"error":   err.Error(),
		}
		c.ServeJSON()
		return
	}

	// 如果没有记录，返回空数组而不是 null
	if records == nil {
		records = make([]RecordResult, 0)
	}

	c.Data["json"] = map[string]interface{}{
		"success": true,
		"records": records,
		"total":   total,
	}
	c.ServeJSON()
}

// AddActivityRecord 添加扣除积分列表
func (c *StudentController) AddActivityRecord() {
	beego.Info("开始添加扣除积分记录")

	// 检查用户是否是社团管理员
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
	err := o.QueryTable("users").Filter("username", userID).One(&user)
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "获取用户信息失败",
		}
		c.ServeJSON()
		return
	}

	if user.Role_name != "社团管理员" {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "您没有权限执行此操作",
		}
		c.ServeJSON()
		return
	}

	activityId, _ := c.GetInt64("activityId")
	summary := c.GetString("summary")
	absentStr := c.GetString("absentStudents")

	// 自动计算实际到场人数（根据报名表中 status 为 1 的记录）
	attendanceCount64, err := o.QueryTable("activity_registrations").
		Filter("activity_id", activityId).
		Filter("status", 1). // 根据业务定义实际到场状态
		Count()

	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "统计实际到场人数失败",
		}
		c.ServeJSON()
		return
	}

	attendanceCount := int(attendanceCount64) // 转换为 int 类型

	// 检查活动是否存在且已结束
	var activity Models.Activities
	err = o.QueryTable("activities").Filter("id", activityId).One(&activity)
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "活动不存在",
		}
		c.ServeJSON()
		return
	}

	if activity.Status != 3 { // 3表示活动已结束
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "只能为已结束的活动添加记录",
		}
		c.ServeJSON()
		return
	}

	// 检查是否已存在该活动的记录
	exist := o.QueryTable("activity_records").Filter("activity_id", activityId).Exist()
	if exist {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "该活动已有记录",
		}
		c.ServeJSON()
		return
	}

	// 👇 处理未到场学生自动扣分
	// 👇 处理未到场学生自动扣分（根据学号处理）
	deductedNames := []string{}
	if absentStr != "" {
		trimmed := strings.TrimSpace(absentStr)

		// 判断是否是多个学号（包含英文逗号）
		if strings.Contains(trimmed, ",") {
			ids := strings.Split(trimmed, ",")
			for _, id := range ids {
				id = strings.TrimSpace(id)

				// 校验学号是否为9位数字
				if len(id) != 9 || !regexp.MustCompile(`^\d{9}$`).MatchString(id) {
					c.Data["json"] = map[string]interface{}{
						"success": false,
						"message": "学生学号有误，请确保学号为9位数字并用英文逗号分隔",
					}
					c.ServeJSON()
					return
				}

				var student Models.Users
				err := o.QueryTable("users").Filter("username", id).One(&student)
				if err != nil {
					c.Data["json"] = map[string]interface{}{
						"success": false,
						"message": fmt.Sprintf("学生 %s 不存在，请验证后重新添加", id),
					}
					c.ServeJSON()
					return
				}

				// 扣除当前活动的积分
				student.Points -= activity.Points
				_, err = o.Update(&student, "Points")
				if err != nil {
					beego.Warn("扣分失败：", id)
					continue
				}

				deductedNames = append(deductedNames, id)
			}

		} else {
			// 处理单个学号
			id := trimmed
			if len(id) != 9 || !regexp.MustCompile(`^\d{9}$`).MatchString(id) {
				c.Data["json"] = map[string]interface{}{
					"success": false,
					"message": "学生学号格式有误，请确保为9位数字",
				}
				c.ServeJSON()
				return
			}

			var student Models.Users
			err := o.QueryTable("users").Filter("username", id).One(&student)
			if err != nil {
				c.Data["json"] = map[string]interface{}{
					"success": false,
					"message": "当前学生不存在，请验证后添加记录",
				}
				c.ServeJSON()
				return
			}

			student.Points -= activity.Points
			_, err = o.Update(&student, "Points")
			if err != nil {
				beego.Warn("扣分失败：", id)
			} else {
				deductedNames = append(deductedNames, id)
			}
		}
	}

	// ⬇️ 把扣分信息写入 summary 备注
	if len(deductedNames) > 0 {
		summary += "\n未到场学生扣分：" + strings.Join(deductedNames, ", ")
	}

	//插入记录
	record := Models.ActivityRecords{
		ActivityId:      activityId,
		AttendanceCount: attendanceCount,
		AbsentStudents:  absentStr,
		CreatedBy:       user.Id,
		CreatedAt:       time.Now(),
	}

	_, err = o.Insert(&record)
	if err != nil {
		beego.Error("添加扣除积分记录失败：", err)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "添加扣除积分记录失败",
		}
		c.ServeJSON()
		return
	}

	c.Data["json"] = map[string]interface{}{
		"success": true,
		"message": "添加扣除积分记录成功",
	}

	pointsRecord := Models.PointsRecord{
		UserId:      user.Id, // 或者上传该活动的用户ID
		ActivityId:  activity.Id,
		Points:      activity.Points, // 通常是你表单中的积分值
		Description: "扣除积分记录",        // 你也可以根据需求自定义描述
		CreatedAt:   time.Now(),
	}
	_, err = o.Insert(&pointsRecord)
	if err != nil {
		beego.Error("创建扣除积分记录失败：", err)
	}

	c.ServeJSON()
}

// GetActivityById 获取指定活动的积分（用于前端自动填充积分信息）
func (c *StudentController) GetActivityById() {
	activityId, err := c.GetInt64("id")
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "参数错误",
		}
		c.ServeJSON()
		return
	}

	o := orm.NewOrm()
	var activity Models.Activities
	err = o.QueryTable("activities").Filter("id", activityId).One(&activity)
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "活动不存在",
		}
		c.ServeJSON()
		return
	}

	c.Data["json"] = map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"points": activity.Points,
		},
	}
	c.ServeJSON()
}

// GetAttendanceCount 获取活动的实际到场人数
func (c *StudentController) GetAttendanceCount() {
	activityId, err := c.GetInt64("id")
	if err != nil || activityId <= 0 {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "活动ID无效",
		}
		c.ServeJSON()
		return
	}

	o := orm.NewOrm()
	count, err := o.QueryTable("activity_registrations").
		Filter("activity_id", activityId).
		Filter("status", 1). // 1 表示实际到场
		Count()

	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "查询到场人数失败",
		}
		c.ServeJSON()
		return
	}

	c.Data["json"] = map[string]interface{}{
		"success": true,
		"count":   count,
	}
	c.ServeJSON()
}

// 在修改积分的函数中添加更新头衔的逻辑
func (c *StudentController) UpdatePoints(points int) {
	userID := c.GetSession("userId")
	if userID == nil {
		return
	}

	o := orm.NewOrm()
	var user Models.Users
	err := o.QueryTable("users").Filter("username", userID).One(&user)
	if err != nil {
		return
	}

	// 更新积分
	user.Points += points
	user.UpdateUserTitle() // 更新用户头衔

	// 读取私钥文件
	files, err := os.ReadDir(`D:/geth student 1009/student/dev-chain/keystore`)
	if err != nil {
		beego.Error("读取keystore目录失败:", err)
		return
	}

	keystoreFile := `D:/geth student 1009/student/dev-chain/keystore/` + files[0].Name()
	password := "12345678" // 钱包密码

	// 打开keystore文件
	file, err := os.Open(keystoreFile)
	if err != nil {
		beego.Error("打开keystore文件失败:", err)
		return
	}
	defer file.Close()

	// 解密keystore文件
	keyjson, err := os.ReadFile(keystoreFile)
	if err != nil {
		beego.Error("读取keystore文件失败:", err)
		return
	}

	// 解密为私钥
	key, err := keystore.DecryptKey(keyjson, password)
	if err != nil {
		beego.Error("解密keystore失败:", err)
		return
	}

	// 获取私钥
	privateKey := key.PrivateKey

	// 先更新数据库
	_, err = o.Update(&user, "Points", "Title")
	if err != nil {
		beego.Error("更新用户积分失败：", err)
		return
	}

	// 重新上链用户信息
	err = blockchain.UploadStudentPoints(privateKey, user.Username)
	if err != nil {
		beego.Error("用户信息重新上链失败:", err)
	} else {
		beego.Info("用户信息成功重新上链")
	}
}

// ShowRanking 显示积分排名页面
func (c *StudentController) ShowRanking() {
	beego.Info("进入积分排名页面")

	// 检查用户是否登录
	userID := c.GetSession("userId")
	if userID == nil {
		c.Redirect("/", 302)
		return
	}

	c.Data["ActivePage"] = "ranking"
	c.Data["IsClubAdmin"] = c.isClubAdmin()

	o := orm.NewOrm()

	// --- 新增：每次进入排名页面时，实时统计并更新total、activity_count、points ---
	// 1. 统计本月所有积分记录，重算total和activity_count
	// 先将所有用户的total和activity_count清零
	o.Raw("UPDATE users SET total = 0, activity_count = 0").Exec()

	// 获取本月开始时间
	now := time.Now()
	firstDayOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	// 查询本月所有积分记录
	var records []Models.PointsRecord
	_, err := o.QueryTable("points_record").
		Filter("created_at__gte", firstDayOfMonth).
		All(&records)
	if err == nil {
		userTotal := make(map[int64]int)
		userCount := make(map[int64]int)
		for _, record := range records {
			userTotal[record.UserId] += record.Points
			userCount[record.UserId] += 1
		}
		for userId, total := range userTotal {
			count := userCount[userId]
			o.QueryTable("users").Filter("id", userId).Update(orm.Params{
				"total":          total,
				"activity_count": count,
			})
		}
	}
	// 2. 实时刷新points字段（积分商城是直接更新users表的Points字段）
	// 这里不需要额外处理，因商城和活动等操作已实时更新Points

	// --- 继续原有排名逻辑 ---
	var users []Models.Users
	_, err = o.QueryTable("users").OrderBy("-points").All(&users)
	if err != nil {
		beego.Error("获取排名数据失败：", err)
		c.Data["Error"] = "获取排名数据失败"
		c.TplName = "student_ranking.html"
		return
	}

	// 处理排名逻辑，实现相同积分相同排名
	if len(users) > 0 {
		rank := 1
		prevPoints := users[0].Points
		for i, user := range users {
			if user.Points < prevPoints {
				rank = i + 1
				prevPoints = user.Points
			}
			users[i].Rank = rank
		}
	}

	// 获取当前用户信息用于导航栏显示
	var user Models.Users
	err = o.QueryTable("users").Filter("username", userID).One(&user)
	if err == nil {
		c.Data["UserPoints"] = user.Points
		c.Data["UserTitle"] = user.Title
	}

	c.Data["Users"] = users
	c.TplName = "student_ranking.html"
}

// GetPointsRecords 获取用户的积分记录
func (c *StudentController) GetPointsRecords() {
	beego.Info("开始获取积分记录")

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
	err := o.QueryTable("users").Filter("username", userID).One(&user)
	if err != nil {
		beego.Error("查询用户信息失败：", err)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "用户信息获取失败",
		}
		c.ServeJSON()
		return
	}

	// 查询用户的积分记录
	var records []Models.PointsRecord
	_, err = o.QueryTable("points_record").
		Filter("user_id", user.Id).
		OrderBy("-created_at").
		All(&records)
	if err != nil {
		beego.Error("查询积分记录失败：", err)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "获取积分记录失败",
		}
		c.ServeJSON()
		return
	}

	// 获取活动名称
	var result []map[string]interface{}
	for _, record := range records {
		var activityName string

		// 根据积分记录的来源类型判断
		switch record.Source {
		case "activity":
			var activity Models.Activities
			err = o.QueryTable("activities").Filter("id", record.ActivityId).One(&activity)
			if err != nil {
				beego.Error("查询活动信息失败：", err)
				activityName = "未知活动"
			} else {
				activityName = activity.Name
			}
		case "exchange":
			activityName = "兑换商品"
		case "admin":
			activityName = "管理员操作"
		default:
			activityName = "其他来源"
		}

		result = append(result, map[string]interface{}{
			"id":            record.Id,
			"activity_name": activityName,
			"points":        record.Points,
			"description":   record.Description,
			"created_at":    record.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	c.Data["json"] = map[string]interface{}{
		"success": true,
		"records": result,
	}
	c.ServeJSON()
}

// ShowPointsRecords 显示积分记录页面
func (c *StudentController) ShowPointsRecords() {
	beego.Info("进入积分记录页面")
	c.Data["ActivePage"] = "points_records"
	c.Data["IsClubAdmin"] = c.isClubAdmin()
	c.TplName = "student_points_records.html"
}

// 兑换商品接口
func (c *StudentController) ExchangeItem() {
	userID := c.GetSession("userId")
	if userID == nil {
		c.Data["json"] = map[string]interface{}{"success": false, "message": "请先登录"}
		c.ServeJSON()
		return
	}

	item := strings.TrimSpace(c.GetString("item"))
	beego.Info("收到兑换商品类型：", item)

	// 商品类型与所需积分映射
	itemCostMap := map[string]int{
		"clothe":      50,
		"book":        30,
		"coupon":      20,
		"badge":       25,
		"canvas_bag":  40,
		"cup":         35,
		"fan":         28,
		"kettle":      60,
		"keychain":    15,
		"mouse_pad":   18,
		"notebook":    22,
		"phone_stand": 20,
		"speaker":     80,
		"tshirt":      55,
	}

	cost, ok := itemCostMap[item]
	if !ok {
		c.Data["json"] = map[string]interface{}{"success": false, "message": "无效的商品类型"}
		c.ServeJSON()
		return
	}

	o := orm.NewOrm()
	var user Models.Users
	err := o.QueryTable("users").Filter("username", userID).One(&user)
	if err != nil {
		c.Data["json"] = map[string]interface{}{"success": false, "message": "用户不存在"}
		c.ServeJSON()
		return
	}

	if user.Points < cost {
		c.Data["json"] = map[string]interface{}{"success": false, "message": "积分不足，无法兑换"}
		c.ServeJSON()
		return
	}

	// 扣除积分
	user.Points -= cost

	// 追加兑换记录
	if user.Exchange == "" {
		user.Exchange = item
	} else {
		user.Exchange = user.Exchange + "," + item
	}

	_, err = o.Update(&user, "Points", "Exchange")
	if err != nil {
		c.Data["json"] = map[string]interface{}{"success": false, "message": "兑换失败，请稍后重试"}
		c.ServeJSON()
		return
	}

	c.Data["json"] = map[string]interface{}{"success": true, "message": "兑换成功", "points": user.Points}

	// 添加兑换商品的积分记录
	pointsRecord := Models.PointsRecord{
		UserId:      user.Id,
		Points:      -cost,
		Description: fmt.Sprintf("兑换商品：%s", item),
		Source:      "exchange",
		CreatedAt:   time.Now(),
	}

	_, err = o.Insert(&pointsRecord)
	if err != nil {
		beego.Error("创建兑换商品积分记录失败：", err)
	}

	c.ServeJSON()
}

// 展示兑换历史页面
func (c *StudentController) ShowExchangeHistory() {
	beego.Info("进入兑换历史页面")
	c.Data["ActivePage"] = "exchange_history"
	c.Data["IsClubAdmin"] = c.isClubAdmin()
	userID := c.GetSession("userId")
	if userID == nil {
		c.Redirect("/", 302)
		return
	}

	o := orm.NewOrm()
	var user Models.Users
	err := o.QueryTable("users").Filter("username", userID).One(&user)
	if err != nil {
		c.Data["ExchangeMap"] = map[string]int{}
		c.TplName = "student_exchange_history.html"
		return
	}

	exchangeMap := map[string]int{}
	if user.Exchange != "" {
		items := strings.Split(user.Exchange, ",")
		for _, item := range items {
			exchangeMap[item]++
		}
	}
	c.Data["ExchangeMap"] = exchangeMap
	c.TplName = "student_exchange_history.html"
}
