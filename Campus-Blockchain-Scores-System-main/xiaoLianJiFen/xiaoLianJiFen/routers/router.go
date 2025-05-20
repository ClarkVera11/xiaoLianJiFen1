package routers

import (
	"xiaoLianJiFen/controllers"

	"github.com/astaxie/beego"
)

func init() {
	//登入
	beego.Router("/", &controllers.MainController{}, "get:ShowLogin;post:HandleLogin")

	// 教师端路由
	beego.Router("/teacher/nav", &controllers.TeacherController{}, "get:ShowTeacherNav")
	beego.Router("/teacher", &controllers.TeacherController{}, "get:ShowDashboard")
	beego.Router("/teacher/activity", &controllers.TeacherController{}, "get:ShowActivity")
	beego.Router("/teacher/student", &controllers.TeacherController{}, "get:ShowStudent")
	beego.Router("/teacher/club", &controllers.TeacherController{}, "get:ShowClub")
	beego.Router("/teacher/requests", &controllers.TeacherController{}, "get:GetAdminRequests")
	beego.Router("/teacher/set-admin", &controllers.TeacherController{}, "post:SetAsAdmin")
	beego.Router("/teacher/students", &controllers.TeacherController{}, "get:ShowStudents")
	beego.Router("/teacher/clubs", &controllers.TeacherController{}, "get:ShowClubs")
	beego.Router("/teacher/get-students", &controllers.TeacherController{}, "get:GetStudents")
	beego.Router("/teacher/update-admin-request", &controllers.TeacherController{}, "post:UpdateAdminRequest")
	beego.Router("/teacher/get-activities", &controllers.TeacherController{}, "get:GetActivities")
	beego.Router("/teacher/approve-activity", &controllers.TeacherController{}, "post:ApproveActivity")
	beego.Router("/teacher/update-activity", &controllers.TeacherController{}, "post:UpdateActivity")
	beego.Router("/teacher/reject-activity", &controllers.TeacherController{}, "post:RejectActivity")
	beego.Router("/teacher/handle-admin-request", &controllers.TeacherController{}, "*:HandleAdminRequest")
	beego.Router("/teacher/revoke-admin", &controllers.TeacherController{}, "post:RevokeAdmin")

	// 教师端扣除积分记录路由
	beego.Router("/teacher/activity-records", &controllers.TeacherController{}, "get:ShowActivityRecords")
	beego.Router("/teacher/activity-records/list", &controllers.TeacherController{}, "get:GetActivityRecords")
	beego.Router("/teacher/activity-records/update", &controllers.TeacherController{}, "post:UpdateActivityRecord")

	// 学生端路由
	beego.Router("/student", &controllers.StudentController{}, "get:ShowStudentNav")
	beego.Router("/student/activities", &controllers.StudentController{}, "get:ShowActivities")
	beego.Router("/student/activities/list", &controllers.StudentController{}, "get:GetActivities")
	beego.Router("/student/shop", &controllers.StudentController{}, "get:ShowShop")
	beego.Router("/student/ranking", &controllers.StudentController{}, "get:ShowRanking")
	beego.Router("/student/club", &controllers.StudentController{}, "get:ShowClub")
	beego.Router("/student/apply", &controllers.StudentController{}, "post:ApplyForAdmin")
	beego.Router("/student/club-activities", &controllers.StudentController{}, "get:GetClubActivities")
	beego.Router("/student/submit-activity", &controllers.StudentController{}, "post:SubmitActivity")
	beego.Router("/student/request-admin", &controllers.StudentController{}, "post:RequestAdmin")
	beego.Router("/student/register-activity", &controllers.StudentController{}, "post:RegisterActivity")
	beego.Router("/student/registration-status", &controllers.StudentController{}, "get:GetRegistrationStatus")
	beego.Router("/student/my-activities", &controllers.StudentController{}, "get:ShowMyActivities")
	beego.Router("/student/my-activities/list", &controllers.StudentController{}, "get:GetMyActivities")
	beego.Router("/student/cancel-registration", &controllers.StudentController{}, "post:CancelRegistration")
	beego.Router("/student/activity-records", &controllers.StudentController{}, "get:ShowActivityRecords")
	beego.Router("/student/activity-records/attendance-count", &controllers.StudentController{}, "get:GetAttendanceCount")
	beego.Router("/student/activity-records/list", &controllers.StudentController{}, "get:GetActivityRecords")
	beego.Router("/student/activity-records/add", &controllers.StudentController{}, "post:AddActivityRecord")
	beego.Router("/student/points-records", &controllers.StudentController{}, "get:ShowPointsRecords")
	beego.Router("/student/points-records/list", &controllers.StudentController{}, "get:GetPointsRecords")
	beego.Router("/student/activities/get", &controllers.StudentController{}, "get:GetActivityById")
	beego.Router("/student/exchange", &controllers.StudentController{}, "post:ExchangeItem")
	beego.Router("/student/exchange-history", &controllers.StudentController{}, "get:ShowExchangeHistory")

	// 区块链信息路由
	beego.Router("/blockchain/info", &controllers.BlockchainController{}, "get:ShowBlockchainInfo")
	beego.Router("/teacher/blockchain-info", &controllers.BlockchainController{}, "get:ShowTeacherBlockchainInfo")

	//注册，实现了get请求方法之后，不会再访问默认方法
	beego.Router("/ZhuCe", &controllers.MainController{}, "get:ZhuCeGet;post:HuoQu")

	//重置密码
	beego.Router("/ChongZhiMiMa", &controllers.MainController{}, "get:ShowChongZhiMiMa;post:HandleChongZhiMiMa")

	// 退出路由
	beego.Router("/logout", &controllers.MainController{}, "get:Logout")

	beego.SetStaticPath("/static", "xiaoLianJiFen/xiaoLianJiFen/static")
}
