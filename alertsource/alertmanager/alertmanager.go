package alertmanager

import (
	"encoding/json"
	"github.com/shdubna/vkteamsng/templates"
	"go.uber.org/zap"
	"io"
	"net/http"
)

//https://prometheus.io/docs/alerting/configuration/#webhook_config

// Alert data
type Alert struct {
	Status       string      `json:"status"`
	Labels       interface{} `json:"labels"`
	Annotations  interface{} `json:"annotations"`
	StartsAt     string      `json:"startsAt"`
	EndsAt       string      `json:"endsAt"`
	GeneratorURL string      `json:"generatorURL"`
}

// Message data struct by Prometheus Alertmanager
type Message struct {
	Version      string      `json:"version"`
	GroupKey     string      `json:"groupKey"`
	Status       string      `json:"status"`
	Receiver     string      `json:"receiver"`
	GroupLabels  interface{} `json:"groupLabels"`
	CommonLabels interface{} `json:"commonLabels"`
	ExternalURL  string      `json:"externalURL"`
	Alerts       []Alert     `json:"alerts"`
}

// Parse implement Payload.Parse()
func (gm Message) Parse(req *http.Request, logger *zap.Logger) (string, error) {
	messageBytes, err := io.ReadAll(req.Body)
	if err != nil {
		return "", err
	}
	logger.Debug(string(messageBytes))
	message := Message{}
	err = json.Unmarshal(messageBytes, &message)
	if err != nil {
		return "", err
	}

	return templates.Render("alertmanager", message)
}
