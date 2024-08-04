package clients

import (
	"github.com/spf13/viper"
	"net/http"
	"time"
)

type HttpClient struct {
	HttpClient *http.Client
}

func NewHttpClient() *HttpClient {
	return &HttpClient{
		HttpClient: &http.Client{
			Transport: &http.Transport{
				MaxIdleConns:    viper.GetInt("http.max.idle.connection.pool"),
				MaxConnsPerHost: viper.GetInt("http.max.connection"),
			},
			Timeout: time.Duration(viper.GetInt("http.max.timeout.seconds")) * time.Second,
		},
	}
}
