package usersDb

import "gorm.io/gorm"
import managersDb "kleos/db/managerDb"
import hrbpDb "kleos/db/hrbp_db"
import departmentDb "kleos/db/department_db"

type UserEntity struct {
	gorm.Model
	SlackUserId   string `gorm:"column:slack_user_id"`
	Name          string `gorm:"column:user_name"`
	Email         string `gorm:"column:email"`
	SlackImageUrl string `gorm:"column:slack_image_url"`
	RealName      string `gorm:"column:real_name"`
	EmployeeId    string `gorm:"column:employee_id"`
	ManagerId     int    `gorm:"column:manager_id"`
	GivenCount    int    `gorm:"column:given_count"`
	ReceivedCount int    `gorm:"column:received_count"`
	HrbpId        int    `gorm:"column:hrbp_id"`
	DepartmentId  int    `gorm:"column:department_id"`
}

type AdvUserEntity struct {
	Id            int    `gorm:"column:id"`
	SlackUserId   string `gorm:"column:slack_user_id"`
	Email         string `gorm:"column:email"`
	Name          string `gorm:"column:user_name"`
	GivenCount    int    `gorm:"column:given_count"`
	ReceivedCount int    `gorm:"column:received_count"`
	ManagerId     int    `gorm:"column:manager_id"`
	HrbpId        int    `gorm:"column:hrbp_id"`
	DepartmentId  int    `gorm:"column:department_id"`
	EmployeeId    string `gorm:"column:employee_id"`
}

type UserIdInfo struct {
	Id            int    `gorm:"column:id"`
	SlackUserId   string `gorm:"column:slack_user_id"`
	Email         string `gorm:"column:email"`
	ManagerId     int    `gorm:"column:manager_id"`
	SlackImageUrl string `gorm:"column:slack_image_url"`
	RealName      string `gorm:"column:real_name"`
}

type UserCount struct {
	Id            int    `gorm:"column:id"`
	SlackUserId   string `gorm:"column:slack_user_id"`
	GivenCount    int    `gorm:"column:given_count"`
	ReceivedCount int    `gorm:"column:received_count"`
}

type UserData struct {
	Email         string `gorm:"column:email"`
	RealName      string `gorm:"column:real_name"`
	SlackImageUrl string `gorm:"column:slack_image_url"`
}

type CustomUserInfo struct {
	Id            int    `gorm:"column:id"`
	Email         string `gorm:"column:email"`
	RealName      string `gorm:"column:real_name"`
	SlackImageUrl string `gorm:"column:slack_image_url"`
}

type UserEmailIdCombo struct {
	Id    int    `gorm:"column:id"`
	Email string `gorm:"column:email"`
}

type UserSlackCombo struct {
	Id          int    `gorm:"column:id"`
	SlackUserId string `gorm:"column:slack_user_id"`
}

type UserGivenReceivedCount struct {
	Id            int `gorm:"column:id"`
	GivenCount    int `gorm:"column:given_count"`
	ReceivedCount int `gorm:"column:received_count"`
}

type UserWithForeignData struct {
	UserEntity
	managersDb.ManagerEntity
	departmentDb.DepartmentEntity
	hrbpDb.HrbpEntity
}

func (UserEntity) TableName() string {
	return "users"
}
