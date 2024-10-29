$(document).ready(function() {
    var allProducts = [];

    loadProducts();

    $('#productForm').submit(function(e) {
        e.preventDefault();
        var productId = $('#productId').val();
        var productName = $('#productName').val();
        var serviceId = $('#serviceId').val();
        var url = productId ? '/api/v1/products/' + productId : '/api/v1/products';
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
                    resetForm();
                    loadProducts();
                });
            },
            error: function(xhr, status, error) {
                console.error('Erro:', xhr.responseText);
                Swal.fire({
                    icon: 'error',
                    title: 'Erro!',
                    text: 'Erro ao processar produto: ' + error,
                });
            }
        });
    });

    $('#cancelBtn').click(function() {
        resetForm();
    });

    function loadProducts() {
        $.ajax({
            url: '/api/v1/products-list',
            type: 'GET',
            success: function(response) {
                allProducts = response;
                renderProducts(allProducts);
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
            list.html('<p>Nenhum produto encontrado.</p>');
        } else {
            products.forEach(function(product) {
                var serviceName = product.Service ? product.Service.name : 'Sem serviço';
                list.append(
                    `<div class="product-item">
                        <div class="d-flex align-items-center">
                            <input type="checkbox" class="product-checkbox mr-2" data-id="${product.ID}">
                            <span>${product.name} (${serviceName})</span>
                        </div>
                        <div>
                            <button class="btn btn-sm btn-warning mr-2 edit-btn" 
                                data-id="${product.ID}" 
                                data-name="${product.name}" 
                                data-service="${product.ServiceID}">
                                <i class="fas fa-edit"></i>
                            </button>
                            <button class="btn btn-sm btn-danger delete-btn" data-id="${product.ID}">
                                <i class="fas fa-trash"></i>
                            </button>
                        </div>
                    </div>`
                );
            });
            setupEditButtons();
            setupDeleteButtons();
            setupCheckboxEvents();
        }
    }

    function setupEditButtons() {
        $('.edit-btn').click(function() {
            var id = $(this).data('id');
            var name = $(this).data('name');
            var serviceId = $(this).data('service');
            $('#productId').val(id);
            $('#productName').val(name);
            $('#serviceId').val(serviceId);
            $('#submitBtn').html('<i class="fas fa-save mr-2"></i>Atualizar Produto');
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

    function setupCheckboxEvents() {
        $('.product-checkbox').change(function() {
            updateDeleteSelectedButton();
        });
    }

    function updateDeleteSelectedButton() {
        var checkedCount = $('.product-checkbox:checked').length;
        var hasProducts = allProducts.length > 0;
        $('#deleteSelectedBtn').toggle(checkedCount > 0);
        $('#selectAllBtn').toggle(hasProducts);
    }

    $('#selectAllBtn').click(function() {
        var isAllSelected = $('.product-checkbox:checked').length === $('.product-checkbox').length;
        $('.product-checkbox').prop('checked', !isAllSelected);
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

    function deleteProduct(id) {
        $.ajax({
            url: '/api/v1/products/' + id,
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

    function resetForm() {
        $('#productId').val('');
        $('#productName').val('');
        $('#serviceId').val('');
        $('#submitBtn').html('<i class="fas fa-save mr-2"></i>Cadastrar Produto');
        $('#cancelBtn').hide();
    }

    $('#productSearch').on('input', function() {
        var searchTerm = $(this).val().toLowerCase();
        var serviceId = $('#serviceFilter').val();
        filterProducts(searchTerm, serviceId);
    });

    $('#serviceFilter').change(function() {
        var searchTerm = $('#productSearch').val().toLowerCase();
        var serviceId = $(this).val();
        filterProducts(searchTerm, serviceId);
    });

    function filterProducts(searchTerm, serviceId) {
        var filteredProducts = allProducts.filter(function(product) {
            var matchesSearch = product.name.toLowerCase().includes(searchTerm);
            var matchesService = !serviceId || product.ServiceID == serviceId;
            return matchesSearch && matchesService;
        });
        renderProducts(filteredProducts);
    }
}); 