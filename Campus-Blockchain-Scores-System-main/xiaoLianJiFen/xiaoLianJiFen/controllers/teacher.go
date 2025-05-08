package controllers

import (
	"encoding/json"
	Models "xiaoLianJiFen/models"

	"fmt"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type TeacherController struct {
	beego.Controller
}

// ShowDashboard 显示教师端首页
func (c *TeacherController) ShowDashboard() {
	beego.Info("进入教师端首页")
	c.Data["ActivePage"] = "home"
	c.TplName = "teacher_nav.html"
}

// ShowActivity 显示活动管理页面
func (c *TeacherController) ShowActivity() {
	beego.Info("进入活动管理页面")
	c.Data["ActivePage"] = "activity"
	c.TplName = "teacher_activity.html"
}

// ShowStudent 显示学生管理页面
func (c *TeacherController) ShowStudent() {
	beego.Info("进入学生管理页面")
	c.Data["ActivePage"] = "student"
	c.TplName = "teacher_student.html"
}

// ShowClub 显示社团管理页面
func (c *TeacherController) ShowClub() {
	beego.Info("进入社团管理页面")
	c.Data["ActivePage"] = "club"
	c.TplName = "teacher_club.html"
}

// ShowProfile 显示个人中心页面
func (c *TeacherController) ShowProfile() {
	beego.Info("进入个人中心页面")
	c.Data["ActivePage"] = "profile"
	c.TplName = "teacher_profile.html"
}

// ShowTeacherNav 显示教师导航页面
func (c *TeacherController) ShowTeacherNav() {
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
	if err != nil || user.Role_name != "教师" {
		c.Redirect("/", 302)
		return
	}

	// 设置导航栏当前页面
	c.Data["ActivePage"] = "home"
	// 渲染教师导航页面
	c.TplName = "teacher_nav.html"
}

// GetAdminRequests 获取所有申请成为管理员的学生列表
func (c *TeacherController) GetAdminRequests() {
	beego.Info("开始获取申请成为管理员的学生列表")

	o := orm.NewOrm()
	var users []Models.Users

	// 查询所有已申请管理员的学生
	_, err := o.QueryTable("users").Filter("is_admin_request", true).All(&users)
	if err != nil {
		beego.Error("查询申请列表失败：", err)
		c.Data["json"] = map[string]interface{}{"success": false, "message": "获取申请列表失败"}
		c.ServeJSON()
		return
	}

	beego.Info("成功获取申请列表，数量：", len(users))
	// 返回申请列表
	c.Data["json"] = map[string]interface{}{"success": true, "data": users}
	c.ServeJSON()
}

// SetAsAdmin 将学生设置为管理员
func (c *TeacherController) SetAsAdmin() {
	beego.Info("开始处理设置管理员请求")

	// 通过表单参数获取学生 ID
	studentID, err := c.GetInt("id")
	if err != nil {
		beego.Error("获取学生 ID 失败：", err)
		c.Data["json"] = map[string]interface{}{"success": false, "message": "无效的学生 ID"}
		c.ServeJSON()
		return
	}

	beego.Info("获取的学生 ID：", studentID)

	o := orm.NewOrm()
	user := Models.Users{Id: int64(studentID)}

	// 查询学生信息
	err = o.Read(&user)
	if err != nil {
		beego.Error("查询学生信息失败：", err)
		c.Data["json"] = map[string]interface{}{"success": false, "message": "学生不存在"}
		c.ServeJSON()
		return
	}

	beego.Info("查询到的学生信息：", user)

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

// ShowStudents 显示学生列表页面
func (c *TeacherController) ShowStudents() {
	beego.Info("进入学生列表页面")
	c.Data["ActivePage"] = "students"
	c.TplName = "teacher_student.html"
}

// ShowClubs 显示社团列表页面
func (c *TeacherController) ShowClubs() {
	beego.Info("进入社团列表页面")
	c.Data["ActivePage"] = "clubs"
	c.TplName = "teacher_club.html"
}

// GetStudents 获取所有学生列表
func (c *TeacherController) GetStudents() {
	beego.Info("开始获取学生列表")

	o := orm.NewOrm()
	var users []Models.Users

	// 查询角色为学生和社团管理员的用户
	_, err := o.QueryTable("users").Filter("role_name__in", []string{"学生", "社团管理员"}).All(&users)
	if err != nil {
		beego.Error("查询学生列表失败：", err)
		c.Data["json"] = map[string]interface{}{"success": false, "message": "获取学生列表失败"}
		c.ServeJSON()
		return
	}

	beego.Info("成功获取学生列表，数量：", len(users))
	// 返回学生列表
	c.Data["json"] = map[string]interface{}{"success": true, "data": users}
	c.ServeJSON()
}

// UpdateAdminRequest 更新学生的管理员申请状态
func (c *TeacherController) UpdateAdminRequest() {
	beego.Info("开始更新学生管理员申请状态")

	// 获取表单数据
	id, err := c.GetInt64("id")
	if err != nil {
		beego.Error("获取学生ID失败：", err)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "无效的学生ID",
		}
		c.ServeJSON()
		return
	}

	status, err := c.GetInt8("status")
	if err != nil {
		beego.Error("获取状态值失败：", err)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "无效的状态值",
		}
		c.ServeJSON()
		return
	}

	// 直接更新数据库
	o := orm.NewOrm()
	num, err := o.QueryTable("users").Filter("id", id).Update(orm.Params{
		"is_admin_request": status,
	})

	if err != nil {
		beego.Error("更新失败：", err)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "更新失败",
		}
		c.ServeJSON()
		return
	}

	if num == 0 {
		beego.Error("未找到要更新的学生记录")
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "未找到学生记录",
		}
		c.ServeJSON()
		return
	}

	beego.Info("成功更新学生管理员申请状态，学生ID：", id, "，新状态：", status)
	c.Data["json"] = map[string]interface{}{
		"success": true,
		"message": "更新成功",
	}
	c.ServeJSON()
}

// GetActivities 获取活动列表
func (c *TeacherController) GetActivities() {
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
	for _, activity := range activities {
		beego.Info("活动信息：", activity.Name, "状态：", activity.Status)
	}
	c.Data["json"] = map[string]interface{}{
		"success": true,
		"data":    activities,
	}
	c.ServeJSON()
}

// ApproveActivity 审批通过活动
func (c *TeacherController) ApproveActivity() {
	beego.Info("开始处理活动审批")

	// 获取活动ID
	id, err := c.GetInt64("id")
	if err != nil {
		beego.Error("获取活动ID失败：", err)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "无效的活动ID",
		}
		c.ServeJSON()
		return
	}

	// 直接更新数据库
	o := orm.NewOrm()
	num, err := o.QueryTable("activities").Filter("id", id).Update(orm.Params{
		"status": 1, // 1表示已通过
	})

	if err != nil {
		beego.Error("更新失败：", err)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "更新失败",
		}
		c.ServeJSON()
		return
	}

	if num == 0 {
		beego.Error("未找到要更新的活动记录")
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "未找到活动记录",
		}
		c.ServeJSON()
		return
	}

	beego.Info("成功审批通过活动，活动ID：", id)
	c.Data["json"] = map[string]interface{}{
		"success": true,
		"message": "审批通过成功",
	}
	c.ServeJSON()
}

// UpdateActivity 更新活动信息
func (c *TeacherController) UpdateActivity() {
	beego.Info("开始更新活动信息")

	// 获取活动ID
	id, err := c.GetInt64("id")
	if err != nil {
		beego.Error("获取活动ID失败：", err)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "无效的活动ID",
		}
		c.ServeJSON()
		return
	}

	// 获取更新的活动信息
	name := c.GetString("name")
	startTime := c.GetString("start_time") // 获取 StartTime
	endTime := c.GetString("end_time")     // 获取 EndTime
	location := c.GetString("location")
	description := c.GetString("description")
	points, err := c.GetInt("points")
	if err != nil {
		beego.Error("获取积分失败：", err)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "无效的积分值",
		}
		c.ServeJSON()
		return
	}
	capacity, err := c.GetInt("capacity")
	if err != nil {
		beego.Error("获取人数失败：", err)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "无效的人数值",
		}
		c.ServeJSON()
		return
	}

	// 更新数据库
	o := orm.NewOrm()
	num, err := o.QueryTable("activities").Filter("id", id).Update(orm.Params{
		"name":        name,
		"start_time":  startTime, // 更新 StartTime
		"end_time":    endTime,   // 更新 EndTime
		"location":    location,
		"description": description,
		"points":      points,
		"capacity":    capacity,
	})

	if err != nil {
		beego.Error("更新活动失败：", err)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "更新失败",
		}
		c.ServeJSON()
		return
	}

	if num == 0 {
		beego.Error("未找到要更新的活动记录")
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "未找到活动记录",
		}
		c.ServeJSON()
		return
	}

	beego.Info("成功更新活动信息，活动ID：", id)
	c.Data["json"] = map[string]interface{}{
		"success": true,
		"message": "更新成功",
	}
	c.ServeJSON()
}

// RejectActivity 拒绝活动申请
func (c *TeacherController) RejectActivity() {
	beego.Info("开始处理活动拒绝")

	// 获取活动ID
	id, err := c.GetInt64("id")
	if err != nil {
		beego.Error("获取活动ID失败：", err)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "无效的活动ID",
		}
		c.ServeJSON()
		return
	}

	// 直接更新数据库
	o := orm.NewOrm()
	num, err := o.QueryTable("activities").Filter("id", id).Update(orm.Params{
		"status": 2, // 2表示已拒绝
	})

	if err != nil {
		beego.Error("更新失败：", err)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "更新失败",
		}
		c.ServeJSON()
		return
	}

	if num == 0 {
		beego.Error("未找到要更新的活动记录")
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "未找到活动记录",
		}
		c.ServeJSON()
		return
	}

	beego.Info("成功拒绝活动，活动ID：", id)
	c.Data["json"] = map[string]interface{}{
		"success": true,
		"message": "已拒绝该活动",
	}
	c.ServeJSON()
}

// HandleAdminRequest 处理学生的管理员申请
func (c *TeacherController) HandleAdminRequest() {
	beego.Info("开始处理管理员申请")

	// 从表单获取数据
	studentId := c.GetString("id")
	action := c.GetString("action")
	beego.Info("收到的请求数据 - 学生ID:", studentId, "操作:", action)

	if studentId == "" || action == "" {
		beego.Error("缺少必要参数 - 学生ID:", studentId, "操作:", action)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "缺少必要参数",
		}
		c.ServeJSON()
		return
	}

	o := orm.NewOrm()
	var student Models.Users
	err := o.QueryTable("users").Filter("id", studentId).One(&student)
	if err != nil {
		beego.Error("查询学生失败：", err)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "找不到该学生",
		}
		c.ServeJSON()
		return
	}

	beego.Info("找到学生：", student)

	// 根据操作类型处理
	switch action {
	case "approve":
		if student.IsAdminRequest != 1 {
			c.Data["json"] = map[string]interface{}{
				"success": false,
				"message": "该学生没有申请管理员权限",
			}
			c.ServeJSON()
			return
		}
		student.Role_name = "社团管理员"
		student.IsAdminRequest = 0
		beego.Info("准备将学生设置为社团管理员")
	case "reject":
		if student.IsAdminRequest != 1 {
			c.Data["json"] = map[string]interface{}{
				"success": false,
				"message": "该学生没有申请管理员权限",
			}
			c.ServeJSON()
			return
		}
		student.IsAdminRequest = 0
		beego.Info("准备拒绝学生的管理员申请")
	case "revoke":
		if student.Role_name != "社团管理员" {
			c.Data["json"] = map[string]interface{}{
				"success": false,
				"message": "该学生不是管理员，无需撤销",
			}
			c.ServeJSON()
			return
		}
		student.Role_name = "学生"
		beego.Info("准备撤销学生的管理员权限")
	default:
		beego.Error("无效的操作类型：", action)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "无效的操作类型：" + action,
		}
		c.ServeJSON()
		return
	}

	// 开启事务
	o.Begin()
	beego.Info("开始更新学生信息")

	// 根据不同的操作类型更新不同的字段
	var updateFields []string
	switch action {
	case "approve":
		updateFields = []string{"role_name", "is_admin_request"}
	case "reject":
		updateFields = []string{"is_admin_request"}
	case "revoke":
		updateFields = []string{"role_name"}
	}

	_, err = o.Update(&student, updateFields...)
	if err != nil {
		beego.Error("更新学生信息失败：", err)
		o.Rollback()
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "处理申请失败：" + err.Error(),
		}
		c.ServeJSON()
		return
	}
	o.Commit()
	beego.Info("成功更新学生信息")

	c.Data["json"] = map[string]interface{}{
		"success": true,
		"message": "处理成功",
	}
	c.ServeJSON()
}

// RevokeAdmin 撤销社团管理员权限
func (c *TeacherController) RevokeAdmin() {
	beego.Info("开始处理撤销管理员请求")

	// 获取学生 ID
	studentID, err := c.GetInt("id")
	if err != nil {
		beego.Error("获取学生 ID 失败：", err)
		c.Data["json"] = map[string]interface{}{"success": false, "message": "无效的学生 ID"}
		c.ServeJSON()
		return
	}

	beego.Info("获取的学生 ID：", studentID)

	o := orm.NewOrm()
	user := Models.Users{Id: int64(studentID)}

	// 查询学生信息
	err = o.Read(&user)
	if err != nil {
		beego.Error("查询学生信息失败：", err)
		c.Data["json"] = map[string]interface{}{"success": false, "message": "学生不存在"}
		c.ServeJSON()
		return
	}

	beego.Info("查询到的学生信息：", user)

	// 更新学生角色回到 "学生"
	user.Role_name = "学生"
	_, err = o.Update(&user, "Role_name")
	if err != nil {
		beego.Error("更新学生角色失败: ", err)
		c.Data["json"] = map[string]interface{}{"success": false, "message": "撤销管理员失败"}
		c.ServeJSON()
		return
	}

	// 返回成功结果
	beego.Info("成功撤销管理员权限, 学生 ID: ", studentID)
	c.Data["json"] = map[string]interface{}{"success": true, "message": "已撤销该学生的管理员权限"}
	c.ServeJSON()
}

// ShowActivityRecords 显示活动记录页面
func (c *TeacherController) ShowActivityRecords() {
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
	if err != nil || user.Role_name != "教师" {
		c.Redirect("/", 302)
		return
	}

	beego.Info("进入活动记录页面")
	c.Data["ActivePage"] = "activity_records"
	c.TplName = "teacher_activity_records.html"
}

// GetActivityRecords 获取活动记录列表
func (c *TeacherController) GetActivityRecords() {
	// 检查用户是否登录
	userID := c.GetSession("userId")
	if userID == nil {
		c.Data["json"] = map[string]interface{}{
			"code": 401,
			"msg":  "请先登录",
		}
		c.ServeJSON()
		return
	}

	// 获取分页参数
	page, _ := c.GetInt("page", 1)
	pageSize, _ := c.GetInt("pageSize", 6)
	searchQuery := c.GetString("searchQuery", "")
	recordId, _ := c.GetInt64("recordId", 0)

	// 构建查询
	o := orm.NewOrm()

	// 首先获取所有已结束的活动
	var activities []Models.Activities
	qs := o.QueryTable("activities").Filter("status", 3)

	// 添加搜索条件
	if searchQuery != "" {
		qs = qs.Filter("name__icontains", searchQuery)
	}

	// 如果指定了recordId，只获取该记录对应的活动
	if recordId > 0 {
		var record Models.ActivityRecords
		err := o.QueryTable("activity_records").Filter("id", recordId).One(&record)
		if err == nil {
			qs = qs.Filter("id", record.Activity.Id)
		}
	}

	// 获取总记录数
	total, err := qs.Count()
	if err != nil {
		beego.Error("获取活动总数失败：", err)
		c.Data["json"] = map[string]interface{}{
			"code": 500,
			"msg":  "获取活动总数失败",
		}
		c.ServeJSON()
		return
	}

	// 获取分页数据
	_, err = qs.OrderBy("-end_time").
		Limit(pageSize, (page-1)*pageSize).
		All(&activities)
	if err != nil {
		beego.Error("获取活动列表失败：", err)
		c.Data["json"] = map[string]interface{}{
			"code": 500,
			"msg":  "获取活动列表失败",
		}
		c.ServeJSON()
		return
	}

	// 构造返回数据
	var recordsData []map[string]interface{}
	for _, activity := range activities {
		// 查找该活动的记录
		var record Models.ActivityRecords
		err := o.QueryTable("activity_records").
			Filter("activity_id", activity.Id).
			One(&record)

		// 构造活动数据
		activityData := map[string]interface{}{
			"activity_id":     activity.Id,
			"activity_name":   activity.Name,
			"activity_points": activity.Points,
			"end_time":        activity.EndTime,
			"has_record":      err == nil, // 是否已有记录
		}

		// 如果已有记录，添加记录信息
		if err == nil {
			activityData["id"] = record.Id
			activityData["attendance_count"] = record.AttendanceCount
			activityData["summary"] = record.Summary
			activityData["created_at"] = record.CreatedAt
		} else {
			// 如果没有记录，设置默认值
			activityData["attendance_count"] = 0
			activityData["summary"] = ""
		}

		recordsData = append(recordsData, activityData)
	}

	c.Data["json"] = map[string]interface{}{
		"code": 0,
		"data": map[string]interface{}{
			"total":   total,
			"records": recordsData,
		},
	}
	c.ServeJSON()
}

// UpdateActivityRecord 更新活动记录
func (c *TeacherController) UpdateActivityRecord() {
	// 检查用户是否登录
	userID := c.GetSession("userId")
	if userID == nil {
		c.Data["json"] = map[string]interface{}{
			"code": 1,
			"msg":  "未登录",
		}
		c.ServeJSON()
		return
	}

	// 检查用户角色
	o := orm.NewOrm()
	var user Models.Users
	err := o.QueryTable("users").Filter("username", userID).One(&user)
	if err != nil || user.Role_name != "教师" {
		c.Data["json"] = map[string]interface{}{
			"code": 2,
			"msg":  "无权限",
		}
		c.ServeJSON()
		return
	}

	// 获取表单数据
	recordId, err := c.GetInt64("id")
	if err != nil {
		beego.Error("获取记录ID失败：", err)
		c.Data["json"] = map[string]interface{}{
			"code": 3,
			"msg":  "无效的记录ID",
		}
		c.ServeJSON()
		return
	}

	attendanceCount, err := c.GetInt("attendance_count")
	if err != nil {
		beego.Error("获取到场人数失败：", err)
		c.Data["json"] = map[string]interface{}{
			"code": 3,
			"msg":  "无效的到场人数",
		}
		c.ServeJSON()
		return
	}

	summary := c.GetString("summary")
	if summary == "" {
		c.Data["json"] = map[string]interface{}{
			"code": 3,
			"msg":  "活动总结不能为空",
		}
		c.ServeJSON()
		return
	}

	// 更新活动记录
	record := Models.ActivityRecords{Id: recordId}
	if err := o.Read(&record); err != nil {
		c.Data["json"] = map[string]interface{}{
			"code": 4,
			"msg":  "活动记录不存在",
		}
		c.ServeJSON()
		return
	}

	record.AttendanceCount = attendanceCount
	record.Summary = summary

	if _, err := o.Update(&record, "attendance_count", "summary"); err != nil {
		c.Data["json"] = map[string]interface{}{
			"code": 5,
			"msg":  "更新失败: " + err.Error(),
		}
		c.ServeJSON()
		return
	}

	c.Data["json"] = map[string]interface{}{
		"code": 0,
		"msg":  "更新成功",
	}
	c.ServeJSON()
}

// CheckActivityEndStatus 检查活动是否已结束并更新状态
func CheckActivityEndStatus() {
	o := orm.NewOrm()
	now := time.Now()

	// 查找所有状态为已通过（1）且结束时间已到或已过的活动
	var activities []Models.Activities
	_, err := o.QueryTable("activities").Filter("status", 1).Filter("end_time__lte", now).All(&activities)
	if err != nil {
		beego.Error("查询活动失败：", err)
		return
	}

	// 更新这些活动的状态为已结束（3）并为已报名学生发放积分
	for _, activity := range activities {
		// 开启事务
		err = o.Begin()
		if err != nil {
			beego.Error("开启事务失败：", err)
			continue
		}

		// 更新活动状态为已结束
		activity.Status = 3
		_, err = o.Update(&activity, "Status")
		if err != nil {
			o.Rollback()
			beego.Error("更新活动状态失败：", err)
			continue
		}

		// 查找所有已报名且未取消的学生
		var registrations []Models.ActivityRegistrations
		_, err = o.QueryTable("activity_registrations").
			Filter("activity_id", activity.Id).
			Filter("status", 1).
			All(&registrations)
		if err != nil {
			o.Rollback()
			beego.Error("查询报名记录失败：", err)
			continue
		}

		// 为每个已报名的学生发放积分
		for _, registration := range registrations {
			var user Models.Users
			err = o.QueryTable("users").Filter("id", registration.UserId).One(&user)
			if err != nil {
				beego.Error("查询用户信息失败：", err)
				continue
			}

			// 检查是否已经为该用户创建过该活动的积分记录
			exist := o.QueryTable("points_record").
				Filter("user_id", user.Id).
				Filter("activity_id", activity.Id).
				Exist()

			if exist {
				beego.Info("该用户已有该活动的积分记录，跳过：", user.Id, activity.Id)
				continue
			}

			// 更新用户积分
			user.Points += activity.Points
			_, err = o.Update(&user, "Points")
			if err != nil {
				beego.Error("更新用户积分失败：", err)
				continue
			}

			// 创建积分记录
			pointsRecord := Models.PointsRecord{
				UserId:      user.Id,
				ActivityId:  activity.Id,
				Points:      activity.Points,
				Description: fmt.Sprintf("参加活动《%s》获得积分", activity.Name),
			}
			_, err = o.Insert(&pointsRecord)
			if err != nil {
				beego.Error("创建积分记录失败：", err)
				continue
			}

			beego.Info("为用户发放积分成功，用户ID：", user.Id, "，活动ID：", activity.Id, "，积分：", activity.Points)
		}

		// 提交事务
		err = o.Commit()
		if err != nil {
			o.Rollback()
			beego.Error("提交事务失败：", err)
			continue
		}

		beego.Info("活动已结束，ID：", activity.Id, "，结束时间：", activity.EndTime)
	}
}

// init 初始化函数
func init() {
	// 启动定时任务，每分钟检查一次活动状态
	go func() {
		for {
			CheckActivityEndStatus()
			time.Sleep(1 * time.Minute)
		}
	}()
}

// CreateActivityRecord 创建新的活动记录
func (c *TeacherController) CreateActivityRecord() {
	// 检查用户是否登录
	userID := c.GetSession("userId")
	if userID == nil {
		c.Data["json"] = map[string]interface{}{
			"code": 401,
			"msg":  "请先登录",
		}
		c.ServeJSON()
		return
	}

	// 检查用户角色
	o := orm.NewOrm()
	var user Models.Users
	err := o.QueryTable("users").Filter("username", userID).One(&user)
	if err != nil || user.Role_name != "教师" {
		c.Data["json"] = map[string]interface{}{
			"code": 403,
			"msg":  "无权限执行此操作",
		}
		c.ServeJSON()
		return
	}

	// 获取表单数据
	activityId, err := c.GetInt64("activityId")
	if err != nil {
		beego.Error("获取活动ID失败：", err)
		c.Data["json"] = map[string]interface{}{
			"code": 400,
			"msg":  "无效的活动ID",
		}
		c.ServeJSON()
		return
	}

	// 获取到场人数
	attendanceCountStr := c.GetString("attendanceCount")
	if attendanceCountStr == "" {
		c.Data["json"] = map[string]interface{}{
			"code": 400,
			"msg":  "请输入到场人数",
		}
		c.ServeJSON()
		return
	}

	attendanceCount, err := strconv.Atoi(attendanceCountStr)
	if err != nil {
		beego.Error("到场人数转换失败：", err)
		c.Data["json"] = map[string]interface{}{
			"code": 400,
			"msg":  "到场人数必须是整数",
		}
		c.ServeJSON()
		return
	}

	if attendanceCount < 0 {
		c.Data["json"] = map[string]interface{}{
			"code": 400,
			"msg":  "到场人数不能为负数",
		}
		c.ServeJSON()
		return
	}

	summary := c.GetString("summary")
	if summary == "" {
		c.Data["json"] = map[string]interface{}{
			"code": 400,
			"msg":  "活动总结不能为空",
		}
		c.ServeJSON()
		return
	}

	// 检查活动是否存在且已结束
	var activity Models.Activities
	err = o.QueryTable("activities").Filter("id", activityId).One(&activity)
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"code": 404,
			"msg":  "活动不存在",
		}
		c.ServeJSON()
		return
	}

	// 检查活动状态是否为已结束（3）
	if activity.Status != 3 {
		c.Data["json"] = map[string]interface{}{
			"code": 400,
			"msg":  "只能为已结束的活动创建记录",
		}
		c.ServeJSON()
		return
	}

	// 检查是否已存在该活动的记录
	exist := o.QueryTable("activity_records").Filter("activity_id", activityId).Exist()
	if exist {
		c.Data["json"] = map[string]interface{}{
			"code": 400,
			"msg":  "该活动已有记录",
		}
		c.ServeJSON()
		return
	}

	// 创建活动记录
	record := Models.ActivityRecords{
		Activity:        &Models.Activities{Id: activityId},
		AttendanceCount: attendanceCount,
		Summary:         summary,
		Creator:         &Models.Users{Id: user.Id},
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// 保存记录
	_, err = o.Insert(&record)
	if err != nil {
		beego.Error("创建活动记录失败：", err)
		c.Data["json"] = map[string]interface{}{
			"code": 500,
			"msg":  "创建记录失败",
		}
		c.ServeJSON()
		return
	}

	// 构建返回数据
	result := map[string]interface{}{
		"id":               record.Id,
		"activity_id":      activityId,
		"activity_name":    activity.Name,
		"activity_points":  activity.Points,
		"attendance_count": attendanceCount,
		"summary":          summary,
		"created_at":       record.CreatedAt,
	}

	c.Data["json"] = map[string]interface{}{
		"code": 0,
		"msg":  "创建成功",
		"data": result,
	}
	c.ServeJSON()
}

// GetActivityParticipants 获取活动参加人列表
func (c *TeacherController) GetActivityParticipants() {
	// 检查用户是否登录
	user := c.GetSession("user")
	if user == nil {
		c.Data["json"] = map[string]interface{}{
			"code": 401,
			"msg":  "请先登录",
		}
		c.ServeJSON()
		return
	}

	// 获取活动ID
	activityId, err := c.GetInt("activity_id")
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"code": 400,
			"msg":  "无效的活动ID",
		}
		c.ServeJSON()
		return
	}

	// 获取活动信息
	var activity Models.Activities
	if err := orm.NewOrm().QueryTable("activities").Filter("id", activityId).One(&activity); err != nil {
		c.Data["json"] = map[string]interface{}{
			"code": 500,
			"msg":  "获取活动信息失败",
		}
		c.ServeJSON()
		return
	}

	// 获取所有报名记录
	var registrations []Models.ActivityRegistrations
	_, err = orm.NewOrm().QueryTable("activity_registrations").
		Filter("activity_id", activityId).
		Filter("status", 1). // 已报名状态
		All(&registrations)
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"code": 500,
			"msg":  "获取报名记录失败",
		}
		c.ServeJSON()
		return
	}

	// 构造参加人列表
	var participants []map[string]interface{}
	for _, reg := range registrations {
		var student Models.Users
		if err := orm.NewOrm().QueryTable("users").Filter("id", reg.UserId).One(&student); err != nil {
			continue
		}

		participants = append(participants, map[string]interface{}{
			"id":                student.Id,
			"username":          student.Username,
			"registration_time": reg.CreatedAt,
			"is_attended":       false, // 默认未签到
		})
	}

	c.Data["json"] = map[string]interface{}{
		"code": 200,
		"data": map[string]interface{}{
			"activity_name": activity.Name,
			"points":        activity.Points,
			"participants":  participants,
		},
	}
	c.ServeJSON()
}

// UpdateAttendanceStatus 更新活动参加人的签到状态
func (c *TeacherController) UpdateAttendanceStatus() {
	// 检查用户是否登录
	user := c.GetSession("user")
	if user == nil {
		c.Data["json"] = map[string]interface{}{
			"code": 1,
			"msg":  "请先登录",
		}
		c.ServeJSON()
		return
	}

	// 获取请求数据
	var requestData struct {
		ActivityId   int64 `json:"activityId"`
		Participants []struct {
			Id         int64 `json:"id"`
			IsAttended bool  `json:"is_attended"`
		} `json:"participants"`
	}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &requestData); err != nil {
		c.Data["json"] = map[string]interface{}{
			"code": 1,
			"msg":  "无效的请求数据",
		}
		c.ServeJSON()
		return
	}

	// 验证活动ID
	if requestData.ActivityId <= 0 {
		c.Data["json"] = map[string]interface{}{
			"code": 1,
			"msg":  "无效的活动ID",
		}
		c.ServeJSON()
		return
	}

	// 开启事务
	o := orm.NewOrm()
	err := o.Begin()
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"code": 1,
			"msg":  "开启事务失败",
		}
		c.ServeJSON()
		return
	}

	// 更新签到状态
	for _, participant := range requestData.Participants {
		// 更新报名记录状态
		_, err := o.QueryTable("activity_registrations").
			Filter("id", participant.Id).
			Update(orm.Params{
				"is_attended": participant.IsAttended,
			})
		if err != nil {
			o.Rollback()
			c.Data["json"] = map[string]interface{}{
				"code": 1,
				"msg":  "更新签到状态失败",
			}
			c.ServeJSON()
			return
		}

		// 如果未签到，扣除积分
		if !participant.IsAttended {
			// 获取活动信息
			var activity Models.Activities
			err := o.QueryTable("activities").
				Filter("id", requestData.ActivityId).
				One(&activity)
			if err != nil {
				o.Rollback()
				c.Data["json"] = map[string]interface{}{
					"code": 1,
					"msg":  "获取活动信息失败",
				}
				c.ServeJSON()
				return
			}

			// 获取用户信息
			var registration Models.ActivityRegistrations
			err = o.QueryTable("activity_registrations").
				Filter("id", participant.Id).
				One(&registration)
			if err != nil {
				o.Rollback()
				c.Data["json"] = map[string]interface{}{
					"code": 1,
					"msg":  "获取报名信息失败",
				}
				c.ServeJSON()
				return
			}

			// 扣除积分
			_, err = o.QueryTable("users").
				Filter("id", registration.UserId).
				Update(orm.Params{
					"points": orm.ColValue(orm.ColMinus, activity.Points),
				})
			if err != nil {
				o.Rollback()
				c.Data["json"] = map[string]interface{}{
					"code": 1,
					"msg":  "扣除积分失败",
				}
				c.ServeJSON()
				return
			}

			// 添加积分记录
			pointsRecord := Models.PointsRecord{
				UserId:      registration.UserId,
				ActivityId:  requestData.ActivityId,
				Points:      -activity.Points,
				Description: fmt.Sprintf("未参加活动：%s", activity.Name),
				CreatedAt:   time.Now(),
			}
			_, err = o.Insert(&pointsRecord)
			if err != nil {
				o.Rollback()
				c.Data["json"] = map[string]interface{}{
					"code": 1,
					"msg":  "添加积分记录失败",
				}
				c.ServeJSON()
				return
			}
		}
	}

	// 更新活动记录的签到人数
	attendedCount := 0
	for _, participant := range requestData.Participants {
		if participant.IsAttended {
			attendedCount++
		}
	}

	_, err = o.QueryTable("activity_records").
		Filter("activity_id", requestData.ActivityId).
		Update(orm.Params{
			"attendance_count": attendedCount,
		})
	if err != nil {
		o.Rollback()
		c.Data["json"] = map[string]interface{}{
			"code": 1,
			"msg":  "更新签到人数失败",
		}
		c.ServeJSON()
		return
	}

	// 提交事务
	err = o.Commit()
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"code": 1,
			"msg":  "提交事务失败",
		}
		c.ServeJSON()
		return
	}

	c.Data["json"] = map[string]interface{}{
		"code": 0,
		"msg":  "更新成功",
	}
	c.ServeJSON()
}
