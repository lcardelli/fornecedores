package schemas

import (
	"gorm.io/gorm"
	"time"
)

type Supplier struct {
	gorm.Model
	ID          int        `json:"id"`
    Nome        string     `json:"nome"`
    CNPJ        string     `json:"cnpj"`
    Email       string     `json:"email"`
    Telefone    string     `json:"telefone"`
    Endereco    string     `json:"endereco"`
    CategoriaID int        `json:"categoria_id"`
    Categoria   SupplierCategory  `json:"categoria"`
    Servicos    []SupplierService  `json:"servicos"`
    DataCadastro time.Time  `json:"data_cadastro"`
}


