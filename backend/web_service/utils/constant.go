package utils

type BlockActionType string

const (
	GiveKleos       BlockActionType = "give_kleos"
	ShowMyKleos     BlockActionType = "show_kleos"
	ShowLeaderBoard BlockActionType = "leader_board"
	HelpKleos       BlockActionType = "help_kleos"
)

type ViewSubmissionType string

const (
	StartGiveKleos       ViewSubmissionType = "start_give_kleos"
	StartshowMyKleos                        = "show_kleos"
	StartshowLeaderBoard                    = "leader_board"
)

const DefaultGivenValue = 0
const DefaultReceivedValue = 0

//const EmployeeIdColumn = "employee id"
//const EmployeeNameColumn = "full name"
//const EmployeeEmailColumn = "company email id"
//const ManagerNameColumn = "direct manager name"
//const ManagerEmailColumn = "direct manager email id"
