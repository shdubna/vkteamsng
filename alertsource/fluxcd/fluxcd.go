package fluxcd

import (
	"encoding/json"
	"github.com/shdubna/vkteamsng/templates"
	"go.uber.org/zap"
	"io"
	"net/http"
)

//https://fluxcd.io/flux/components/notification/providers/#generic-webhook

// Alert data
type InvolvedObject struct {
	ApiVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Name       string `json:"name"`
	Namespace  string `json:"namespace"`
	Uid        string `json:"uid"`
}

// Message data struct by Flux  notification controller
type Message struct {
	Severity            string         `json:"severity"`
	Reason              string         `json:"reason"`
	Message             string         `json:"message"`
	ReportingController string         `json:"reportingController"`
	Metadata            interface{}    `json:"metadata"`
	InvolvedObject      InvolvedObject `json:"involvedObject"`
	ReportingInstance   string         `json:"reportingInstance"`
	Timestamp           string         `json:"timestamp"`
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

	return templates.Render("fluxcd", message)
}
