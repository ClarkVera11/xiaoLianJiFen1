package controllers

import (
	Models "xiaoLianJiFen/models"

	"fmt"
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
		beego.Error("未找到要更新的扣除积分记录")
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "未找到扣除积分记录",
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
		beego.Error("未找到要更新的扣除积分记录")
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "未找到扣除积分记录",
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
		beego.Error("未找到要更新的扣除积分记录")
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"message": "未找到扣除积分记录",
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

// ShowActivityRecords 显示扣除积分页面
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

	beego.Info("进入扣除积分页面")
	c.Data["ActivePage"] = "activity_records"
	c.TplName = "teacher_activity_records.html"
}

// GetActivityRecords 获取扣除积分列表
func (c *TeacherController) GetActivityRecords() {
	page, _ := c.GetInt("page", 1)
	pageSize, _ := c.GetInt("pageSize", 6)
	search := c.GetString("search")

	// 构建查询条件
	query := orm.NewOrm().QueryTable("activity_records")
	if search != "" {
		var activityIds orm.ParamsList
		orm.NewOrm().QueryTable("activities").Filter("name__icontains", search).ValuesFlat(&activityIds, "id")
		if len(activityIds) > 0 {
			query = query.Filter("activity_id__in", activityIds)
		}
	}

	// 获取总记录数
	total, _ := query.Count()

	// 获取分页数据
	var records []*Models.ActivityRecords
	query = query.OrderBy("-created_at").
		Limit(pageSize, (page-1)*pageSize)
	query.All(&records)

	// 获取活动名称和创建者名称
	var result []map[string]interface{}
	for _, record := range records {
		// 获取活动信息
		var activity Models.Activities
		orm.NewOrm().QueryTable("activities").Filter("id", record.ActivityId).One(&activity)

		// 获取创建者信息
		var user Models.Users
		orm.NewOrm().QueryTable("users").Filter("id", record.CreatedBy).One(&user)

		result = append(result, map[string]interface{}{
			"id":               record.Id,
			"activity_name":    activity.Name,
			"attendance_count": record.AttendanceCount,
			"created_by_name":  user.Username,
			"created_at":       record.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	c.Data["json"] = map[string]interface{}{
		"code": 0,
		"msg":  "success",
		"data": map[string]interface{}{
			"records":  result,
			"total":    total,
			"page":     page,
			"pageSize": pageSize,
		},
	}
	c.ServeJSON()
}

// UpdateActivityRecord 更新扣除积分
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

	// 更新扣除积分
	record := Models.ActivityRecords{Id: recordId}
	if err := o.Read(&record); err != nil {
		c.Data["json"] = map[string]interface{}{
			"code": 4,
			"msg":  "扣除积分记录不存在",
		}
		c.ServeJSON()
		return
	}

	record.AttendanceCount = attendanceCount


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
