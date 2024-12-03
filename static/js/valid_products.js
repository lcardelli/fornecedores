$(document).ready(function() {
    var allProducts = [];

    loadProducts();

    $('#newProductBtn').click(function() {
        $('#formSection').slideDown();
        resetForm();
    });

    $('#productForm').submit(function(e) {
        e.preventDefault();
        var productId = $('#productId').val();
        var productName = $('#productName').val().trim();
        var serviceId = $('#serviceId').val();

        if (!productName) {
            Swal.fire({
                icon: 'error',
                title: 'Erro!',
                text: 'O nome do produto não pode estar vazio'
            });
            return;
        }

        if (!serviceId) {
            Swal.fire({
                icon: 'error',
                title: 'Erro!',
                text: 'Selecione um serviço'
            });
            return;
        }

        var url = productId ? `/api/v1/products/${productId}` : '/api/v1/products';
        var method = productId ? 'PUT' : 'POST';

        var data = {
            name: productName,
            service_id: parseInt(serviceId)
        };

        $.ajax({
            url: url,
            type: method,
            data: JSON.stringify(data),
            contentType: 'application/json',
            success: function(response) {
                Swal.fire({
                    icon: 'success',
                    title: 'Sucesso!',
                    text: productId ? 'Produto atualizado com sucesso!' : 'Produto cadastrado com sucesso!',
                }).then(() => {
                    $('#formSection').slideUp();
                    resetForm();
                    loadProducts();
                });
            },
            error: function(xhr, status, error) {
                let errorMessage = 'Erro ao processar produto';
                if (xhr.responseJSON && xhr.responseJSON.error) {
                    errorMessage = xhr.responseJSON.error;
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

    function loadProducts() {
        $.ajax({
            url: '/api/v1/products-list',
            type: 'GET',
            success: function(response) {
                allProducts = response;
                filterProducts();
            },
            error: function(xhr, status, error) {
                console.error('Erro ao carregar produtos:', error);
            }
        });
    }

    function renderProducts(products) {
        var list = $('#productsList');
        list.empty();
        
        if (products.length === 0) {
            list.html('<div class="empty-state"><i class="fas fa-box-open"></i><h4>Nenhum produto encontrado</h4><p>Não há produtos cadastrados com os filtros atuais.</p></div>');
            return;
        }

        var table = `
            <table class="table table-hover mb-0">
                <thead>
                    <tr>
                        <th width="40px">
                            <input type="checkbox" id="selectAll" class="select-all-checkbox">
                        </th>
                        <th>Nome do Produto</th>
                        <th>Categoria</th>
                        <th width="120px">Ações</th>
                    </tr>
                </thead>
                <tbody>
        `;
        
        products.forEach(function(product) {
            var serviceName = product.Service ? product.Service.name : 'Sem categoria';
            table += `
                <tr>
                    <td>
                        <input type="checkbox" class="product-checkbox" data-id="${product.ID}">
                    </td>
                    <td>${product.name}</td>
                    <td>${serviceName}</td>
                    <td>
                        <div class="btn-group-actions">
                            <button class="btn btn-sm btn-warning edit-btn" data-id="${product.ID}" data-name="${product.name}" data-service="${product.ServiceID}">
                                <i class="fas fa-edit"></i>
                            </button>
                            <button class="btn btn-sm btn-danger delete-btn" data-id="${product.ID}">
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
        updateDeleteSelectedButton();

        // Após renderizar a tabela, aplique os delays de animação
        $('.table tbody tr').each(function(index) {
            $(this).css({
                'animation-delay': `${index * 0.1}s`
            });
        });
    }

    function setupEditButtons() {
        $('.edit-btn').click(function() {
            var id = $(this).data('id');
            var name = $(this).data('name');
            var serviceId = $(this).data('service');
            $('#productId').val(id);
            $('#productName').val(name);
            $('#serviceId').val(serviceId);
            $('#submitBtn').html('<i class="fas fa-save mr-2"></i>Atualizar');
            $('#formSection').slideDown();
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
                    deleteProduct(id);
                }
            });
        });
    }

    function deleteProduct(id) {
        $.ajax({
            url: `/api/v1/products/${id}`,
            type: 'DELETE',
            success: function(response) {
                Swal.fire(
                    'Deletado!',
                    'O produto foi deletado com sucesso.',
                    'success'
                );
                loadProducts();
            },
            error: function(xhr, status, error) {
                console.error('Erro ao deletar produto:', xhr.responseText);
                Swal.fire(
                    'Erro!',
                    'Não foi possível deletar o produto: ' + (xhr.responseJSON ? xhr.responseJSON.error : error),
                    'error'
                );
            }
        });
    }

    function setupCheckboxEvents() {
        $('.product-checkbox').change(function() {
            updateDeleteSelectedButton();
        });
    }

    function updateDeleteSelectedButton() {
        var checkedCount = $('.product-checkbox:checked').length;
        $('#deleteSelectedBtn').toggle(checkedCount > 0);
    }

    $('#selectAllBtn').click(function() {
        var isAllSelected = $('.product-checkbox:checked').length === $('.product-checkbox').length;
        $('.product-checkbox').prop('checked', !isAllSelected);
        $(this).html(
            isAllSelected ? 
            '<i class="fas fa-check-square mr-2"></i>Selecionar Todos' : 
            '<i class="fas fa-square mr-2"></i>Desmarcar Todos'
        );
        updateDeleteSelectedButton();
    });

    $('#deleteSelectedBtn').click(function() {
        var selectedIds = $('.product-checkbox:checked').map(function() {
            return $(this).data('id');
        }).get();

        if (selectedIds.length === 0) return;

        Swal.fire({
            title: 'Tem certeza?',
            text: `Você está prestes a excluir ${selectedIds.length} produto(s). Esta ação não pode ser revertida!`,
            icon: 'warning',
            showCancelButton: true,
            confirmButtonColor: '#3085d6',
            cancelButtonColor: '#d33',
            confirmButtonText: 'Sim, excluir!',
            cancelButtonText: 'Cancelar'
        }).then((result) => {
            if (result.isConfirmed) {
                deleteMultipleProducts(selectedIds);
            }
        });
    });

    function deleteMultipleProducts(ids) {
        $.ajax({
            url: '/api/v1/products/batch',
            type: 'DELETE',
            data: JSON.stringify({ ids: ids }),
            contentType: 'application/json',
            success: function(response) {
                Swal.fire(
                    'Excluídos!',
                    'Os produtos selecionados foram excluídos com sucesso.',
                    'success'
                );
                loadProducts();
            },
            error: function(xhr, status, error) {
                console.error('Erro ao excluir produtos:', xhr.responseText);
                Swal.fire(
                    'Erro!',
                    'Não foi possível excluir os produtos selecionados: ' + (xhr.responseJSON ? xhr.responseJSON.error : error),
                    'error'
                );
            }
        });
    }

    function resetForm() {
        $('#productId').val('');
        $('#productName').val('');
        $('#serviceId').val('');
        $('#submitBtn').html('<i class="fas fa-save mr-2"></i>Salvar');
        $('#cancelBtn').show();
    }

    $('#productSearch, #serviceFilter').on('input change', filterProducts);

    function filterProducts() {
        var searchTerm = $('#productSearch').val().toLowerCase();
        var serviceId = $('#serviceFilter').val();

        var filteredProducts = allProducts;
        
        if (searchTerm || serviceId) {
            filteredProducts = allProducts.filter(function(product) {
                var matchesSearch = !searchTerm || product.name.toLowerCase().includes(searchTerm);
                var matchesService = !serviceId || product.ServiceID == serviceId;
                return matchesSearch && matchesService;
            });
        }

        renderProducts(filteredProducts);
    }

    $(document).on('change', '.select-all-checkbox', function() {
        var isChecked = $(this).prop('checked');
        $('.product-checkbox').prop('checked', isChecked);
        updateDeleteSelectedButton();
    });
}); 