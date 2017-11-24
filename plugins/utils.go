package plugins

import (
	"html/template"
	"log"
	"github.com/Masterminds/sprig"
	"io/ioutil"
	"bytes"
	"path/filepath"
)

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
