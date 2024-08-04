package service

type AdminAllUserData struct {
	AdminAllUser []UserAdminData `json:"adminAllUser"`
}

type UserAdminData struct {
	KleosGiven    int      `json:"kleosGiven"`
	KleosReceived int      `json:"kleosReceived"`
	User          UserData `json:"user"`
}

type AdvAdminAllUserData struct {
	AdminAllUser []AdvUserAdminData `json:"adminAllUser"`
}

type AdvUserAdminData struct {
	KleosGiven    int         `json:"kleosGiven"`
	KleosReceived int         `json:"kleosReceived"`
	User          AdvUserData `json:"user"`
}

type AdvUserData struct {
	Email      string         `json:"email"`
	UserName   string         `json:"userName"`
	Hrbp       HrbpData       `json:"hrbp"`
	Manager    ManagerData    `json:"manager"`
	Department DepartmentData `json:"department"`
	EmployeeId string         `json:"employeeId"`
}

type ManagerData struct {
	Email string `json:"email"`
	Name  string `json:"userName"`
}

type HrbpData struct {
	Email string `json:"email"`
	Name  string `json:"userName"`
}

type DepartmentData struct {
	Name string `json:"name"`
}

type AchievementData struct {
	Achievement string `json:"achievement"`
}

type KleosAllData struct {
	SenderData      UserData        `json:"senderData"`
	ReceiverData    UserData        `json:"receiverData"`
	AchievementData AchievementData `json:"achievementData"`
	Message         string          `json:"message"`
}
