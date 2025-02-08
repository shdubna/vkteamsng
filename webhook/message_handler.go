package webhook

import (
	"github.com/gorilla/mux"
	"github.com/mail-ru-im/bot-golang"
	"go.uber.org/zap"
	"net/http"
)

var parseMode botgolang.ParseMode

func SetParseMode(mode string) {
	parseMode = botgolang.ParseMode(mode)
}

func (p *Provider) handleMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	payload, err := p.payloadBySourceName(vars["source"])
	if err != nil {
		p.Logger.Error(err.Error(), zap.String("url %s", r.RequestURI))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	messageString, err := payload.Parse(r, p.Logger)
	if err != nil {
		p.Logger.Error(err.Error(), zap.String("url", r.RequestURI))
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	message := p.Bot.NewTextMessage(vars["target"], messageString)
	message.ParseMode = parseMode
	err = message.Send()
	if err != nil {
		p.Logger.Error(err.Error(), zap.String("url", r.RequestURI))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

}
