$(document).ready(function() {
    var allCategories = [];

    loadCategories();

    $('#categoryForm').submit(function(e) {
        e.preventDefault();
        var categoryId = $('#categoryId').val();
        var categoryName = $('#categoryName').val();
        var url = categoryId ? '/api/v1/categories/' + categoryId : '/api/v1/categories';
        var method = categoryId ? 'PUT' : 'POST';

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
                });
                resetForm();
                loadCategories();
            },
            error: function(xhr, status, error) {
                Swal.fire({
                    icon: 'error',
                    title: 'Erro!',
                    text: 'Erro ao processar categoria: ' + error,
                });
            }
        });
    });

    $('#cancelBtn').click(function() {
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
                console.error('Erro ao carregar categorias:', error);
            }
        });
    }

    function renderCategories(categories) {
        var list = $('#categoriesList');
        list.empty();
        if (categories.length === 0) {
            list.html('<p>Nenhuma categoria encontrada.</p>');
        } else {
            categories.forEach(function(category) {
                list.append(
                    `<div class="category-item">
                        <div class="d-flex align-items-center">
                            <input type="checkbox" class="category-checkbox mr-2" data-id="${category.ID}">
                            <span>${category.name}</span>
                        </div>
                        <div>
                            <button class="btn btn-sm btn-warning mr-2 edit-btn" data-id="${category.ID}" data-name="${category.name}"><i class="fas fa-edit"></i></button>
                            <button class="btn btn-sm btn-danger delete-btn" data-id="${category.ID}"><i class="fas fa-trash"></i></button>
                        </div>
                    </div>`
                );
            });
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
            $('#submitBtn').html('<i class="fas fa-save mr-2"></i>Atualizar Categoria');
            $('#cancelBtn').show();
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
            text: `Você está prestes a excluir ${selectedIds.length} categoria(s). Esta ação não pode ser revertida!`,
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
                    'As categorias selecionadas foram excluídas com sucesso.',
                    'success'
                );
                // Remover as categorias excluídas da lista local
                allCategories = allCategories.filter(category => !ids.includes(category.ID));
                // Atualizar a visualização
                renderCategories(allCategories);
            },
            error: function(xhr, status, error) {
                console.error('Erro ao excluir categorias:', xhr.responseText);
                Swal.fire(
                    'Erro!',
                    'Não foi possível excluir as categorias selecionadas: ' + (xhr.responseJSON ? xhr.responseJSON.error : error),
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
                    'A categoria foi deletada com sucesso.',
                    'success'
                );
                // Remover a categoria excluída da lista local
                allCategories = allCategories.filter(category => category.ID !== id);
                // Atualizar a visualização
                renderCategories(allCategories);
            },
            error: function(xhr, status, error) {
                console.error('Erro ao deletar categoria:', xhr.responseText);
                Swal.fire(
                    'Erro!',
                    'Não foi possível deletar a categoria: ' + (xhr.responseJSON ? xhr.responseJSON.error : error),
                    'error'
                );
            }
        });
    }

    function resetForm() {
        $('#categoryId').val('');
        $('#categoryName').val('');
        $('#submitBtn').html('<i class="fas fa-save mr-2"></i>Cadastrar Categoria');
        $('#cancelBtn').hide();
    }
});