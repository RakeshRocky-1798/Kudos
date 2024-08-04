package web_service

import (
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
	"github.com/spf13/viper"
	"github.com/xuri/excelize/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"io"
	slackbotconfig "kleos/config/slack_config"
	"kleos/cron"
	achievementdb "kleos/db/achievement_db"
	departmentDb "kleos/db/department_db"
	hrbpDb "kleos/db/hrbp_db"
	"kleos/db/kleosDb"
	managersDb "kleos/db/managerDb"
	userCountDb "kleos/db/user_count_db"
	"kleos/db/usersDb"
	"kleos/slack_service/processor/action"
	"kleos/web_service/request"
	service "kleos/web_service/response"
	common "kleos/web_service/response/common"
	"kleos/web_service/utils"
	"math"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type KleosService struct {
	gin                *gin.Engine
	logger             *zap.Logger
	db                 *gorm.DB
	socketModeClient   *socketmode.Client
	hNudgeSocketMode   *socketmode.Client
	slackbotClient     *slackbotconfig.Client
	kleosDataService   *kleosDb.KleosRepository
	achievementService *achievementdb.AchievementRepository
	managerService     *managersDb.ManagerRepository
	hrbpService        *hrbpDb.HrbpRepository
	departmentService  *departmentDb.DepartmentRepository
	userDataService    *usersDb.UserRepository
	userCountService   *userCountDb.UserCountRepository
	giveKleosAction    *action.GiveKleosAction
}

func NewKleosService(gin *gin.Engine, logger *zap.Logger, db *gorm.DB,
	socketModeClient *socketmode.Client,
	slackbotClient *slackbotconfig.Client, socketModeClient2 *socketmode.Client) *KleosService {

	kleosDataService := kleosDb.NewKleosRepository(logger, db)
	achievementService := achievementdb.NewAchievementRepository(logger, db)
	userDataService := usersDb.NewUserRepository(logger, db)
	managerDataService := managersDb.NewManagerRepository(logger, db)
	userCountService := userCountDb.NewUserCountRepository(logger, db)
	hrbpService := hrbpDb.NewHrbpRepository(logger, db)
	departmentService := departmentDb.NewDepartmentRepository(logger, db)

	return &KleosService{
		gin:                gin,
		logger:             logger,
		db:                 db,
		socketModeClient:   socketModeClient,
		hNudgeSocketMode:   socketModeClient2,
		slackbotClient:     slackbotClient,
		kleosDataService:   kleosDataService,
		achievementService: achievementService,
		userDataService:    userDataService,
		managerService:     managerDataService,
		hrbpService:        hrbpService,
		departmentService:  departmentService,
		userCountService:   userCountService,

		giveKleosAction: action.NewGiveKleosAction(socketModeClient, logger,
			slackbotClient, kleosDataService, achievementService, userDataService, managerDataService, userCountService),
	}
}

//To be used for UI Dashboard User

func (k *KleosService) GetKleosReceived(c *gin.Context) {
	UserEmail := c.Param("userid")
	var userResponse service.KleosReceivedResponse
	userData, err := k.userDataService.GetUserDataFromGmail(UserEmail)
	if err != nil {
		k.logger.Error("Error while fetching user data", zap.Error(err))
		c.JSON(http.StatusBadGateway, common.ErrorResponse(err, http.StatusBadGateway, nil))
		return
	}

	user := strconv.Itoa(int(userData.ID))
	//givenCount, err := k.kleosDataService.KleosGivenCountPerUser(user)
	//if err != nil {
	//	k.logger.Error("Error while fetching given count", zap.Error(err))
	//	c.JSON(http.StatusBadGateway, common.ErrorResponse(err, http.StatusBadGateway, nil))
	//	return
	//}
	//receivedCount, err := k.kleosDataService.KleosReceivedPerUser(user)
	//if err != nil {
	//	k.logger.Error("Error while fetching received count", zap.Error(err))
	//	c.JSON(http.StatusBadGateway, common.ErrorResponse(err, http.StatusBadGateway, nil))
	//	return
	//}

	myCount, err := k.userDataService.GetUserCount(user)
	if err != nil {
		k.logger.Error("Error while fetching user count", zap.Error(err))
		c.JSON(http.StatusBadGateway, common.ErrorResponse(err, http.StatusBadGateway, nil))
		return
	}

	userResponse = service.KleosReceivedResponse{
		UserId:        strconv.Itoa(int(userData.ID)),
		UserEmail:     userData.Email,
		GivenCount:    strconv.Itoa(myCount.GivenCount),
		ReceivedCount: strconv.Itoa(myCount.ReceivedCount),
	}

	c.JSON(http.StatusOK, common.SuccessResponse(userResponse, http.StatusOK))
}

func (k *KleosService) GetKleosDashboard(c *gin.Context) {
	UserEmail := c.Param("userid")
	var dashboardResponse service.DashboardResponse
	userData, err := k.userDataService.GetUserDataFromGmail(UserEmail)
	if err != nil {
		k.logger.Error("Error while fetching user data", zap.Error(err))
		c.JSON(http.StatusBadGateway, common.ErrorResponse(err, http.StatusBadGateway, nil))
		return
	}

	if userData.ID == 0 {

		slackId := ""
		userName := ""
		slackImage := ""

		api := k.slackbotClient.SocketModeClient
		info, err := api.GetUserInfo(UserEmail)
		if err != nil {
			if err.Error() == "user_not_found" {
				userName = strings.Split(UserEmail, "@")[0]
			}
			userName = strings.Split(UserEmail, "@")[0]
		} else {
			userName = info.RealName
			slackId = info.ID
			slackImage = info.Profile.Image192
		}

		var userNewData = &usersDb.UserRequest{
			SlackUserId:   slackId,
			UserName:      userName,
			Email:         UserEmail,
			SlackImageUrl: slackImage,
			RealName:      userName,
			GivenCount:    0,
			ReceivedCount: 0,
		}

		_, err = k.userDataService.AddUser(userNewData)
		if err != nil {
			k.logger.Error("Error while adding user", zap.Error(err))
			c.JSON(http.StatusBadGateway, common.ErrorResponse(err, http.StatusBadGateway, nil))
			return
		}

		userData, err = k.userDataService.GetUserDataFromGmail(UserEmail)
		if err != nil {
			k.logger.Error("Error while fetching user data", zap.Error(err))
			c.JSON(http.StatusBadGateway, common.ErrorResponse(err, http.StatusBadGateway, nil))
			return
		}
	}

	myDataResponse := service.UserData{
		Email:      userData.Email,
		UserName:   userData.RealName,
		ProfileUrl: userData.SlackImageUrl,
	}
	dashboardResponse.MyData = myDataResponse

	user := strconv.Itoa(int(userData.ID))

	//givenCount, err := k.kleosDataService.KleosGivenCountPerUser(user)
	//if err != nil {
	//	k.logger.Error("Error while fetching given count", zap.Error(err))
	//	c.JSON(http.StatusBadGateway, common.ErrorResponse(err, http.StatusBadGateway, nil))
	//	return
	//}
	//receivedCount, err := k.kleosDataService.KleosReceivedPerUser(user)
	//if err != nil {
	//	k.logger.Error("Error while fetching received count", zap.Error(err))
	//	c.JSON(http.StatusBadGateway, common.ErrorResponse(err, http.StatusBadGateway, nil))
	//	return
	//}

	myCount, err := k.userDataService.GetUserCount(user)
	if err != nil {
		k.logger.Error("Error while fetching user count", zap.Error(err))
		c.JSON(http.StatusBadGateway, common.ErrorResponse(err, http.StatusBadGateway, nil))
		return
	}

	dashboardResponse.KleosMetrics = service.KleosMetrics{
		GivenCount:    strconv.Itoa(myCount.GivenCount),
		ReceivedCount: strconv.Itoa(myCount.ReceivedCount),
	}

	achievementData, err := k.achievementService.GetAllAchievementName()
	if err != nil {
		k.logger.Error("Error while fetching achievement data", zap.Error(err))
		c.JSON(http.StatusBadGateway, common.ErrorResponse(err, http.StatusBadGateway, nil))
		return
	}

	for _, achData := range achievementData {
		achievementAllData := service.Options{
			Label: achData.DisplayName,
			Value: achData.Emoji,
		}
		dashboardResponse.AchievementDropDown = append(dashboardResponse.AchievementDropDown, achievementAllData)
	}

	achievementCount, err := k.kleosDataService.GetMyKleosPerAchievement(user)
	if err != nil {
		k.logger.Error("Error while fetching achievement data", zap.Error(err))
		c.JSON(http.StatusBadGateway, common.ErrorResponse(err, http.StatusBadGateway, nil))
		return
	}

	for _, achData := range *achievementCount {
		achievementId, _ := strconv.Atoi(achData.Achievement)
		for _, temp := range achievementData {
			if int(temp.ID) == achievementId {
				achievementData := service.AchievementOptions{
					AchievementName: temp.DisplayName,
					Count:           achData.Count,
					Emoji:           temp.Emoji,
				}
				dashboardResponse.TotalAchievement = append(dashboardResponse.TotalAchievement, achievementData)
			}
		}
	}

	// get last 3 my kleos data

	lastThreeData, err := k.kleosDataService.LastThreeKleosData(user)
	if err != nil {
		k.logger.Error("Error while fetching last three data", zap.Error(err))
		c.JSON(http.StatusBadGateway, common.ErrorResponse(err, http.StatusBadGateway, nil))
		return
	}
	for _, threeData := range *lastThreeData {
		message := threeData.Message
		achievementId, _ := strconv.Atoi(threeData.Achievement)

		for _, temp := range achievementData {
			if int(temp.ID) == achievementId {
				achievementFrom := threeData.SenderID
				userData, err := k.userDataService.GetUserInfoFromId(achievementFrom)
				if err != nil {
					k.logger.Error("Error while fetching user data", zap.Error(err))
					c.JSON(http.StatusBadGateway, common.ErrorResponse(err, http.StatusBadGateway, nil))
					return
				}
				receivedData := service.RecentKleosResponse{
					Message: message,
					Achievement: service.AdvancedAchievementResponse{
						ALabel: temp.DisplayName,
						AEmoji: temp.Emoji,
						AFrom: service.UserData{
							Email:      userData.Email,
							UserName:   userData.RealName,
							ProfileUrl: userData.SlackImageUrl,
						},
					},
				}
				dashboardResponse.RecentRecognition = append(dashboardResponse.RecentRecognition, receivedData)
			}
		}
	}
	c.JSON(http.StatusOK, common.SuccessResponse(dashboardResponse, http.StatusOK))
}

func (k *KleosService) GetAllUsers(c *gin.Context) {

	UserEmail := c.Query("userid")

	var UserLabelValue []service.Options
	userData, err := k.userDataService.GetAllUserData(UserEmail)
	if err != nil {
		k.logger.Error("Error while fetching user data", zap.Error(err))
		c.JSON(http.StatusBadGateway, common.ErrorResponse(err, http.StatusBadGateway, nil))
		return
	}
	for _, userData := range *userData {
		userData := service.Options{
			Label: userData.Email,
			Value: strconv.Itoa(userData.Id),
		}
		UserLabelValue = append(UserLabelValue, userData)
	}
	c.JSON(http.StatusOK, common.SuccessResponse(UserLabelValue, http.StatusOK))
}

func (k *KleosService) GetPaginatedInfo(c *gin.Context) {

	UserEmail := c.Param("userid")
	dataType := c.Query("data_type")
	paramPageSize := c.Query("page_size")
	paramPageNumber := c.Query("page_number")
	pageSize, pageNumber, err := utils.ValidatePage(paramPageSize, paramPageNumber)
	if err != nil {
		k.logger.Error("Error while validating page", zap.Error(err))
		c.JSON(http.StatusBadRequest, common.ErrorResponse(err, http.StatusBadRequest, nil))
		return
	}

	userData, err := k.userDataService.GetUserDataFromGmail(UserEmail)
	if err != nil {
		k.logger.Error("Error while fetching user data", zap.Error(err))
		c.JSON(http.StatusBadGateway, common.ErrorResponse(err, http.StatusBadGateway, nil))
		return
	}
	userId := strconv.Itoa(int(userData.ID))
	kleos, totalElement, err := k.kleosDataService.FetchPaginatedKleos(dataType, pageNumber, pageSize, userId)
	if err != nil {
		k.logger.Error("Error while fetching kleos data", zap.Error(err))
		c.JSON(http.StatusBadGateway, common.ErrorResponse(err, http.StatusBadGateway, nil))
		return
	}

	data, err := k.GetAllUserKleosData(kleos, dataType, c)
	if err != nil {
		k.logger.Error("Error while fetching kleos data", zap.Error(err))
		c.JSON(http.StatusBadGateway, common.ErrorResponse(err, http.StatusBadGateway, nil))
		return
	}

	totalPages := int64(math.Ceil(float64(totalElement) / float64(pageSize)))
	hasData := pageNumber < totalPages-1

	pageData := common.Page{
		PageSize:      pageSize,
		TotalPages:    int64(math.Ceil(float64(totalElement) / float64(pageSize))),
		PageNumber:    pageNumber,
		TotalElements: totalElement,
		HasData:       hasData,
	}
	c.JSON(http.StatusOK, common.SuccessPaginatedResponse(data, pageData, http.StatusOK))
}

func (k *KleosService) GetAllUserKleosData(kleosData []kleosDb.FilteredKleosEntity,
	dataType string, c *gin.Context) ([]service.FilteredKleosResponse, error) {

	//fetch all achievement data
	var filteredKleosResponse []service.FilteredKleosResponse
	achievementData, err := k.achievementService.GetAllAchievementName()
	if err != nil {
		k.logger.Error("Error while fetching achievement data", zap.Error(err))
		c.JSON(http.StatusBadGateway, common.ErrorResponse(err, http.StatusBadGateway, nil))
		return nil, err
	}
	for _, kleos := range kleosData {
		Id := kleos.Id
		message := kleos.Message
		achievementId, _ := strconv.Atoi(kleos.Achievement)
		createdAt := kleos.CreatedAt

		for _, temp := range achievementData {
			switch dataType {
			case "given":
				{
					if int(temp.ID) == achievementId {
						userData, err := k.userDataService.GetUserInfoFromId(kleos.ReceiverID)
						if err != nil {
							k.logger.Error("Error while fetching user data", zap.Error(err))
							c.JSON(http.StatusBadGateway,
								common.ErrorResponse(err, http.StatusBadGateway, nil))
							return nil, err
						}
						receivedData := service.FilteredKleosResponse{
							Id:      Id,
							Message: message,
							Achievement: service.AchievementOption{
								AType:      temp.DisplayName,
								ACreatedAt: createdAt.Format("2006-01-02 15:04:05.999999"),
								AEmoji:     temp.Emoji,
								User: service.UserData{
									Email:      userData.Email,
									UserName:   userData.RealName,
									ProfileUrl: userData.SlackImageUrl,
								},
							},
						}
						filteredKleosResponse = append(filteredKleosResponse, receivedData)
					}
				}
			case "received":
				{
					if int(temp.ID) == achievementId {
						userData, err := k.userDataService.GetUserInfoFromId(kleos.SenderID)
						if err != nil {
							k.logger.Error("Error while fetching user data", zap.Error(err))
							c.JSON(http.StatusBadGateway, common.ErrorResponse(err,
								http.StatusBadGateway, nil))
							return nil, err
						}
						receivedData := service.FilteredKleosResponse{
							Id:      Id,
							Message: message,
							Achievement: service.AchievementOption{
								AType:      temp.DisplayName,
								ACreatedAt: createdAt.Format("2006-01-02 15:04:05.999999"),
								AEmoji:     temp.Emoji,
								User: service.UserData{
									Email:      userData.Email,
									UserName:   userData.RealName,
									ProfileUrl: userData.SlackImageUrl,
								},
							},
						}
						filteredKleosResponse = append(filteredKleosResponse, receivedData)
					}
				}

			default:
				{
					return nil, errors.New("not a valid data type")
				}
			}
		}
	}
	return filteredKleosResponse, nil
}

func (k *KleosService) ValidatePayload(kleosRequest request.KleosRequest) (bool, error) {

	senderId := kleosRequest.SenderId

	achievementId, _ := strconv.Atoi(kleosRequest.Achievement)

	receiverIdList := kleosRequest.ReceiverId

	message := kleosRequest.Message

	if kleosRequest.ReceiverId == nil {
		return false, errors.New("receiver id is empty")
	}
	if achievementId == 0 {
		return false, errors.New("achievement is empty")
	}
	if message == "" {
		return false, errors.New("message is empty")
	}

	for _, receiverId := range receiverIdList {
		id, err := strconv.Atoi(receiverId)
		if err != nil {
			return false, errors.New("receiver id is not a valid integer")
		}

		exist, err := k.userDataService.CheckIfUserExists(id)
		if err != nil || !exist {
			return false, errors.New(fmt.Sprintf("user with id %d does not exist", id))
		}

		if receiverId == senderId {
			return false, errors.New("cannot give kleos to yourself")
		}
	}

	exist, err := k.achievementService.CheckIfAchievementExists(achievementId)
	if err != nil || !exist {
		return false, errors.New("achievement given does not exist")
	}

	return true, nil
}

func (k *KleosService) GiveKleosFromWeb(c *gin.Context) {

	var kleosRequest request.KleosRequest

	email := c.GetString("email")

	if email == "" {
		k.logger.Error("Error while fetching email", zap.String("invalid Email", email))
		c.JSON(http.StatusUnauthorized, common.ErrorResponse(errors.New("unauthorized"), http.StatusUnauthorized, nil))
		return
	}

	err := c.ShouldBindJSON(&kleosRequest)
	if err != nil {
		k.logger.Error("Error while binding request", zap.Error(err))
		c.JSON(http.StatusInternalServerError, common.ErrorResponse(err,
			http.StatusInternalServerError, nil))
		return
	}

	kleosFromId, _ := strconv.Atoi(kleosRequest.SenderId)
	loggedUserId, err := k.userDataService.GetSenderIdFromEmail(email)

	k.logger.Info("Auth Check", zap.Int("Logged in user:", loggedUserId), zap.Int("Kleos from user:", kleosFromId))

	if kleosFromId != loggedUserId {
		k.logger.Error("Error while validating user", zap.Error(err))
		c.JSON(http.StatusUnauthorized, common.ErrorResponse(errors.New("unauthorized"), http.StatusUnauthorized, nil))
		return
	}

	validPayload, err := k.ValidatePayload(kleosRequest)

	if !validPayload {
		k.logger.Error("Invalid payload", zap.Error(err))
		c.JSON(http.StatusBadRequest, common.ErrorResponse(err, http.StatusBadRequest, nil))
		return
	}

	var multiselect bool
	receiverId := ""
	kleosToId, _ := strconv.Atoi(receiverId)
	achievementId, _ := strconv.Atoi(kleosRequest.Achievement)
	message := kleosRequest.Message
	slackFlag := kleosRequest.NeedSlack
	currentTime := time.Now()
	_, currentWeek := currentTime.ISOWeek()

	limitPerWeekGiven := viper.GetInt("limit.per.week.g")

	givenDetails, err := k.userCountService.GetUserGiveCountCurrentWeek(strconv.Itoa(kleosFromId), fmt.Sprintf("%d", currentWeek))
	if err != nil {
		k.logger.Error("Error while fetching received user count", zap.Error(err))
		c.JSON(http.StatusBadGateway, common.ErrorResponse(err, http.StatusBadGateway, nil))
		return
	}

	givenFrom, err := k.userDataService.GetSlackIdFromId(kleosFromId)
	if err != nil {
		k.logger.Error("Error while fetching user data", zap.Error(err))
		c.JSON(http.StatusBadGateway, common.ErrorResponse(err, http.StatusBadGateway, nil))
		return
	}

	if limitPerWeekGiven <= givenDetails.GivenCount {
		err := errors.New("[given] limit exceeded")
		c.JSON(http.StatusBadGateway, common.ErrorResponse(err, http.StatusBadGateway, nil))
		return
	}

	if len(kleosRequest.ReceiverId) == 1 {
		multiselect = false
	} else {
		multiselect = true
	}

	if !multiselect {
		receiverId = kleosRequest.ReceiverId[0]
	} else {
		if limitPerWeekGiven <= givenDetails.GivenCount {
			err := errors.New("[given] limit exceeded")
			c.JSON(http.StatusBadGateway, common.ErrorResponse(err, http.StatusBadGateway, nil))
			return
		}

		for _, receiver := range kleosRequest.ReceiverId {
			kleosToId, _ := strconv.Atoi(receiver)
			givenTo, err := k.userDataService.GetSlackIdFromId(kleosToId)
			if err != nil {
				k.logger.Error("Error while fetching user data", zap.Error(err))
				c.JSON(http.StatusBadGateway, common.ErrorResponse(err, http.StatusBadGateway, nil))
				return
			}
			if givenFrom.Id == givenTo.Id {
				err := errors.New("cannot give kleos to yourself")
				c.JSON(http.StatusBadGateway, common.ErrorResponse(err, http.StatusBadGateway, nil))
				return
			}
			k.GiveKleosOneByOne(multiselect, givenFrom, givenTo, achievementId, message, currentTime, kleosRequest, slackFlag, currentWeek)
		}

		k.userCountService.UpdateGivenCount(strconv.Itoa(givenFrom.Id), fmt.Sprintf("%d", currentWeek))
		err = k.userDataService.UpdateUserGivenCount(strconv.Itoa(givenFrom.Id))
		if err != nil {
			k.logger.Error("Error while updating user given count", zap.Error(err))
			c.JSON(http.StatusBadGateway, common.ErrorResponse(err, http.StatusBadGateway, nil))
			return
		}

		c.JSON(http.StatusOK, common.SuccessResponse("Data Successfully Added", http.StatusOK))
		return
	}

	kleosToId, _ = strconv.Atoi(receiverId)
	givenTo, err := k.userDataService.GetSlackIdFromId(kleosToId)
	if err != nil {
		k.logger.Error("Error while fetching user data", zap.Error(err))
		c.JSON(http.StatusBadGateway, common.ErrorResponse(err, http.StatusBadGateway, nil))
		return
	}

	limitOneOnOne := viper.GetInt("limit.one.on.one")

	if givenTo.Id == givenFrom.Id {
		err := errors.New("cannot give kleos to yourself")
		c.JSON(http.StatusBadGateway, common.ErrorResponse(err, http.StatusBadGateway, nil))
		return
	}

	oneOnOneDetails, err := k.kleosDataService.GetOneOnOneCount(strconv.Itoa(kleosFromId), strconv.Itoa(kleosToId), strconv.Itoa(currentWeek))
	if err != nil {
		k.logger.Error("Error while fetching one on one count", zap.Error(err))
		c.JSON(http.StatusBadGateway, common.ErrorResponse(err, http.StatusBadGateway, nil))
		return
	}

	if limitOneOnOne <= oneOnOneDetails.Count {
		err := errors.New("[one on one] limit exceeded")
		c.JSON(http.StatusBadGateway, common.ErrorResponse(err, http.StatusBadGateway, nil))
		return
	}

	k.GiveKleosOneByOne(multiselect, givenFrom, givenTo, achievementId, message, currentTime, kleosRequest, slackFlag, currentWeek)
	c.JSON(http.StatusOK, common.SuccessResponse("Data Successfully Added", http.StatusOK))
	return
}

func (k *KleosService) GiveKleosOneByOne(multiSelect bool, kleosFrom *usersDb.UserIdInfo, kleosTo *usersDb.UserIdInfo,
	achievementId int, message string, currentTime time.Time, kleosRequest request.KleosRequest, slackFlag bool,
	currentWeek int) *gin.Context {

	var c *gin.Context

	managerInfo, err := k.managerService.GetManagerDataFromId(kleosTo.ManagerId)
	if err != nil {
		k.logger.Error("Error while fetching manager data", zap.Error(err))
		c.JSON(http.StatusBadGateway, common.ErrorResponse(err, http.StatusBadGateway, nil))
		return c
	}
	k.logger.Info("[Kleos_Web] Kleos data", zap.String("givenFrom", managerInfo.SlackUserId))

	achievement, err := k.achievementService.GetAchievementNameFromID(achievementId)
	if err != nil {
		k.logger.Error("Error while fetching achievement data", zap.Error(err))
		c.JSON(http.StatusBadGateway, common.ErrorResponse(err, http.StatusBadGateway, nil))
		return c
	}

	achievementName := achievement.AchievementName
	k.logger.Info("[Kleos_Web] Kleos data", zap.String("givenFrom", kleosFrom.SlackUserId),
		zap.String("givenTo", kleosTo.SlackUserId),
		zap.String("achievement", achievementName),
		zap.String("message", message))

	var data = &kleosDb.CreateKleosRequest{
		SenderID:    kleosRequest.SenderId,
		Message:     message,
		Achievement: kleosRequest.Achievement,
		ReceiverID:  strconv.Itoa(kleosTo.Id),
		Year:        fmt.Sprintf("%d", currentTime.Year()),
		Month:       fmt.Sprintf("%d", currentTime.Month()),
		Week:        fmt.Sprintf("%d", currentWeek),
		Day:         fmt.Sprintf("%d", currentTime.Day()),
	}

	_, err11 := k.kleosDataService.GiveKleos(data)
	if err11 != nil {
		k.logger.Error("[CIP] Error while giving Kleos", zap.Error(err11))
		c.JSON(http.StatusBadGateway, common.ErrorResponse(err11, http.StatusBadGateway, nil))
		return c
	}

	ReceiverCountData, err := k.userCountService.CurrentWeekCount(strconv.Itoa(kleosTo.Id), fmt.Sprintf("%d", currentWeek))
	if err != nil {
		k.logger.Error("Error while fetching received user count", zap.Error(err))
		c.JSON(http.StatusBadGateway, common.ErrorResponse(err, http.StatusBadGateway, nil))
		return c
	}

	if ReceiverCountData.UserId == "" {
		var userCountData = &userCountDb.CreateUserCountRequest{
			UserId:        strconv.Itoa(kleosTo.Id),
			GivenCount:    0,
			ReceivedCount: 1,
			Month:         fmt.Sprintf("%d", currentTime.Month()),
			Week:          fmt.Sprintf("%d", currentWeek),
		}
		_, err := k.userCountService.AddUserCountEntity(userCountData)
		if err != nil {
			k.logger.Error("Error while adding user count", zap.Error(err))
			c.JSON(http.StatusBadGateway, common.ErrorResponse(err, http.StatusBadGateway, nil))
			return c
		}
	} else {
		k.userCountService.UpdateReceivedCount(strconv.Itoa(kleosTo.Id), fmt.Sprintf("%d", currentWeek))
	}

	err = k.userDataService.UpdateUserReceivedCount(strconv.Itoa(kleosTo.Id))
	if err != nil {
		k.logger.Error("Error while updating user received count", zap.Error(err))
		c.JSON(http.StatusBadGateway, common.ErrorResponse(err, http.StatusBadGateway, nil))
		return c
	}

	givenCountData, err := k.userCountService.CurrentWeekCount(strconv.Itoa(kleosFrom.Id), fmt.Sprintf("%d", currentWeek))
	if err != nil {
		k.logger.Error("Error while fetching given user count", zap.Error(err))
		c.JSON(http.StatusBadGateway, common.ErrorResponse(err, http.StatusBadGateway, nil))
		return c
	}

	if givenCountData.UserId == "" {
		var userCountData = &userCountDb.CreateUserCountRequest{
			UserId:        strconv.Itoa(kleosFrom.Id),
			GivenCount:    1,
			ReceivedCount: 0,
			Month:         fmt.Sprintf("%d", currentTime.Month()),
			Week:          fmt.Sprintf("%d", currentWeek),
		}
		_, err := k.userCountService.AddUserCountEntity(userCountData)
		if err != nil {
			k.logger.Error("Error while adding user count", zap.Error(err))
			c.JSON(http.StatusBadGateway, common.ErrorResponse(err, http.StatusBadGateway, nil))
			return c
		}
	} else {
		if !multiSelect {
			k.userCountService.UpdateGivenCount(strconv.Itoa(kleosFrom.Id), fmt.Sprintf("%d", currentWeek))
		}
	}

	if !multiSelect {
		err = k.userDataService.UpdateUserGivenCount(strconv.Itoa(kleosFrom.Id))
		if err != nil {
			k.logger.Error("Error while updating user given count", zap.Error(err))
			c.JSON(http.StatusBadGateway, common.ErrorResponse(err, http.StatusBadGateway, nil))
			return c
		}
	}

	receiver, sender := k.fetchSlackIdsIfNotPresent(kleosTo, kleosFrom)

	receiverSlackId := receiver.SlackUserId
	receiverEmail := receiver.Email
	receiverName := receiver.RealName

	managerEmail := managerInfo.Email

	senderSlackId := sender.SlackUserId
	senderName := sender.RealName
	senderSlackImage := sender.SlackImageUrl
	senderEmail := sender.Email

	k.giveKleosAction.PostSlackCommunicationsAndTriggerEmail(senderSlackId, receiverSlackId, achievementName, message,
		slackFlag, receiverEmail, receiverName, managerEmail, senderName, senderSlackImage, senderEmail)

	//
	//go func() {
	//	err := k.postSlackBotMessageWeb(kleosFrom.SlackUserId, kleosTo.SlackUserId, achievementName, message)
	//	if err != nil {
	//		k.logger.Error("Error while posting slack message", zap.Error(err))
	//	}
	//
	//	comms.TriggerComms(k.logger, kleosTo.Email, managerInfo.Email, kleosTo.RealName, message, kleosFrom.SlackImageUrl, kleosFrom.RealName, kleosFrom.Email)
	//
	//	if slackFlag == true {
	//		_, err := k.giveKleosAction.PostCardFromWeb(kleosTo.SlackUserId, message, achievementName, kleosTo.Email, kleosFrom.SlackUserId)
	//		if err != nil {
	//			k.logger.Error("[Web] [Slack Channel] Error while posting slack message", zap.Error(err))
	//		}
	//	}
	//	return
	//}()
	return c
}

func (k *KleosService) fetchSlackIdsIfNotPresent(receiver *usersDb.UserIdInfo, sender *usersDb.UserIdInfo) (*usersDb.UserIdInfo, *usersDb.UserIdInfo) {

	senderEmail := sender.Email
	senderSlackId := sender.SlackUserId
	receiverSlackId := receiver.SlackUserId
	receiverEmail := receiver.Email

	if receiverSlackId == "" {

		k.logger.Info("Receiver slack id is not present", zap.String("receiverEmail", receiverEmail))

		api := k.slackbotClient.SocketModeClient
		info, err := api.GetUserByEmail(receiverEmail)
		if err != nil {
			if err.Error() == "users_not_found" {
				err = k.userDataService.SetSlackIdFromWeb(receiverEmail, "nil")
				if err != nil {
					k.logger.Error("Error while setting the slack id", zap.Error(err))
				}
				k.logger.Error("User is not registered on slack", zap.Error(err))
			}
		}

		if info != nil {
			receiverSlackId = info.ID
			receiver.SlackUserId = receiverSlackId
		} else {
			k.logger.Info("Receiver slack id not retrieved", zap.String("receiverEmail", receiverEmail))
		}

		err = k.userDataService.SetSlackIdFromWeb(receiverEmail, receiverSlackId)

		if err != nil {
			k.logger.Error("Error while setting the slack id", zap.Error(err))
		}
	}

	if senderSlackId == "" {

		k.logger.Info("Sender slack id is not present", zap.String("senderEmail", senderEmail))

		api := k.slackbotClient.SocketModeClient
		info, err := api.GetUserByEmail(senderEmail)
		if err != nil {
			if err.Error() == "users_not_found" {
				err = k.userDataService.SetSlackIdFromWeb(senderEmail, "nil")
				if err != nil {
					k.logger.Error("Error while setting the slack id", zap.Error(err))
				}
				k.logger.Error("User is not registered on slack", zap.Error(err))
			}
		}

		if info != nil {
			senderSlackId = info.ID
			sender.SlackUserId = senderSlackId
		} else {
			k.logger.Info("Receiver slack id not retrieved", zap.String("senderEmail", senderEmail))
		}

		err = k.userDataService.SetSlackIdFromWeb(senderEmail, senderSlackId)

		if err != nil {
			k.logger.Error("Error while setting the slack id", zap.Error(err))
		}
	}

	return receiver, sender

}

func (k *KleosService) GetLeaderBoard(c *gin.Context) {

	userEmail := c.Query("userid")
	isReceived := c.Query("isReceived")
	isReceivedBool, err := strconv.ParseBool(isReceived)

	k.logger.Info("Inside LeaderBoard", zap.String("userEmail", userEmail))
	currentUser, err := k.userDataService.GetUserDataFromGmail(userEmail)

	k.logger.Info("Inside LeaderBoard", zap.String("userEmail", currentUser.Email))

	if err != nil {
		k.logger.Error("Error while fetching user data", zap.Error(err))
		c.JSON(http.StatusBadGateway, common.ErrorResponse(err, http.StatusBadGateway, nil))
		return
	}

	var leaderBoardResponse service.LeaderBoardResponse

	currentTime := time.Now()
	mST := fmt.Sprintf("%d", currentTime.Month())

	leaderBoardDataList, err := k.userCountService.GetUserLeaderBoardData(isReceivedBool, mST)

	if err != nil {
		k.logger.Error("Error while fetching leaderboard data", zap.Error(err))
		c.JSON(http.StatusBadGateway, common.ErrorResponse(err, http.StatusBadGateway, nil))
		return
	}

	currentUserFound := false

	for i, leaderBoard := range *leaderBoardDataList {

		if i < 10 {
			userDataInfo, err := k.userDataService.GetUserInfoFromId(leaderBoard.UserId)
			if err != nil {
				k.logger.Error("Error while fetching user data", zap.Error(err))
				c.JSON(http.StatusBadGateway, common.ErrorResponse(err, http.StatusBadGateway, nil))
				return
			}

			userMeta := service.UserData{
				Email:      userDataInfo.Email,
				UserName:   userDataInfo.RealName,
				ProfileUrl: userDataInfo.SlackImageUrl,
			}

			leaderBoardData := service.LeaderBoardData{
				UserMeta: userMeta,
				Rank:     leaderBoard.Rank,
				Count:    leaderBoard.Count,
			}

			if strconv.FormatUint(uint64(currentUser.ID), 10) == leaderBoard.UserId {
				currentUserFound = true
				leaderBoardData.IsCurrentUser = true
			}

			leaderBoardResponse.LeaderBoardData = append(leaderBoardResponse.LeaderBoardData, leaderBoardData)
		} else {
			if currentUserFound {
				break
			}

			//currentUserPosition, err := k.kleosDataService.GetCurrentUserData(isReceivedBool, strconv.Itoa(int(currentUser.ID)))
			//
			//if err != nil {
			//	k.logger.Error("Error while fetching current user data", zap.Error(err))
			//	c.JSON(http.StatusBadGateway, common.ErrorResponse(err, http.StatusBadGateway, nil))
			//	return
			//}

			if strconv.FormatUint(uint64(currentUser.ID), 10) == leaderBoard.UserId {
				currentUserFound = true

				userDataInfo, err := k.userDataService.GetUserInfoFromId(leaderBoard.UserId)
				if err != nil {
					k.logger.Error("Error while fetching user data", zap.Error(err))
					c.JSON(http.StatusBadGateway, common.ErrorResponse(err, http.StatusBadGateway, nil))
					return
				}

				userMeta := service.UserData{
					Email:      userDataInfo.Email,
					UserName:   userDataInfo.RealName,
					ProfileUrl: userDataInfo.SlackImageUrl,
				}

				leaderBoardData := service.LeaderBoardData{
					UserMeta:      userMeta,
					Rank:          leaderBoard.Rank,
					Count:         leaderBoard.Count,
					IsCurrentUser: true,
				}

				leaderBoardResponse.LeaderBoardData = append(leaderBoardResponse.LeaderBoardData, leaderBoardData)
			}
		}
	}

	c.JSON(http.StatusOK, common.SuccessResponse(leaderBoardResponse, http.StatusOK))
}

func (k *KleosService) GetAllAchievement(c *gin.Context) {

	var AchievementLabelValue []service.AchievementOptionOne
	k.logger.Info("Inside All Achievement")
	aEntity, err := k.achievementService.GetAllAchievementName()
	if err != nil {
		k.logger.Error("Error while fetching achievement data", zap.Error(err))
		c.JSON(http.StatusBadGateway, common.ErrorResponse(err, http.StatusBadGateway, nil))
		return
	}
	for _, achData := range aEntity {
		achievementAllData := service.AchievementOptionOne{
			Label:  achData.DisplayName,
			Value:  strconv.Itoa(int(achData.ID)),
			AEmoji: achData.Emoji,
		}
		AchievementLabelValue = append(AchievementLabelValue, achievementAllData)
	}

	c.JSON(http.StatusOK, common.SuccessResponse(AchievementLabelValue, http.StatusOK))
}

func (k *KleosService) GetAdminData(c *gin.Context) {

	var adminDataResponse service.AdminAllUserData
	adminData, err := k.kleosDataService.GetAdminDataTest()
	if err != nil {
		k.logger.Error("Error while fetching admin data", zap.Error(err))
		c.JSON(http.StatusBadGateway, common.ErrorResponse(err, http.StatusBadGateway, nil))
		return
	}

	for _, aData := range *adminData {
		userInfo, err := k.userDataService.GetUserInfoFromId(aData.UserId)
		myCount, err := k.userDataService.GetUserCount(aData.UserId)
		if err != nil {
			k.logger.Error("Error while fetching user count", zap.Error(err))
			c.JSON(http.StatusBadGateway, common.ErrorResponse(err, http.StatusBadGateway, nil))
			return
		}
		if err != nil {
			k.logger.Error("Error while fetching user data", zap.Error(err))
			c.JSON(http.StatusBadGateway, common.ErrorResponse(err, http.StatusBadGateway, nil))
			return
		}
		userData := service.UserData{
			Email:      userInfo.Email,
			UserName:   userInfo.RealName,
			ProfileUrl: userInfo.SlackImageUrl,
		}
		tempAdminData := service.UserAdminData{
			KleosGiven:    myCount.GivenCount,
			KleosReceived: myCount.ReceivedCount,
			User:          userData,
		}
		adminDataResponse.AdminAllUser = append(adminDataResponse.AdminAllUser, tempAdminData)
	}
	c.JSON(http.StatusOK, common.SuccessResponse(adminDataResponse, http.StatusOK))
}

func (k *KleosService) NewAdminData(c *gin.Context) {

	userRoles := c.MustGet("roles").([]string)

	isAdmin := strings.Contains(strings.Join(userRoles, ","), "admin")

	if !isAdmin {
		err := errors.New("unauthorized")
		c.JSON(http.StatusUnauthorized, common.ErrorResponse(err, http.StatusUnauthorized, nil))
		return
	}

	var adminDataResponse service.AdminAllUserData
	userAdminData, err := k.userDataService.GetAllUserCount()
	if err != nil {
		k.logger.Error("Error while fetching user count", zap.Error(err))
		c.JSON(http.StatusBadGateway, common.ErrorResponse(err, http.StatusBadGateway, nil))
		return
	}

	for _, aData := range *userAdminData {
		userInfo, err := k.userDataService.GetUserInfoFromId(strconv.Itoa(aData.Id))
		if err != nil {
			k.logger.Error("Error while fetching user data", zap.Error(err))
			c.JSON(http.StatusBadGateway, common.ErrorResponse(err, http.StatusBadGateway, nil))
			return
		}
		userData := service.UserData{
			Email:      userInfo.Email,
			UserName:   userInfo.RealName,
			ProfileUrl: userInfo.SlackImageUrl,
		}
		tempAdminData := service.UserAdminData{
			KleosGiven:    aData.GivenCount,
			KleosReceived: aData.ReceivedCount,
			User:          userData,
		}
		adminDataResponse.AdminAllUser = append(adminDataResponse.AdminAllUser, tempAdminData)
	}
	c.JSON(http.StatusOK, common.SuccessResponse(adminDataResponse, http.StatusOK))
}

func (k *KleosService) GetPaginatedAdminData(c *gin.Context) {

	paramPageSize := c.Query("page_size")
	paramPageNumber := c.Query("page_number")
	var adminDataResponse service.AdminAllUserData

	pageSize, pageNumber, err := utils.ValidatePage(paramPageSize, paramPageNumber)
	if err != nil {
		k.logger.Error("Error while validating page", zap.Error(err))
		c.JSON(http.StatusBadRequest, common.ErrorResponse(err, http.StatusBadRequest, nil))
		return
	}

	adminData, totalElement, err := k.kleosDataService.GetPaginatedAdminData(pageNumber, pageSize)
	if err != nil {
		k.logger.Error("Error while fetching admin data", zap.Error(err))
		c.JSON(http.StatusBadGateway, common.ErrorResponse(err, http.StatusBadGateway, nil))
		return
	}

	for _, aData := range adminData {
		userInfo, err := k.userDataService.GetUserInfoFromId(aData.UserId)
		if err != nil {
			k.logger.Error("Error while fetching user data", zap.Error(err))
			c.JSON(http.StatusBadGateway, common.ErrorResponse(err, http.StatusBadGateway, nil))
			return
		}
		userData := service.UserData{
			Email:      userInfo.Email,
			UserName:   userInfo.RealName,
			ProfileUrl: userInfo.SlackImageUrl,
		}
		tempAdminData := service.UserAdminData{
			KleosGiven:    aData.KleosGiven,
			KleosReceived: aData.KleosReceived,
			User:          userData,
		}
		adminDataResponse.AdminAllUser = append(adminDataResponse.AdminAllUser, tempAdminData)
	}
	totalPages := int64(math.Ceil(float64(totalElement) / float64(pageSize)))
	hasData := pageNumber < totalPages-1
	pageData := common.Page{
		PageSize:      pageSize,
		TotalPages:    int64(math.Ceil(float64(totalElement) / float64(pageSize))),
		PageNumber:    pageNumber,
		TotalElements: totalElement,
		HasData:       hasData,
	}
	c.JSON(http.StatusOK, common.SuccessPaginatedResponse(adminDataResponse, pageData, http.StatusOK))
}

func (k *KleosService) TriggerLeaderboardData(c *gin.Context) {

	userRoles := c.MustGet("roles").([]string)
	UserEmail := c.GetString("email")

	isAdmin := strings.Contains(strings.Join(userRoles, ","), "admin")

	if !isAdmin {
		err := errors.New("unauthorized")
		c.JSON(http.StatusUnauthorized, common.ErrorResponse(err, http.StatusUnauthorized, nil))
		return
	}

	k.logger.Info("Triggering Leaderboard data", zap.String("UserEmail", UserEmail))

	cron.SendLeaderBoardData(k.socketModeClient, k.logger, k.kleosDataService, k.userDataService)
	c.JSON(http.StatusOK, common.SuccessResponse("Leaderboard data triggered", http.StatusOK))
}

func (k *KleosService) TriggerWeeklyComms(c *gin.Context) {

	userRoles := c.MustGet("roles").([]string)
	isAdmin := strings.Contains(strings.Join(userRoles, ","), "admin")

	if !isAdmin {
		err := errors.New("unauthorized")
		c.JSON(http.StatusUnauthorized, common.ErrorResponse(err, http.StatusUnauthorized, nil))
		return
	}

	k.logger.Info("Triggering Weekly Comms data")
	go func() {
		cron.SendCommsInSlackWeekly(k.hNudgeSocketMode, k.logger, k.userDataService)
		return
	}()
	c.JSON(http.StatusOK, common.SuccessResponse("Weekly Comms data triggered", http.StatusOK))
}

func (k *KleosService) GetAdminDataCsv(c *gin.Context) {

	userRoles := c.MustGet("roles").([]string)
	UserEmail := c.GetString("email")

	isAdmin := strings.Contains(strings.Join(userRoles, ","), "admin")

	if !isAdmin {
		err := errors.New("unauthorized")
		c.JSON(http.StatusUnauthorized, common.ErrorResponse(err, http.StatusUnauthorized, nil))
		return
	}

	k.logger.Info("Received request to get admin CSV data", zap.String("UserEmail", UserEmail))

	userData, err := k.userDataService.GetUserDataFromGmail(UserEmail)

	if err != nil {
		k.logger.Error("Error while fetching user data", zap.Error(err))
		c.JSON(http.StatusBadGateway, common.ErrorResponse(err, http.StatusBadGateway, nil))
		return
	}

	slackChannel := userData.SlackUserId

	var adminDataResponse service.AdminAllUserData
	adminData, err := k.kleosDataService.GetAdminDataTest()
	if err != nil {
		k.logger.Error("Error while fetching admin data", zap.Error(err))
		c.JSON(http.StatusBadGateway, common.ErrorResponse(err, http.StatusBadGateway, nil))
		return
	}

	for _, aData := range *adminData {
		userInfo, err := k.userDataService.GetUserInfoFromId(aData.UserId)
		if err != nil {
			k.logger.Error("Error while fetching user data", zap.Error(err))
			c.JSON(http.StatusBadGateway, common.ErrorResponse(err, http.StatusBadGateway, nil))
			return
		}
		userData := service.UserData{
			Email:      userInfo.Email,
			UserName:   userInfo.RealName,
			ProfileUrl: userInfo.SlackImageUrl,
		}
		tempAdminData := service.UserAdminData{
			KleosGiven:    aData.KleosGiven,
			KleosReceived: aData.KleosReceived,
			User:          userData,
		}
		adminDataResponse.AdminAllUser = append(adminDataResponse.AdminAllUser, tempAdminData)
	}

	file, err := os.Create("output.csv")
	if err != nil {
		k.logger.Error("Error while creating csv file", zap.Error(err))
		c.JSON(http.StatusBadGateway, common.ErrorResponse(err, http.StatusBadGateway, nil))
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			return
		}
		err = os.Remove("output.csv")
		if err != nil {
			return
		}
	}(file)
	csvWriter := csv.NewWriter(file)

	header := []string{"email", "kleosGiven", "kleosReceived", "userName"}
	if err := csvWriter.Write(header); err != nil {
		k.logger.Error("Error while writing csv file", zap.Error(err))
		c.JSON(http.StatusBadGateway, common.ErrorResponse(err, http.StatusBadGateway, nil))
		return
	}

	for _, data := range adminDataResponse.AdminAllUser {
		row := []string{
			data.User.Email,
			strconv.Itoa(data.KleosGiven),
			strconv.Itoa(data.KleosReceived),
			data.User.UserName,
		}
		if err := csvWriter.Write(row); err != nil {
			k.logger.Error("Error while writing csv file", zap.Error(err))
			c.JSON(http.StatusBadGateway, common.ErrorResponse(err, http.StatusBadGateway, nil))
			return
		}
	}
	csvWriter.Flush()

	k.logger.Info("CSV file created and being sent to the user", zap.String("slackChannel", slackChannel))

	params := slack.FileUploadParameters{
		Title:    "kleosData.csv",
		File:     file.Name(),
		Channels: []string{slackChannel},
	}

	_, err = k.socketModeClient.UploadFile(params)

	if err != nil {
		k.logger.Error("Error while uploading file", zap.Error(err))
		c.JSON(http.StatusBadGateway, common.ErrorResponse(err, http.StatusBadGateway, nil))
		return
	}

	c.JSON(http.StatusOK, common.SuccessResponse("CSV Created", http.StatusOK))
}

func (k *KleosService) GetAdminXlsxFile(c *gin.Context) {

	userRoles := c.MustGet("roles").([]string)
	UserEmail := c.GetString("email")

	isAdmin := strings.Contains(strings.Join(userRoles, ","), "admin")

	if !isAdmin {
		err := errors.New("unauthorized")
		c.JSON(http.StatusUnauthorized, common.ErrorResponse(err, http.StatusUnauthorized, nil))
		return
	}

	k.logger.Info("Received request to get admin XLSX data", zap.String("UserEmail", UserEmail))

	userData, err := k.userDataService.GetUserDataFromGmail(UserEmail)
	if err != nil {
		k.logger.Error("Error while fetching user data", zap.Error(err))
		c.JSON(http.StatusBadGateway, common.ErrorResponse(err, http.StatusBadGateway, nil))
		return
	}

	var advAdminData service.AdvAdminAllUserData
	adminData, err := k.kleosDataService.GetAdminDataTest()
	if err != nil {
		k.logger.Error("Error while fetching admin data", zap.Error(err))
		c.JSON(http.StatusBadGateway, common.ErrorResponse(err, http.StatusBadGateway, nil))
		return
	}
	for _, aData := range *adminData {
		userInfo, err := k.userDataService.GetAdminAdvData(aData.UserId)
		if err != nil {
			k.logger.Error("Error while fetching user data", zap.Error(err))
			c.JSON(http.StatusBadGateway, common.ErrorResponse(err, http.StatusBadGateway, nil))
			return
		}

		hrbpInfo, err := k.hrbpService.GetHrbpData(userInfo.HrbpId)
		if err != nil {
			k.logger.Error("Error while fetching HRBP data", zap.Error(err))
			c.JSON(http.StatusBadGateway, common.ErrorResponse(err, http.StatusBadGateway, nil))
			return
		}

		departmentInfo, err := k.departmentService.GetDepartmentData(userInfo.DepartmentId)
		if err != nil {
			k.logger.Error("Error while fetching Department data", zap.Error(err))
			c.JSON(http.StatusBadGateway, common.ErrorResponse(err, http.StatusBadGateway, nil))
			return
		}

		managerInfo, err := k.managerService.GetManagerDataFromId(userInfo.ManagerId)
		if err != nil {
			k.logger.Error("Error while fetching Manager data", zap.Error(err))
			c.JSON(http.StatusBadGateway, common.ErrorResponse(err, http.StatusBadGateway, nil))
			return
		}

		hrbpData := service.HrbpData{
			Email: hrbpInfo.Email,
			Name:  hrbpInfo.Name,
		}

		managerData := service.ManagerData{
			Email: managerInfo.Email,
			Name:  managerInfo.Name,
		}

		departmentData := service.DepartmentData{
			Name: departmentInfo.Name,
		}

		userData := service.AdvUserData{
			Email:      userInfo.Email,
			UserName:   userInfo.Name,
			Hrbp:       hrbpData,
			Manager:    managerData,
			Department: departmentData,
			EmployeeId: userInfo.EmployeeId,
		}

		tempData := service.AdvUserAdminData{
			KleosGiven:    aData.KleosGiven,
			KleosReceived: aData.KleosReceived,
			User:          userData,
		}

		advAdminData.AdminAllUser = append(advAdminData.AdminAllUser, tempData)
	}

	f := excelize.NewFile()
	index, _ := f.NewSheet("Sheet1")
	f.SetActiveSheet(index)
	err = f.SetSheetName("Sheet1", "Consolidated_Data")
	if err != nil {
		k.logger.Error("Error while setting sheet name", zap.Error(err))
		c.JSON(http.StatusBadGateway, common.ErrorResponse(err, http.StatusBadGateway, nil))
		return
	}

	header := []string{"EmployeeId", "email", "userName", "kleosGiven", "kleosReceived", "hrbpEmail", "hrbpName", "managerEmail", "managerName", "departmentName"}
	err = f.SetSheetRow("Consolidated_Data", "A1", &header)
	if err != nil {
		k.logger.Error("Error while setting sheet row", zap.Error(err))
		c.JSON(http.StatusBadGateway, common.ErrorResponse(err, http.StatusBadGateway, nil))
		return
	}

	for i, data := range advAdminData.AdminAllUser {
		row := []string{
			data.User.EmployeeId,
			data.User.Email,
			data.User.UserName,
			strconv.Itoa(data.KleosGiven),
			strconv.Itoa(data.KleosReceived),
			data.User.Hrbp.Email,
			data.User.Hrbp.Name,
			data.User.Manager.Email,
			data.User.Manager.Name,
			data.User.Department.Name,
		}

		startCell, err := excelize.JoinCellName("A", i+2)
		if err != nil {
			k.logger.Error("Error while joining cell name", zap.Error(err))
			c.JSON(http.StatusBadGateway, common.ErrorResponse(err, http.StatusBadGateway, nil))
			return
		}
		if err = f.SetSheetRow("Consolidated_Data", startCell, &row); err != nil {
			k.logger.Error("Error while setting sheet row", zap.Error(err))
			c.JSON(http.StatusBadGateway, common.ErrorResponse(err, http.StatusBadGateway, nil))
			return
		}
	}

	allKleosData, err := k.kleosDataService.GetAllKleosData()
	if err != nil {
		k.logger.Error("Error while fetching all kleos data", zap.Error(err))
		c.JSON(http.StatusBadGateway, common.ErrorResponse(err, http.StatusBadGateway, nil))
		return

	}

	index2, _ := f.NewSheet("Sheet2")
	f.SetActiveSheet(index2)
	err = f.SetSheetName("Sheet2", "All_Data")
	if err != nil {
		k.logger.Error("Error while setting sheet name", zap.Error(err))
		c.JSON(http.StatusBadGateway, common.ErrorResponse(err, http.StatusBadGateway, nil))
		return
	}

	headerTwo := []string{"Sender_Email", "Sender_Name", "Receiver_Email", "Receiver_Name", "Achievement_Name", "Message", "Date", "Month"}

	err = f.SetSheetRow("All_Data", "A1", &headerTwo)

	for i, data := range *allKleosData {
		senderInfo, err := k.userDataService.GetUserInfoFromId(data.SenderID)
		if err != nil {
			k.logger.Error("Error while fetching user data", zap.Error(err))
			c.JSON(http.StatusBadGateway, common.ErrorResponse(err, http.StatusBadGateway, nil))
			return
		}
		receiverInfo, err := k.userDataService.GetUserInfoFromId(data.ReceiverID)
		if err != nil {
			k.logger.Error("Error while fetching user data", zap.Error(err))
			c.JSON(http.StatusBadGateway, common.ErrorResponse(err, http.StatusBadGateway, nil))
			return
		}

		achievementData, err := k.achievementService.GetAllAchievementName()
		if err != nil {
			k.logger.Error("Error while fetching achievement data", zap.Error(err))
			c.JSON(http.StatusBadGateway, common.ErrorResponse(err, http.StatusBadGateway, nil))
			return
		}

		achievementName := ""

		for _, achData := range achievementData {

			achievementId, _ := strconv.Atoi(data.Achievement)
			if int(achData.ID) == achievementId {
				achievementName = achData.DisplayName
			}
		}

		Date := data.Day + "-" + data.Month + "-" + data.Year

		row := []string{
			senderInfo.Email,
			senderInfo.RealName,
			receiverInfo.Email,
			receiverInfo.RealName,
			achievementName,
			data.Message,
			Date,
			data.Month,
		}

		startCell, err := excelize.JoinCellName("A", i+2)
		if err != nil {
			k.logger.Error("Error while joining cell name", zap.Error(err))
			c.JSON(http.StatusBadGateway, common.ErrorResponse(err, http.StatusBadGateway, nil))
			return
		}
		if err = f.SetSheetRow("All_Data", startCell, &row); err != nil {
			k.logger.Error("Error while setting sheet row", zap.Error(err))
			c.JSON(http.StatusBadGateway, common.ErrorResponse(err, http.StatusBadGateway, nil))
			return
		}
	}

	f.SetActiveSheet(index)

	err = f.SaveAs("output.xlsx")
	if err != nil {
		k.logger.Error("Error while saving file", zap.Error(err))
		c.JSON(http.StatusBadGateway, common.ErrorResponse(err, http.StatusBadGateway, nil))
		return
	}

	slackChannelToPost := userData.SlackUserId

	params := slack.FileUploadParameters{
		Title:    "kleosData.xlsx",
		File:     "output.xlsx",
		Channels: []string{slackChannelToPost},
	}

	_, err = k.socketModeClient.UploadFile(params)
	if err != nil {
		k.logger.Error("Error while uploading file", zap.Error(err))
		c.JSON(http.StatusBadGateway, common.ErrorResponse(err, http.StatusBadGateway, nil))
		return
	}

	err = os.Remove("output.xlsx")
	if err != nil {
		k.logger.Error("Error while deleting file", zap.Error(err))
		c.JSON(http.StatusBadGateway, common.ErrorResponse(err, http.StatusBadGateway, nil))
		return
	}

	c.JSON(http.StatusOK, common.SuccessResponse("CSV Created", http.StatusOK))
}

func (k *KleosService) CreateUser(c *gin.Context) {

	var user request.CreateUserRequest
	err := c.ShouldBindJSON(&user)

	if err != nil {
		k.logger.Error("Error while binding request", zap.Error(err))
		c.JSON(http.StatusInternalServerError, common.ErrorResponse(err,
			http.StatusInternalServerError, nil))
		return
	}

	_, err = k.createUserEntityAndSaveToDB(user)

	if err != nil {
		k.logger.Error("Error while creating user", zap.Error(err))
		c.JSON(http.StatusBadGateway, common.ErrorResponse(err, http.StatusBadGateway, nil))
		return
	}

	c.JSON(http.StatusOK, common.SuccessResponse("User Created Successfully",
		http.StatusOK))
}

func (k *KleosService) CreateUserBulk(c *gin.Context) {

	mandatoryColumns := strings.Split(viper.GetString("csv.mandatory.columns"), ",")

	fmt.Println("Mandatory Columns: ", mandatoryColumns)

	file, err := c.FormFile("file")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error reading file: " + err.Error()})
		return
	}

	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error opening file: " + err.Error()})
		return
	}
	defer func(f multipart.File) {
		err := f.Close()
		if err != nil {

		}
	}(f)
	reader := csv.NewReader(f)

	headers, err := reader.Read()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error reading CSV: " + err.Error()})
		return
	}

	if !utils.ValidateHeaders(headers, mandatoryColumns) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Mandatory Columns are missing or not in correct order", "Mandatory Columns: ": mandatoryColumns})
		return
	}

	var successCount int
	var failureCount int
	var failedRows []map[string]string

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error reading CSV: " + err.Error()})
			return
		}

		newUser := createUserRequestFromRecord(record, headers, k.slackbotClient.SocketModeClient)

		existingUser, _ := k.userDataService.GetUserDataWithForeignJoins(newUser.Email)

		if existingUser != nil {
			err := k.validateAndUpdateExistingUser(existingUser, newUser)

			if err != nil {
				failedRows = append(failedRows, map[string]string{"row": strings.Join(record, ","), "reason": err.Error()})
				failureCount++
				continue
			}

			successCount++
			continue
		}

		newUser, err = k.addSlackUserIdAndImageUrl(newUser, k.slackbotClient.SocketModeClient)

		if err != nil {
			failedRows = append(failedRows, map[string]string{"row": strings.Join(record, ","), "reason": err.Error()})
			failureCount++
			continue
		}

		_, err = k.createUserEntityAndSaveToDB(newUser)

		if err != nil {
			failedRows = append(failedRows, map[string]string{"row": strings.Join(record, ","), "reason": err.Error()})
			failureCount++
			continue
		}

		successCount++
	}

	response := &service.CreateUserBulkResponse{
		SuccessCount: successCount,
		FailureCount: failureCount,
		FailedRows:   failedRows,
	}

	if failureCount > 0 {
		c.JSON(http.StatusMultiStatus, response)
	} else {
		c.JSON(http.StatusOK, response)
	}
}

func (k *KleosService) validateAndUpdateExistingUser(existingUser *usersDb.UserWithForeignData, newUser request.CreateUserRequest) error {

	managerData, departmentData, hrbpData, err := k.fetchData(existingUser)

	if err != nil {
		return fmt.Errorf("error fetching user data: %w", err)
	}

	if hrbpData.Name == newUser.HrbpName &&
		managerData.Email == newUser.ManagerEmail &&
		departmentData.Name == newUser.DepartmentName {
		k.logger.Info("User already exists", zap.String("email", newUser.Email))
		return nil
	}

	return k.updateUserDataForExistingUser(existingUser, newUser)
}

func (k *KleosService) fetchData(existingUser *usersDb.UserWithForeignData) (*managersDb.ManagerEntity, *departmentDb.DepartmentEntity, *hrbpDb.HrbpEntity, error) {
	managerData, err := k.managerService.GetManagerDataFromId(existingUser.ManagerId)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error fetching manager data: %w", err)
	}

	departmentData, err := k.departmentService.GetDepartmentDataById(existingUser.DepartmentId)
	if err != nil {
		return managerData, nil, nil, fmt.Errorf("error fetching department data: %w", err)
	}

	hrbpData, err := k.hrbpService.GetHrbpDataById(existingUser.HrbpId)
	if err != nil {
		return managerData, departmentData, nil, fmt.Errorf("error fetching hrbp data: %w", err)
	}

	return managerData, departmentData, hrbpData, nil
}

func (k *KleosService) updateUserDataForExistingUser(existingUser *usersDb.UserWithForeignData, newUser request.CreateUserRequest) error {
	savedManagerEntity, err := k.createOrGetManager(newUser)
	if err != nil {
		return err
	}

	savedHRBPEntity, err := k.createOrGetHRBP(newUser)
	if err != nil {
		return err
	}

	savedDepartmentEntity, err := k.createOrGetDepartment(newUser)
	if err != nil {
		return err
	}

	updatedUser := existingUser.UserEntity

	updatedUser.HrbpId = int(savedHRBPEntity.ID)
	updatedUser.DepartmentId = int(savedDepartmentEntity.ID)
	updatedUser.ManagerId = int(savedManagerEntity.ID)

	_, err = k.userDataService.UpdateUser(&updatedUser)
	if err != nil {
		return err
	}

	return nil
}

func (k *KleosService) addSlackUserIdAndImageUrl(user request.CreateUserRequest, api *socketmode.Client) (request.CreateUserRequest, error) {
	info, err := api.GetUserByEmail(user.Email)
	if err != nil {
		return request.CreateUserRequest{}, err
	}

	user.SlackUserId = info.ID
	user.SlackImageUrl = info.Profile.Image192

	return user, nil
}

func (k *KleosService) createUserEntityAndSaveToDB(userCreateRequest request.CreateUserRequest) (*usersDb.UserEntity, error) {

	savedManagerEntity, err := k.createOrGetManager(userCreateRequest)
	savedHRBPEntity, err := k.createOrGetHRBP(userCreateRequest)
	savedDepartmentEntity, err := k.createOrGetDepartment(userCreateRequest)

	userEntity := &usersDb.UserEntity{
		Name:          utils.ExtractUsername(userCreateRequest.Email),
		Email:         userCreateRequest.Email,
		RealName:      userCreateRequest.RealName,
		SlackUserId:   userCreateRequest.SlackUserId,
		SlackImageUrl: userCreateRequest.SlackImageUrl,
		EmployeeId:    userCreateRequest.EmployeeId,
		HrbpId:        int(savedHRBPEntity.ID),
		DepartmentId:  int(savedDepartmentEntity.ID),
		ManagerId:     int(savedManagerEntity.ID),
		GivenCount:    utils.DefaultGivenValue,
		ReceivedCount: utils.DefaultReceivedValue,
	}

	newUserEntity, err := k.userDataService.CreateUser(userEntity)

	if err != nil {
		return nil, err
	}
	return newUserEntity, nil
}

func createUserRequestFromRecord(record []string, headers []string, api *socketmode.Client) request.CreateUserRequest {

	var user request.CreateUserRequest
	headerMap := map[string]*string{
		"Employee Id":             &user.EmployeeId,
		"Full Name":               &user.RealName,
		"Official Email Id":       &user.Email,
		"Direct Manager Name":     &user.ManagerName,
		"Direct Manager Email Id": &user.ManagerEmail,
		"Current Department":      &user.DepartmentName,
		"Hrbp Name":               &user.HrbpName,
	}

	for i, header := range headers {
		if fieldPointer, exists := headerMap[header]; exists {
			*fieldPointer = record[i]
		}
	}

	return user
}

func (k *KleosService) createOrGetManager(userCreateRequest request.CreateUserRequest) (*managersDb.ManagerEntity, error) {

	existingManager, err := k.managerService.GetManagerDataFromEmail(userCreateRequest.ManagerEmail)

	if err != nil {
		return nil, fmt.Errorf("error checking for existing manager: %w", err)
	}

	if existingManager == nil {
		managerEntity := &managersDb.ManagerEntity{
			Email:    userCreateRequest.ManagerEmail,
			Name:     utils.ExtractUsername(userCreateRequest.ManagerEmail),
			RealName: userCreateRequest.ManagerName,
		}

		savedManagerEntity, err := k.managerService.CreateManager(managerEntity)
		if err != nil {
			return nil, fmt.Errorf("error creating manager: %w", err)
		}

		return savedManagerEntity, nil
	}

	return existingManager, nil
}

func (k *KleosService) createOrGetHRBP(userCreateRequest request.CreateUserRequest) (*hrbpDb.HrbpEntity, error) {

	existingHrbp, err := k.hrbpService.GetHrbpDataFromName(userCreateRequest.HrbpName)

	if err != nil {
		return nil, fmt.Errorf("error checking for existing manager: %w", err)
	}

	if existingHrbp == nil {
		savedHrbp, err := k.hrbpService.CreateHrbp(userCreateRequest.HrbpName, "")
		if err != nil {
			return nil, fmt.Errorf("error creating new HRBP: %w", err)
		}
		return savedHrbp, nil
	}

	return existingHrbp, nil
}

func (k *KleosService) createOrGetDepartment(userCreateRequest request.CreateUserRequest) (*departmentDb.DepartmentEntity, error) {

	existingDepartment, err := k.departmentService.GetDepartmentDataFromName(userCreateRequest.DepartmentName)

	if err != nil {
		return nil, fmt.Errorf("error checking for existing manager: %w", err)
	}

	if existingDepartment == nil {
		newDepartment, err := k.departmentService.CreateDepartment(userCreateRequest.DepartmentName)
		if err != nil {
			return nil, fmt.Errorf("error creating new department: %w", err)
		}
		return newDepartment, nil
	}

	return existingDepartment, nil

}

//func (k *KleosService) postSlackBotMessageWeb(senderId string, receiverId string, level string, title string) error {
//
//	senderBlock := slack.MsgOptionCompose(
//		slack.MsgOptionText(
//			fmt.Sprintf("Yay <@%s>, Kudos has been sent to <@%s>", senderId, receiverId),
//			false),
//		slack.MsgOptionPostMessageParameters(slack.PostMessageParameters{
//			AsUser: true,
//		}),
//	)
//
//	receiverBlock := view.ReceiverSlackBotSection(receiverId, senderId, level, title)
//
//	isSuccess := k.postBotForSenderAndReceiver(senderId, receiverId, senderBlock, receiverBlock)
//
//	if isSuccess {
//		k.logger.Info("[Web] Successfully posted bot messages", zap.String("senderId", senderId), zap.String("senderId", receiverId), zap.Any("receiverId", time.Now()))
//	}
//
//	return nil
//}
//
//func (k *KleosService) postBotForSenderAndReceiver(senderId string, receiverId string, senderBlock slack.MsgOption, receiverBlock slack.Attachment) bool {
//
//	isSuccess := true
//
//	//Sending message to the sender
//	_, _, err := k.socketModeClient.PostMessage(senderId, senderBlock)
//	if err != nil {
//		k.logger.Error("[Web] [Kleos Bot] Error will posting message to bot for sender", zap.Error(err))
//		k.logger.Info("Posting message to bot failed, continuing with the channel & email communication")
//		isSuccess = false
//	}
//
//	//Sending message to the receiver
//	_, _, err = k.socketModeClient.PostMessage(receiverId, slack.MsgOptionAttachments(receiverBlock))
//	if err != nil {
//		k.logger.Error("[Web] [Kleos Bot] Error will posting message to bot receiver", zap.Error(err))
//		k.logger.Info("Posting message to bot failed, continuing with the channel & email communication")
//		isSuccess = false
//	}
//
//	return isSuccess
//}
