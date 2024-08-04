package request

type CreateUserRequest struct {
	Email          string `json:"email"`
	RealName       string `json:"RealName"`
	SlackUserId    string `json:"slackUserId"`
	SlackImageUrl  string `json:"slackImageUrl"`
	ManagerEmail   string `json:"managerEmail"`
	ManagerName    string `json:"managerName"`
	EmployeeId     string `json:"employeeId"`
	HrbpName       string `json:"hrbpName"`
	DepartmentName string `json:"departmentName"`
}
