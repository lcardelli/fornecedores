package handler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/lcardelli/fornecedores/config"
	"github.com/lcardelli/fornecedores/schemas"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gorm.io/gorm"
)

// AuthHandler é o manipulador para autenticação
var (
	oauthConfig *oauth2.Config
	db          *gorm.DB
)

// init inicializa as variáveis de ambiente e a configuração do OAuth2
func init() {
	// Carregar variáveis de ambiente do arquivo .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// Configuração do OAuth2
	oauthConfig = &oauth2.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		RedirectURL:  os.Getenv("REDIRECT_URL"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
	db = config.GetMysql()
}

// GoogleLogin redireciona para a página de login do Google
func GoogleLogin(c *gin.Context) {
	url := oauthConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// GoogleCallback é o callback do Google após o login
func GoogleCallback(c *gin.Context) {
	code := c.Query("code")
	token, err := oauthConfig.Exchange(c, code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to exchange token"})
		return
	}

	// Obter informações do usuário
	client := oauthConfig.Client(c, token)
	response, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get user info"})
		return
	}
	defer response.Body.Close()

	// Verificar se a resposta do Google é OK
	if response.StatusCode != http.StatusOK {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get user info, status: " + response.Status})
		return
	}

	// Ler o corpo da resposta
	userInfo, err := io.ReadAll(response.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read user info"})
		return
	}

	// Estrutura para deserializar a resposta do Google
	type GoogleUserInfo struct {
		ID            string `json:"id"`
		Email         string `json:"email"`
		Name          string `json:"name"`
		Picture       string `json:"picture"`
		VerifiedEmail bool   `json:"verified_email"`
	}

	// Deserializar o JSON
	var googleUser GoogleUserInfo
	if err := json.Unmarshal(userInfo, &googleUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unmarshal user info"})
		return
	}

	// Declarar a variável user no início da função
	var user schemas.User

	// Verificar se o usuário já existe no banco de dados
	result := db.Where("email = ?", googleUser.Email).First(&user)

	if result.Error == nil {
		// Usuário já existe, atualizar informações se necessário
		user.Name = googleUser.Name
		user.Avatar = googleUser.Picture
		if err := db.Save(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user in database"})
			return
		}
	} else if result.Error == gorm.ErrRecordNotFound {
		// Usuário não existe, criar novo registro
		user = schemas.User{
			Name:   googleUser.Name,
			Email:  googleUser.Email,
			Avatar: googleUser.Picture,
		}
		if err := db.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save user to database"})
			return
		}
	} else {
		// Outro erro ocorreu ao consultar o banco de dados
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check user in database"})
		return
	}

	// Armazenar o ID do usuário na sessão
	session := sessions.Default(c)
	session.Set("userID", user.ID)
	session.Save()

	// Redirecionar para a dashboard após o login
	c.Redirect(http.StatusFound, "/api/v1/dashboard")
}


// GoogleLogout trata o logout do usuário
func GoogleLogout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear() // Limpa a sessão do usuário
	session.Save()  // Salva as alterações na sessão

	// Redireciona para a página de login
	c.Redirect(http.StatusFound, "/api/v1/index") // Corrigido para a URL da sua página de login
}
