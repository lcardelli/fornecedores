package handler

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/lcardelli/fornecedores/schemas"
)

type Compra struct {
	CODFILIAL  int
	FORNECEDOR sql.NullString
	TIPO       sql.NullString
}

func ListaFornecedoresHandler(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("userID") // Obtém o ID do usuário da sessão

	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var user schemas.User
	if err := db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// Configuração da conexão com o SQL Server
	connString := os.Getenv("DATABASE_SQL")
	db, err := sql.Open("mssql", connString)
	if err != nil {
		log.Fatal("Erro ao conectar ao banco de dados:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao conectar ao banco de dados"})
		return
	}
	defer db.Close()

	// Query SQL simplificada
	query := `
		SELECT 
			F.CODFILIAL, 
			F.FORNECEDOR, 
			CASE 
				WHEN SUBSTRING(F.CODTBORCAMENTO, 1, 1) = '4' THEN 'Despesas' 
				WHEN SUBSTRING(F.CODTBORCAMENTO, 1, 1) = '1' THEN 'Investimento'
				ELSE 'Rec.Diversas' 
			END AS TIPO
		FROM 
			RESUMO_COMPRAS_CND AS F
		WHERE 
			F.CODFILIAL IN (1, 3, 5)
	`

	rows, err := db.Query(query)
	if err != nil {
		log.Fatal("Erro ao executar a query:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao executar a query"})
		return
	}
	defer rows.Close()

	var compras []Compra
	for rows.Next() {
		var codFilial int
		var fornecedor sql.NullString
		var tipo sql.NullString

		err := rows.Scan(&codFilial, &fornecedor, &tipo)
		if err != nil {
			log.Fatal("Erro ao ler dados da compra:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao ler dados da compra: " + err.Error()})
			return
		}

		// Converta os valores nulos de sql.NullString para string
		fornecedorStr := ""
		if fornecedor.Valid {
			fornecedorStr = fornecedor.String
		}
		tipoStr := ""
		if tipo.Valid {
			tipoStr = tipo.String
		}

		// Crie um objeto Compra com os valores validados
		compras = append(compras, Compra{
			CODFILIAL:  codFilial,
			FORNECEDOR: sql.NullString{String: fornecedorStr, Valid: fornecedor.Valid}, // Corrigido para usar sql.NullString
			TIPO:       sql.NullString{String: tipoStr, Valid: tipo.Valid},             // Corrigido para usar sql.NullString
		})
	}

	// Carregar todos os templates
	templates, err := template.ParseGlob(filepath.Join("templates", "*.html"))
	if err != nil {
		log.Fatal("Erro ao carregar os templates:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao carregar os templates"})
		return
	}

	// Renderizar o template
	err = templates.ExecuteTemplate(c.Writer, "lista_fornecedores.html", gin.H{
		"user":       user,
		"Compras":    compras,
		"activeMenu": "lista-fornecedores", // Adicione isso para destacar o item de menu correto
	})
	if err != nil {
		log.Fatal("Erro ao renderizar o template:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao renderizar o template"})
		return
	}
}
