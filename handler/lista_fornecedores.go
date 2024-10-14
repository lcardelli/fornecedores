package handler

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/denisenkom/go-mssqldb"
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
	// Obter o usuário do contexto
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}
	user, ok := userInterface.(schemas.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao obter informações do usuário"})
		return
	}

	compras, err := getComprasFromDatabase()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar dados de compras: " + err.Error()})
		return
	}

	// Renderizar o template
	c.HTML(http.StatusOK, "lista_fornecedores.html", gin.H{
		"user":       user,
		"Compras":    compras,
		"activeMenu": "lista-fornecedores",
	})
}

func getComprasFromDatabase() ([]Compra, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// Configuração da conexão com o SQL Server
	connString := os.Getenv("DATABASE_SQL")
	db, err := sql.Open("mssql", connString)
	if err != nil {
		return nil, err
	}
	defer db.Close()

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
		return nil, err
	}
	defer rows.Close()

	var compras []Compra
	for rows.Next() {
		var compra Compra
		err := rows.Scan(&compra.CODFILIAL, &compra.FORNECEDOR, &compra.TIPO)
		if err != nil {
			log.Printf("Erro ao ler dados da compra: %v", err)
			continue
		}
		compras = append(compras, compra)
	}

	return compras, nil
}
