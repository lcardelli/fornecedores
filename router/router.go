package router

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
)

// Initialize configura o roteador e as sessões
func Initialize() {
	// Cria um novo roteador Gin
	router := gin.Default()

	// Configura o middleware de sessão
	store := cookie.NewStore([]byte("secret")) // Use uma chave secreta segura
	router.Use(sessions.Sessions("mysession", store))

	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "./static")

	// Inicializa as rotas
	InitializeRoutes(router)

	// Inicia o servidor Gin na porta 8080
	router.Run(":8080")
}
