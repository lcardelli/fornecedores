package utils

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/url"
	"strings"
	"time"
	
)

// FormatMoney formata um valor float64 para o formato monetário brasileiro
func FormatMoney(value float64) string {
	// Converte o número para string com 2 casas decimais
	str := fmt.Sprintf("%.2f", value)

	// Separa a parte inteira da decimal
	parts := strings.Split(str, ".")

	// Formata a parte inteira com pontos para milhares
	intPart := parts[0]
	var formatted []string
	for i := len(intPart); i > 0; i -= 3 {
		start := i - 3
		if start < 0 {
			start = 0
		}
		formatted = append([]string{intPart[start:i]}, formatted...)
	}

	// Junta tudo no formato brasileiro
	return fmt.Sprintf("R$ %s,%s", strings.Join(formatted, "."), parts[1])
}

// TemplateFuncs retorna um map com todas as funções auxiliares para os templates
func TemplateFuncs() template.FuncMap {
	return template.FuncMap{
		"formatMoney": FormatMoney,
		"formatDate": func(t time.Time) string {
			return t.Format("02/01/2006")
		},
		"lower": strings.ToLower,
		"toJSON":  toJSON,
	}
}

func toJSON(v interface{}) template.JS {
	b, err := json.Marshal(v)
	if err != nil {
		return template.JS("[]")
	}
	return template.JS(url.QueryEscape(string(b)))
}
