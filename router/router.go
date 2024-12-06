package router

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"html/template"
	"fmt"
	"time"
	"strings"
)

// Funções auxiliares para os templates
var templateFuncs = template.FuncMap{
	"formatMoney": func(value float64) string {
		return fmt.Sprintf("R$ %.2f", value)
	},
	"formatDate": func(t time.Time) string {
		return t.Format("02/01/2006")
	},
	"lower": strings.ToLower,
}

// Initialize configura o roteador e as sessões
func Initialize() {
	// Cria um novo roteador Gin
	router := gin.Default()

	// Configura o middleware de sessão
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("mysession", store))

	// Registra as funções auxiliares antes de carregar os templates
	router.SetFuncMap(templateFuncs)

	// Carrega os templates
	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "./static")

	// Inicializa as rotas
	InitializeRoutes(router)

	// Inicia o servidor Gin na porta 8080
	router.Run(":8080")
}
