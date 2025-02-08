package raw

import (
	"github.com/shdubna/vkteamsng/templates"
	"go.uber.org/zap"
	"io"
	"net/http"
)

// Message represent data from any data
type Message []byte

func transformMessage(data io.ReadCloser, logger *zap.Logger) (string, error) {
	messageBytes, err := io.ReadAll(data)
	if err != nil {
		return "", err
	}
	logger.Debug(string(messageBytes))
	return string(messageBytes), nil
}

// Parse implement Payload.Parse()
func (m Message) Parse(req *http.Request, logger *zap.Logger) (string, error) {
	messageBytes, err := io.ReadAll(req.Body)
	if err != nil {
		return "", err
	}
	logger.Debug(string(messageBytes))
	return templates.Render("raw", string(messageBytes))
}
