$(document).ready(function() {
    var allCategories = [];

    loadCategories();

    $('#categoryForm').submit(function(e) {
        e.preventDefault();
        var categoryId = $('#categoryId').val();
        var categoryName = $('#categoryName').val().trim();
        var url = categoryId ? '/api/v1/categories/' + categoryId : '/api/v1/categories';
        var method = categoryId ? 'PUT' : 'POST';

        if (!categoryName) {
            Swal.fire({
                icon: 'error',
                title: 'Erro!',
                text: 'O nome da categoria não pode estar vazio'
            });
            return;
        }

        $.ajax({
            url: url,
            type: method,
            data: JSON.stringify({name: categoryName}),
            contentType: 'application/json',
            success: function(response) {
                Swal.fire({
                    icon: 'success',
                    title: 'Sucesso!',
                    text: categoryId ? 'Categoria atualizada com sucesso!' : 'Categoria cadastrada com sucesso!',
                }).then(() => {
                    resetForm();
                    loadCategories();
                });
            },
            error: function(xhr, status, error) {
                console.error('Erro:', xhr.responseText);
                let errorMessage = 'Erro ao processar categoria';
                
                if (xhr.responseJSON) {
                    if (xhr.responseJSON.error === "Category already exists") {
                        errorMessage = 'Esta categoria já está cadastrada no sistema.';
                    } else {
                        errorMessage = xhr.responseJSON.error;
                    }
                }

                Swal.fire({
                    icon: 'error',
                    title: 'Erro!',
                    text: errorMessage
                });
            }
        });
    });

    $('#cancelBtn').click(function() {
        $('#formSection').slideUp();
        resetForm();
    });

    $('#newCategoryBtn').click(function() {
        $('#formSection').slideDown();
        resetForm();
    });

    function loadCategories() {
        $.ajax({
            url: '/api/v1/categories',
            type: 'GET',
            success: function(categories) {
                allCategories = categories;
                renderCategories(allCategories);
            },
            error: function(xhr, status, error) {
                console.error('Erro ao carregar áreas:', error);
            }
        });
    }

    function renderCategories(categories) {
        var list = $('#categoriesList');
        list.empty();
        
        if (categories.length === 0) {
            list.html('<div class="p-3 text-center">Nenhuma área encontrada.</div>');
        } else {
            var table = `
                <table class="table table-hover mb-0">
                    <thead>
                        <tr>
                            <th width="40px"></th>
                            <th>Nome da Área</th>
                            <th width="120px">Ações</th>
                        </tr>
                    </thead>
                    <tbody>
            `;
            
            categories.forEach(function(category) {
                table += `
                    <tr>
                        <td>
                            <input type="checkbox" class="category-checkbox" data-id="${category.ID}">
                        </td>
                        <td>${category.name}</td>
                        <td>
                            <div class="btn-group-actions">
                                <button class="btn btn-sm btn-warning edit-btn" data-id="${category.ID}" data-name="${category.name}">
                                    <i class="fas fa-edit"></i>
                                </button>
                                <button class="btn btn-sm btn-danger delete-btn" data-id="${category.ID}">
                                    <i class="fas fa-trash"></i>
                                </button>
                            </div>
                        </td>
                    </tr>
                `;
            });
            
            table += '</tbody></table>';
            list.html(table);
            
            setupEditButtons();
            setupDeleteButtons();
            setupCheckboxEvents();
        }
        updateDeleteSelectedButton();
    }

    function setupEditButtons() {
        $('.edit-btn').click(function() {
            var id = $(this).data('id');
            var name = $(this).data('name');
            $('#categoryId').val(id);
            $('#categoryName').val(name);
            $('#submitBtn').html('<i class="fas fa-save mr-2"></i>Atualizar');
            $('#formSection').slideDown();
        });
    }

    function setupDeleteButtons() {
        $('.delete-btn').click(function() {
            var id = $(this).data('id');
            Swal.fire({
                title: 'Tem certeza?',
                text: "Você não poderá reverter esta ação!",
                icon: 'warning',
                showCancelButton: true,
                confirmButtonColor: '#3085d6',
                cancelButtonColor: '#d33',
                confirmButtonText: 'Sim, delete!'
            }).then((result) => {
                if (result.isConfirmed) {
                    deleteCategory(id);
                }
            });
        });
    }

    function setupCheckboxEvents() {
        $('.category-checkbox').change(function() {
            updateDeleteSelectedButton();
        });
    }

    function updateDeleteSelectedButton() {
        var checkedCount = $('.category-checkbox:checked').length;
        var hasCategories = allCategories.length > 0;
        $('#deleteSelectedBtn').toggle(checkedCount > 0);
        $('#selectAllBtn').toggle(hasCategories);
    }

    $('#selectAllBtn').click(function() {
        var isAllSelected = $('.category-checkbox:checked').length === $('.category-checkbox').length;
        $('.category-checkbox').prop('checked', !isAllSelected);
        $(this).html(
            isAllSelected ? 
            '<i class="fas fa-check-square mr-2"></i>Selecionar Todos' : 
            '<i class="fas fa-square mr-2"></i>Desmarcar Todos'
        );
        updateDeleteSelectedButton();
    });

    $('#deleteSelectedBtn').click(function() {
        var selectedIds = $('.category-checkbox:checked').map(function() {
            return $(this).data('id');
        }).get();

        if (selectedIds.length === 0) return;

        Swal.fire({
            title: 'Tem certeza?',
            text: `Você está prestes a excluir ${selectedIds.length} área(s). Esta ação não pode ser revertida!`,
            icon: 'warning',
            showCancelButton: true,
            confirmButtonColor: '#3085d6',
            cancelButtonColor: '#d33',
            confirmButtonText: 'Sim, excluir!',
            cancelButtonText: 'Cancelar'
        }).then((result) => {
            if (result.isConfirmed) {
                deleteMultipleCategories(selectedIds);
            }
        });
    });

    function deleteMultipleCategories(ids) {
        $.ajax({
            url: '/api/v1/categories/batch',
            type: 'DELETE',
            data: JSON.stringify({ ids: ids }),
            contentType: 'application/json',
            success: function(response) {
                Swal.fire(
                    'Excluídas!',
                    'As áreas selecionadas foram excluídas com sucesso.',
                    'success'
                );
                // Remover as categorias excluídas da lista local
                allCategories = allCategories.filter(category => !ids.includes(category.ID));
                // Atualizar a visualização
                renderCategories(allCategories);
            },
            error: function(xhr, status, error) {
                console.error('Erro ao excluir áreas:', xhr.responseText);
                Swal.fire(
                    'Erro!',
                    'Não foi possível excluir as áreas selecionadas: ' + (xhr.responseJSON ? xhr.responseJSON.error : error),
                    'error'
                );
            }
        });
    }

    function deleteCategory(id) {
        $.ajax({
            url: '/api/v1/categories/' + id,
            type: 'DELETE',
            success: function(response) {
                Swal.fire(
                    'Deletado!',
                    'A área foi deletada com sucesso.',
                    'success'
                );
                // Remover a categoria excluída da lista local
                allCategories = allCategories.filter(category => category.ID !== id);
                // Atualizar a visualização
                renderCategories(allCategories);
            },
            error: function(xhr, status, error) {
                console.error('Erro ao deletar área:', xhr.responseText);
                Swal.fire(
                    'Erro!',
                    'Não foi possível deletar a área: ' + (xhr.responseJSON ? xhr.responseJSON.error : error),
                    'error'
                );
            }
        });
    }

    function resetForm() {
        $('#categoryId').val('');
        $('#categoryName').val('');
        $('#submitBtn').html('<i class="fas fa-save mr-2"></i>Salvar');
    }

    // Adicionar evento de pesquisa
    $('#categorySearch').on('input', function() {
        var searchTerm = $(this).val().toLowerCase();
        
        // Filtrar as categorias baseado no termo de busca
        var filteredCategories = allCategories.filter(function(category) {
            return category.name.toLowerCase().includes(searchTerm);
        });
        
        // Renderizar apenas as categorias filtradas
        renderCategories(filteredCategories);
    });
});