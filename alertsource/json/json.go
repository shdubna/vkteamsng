package json

import (
	"bytes"
	"encoding/json"
	"github.com/shdubna/vkteamsng/templates"
	"go.uber.org/zap"
	"io"
	"net/http"
)

// Message represent data from any JSON
type Message map[string]interface{}

// Parse implement Payload.Parse()
func (m Message) Parse(req *http.Request, logger *zap.Logger) (string, error) {
	messageBytes, err := io.ReadAll(req.Body)
	if err != nil {
		return "", err
	}
	logger.Debug(string(messageBytes))
	var buffer bytes.Buffer
	err = json.Indent(&buffer, messageBytes, "", "  ")

	if err != nil {
		return "", err
	}

	return templates.Render("json", buffer.String())
}
