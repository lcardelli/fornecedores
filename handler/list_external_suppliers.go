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

type Fornecedor struct {
	CODCOLIGADA   int
	CODCFO        string
	NOMEFANTASIA  sql.NullString
	NOME          sql.NullString
	CGCCFO        sql.NullString
	RUA           sql.NullString
	NUMERO        sql.NullString
	COMPLEMENTO   sql.NullString
	BAIRRO        sql.NullString
	CIDADE        sql.NullString
	CEP           sql.NullString
	TELEFONE      sql.NullString
	EMAIL         sql.NullString
	CONTATO       sql.NullString
	UF            sql.NullString
	ATIVO         sql.NullString
	TIPO          sql.NullString
}

// ListaFornecedoresExternosHandler lista todos os fornecedores externos
func ListaFornecedoresExternosHandler(c *gin.Context) {
	// Obter o usuário do contexto
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}

	// Obter informações do usuário
	user, ok := userInterface.(schemas.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao obter informações do usuário"})
		return
	}

	// Buscar fornecedores externos do banco de dados
	fornecedores, err := getFornecedoresExternosFromDatabase()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar dados de fornecedores: " + err.Error()})
		return
	}

	// Renderizar o template
	c.HTML(http.StatusOK, "lista_fornecedores.html", gin.H{
		"user":         user,
		"Fornecedores": fornecedores,
		"activeMenu":   "lista-fornecedores",
	})
}

// getFornecedoresExternosFromDatabase busca fornecedores externos do banco de dados
func getFornecedoresExternosFromDatabase() ([]Fornecedor, error) {
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

	// Query para buscar fornecedores externos
	query := `
		Select
			FCFO.CODCOLIGADA,
			FCFO.CODCFO,
			FCFO.NOMEFANTASIA,
			FCFO.NOME,
			FCFO.CGCCFO,
			FCFO.RUA,
			FCFO.NUMERO,
			FCFO.COMPLEMENTO,
			FCFO.BAIRRO,
			FCFO.CIDADE,
			FCFO.CEP,
			FCFO.TELEFONE,
			FCFO.EMAIL,
			FCFO.CONTATO,
			FCFO.CODETD As UF,
			FCFO.ATIVO As ATIVO,
			FTCF.DESCRICAO AS TIPO
		From
			FCFO Inner Join
			FTCF On FCFO.CODCOLTCF = FTCF.CODCOLIGADA
					And FCFO.CODTCF = FTCF.CODTCF
		Where
			FCFO.CODCOLIGADA = '0' And
			FCFO.ATIVO = '1' And
			FCFO.PESSOAFISOUJUR = 'J' And
			FCFO.PAGREC in ('2', '3')
	`
	// Executar a query
	rows, err := db.Query(query)
	if err != nil {
		log.Printf("Erro ao executar query: %v", err)
		return nil, err
	}
	defer rows.Close()
	// Ler os dados do banco de dados
	var fornecedores []Fornecedor
	for rows.Next() {
		var f Fornecedor
		err := rows.Scan(
				&f.CODCOLIGADA, &f.CODCFO, &f.NOMEFANTASIA, &f.NOME, &f.CGCCFO,
				&f.RUA, &f.NUMERO, &f.COMPLEMENTO, &f.BAIRRO, &f.CIDADE,
				&f.CEP, &f.TELEFONE, &f.EMAIL, &f.CONTATO, &f.UF,
				&f.ATIVO, &f.TIPO,
		)
		if err != nil {
			log.Printf("Erro ao ler dados do fornecedor: %v", err)
			continue
		}
		// Adicionar o fornecedor à lista
		fornecedores = append(fornecedores, f)
	}

	// Retornar a lista de fornecedores
	return fornecedores, nil
}
