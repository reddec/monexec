package plugins

import (
	"html/template"
	"log"
	"github.com/Masterminds/sprig"
	"io/ioutil"
	"bytes"
	"path/filepath"
	"github.com/pkg/errors"
	"time"
	"os"
)

type withTemplate struct {
	Template     string `yaml:"template"`
	TemplateFile string `yaml:"templateFile"` // template file (relative to config dir) has priority. Template supports basic utils
}

func (wt *withTemplate) renderDefault(action, id, label string, err error, logger *log.Logger) (string, error) {
	s, _, err := wt.renderDefaultParams(action, id, label, err, logger)
	return s, err
}

func (wt *withTemplate) renderDefaultParams(action, id, label string, err error, logger *log.Logger) (string, map[string]interface{}, error) {
	hostname, _ := os.Hostname()
	params := map[string]interface{}{
		"id":       id,
		"label":    label,
		"error":    err,
		"action":   action,
		"hostname": hostname,
		"time":     time.Now().String(),
	}
	s, err := wt.render(params, logger)
	return s, params, err
}

func (wt *withTemplate) render(params map[string]interface{}, logger *log.Logger) (string, error) {
	parser, err := parseFileOrTemplate(wt.TemplateFile, wt.Template, logger)
	if err != nil {
		return "", errors.Wrap(err, "parse template")
	}
	message := &bytes.Buffer{}

	renderErr := parser.Execute(message, params)
	if renderErr != nil {
		logger.Println("failed render:", renderErr, "; params:", params)
		return "", err
	}
	return message.String(), nil
}

func (wt *withTemplate) resolvePath(workDir string) {
	wt.TemplateFile = realPath(wt.TemplateFile, workDir)
}

func (wt *withTemplate) MergeFrom(other *withTemplate) error {
	if wt.TemplateFile == "" {
		wt.TemplateFile = other.TemplateFile
	}
	if wt.TemplateFile != other.TemplateFile {
		return errors.New("template files are different")
	}
	if wt.Template == "" {
		wt.Template = other.Template
	}
	if wt.Template != other.Template {
		return errors.New("different templates")
	}
	return nil
}

func unique(names []string) []string {
	var hash = make(map[string]struct{})
	for _, name := range names {
		hash[name] = struct{}{}
	}
	var ans = make([]string, 0, len(hash))
	for name := range hash {
		ans = append(ans, name)
	}
	return ans
}

func makeSet(names []string) map[string]bool {
	var hash = make(map[string]bool)
	for _, name := range names {
		hash[name] = true
	}
	return hash
}

func parseFileOrTemplate(file string, fallbackContent string, logger *log.Logger) (*template.Template, error) {
	templateContent, err := ioutil.ReadFile(file)
	if err != nil {
		logger.Println("read template:", err)
		templateContent = []byte(fallbackContent)
	}
	return template.New("").Funcs(sprig.FuncMap()).Parse(string(templateContent))
}

func renderOrFallback(templateText string, params map[string]interface{}, fallback string, logger *log.Logger) string {
	if templateText == "" {
		return fallback
	}
	t, err := template.New("").Funcs(sprig.FuncMap()).Parse(string(templateText))
	if err != nil {
		logger.Println("failed parse:", err)
		return fallback
	}
	message := &bytes.Buffer{}
	err = t.Execute(message, params)
	if err != nil {
		logger.Println("failed render:", err)
		return fallback
	}
	return message.String()
}

func realPath(path string, workDir string) string {
	if path == "" {
		return ""
	}
	if filepath.IsAbs(path) {
		return path
	}
	p, _ := filepath.Abs(filepath.Join(workDir, path))
	return p
}
