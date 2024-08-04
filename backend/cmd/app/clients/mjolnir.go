package clients

import (
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"io"
	clients "kleos/cmd/app/clients/models"
	"net/http"
)

type MjolnirClient struct {
	HttpClient *http.Client
	baseUrl    string
	realmId    string
	logger     *zap.Logger
}

func NewMjolnirClient(httpClient *http.Client, baseUrl, realmId string, logger *zap.Logger) *MjolnirClient {
	return &MjolnirClient{
		HttpClient: httpClient,
		baseUrl:    baseUrl,
		realmId:    realmId,
		logger:     logger,
	}
}

const (
	url = "%s/session/%s"
)

func (m *MjolnirClient) GetSessionResponse(sessionToken string) (*clients.MjolnirSessionResponse, error) {
	if sessionToken == "null" {
		return nil, errors.New("unauthorized request")
	}
	client := m.HttpClient
	req, _ := http.NewRequest("GET", fmt.Sprintf(url, m.baseUrl, m.realmId), nil)
	req.Header.Add("X-Session-Token", sessionToken)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response clients.MjolnirSessionResponse
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return nil, err
	}

	response.StatusCode = resp.StatusCode

	m.logger.Info(fmt.Sprintf("%v", response))

	return &response, nil
}
