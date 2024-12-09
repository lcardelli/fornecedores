package router

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/lcardelli/fornecedores/utils"
)

// Initialize configura o roteador e as sessões
func Initialize() {
	// Cria um novo roteador Gin
	router := gin.Default()

	// Configura o middleware de sessão
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("mysession", store))

	// Registra as funções auxiliares antes de carregar os templates
	router.SetFuncMap(utils.TemplateFuncs())

	// Carrega os templates
	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "./static")

	// Inicializa as rotas
	InitializeRoutes(router)

	// Inicia o servidor Gin na porta 8080
	router.Run(":8080")
}
