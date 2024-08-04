package comms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strings"
	"time"
)

const INITIATED = "INITIATED"

func TriggerComms(logger *zap.Logger, eEmail string, mEmail string, name string, message string, slackImage string, givenBy string, giverEmail string) {

	fromEmail := viper.GetString("bot.sender.email")
	dashboardUrl := viper.GetString("dashboard.url")
	commsBaseUrl := viper.GetString("comms.base.url")
	templateId := viper.GetString("email.template.id")
	currentTime := time.Now()
	formattedDate := currentTime.Format("02 Jan 2006")

	sanitizedMessage := sanitizedMessage(message)

	apiURL := commsBaseUrl

	data := fmt.Sprintf(`{
    	"templateId": "%s",
    	"data": {
			"from": "%s",
        	"to": "%s",
        	"ccList": ["%s", "%s"],
        	"date": "%s",
        	"name": "%s",
        	"message": "%s",
        	"user_name": "%s",
			"slack_image": "%s",
			"url": "%s"
    	}
	}`, templateId, fromEmail, eEmail, mEmail, giverEmail, formattedDate, name, sanitizedMessage, givenBy, slackImage, dashboardUrl)

	jsonPayload := []byte(data)

	newRequest, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return
	}
	newRequest.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(newRequest)
	if err != nil {
		logger.Error("Error while sending email comms", zap.Error(err))
		return
	}

	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		logger.Error("[Communication] Error while reading response body", zap.Error(err))
	}

	var responseJson CommunicationResponse

	err = json.Unmarshal(responseBody, &responseJson)
	if err != nil {
		logger.Error("[Communication] Error while unmarshalling response body", zap.Error(err))
	}

	//logger.Info("Email comms sent successfully", zap.String("payload", data),
	//	zap.Any("response", responseJson))

	if responseJson.Status != INITIATED {
		logger.Error("[Communication] Error while sending email communication to user",
			zap.String("templateId", templateId),
			zap.String("fromEmail", fromEmail),
			zap.String("eEmail", eEmail),
			zap.String("mEmail", mEmail),
			zap.String("giverEmail", giverEmail),
			zap.Any("response", responseJson))

		return
	}

	logger.Info("[Communication] Email communication sent successfully", zap.String("sender", eEmail), zap.String("receiver", giverEmail), zap.String("sentAt", responseJson.SentAt))

}

func sanitizedMessage(message string) string {

	sanitizedMessage := strings.ReplaceAll(message, "\n", "<br>")

	return sanitizedMessage
}
