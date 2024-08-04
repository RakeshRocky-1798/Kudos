package kleosDb

import (
	"gorm.io/gorm"
	"time"
)

type MyKleosEntity struct {
	gorm.Model
	SenderID    string `gorm:"column:sender_id"`
	Message     string `gorm:"column:message"`
	Achievement string `gorm:"column:achievement"`
	ReceiverID  string `gorm:"column:receiver_id"`
	Year        string `gorm:"column:year"`
	Month       string `gorm:"column:month"`
	Week        string `gorm:"column:week"`
	Day         string `gorm:"column:day"`
}

type FilteredKleosEntity struct {
	Id          int       `gorm:"column:id"`
	SenderID    string    `gorm:"column:sender_id"`
	ReceiverID  string    `gorm:"column:receiver_id"`
	Achievement string    `gorm:"column:achievement"`
	Message     string    `gorm:"column:message"`
	CreatedAt   time.Time `gorm:"column:created_at"`
}

func (FilteredKleosEntity) TableName() string {
	return "kleos"
}

func (MyKleosEntity) TableName() string {
	return "kleos"
}

type KleosAchievementStruct struct {
	Achievement string `gorm:"column:achievement"`
	Count       int    `gorm:"column:count__"`
}

type LeaderBoardGroupStruct struct {
	ReceiverID string `gorm:"column:receiver_id"`
	Count      int    `gorm:"column:count__"`
}

type KleosGiven struct {
	SenderID string `gorm:"column:sender_id"`
	Count    int    `gorm:"column:count__"`
}

type KleosReceived struct {
	ReceiverID string `gorm:"column:receiver_id"`
	Count      int    `gorm:"column:count__"`
}

type LastThreeKleos struct {
	SenderID    string `gorm:"column:sender_id"`
	Message     string `gorm:"column:message"`
	Achievement string `gorm:"column:achievement"`
}

type MyKleosData struct {
	MyUserID      string
	ReceivedCount int
	GivenCount    int
}

type LeaderBoardData struct {
	UserId string `gorm:"column:user_id"`
	Count  int    `gorm:"column:count__"`
	Rank   int    `gorm:"column:rank"`
}

type AdminAllUser struct {
	UserId        string `gorm:"column:user_id"`
	KleosGiven    int    `gorm:"column:kleos_given"`
	KleosReceived int    `gorm:"column:kleos_received"`
}
