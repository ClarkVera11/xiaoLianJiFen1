package models

import (
	"time"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

// Users 用户模型
type Users struct {
	Id             int64  `orm:"column(id);pk;auto"`
	Username       string `orm:"column(username);size(255);unique"`
	Password       string `orm:"column(password);size(255)"`
	Email          string `orm:"column(email);size(255);null"`
	Phone          string `orm:"column(phone);size(20);null"`
	Role_name      string `orm:"column(role_name);size(255)"`
	IsAdminRequest int8   `orm:"column(is_admin_request);default(0)"`
	Points         int    `orm:"column(points);default(100)"`
	Title          string `orm:"column(title);size(50);default(倔强青铜)"` // 用户头衔
	TxHash         string `orm:"column(tx_hash);size(66);null"`        // 区块链交易哈希
	BlockTimestamp int64  `orm:"column(block_timestamp);default(0)"`   // 区块链交易时间戳
	Total          int    `orm:"column(total);default(0)"`             // 总积分
	ActivityCount  int    `orm:"column(activity_count);default(0)"`    // 参与活动次数
	Rank           int    `orm:"-"`                                    // orm:"-" 表示该字段不映射到数据库，仅用于程序逻辑
	Exchange       string `orm:"type(text);null" json:"exchange"`
}

// Activities 活动模型
type Activities struct {
	Id          int64     `orm:"column(id);pk;auto"`                // 活动ID
	Name        string    `orm:"column(name);size(255)"`            // 活动名称
	StartTime   time.Time `orm:"column(start_time);type(datetime)"` // 活动开始时间
	EndTime     time.Time `orm:"column(end_time);type(datetime)"`   // 活动结束时间
	Location    string    `orm:"column(location);size(255)"`        // 活动地点
	Description string    `orm:"column(description);type(text)"`    // 活动描述
	Points      int       `orm:"column(points)"`                    // 活动积分
	Capacity    int       `orm:"column(capacity)"`                  // 人员数量
	CreatedAt   time.Time `orm:"column(created_at);auto_now_add"`   // 创建时间
	UpdatedAt   time.Time `orm:"column(updated_at);auto_now"`       // 更新时间
	Status      int8      `orm:"column(status);default(0)"`         // 活动状态：0-待审核 1-已通过 2-已拒绝 3-已结束
}

// ActivityRegistrations 活动报名模型
type ActivityRegistrations struct {
	Id         int64     `orm:"column(id);pk;auto"`              // 报名ID
	ActivityId int64     `orm:"column(activity_id)"`             // 活动ID
	UserId     int64     `orm:"column(user_id)"`                 // 用户ID
	Status     int8      `orm:"column(status);default(1)"`       // 报名状态：1-已报名 2-已取消
	CreatedAt  time.Time `orm:"column(created_at);auto_now_add"` // 报名时间
	UpdatedAt  time.Time `orm:"column(updated_at);auto_now"`     // 更新时间
}

// ActivityRecords 扣除积分记录模型
type ActivityRecords struct {
	Id              int64     `orm:"column(id);pk;auto"`              // 记录ID
	ActivityId      int64     `orm:"column(activity_id)"`             // 活动ID
	AttendanceCount int       `orm:"column(attendance_count)"`        // 实际到场人数
	AbsentStudents  string       `orm:"type(text)"` // 存储逗号分隔字符串
	CreatedBy       int64     `orm:"column(created_by)"`              // 记录创建者ID
	CreatedAt       time.Time    `orm:"auto_now_add;type(datetime)"`
}


// PointsRecord 积分记录模型
type PointsRecord struct {
	Id          int64     `orm:"column(id);pk;auto"`
	UserId      int64     `orm:"column(user_id)"`
	ActivityId  int64     `orm:"column(activity_id);null"`        // 活动ID，可为 NULL
	Points      int       `orm:"column(points)"`
	Description string    `orm:"column(description);size(255)"`
	Source      string    `orm:"column(source);size(50)"`         // 新增：来源，如 activity, exchange, admin
	CreatedAt   time.Time `orm:"column(created_at);auto_now_add"`
}


// UpdateUserTotalAndCount 统计本月所有积分记录，重算total和activity_count
func UpdateUserTotalAndCount() error {
	o := orm.NewOrm()

	// 先将所有用户的total和activity_count清零
	o.Raw("UPDATE users SET total = 0, activity_count = 0").Exec()

	// 获取本月开始时间
	now := time.Now()
	firstDayOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	// 查询本月所有积分记录
	var records []PointsRecord
	_, err := o.QueryTable("points_record").
		Filter("created_at__gte", firstDayOfMonth).
		All(&records)
	if err != nil {
		return err
	}

	// 用map统计每个用户的积分和次数
	userTotal := make(map[int64]int)
	userCount := make(map[int64]int)

	for _, record := range records {
		userTotal[record.UserId] += record.Points
		userCount[record.UserId] += 1
	}

	// 更新到users表
	for userId, total := range userTotal {
		count := userCount[userId]
		o.QueryTable("users").Filter("id", userId).Update(orm.Params{
			"total":          total,
			"activity_count": count,
		})
	}

	return nil
}

// ResetMonthlyTotal 每月重置用户总积分
func ResetMonthlyTotal() error {
	o := orm.NewOrm()
	_, err := o.Raw("UPDATE users SET total = 0").Exec()
	return err
}

// 初始化数据库
func init() {
	// 注册数据库驱动
	orm.RegisterDriver("mysql", orm.DRMySQL)

	// 注册数据库连接
	orm.RegisterDataBase("default", "mysql", "root:123456@tcp(127.0.0.1:3306)/xiaolianjifen?charset=utf8&parseTime=True&loc=Asia%2FShanghai")

	// 注册模型
	orm.RegisterModel(new(Users), new(Activities), new(ActivityRegistrations), new(ActivityRecords), new(PointsRecord))

	// 自动创建或更新表
	orm.RunSyncdb("default", false, true)

	// 执行区块链字段迁移脚本
	o := orm.NewOrm()
	_, err := o.Raw("ALTER TABLE users ADD COLUMN IF NOT EXISTS tx_hash VARCHAR(66) NULL COMMENT '区块链交易哈希'").Exec()
	if err != nil {
		// println("添加tx_hash字段失败:", err)
	}

	_, err = o.Raw("ALTER TABLE users ADD COLUMN IF NOT EXISTS block_timestamp BIGINT DEFAULT 0 COMMENT '区块链交易时间戳'").Exec()
	if err != nil {
		// println("添加block_timestamp字段失败:", err)
	}

	// 添加新字段
	_, err = o.Raw("ALTER TABLE users ADD COLUMN IF NOT EXISTS total INT DEFAULT 0 COMMENT '总积分'").Exec()
	if err != nil {
		// println("添加total字段失败:", err)
	}

	_, err = o.Raw("ALTER TABLE users ADD COLUMN IF NOT EXISTS activity_count INT DEFAULT 0 COMMENT '参与活动次数'").Exec()
	if err != nil {
		// println("添加activity_count字段失败:", err)
	}

	// 启动定时任务
	go func() {
		for {
			now := time.Now()
			// 计算下个月1号的时间
			nextMonth := now.AddDate(0, 1, 0)
			nextMonthFirstDay := time.Date(nextMonth.Year(), nextMonth.Month(), 1, 0, 0, 0, 0, nextMonth.Location())

			// 等待到下个月1号
			time.Sleep(nextMonthFirstDay.Sub(now))

			// 重置总积分
			err := ResetMonthlyTotal()
			if err != nil {
				println("重置总积分失败:", err)
			}
		}
	}()

	// 启动定时更新用户总积分和活动次数的任务
	go func() {
		for {
			// 每小时更新一次
			err := UpdateUserTotalAndCount()
			if err != nil {
				println("更新用户总积分和活动次数失败:", err)
			}
			time.Sleep(time.Hour)
		}
	}()

	// 插入初始数据
	insertInitialActivities()
	insertInitialUsers()
}

// 插入初始活动数据
func insertInitialActivities() {
	o := orm.NewOrm()

	// 定义活动数据
	activities := []Activities{
		{
			Name:        "篮球赛",
			StartTime:   parseTime("2025-02-02 08:00:00"), // 活动开始时间
			EndTime:     parseTime("2025-02-02 09:00:00"), // 假设活动结束时间
			Location:    "篮球场",
			Description: "5v5男生全场篮球赛",
			Points:      7,
			Capacity:    10,
			Status:      3, // 已结束
		},
		{
			Name:        "围棋比赛",
			StartTime:   parseTime("2025-03-19 08:00:00"), // 活动开始时间
			EndTime:     parseTime("2025-03-19 09:00:00"), // 假设活动结束时间
			Location:    "创客大厦四楼",
			Description: "擂台赛形式",
			Points:      8,
			Capacity:    20,
			Status:      3, // 已结束
		},
		{
			Name:        "义卖活动",
			StartTime:   parseTime("2025-03-07 08:00:00"), // 活动开始时间
			EndTime:     parseTime("2025-03-07 09:00:00"), // 假设活动结束时间
			Location:    "创操场",
			Description: "义卖不需要的物品",
			Points:      9,
			Capacity:    15,
			Status:      2, // 已拒绝
		},
		{
			Name:        "乒乓球赛",
			StartTime:   parseTime("2025-03-12 08:00:00"), // 活动开始时间
			EndTime:     parseTime("2025-03-12 09:00:00"), // 假设活动结束时间
			Location:    "乒乓球场",
			Description: "2v2双打",
			Points:      8,
			Capacity:    14,
			Status:      3, // 已结束
		},
	}

	// 插入数据
	for _, activity := range activities {
		// 检查活动是否已存在
		exist := o.QueryTable("activities").Filter("name", activity.Name).Exist()
		if !exist {
			_, err := o.Insert(&activity)
			if err != nil {
				println("Error inserting activity:", activity.Name, err)
			} else {
				println("Successfully inserted activity:", activity.Name)
			}
		}
	}
}

// 插入初始用户数据
// 插入初始用户数据
func insertInitialUsers() {
	o := orm.NewOrm()

	// 定义用户数据
	users := []Users{
		{
			Username:       "202201001",
			Password:       "Htc123456！",
			Email:          "3331245121@email.com",
			Phone:          "15478954211",
			Role_name:      "学生",
			IsAdminRequest: 0,
			Points:         100,
			Title:          "秩序白银",
		},
		{
			Username:       "202201002",
			Password:       "Lr123456！",
			Email:          "3331245122@email.com",
			Phone:          "15478954212",
			Role_name:      "学生",
			IsAdminRequest: 0,
			Points:         100,
			Title:          "秩序白银",
		},
		{
			Username:       "202201003",
			Password:       "Lxt123456！",
			Email:          "3331245123@email.com",
			Phone:          "15478954213",
			Role_name:      "学生",
			IsAdminRequest: 1,
			Points:         100,
			Title:          "秩序白银",
		},
		{
			Username:       "JX666666",
			Password:       "Abc123456！",
			Email:          "3331245111@email.com",
			Phone:          "15478954221",
			Role_name:      "教师",
			IsAdminRequest: 0,
			Points:         100,
			Title:          "秩序白银",
		},
		{
			Username:       "JX777777",
			Password:       "Bcd123456！",
			Email:          "3331245112@email.com",
			Phone:          "15478954222",
			Role_name:      "教师",
			IsAdminRequest: 0,
			Points:         100,
			Title:          "秩序白银",
		},
	}

	// 插入数据
	for _, user := range users {
		// 检查用户是否已存在
		exist := o.QueryTable("users").Filter("username", user.Username).Exist()
		if !exist {
			// 对密码进行加密
			hashedPwd, err := HashPassword(user.Password)
			if err != nil {
				println("Error hashing password:", user.Username, err)
				continue
			}
			user.Password = hashedPwd

			_, err = o.Insert(&user)
			if err != nil {
				println("Error inserting user:", user.Username, err)
			} else {
				println("Successfully inserted user:", user.Username)
			}
		}
	}
}

// 辅助函数：解析时间字符串
func parseTime(timeStr string) time.Time {
	t, err := time.Parse("2006-01-02 15:04:05", timeStr)
	if err != nil {
		println("Error parsing time:", err)
		return time.Now()
	}
	return t
}

// UpdateUserTitle 根据积分更新用户头衔
func (u *Users) UpdateUserTitle() {
	switch {
	case u.Points >= 400:
		u.Title = "荣耀王者"
	case u.Points >= 350:
		u.Title = "最强王者"
	case u.Points >= 300:
		u.Title = "至尊星耀"
	case u.Points >= 250:
		u.Title = "永恒钻石"
	case u.Points >= 200:
		u.Title = "尊贵铂金"
	case u.Points >= 150:
		u.Title = "荣耀黄金"
	case u.Points >= 100:
		u.Title = "秩序白银"
	case u.Points >= 49:
		u.Title = "倔强青铜"
	default:
		u.Title = "倔强青铜"
	}
}

// 在ActivityRecords结构体后添加AfterInsert方法
func (ar *ActivityRecords) AfterInsert() {
	o := orm.NewOrm()

	// 获取活动信息
	activity := Activities{Id: ar.ActivityId}
	err := o.Read(&activity)
	if err != nil {
		return
	}

	// 更新用户数据
	user := Users{Id: ar.CreatedBy}
	err = o.Read(&user)
	if err != nil {
		return
	}

	user.Total += activity.Points
	user.ActivityCount += 1
	o.Update(&user)
}

// 在PointsRecord结构体后添加AfterInsert方法
func (pr *PointsRecord) AfterInsert() {
	o := orm.NewOrm()

	// 获取本月开始时间
	now := time.Now()
	firstDayOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	// 统计该用户本月的总积分和次数
	var total, count int
	o.Raw("SELECT IFNULL(SUM(points),0), COUNT(*) FROM points_record WHERE user_id = ? AND created_at >= ?", pr.UserId, firstDayOfMonth).QueryRow(&total, &count)

	// 更新users表
	o.QueryTable("users").Filter("id", pr.UserId).Update(orm.Params{
		"total":          total,
		"activity_count": count,
	})
}
