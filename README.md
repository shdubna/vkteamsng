# vkteamsng
[![Go Report Card](https://goreportcard.com/badge/github.com/shdubna/vkteamsng)](https://goreportcard.com/report/github.com/shdubna/vkteamsng)
[![GitHub CodeQL](https://github.com/shdubna/vkteamsng/workflows/CodeQL/badge.svg)](https://github.com/shdubna/vkteamsng/actions?query=workflow%3CodeQL)
[![GitHub Release](https://github.com/shdubna/vkteamsng/workflows/Release/badge.svg)](https://github.com/shdubna/vkteamsng/actions?query=workflow%3ARelease)
[![GitHub license](https://img.shields.io/github/license/shdubna/vkteamsng.svg)](https://github.com/shdubna/vkteamsng/blob/main/LICENSE)
[![GitHub tag](https://img.shields.io/github/v/tag/shdubna/vkteamsng?label=latest)](https://github.com/shdubna/vkteamsng/releases)

Vk Teams Notification Gateway - шлюз для отправки нотификаций в [VK Teams](https://biz.mail.ru/messenger/) 
с помощью [bot api](https://teams.vk.com/botapi/).

## Поддерживаемые коннекторы

- [Prometheus Alertmanager Webhook](https://prometheus.io/docs/alerting/latest/configuration/#webhook_config)
- [Flux CD Notification Controller Generic webhook](https://fluxcd.io/flux/components/notification/providers/#generic-webhook)
- json
- raw

## Пример настройки оповещений alertmanager
- получить токен бота, добавить бота в чат/канал, получит id чата по [инструкции](https://teams.vk.com/botapi/tutorial);
- Скачать vkteamsng со [страницы релизов](https://github.com/shdubna/vktemasng/releases) или [contaner image](https://github.com/shdubna/vktemasng/pkgs/container/vktemasng)
- запустить vkteamsng

в контенере:
```bash
docker run -d --name vkteamsng -e BOT_TOKEN=<токен бота vkteams> -p 8080:8080 ghcr.io/shdubna/vkteamsng
```
или из бинарного файлаЖ
```bash
export BOT_TOKEN=<токен бота vkteams>
vkteamsng 
```
- настроить webhook в alertmanager:
```
...
receivers:
- name: vkteams
  webhook_configs:
  - url: http://<vkteams ip>:8080/webhook/alertmanager/<vkteams chat id>
...
```

## Поддерживаемые опции

```bash
Usage of vkteamsng:
  -debug
        Enable debug logging.
  -listen_address string
        Address to listen proxy requests. (default ":8080")
  -parse_mode string
        Bot parse mode/. Allowed values: MarkdownV2, HTML (default "MarkdownV2")
  -template_path string
        Path to message template file, if not specified use embeded
  -version
        Show version number and quit.
  -vkteams_url string
        VKTeams api url. (default "https://myteam.mail.ru/bot/v1")
```

Токен авторизации пользователя bot задается переменной окружения `BOT_TOKEN`.

## Кастомные шаблоны

Для подключения кастомных шаблонов сообщений требуется создать файл с шаблонами go template и указать путь до файла флагом `-template_path`.

Шаблоны по умолчанию находятся в файле [default.tmpl](./templates/default.tmpl).

Файл шаблона должен содержать определение шаблонов для всех поддерживаемых интеграций.