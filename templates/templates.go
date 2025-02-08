package templates

import (
	"bytes"
	"embed"
	"go.uber.org/zap"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	tmplhtml "html/template"
	"regexp"
	"strings"
	tmpltext "text/template"
	"time"
)

type FuncMap map[string]interface{}

var tmpl *tmpltext.Template

//go:embed default.tmpl
var defaultTemplate embed.FS

const SpecialChar string = `\\` + "`*_{[]}()#+-.!"

//template func
var DefaultFuncs = tmpltext.FuncMap{
	"toUpper": strings.ToUpper,
	"toLower": strings.ToLower,
	"title": func(text string) string {
		// Casers should not be shared between goroutines, instead
		// create a new caser each time this function is called.
		return cases.Title(language.AmericanEnglish).String(text)
	},
	"trimSpace": strings.TrimSpace,
	// join is equal to strings.Join but inverts the argument order
	// for easier pipelining in templates.
	"join": func(sep string, s []string) string {
		return strings.Join(s, sep)
	},
	"match": regexp.MatchString,
	"safeHtml": func(text string) tmplhtml.HTML {
		return tmplhtml.HTML(text)
	},
	"reReplaceAll": func(pattern, repl, text string) string {
		re := regexp.MustCompile(pattern)
		return re.ReplaceAllString(text, repl)
	},
	"stringSlice": func(s ...string) []string {
		return s
	},
	// date returns the text representation of the time in the specified format.
	"date": func(fmt string, t time.Time) string {
		return t.Format(fmt)
	},
	// tz returns the time in the timezone.
	"tz": func(name string, t time.Time) (time.Time, error) {
		loc, err := time.LoadLocation(name)
		if err != nil {
			return time.Time{}, err
		}
		return t.In(loc), nil
	},
	"since": time.Since,
	// escape markdown
	"escapeMarkdownV2": func(text string) string {
		for _, char := range SpecialChar {
			text = strings.ReplaceAll(text, string(char), "\\"+string(char))
		}
		return text
	},
}

// read template
func Load(templatePath string) {
	var err error

	if templatePath == "" {
		d, _ := defaultTemplate.ReadFile("default.tmpl")
		templatePath = "default.tmpl"
		tmpl, err = tmpltext.New(templatePath).Funcs(DefaultFuncs).Parse(string(d))
	} else {
		tmpl, err = tmpltext.New(templatePath).Funcs(DefaultFuncs).ParseFiles(templatePath)
	}

	if err != nil {
		zap.L().Sugar().Fatalf("Problem reading parsing template file: %v", err)
	} else {
		zap.L().Sugar().Infof("Load template file:%s", templatePath)
	}
}

//render template
func Render(name string, data interface{}) (string, error) {
	var result bytes.Buffer
	err := tmpl.ExecuteTemplate(&result, name, data)
	if err != nil {
		zap.L().Sugar().Errorf("Error rendering template: %s", err.Error())
	}
	return result.String(), err
}
