package main

import (
	"flag"
	"fmt"
	"os"
	"github.com/mail-ru-im/bot-golang"
	"github.com/shdubna/vkteamsng/templates"
	"github.com/shdubna/vkteamsng/webhook"
	"go.uber.org/zap"
)

var gitTag string

var (
	listenAddress = flag.String("listen_address", ":8080", "Address to listen proxy requests.")
	apiUrl        = flag.String("vkteams_url", "https://myteam.mail.ru/bot/v1", "VKTeams api url.")
	parseMode     = flag.String("parse_mode", "MarkdownV2", "Bot parse mode/. Allowed values: MarkdownV2, HTML")
	Debug         = flag.Bool("debug", false, "Enable debug logging.")
	templatePath  = flag.String("template_path", "", "Path to message template file, if not specified use embeded")
	version       = flag.Bool("version", false, "Show version number and quit.")
)

func init() {
	flag.Parse()
	config := zap.NewProductionConfig()
	if *Debug {
		config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	}
	zap.ReplaceGlobals(zap.Must(config.Build()))
}

func main() {
	flag.Parse()

	if *version {
		fmt.Println(gitTag)
		os.Exit(0)
	}
	logger := zap.L()

	templates.Load(*templatePath)
	webhook.SetParseMode(*parseMode)

	bot, err := botgolang.NewBot(os.Getenv("BOT_TOKEN"), botgolang.BotDebug(*Debug), botgolang.BotApiURL(*apiUrl))
	if err != nil {
		logger.Fatal("Wrong VK Teams token: ", zap.Error(err))
	}
	logger.Sugar().Infof("Usig bot %s for sending alerts", bot.Info.Nick)

	webProvider := webhook.Provider{Bot: bot, Logger: logger}
	logger.Fatal(webProvider.Start().Error())
}
