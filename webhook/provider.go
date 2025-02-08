package webhook

import (
	"errors"
	"flag"
	"net/http"
	"time"
	"github.com/gorilla/mux"
	"github.com/mail-ru-im/bot-golang"
	"github.com/shdubna/vkteamsng/alertsource/alertmanager"
	"github.com/shdubna/vkteamsng/alertsource/fluxcd"
	"github.com/shdubna/vkteamsng/alertsource/json"
	"github.com/shdubna/vkteamsng/alertsource/raw"
	"go.uber.org/zap"
)

var payloadSourceMap = map[string]Payload{
	"json":         json.Message{},
	"raw":          raw.Message{},
	"alertmanager": alertmanager.Message{},
	"fluxcd":       fluxcd.Message{},
}

// Payload interface for any data from any alert systems
type Payload interface {
	Parse(r *http.Request, logger *zap.Logger) (string, error)
}

// Provider represent single instances of bot and echo
type Provider struct {
	Bot    *botgolang.Bot
	Logger *zap.Logger
}

func (Provider) payloadBySourceName(sourceName string) (Payload, error) {
	payload, ok := payloadSourceMap[sourceName]
	if !ok {
		return nil, errors.New("unknown alert source")
	}
	return payload, nil
}

// Start prepare routes and serve them
func (p *Provider) Start() error {
	router := mux.NewRouter()
	router.HandleFunc("/webhook/{source}/{target}", p.handleMessage)
	router.HandleFunc("/health",
		func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "", http.StatusOK)
		},
	)

	server := &http.Server{
		Addr:         flag.Lookup("listen_address").Value.(flag.Getter).Get().(string),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  15 * time.Second,
	}
	server.SetKeepAlivesEnabled(true)

	return http.ListenAndServe(server.Addr, router)
}
