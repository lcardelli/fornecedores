$(document).ready(function() {
    $('#category').change(function() {
        var categoryId = $(this).val();
        $('#service').val('');
        if (categoryId) {
            $.ajax({
                url: '/api/v1/services-by-category/' + categoryId,
                type: 'GET',
                success: function(services) {
                    updateServiceSelect(services);
                },
                error: function() {
                    console.error('Erro ao carregar serviços');
                }
            });
        } else {
            // Se nenhuma categoria for selecionada, carregue todos os serviços
            $.ajax({
                url: '/api/v1/service-list',
                type: 'GET',
                success: function(services) {
                    updateServiceSelect(services);
                },
                error: function() {
                    console.error('Erro ao carregar serviços');
                }
            });
        }
        $('#filterForm').submit();
    });

    function updateServiceSelect(services) {
        var serviceSelect = $('#service');
        serviceSelect.empty();
        serviceSelect.append($('<option>', {
            value: '',
            text: 'Selecione o serviço'
        }));
        $.each(services, function(i, service) {
            serviceSelect.append($('<option>', {
                value: service.ID,
                text: service.name
            }));
        });
    }

    $('#service').change(function() {
        var serviceId = $(this).val();
        $('#product').empty().append($('<option>', {
            value: '',
            text: 'Selecione o produto'
        }));

        if (serviceId) {
            $.ajax({
                url: '/api/v1/products-by-service/' + serviceId,
                type: 'GET',
                success: function(products) {
                    console.log('Produtos recebidos:', products);
                    if (Array.isArray(products) && products.length > 0) {
                        products.forEach(function(product) {
                            $('#product').append($('<option>', {
                                value: product.ID,
                                text: product.name
                            }));
                        });
                    }
                },
                error: function(xhr, status, error) {
                    console.error('Erro ao carregar produtos:', error);
                }
            });
        } else {
            // Se nenhum serviço selecionado, carrega todos os produtos
            $.ajax({
                url: '/api/v1/products',
                type: 'GET',
                success: function(products) {
                    console.log('Todos os produtos recebidos:', products);
                    if (Array.isArray(products) && products.length > 0) {
                        products.forEach(function(product) {
                            $('#product').append($('<option>', {
                                value: product.ID,
                                text: product.name
                            }));
                        });
                    }
                },
                error: function(xhr, status, error) {
                    console.error('Erro ao carregar todos os produtos:', error);
                }
            });
        }
        $('#filterForm').submit();
    });

    function updateProductSelect(products) {
        var productSelect = $('#product');
        productSelect.empty();
        productSelect.append($('<option>', {
            value: '',
            text: 'Selecione o produto'
        }));
        $.each(products, function(i, product) {
            productSelect.append($('<option>', {
                value: product.ID,
                text: product.name
            }));
        });
    }

    // Manipulador de clique para o botão de edição
    $('.edit-supplier').click(function() {
        var supplierId = $(this).data('id');
        var supplierCNPJ = $(this).data('cnpj');
        
        // Carregar dados do fornecedor
        $.ajax({
            url: '/api/v1/suppliers-by-id?id=' + supplierId,
            type: 'GET',
            success: function(supplier) {
                $('#editSupplierId').val(supplier.ID);
                
                // Preencher categorias
                $.ajax({
                    url: '/api/v1/categories',
                    type: 'GET',
                    success: function(categories) {
                        var categorySelect = $('#editCategory');
                        categorySelect.empty();
                        $.each(categories, function(i, category) {
                            categorySelect.append($('<option>', {
                                value: category.ID,
                                text: category.name,
                                selected: category.ID === supplier.Category.ID
                            }));
                        });
                        // Carregar serviços da categoria selecionada
                        loadServicesForCategory(supplier.Category.ID, supplier.Services);
                    }
                });
                
                $('#editSupplierModal').modal('show');
            },
            error: function(xhr, status, error) {
                alert('Erro ao carregar dados do fornecedor: ' + error);
            }
        });
    });

    // Função para carregar serviços de uma categoria
    function loadServicesForCategory(categoryId, supplierServices) {
        $.ajax({
            url: '/api/v1/services-by-category/' + categoryId,
            type: 'GET',
            success: function(services) {
                updateServicesCheckboxes(services, supplierServices);
            },
            error: function(xhr, status, error) {
                console.error('Erro ao carregar serviços:', error);
            }
        });
    }

    // Função para atualizar checkboxes de serviços
    function updateServicesCheckboxes(services, supplierServices) {
        const servicesDiv = document.getElementById('editServices');
        servicesDiv.innerHTML = ''; // Limpar serviços existentes

        $.each(services, function(i, service) {
            var serviceDiv = document.createElement('div');
            serviceDiv.className = 'form-check';

            var input = document.createElement('input');
            input.className = 'form-check-input';
            input.type = 'checkbox';
            input.name = 'services[]';
            input.value = service.ID;
            input.id = 'editService' + service.ID;

            // Verificar se o serviço está atribuído ao fornecedor
            var isChecked = supplierServices.some(function(supplierService) {
                return supplierService.ServiceID === service.ID;
            });
            input.checked = isChecked;

            var label = document.createElement('label');
            label.className = 'form-check-label';
            label.htmlFor = 'editService' + service.ID;
            label.textContent = service.name;

            serviceDiv.appendChild(input);
            serviceDiv.appendChild(label);

            servicesDiv.appendChild(serviceDiv);
        });
    }

    // Adicionar evento de mudança para o select de categoria no modal de edição
    $('#editCategory').change(function() {
        var categoryId = $(this).val();
        loadServicesForCategory(categoryId, []);
    });

    // Manipulador de clique para salvar alterações
    $('#saveSupplierChanges').click(function() {
        var supplierId = $('#editSupplierId').val();
        var categoryId = $('#editCategory').val();
        var selectedServices = $('input[name="services[]"]:checked').map(function() {
            return parseInt(this.value);
        }).get();

        var data = {
            category_id: parseInt(categoryId),
            service_ids: selectedServices
        };

        console.log('Dados a serem enviados:', data);

        $.ajax({
            url: '/api/v1/suppliers/' + supplierId,
            type: 'PUT',
            contentType: 'application/json',
            data: JSON.stringify(data),
            success: function(result) {
                console.log('Resposta do servidor:', result);
                Swal.fire({
                    icon: 'success',
                    title: 'Sucesso!',
                    text: 'Fornecedor atualizado com sucesso!',
                    confirmButtonText: 'OK'
                }).then((result) => {
                    if (result.isConfirmed) {
                        $('#editSupplierModal').modal('hide');
                        location.reload();
                    }
                });
            },
            error: function(xhr, status, error) {
                console.error('Erro ao atualizar fornecedor:', xhr.responseText);
                Swal.fire({
                    icon: 'error',
                    title: 'Erro!',
                    text: 'Erro ao atualizar fornecedor: ' + error,
                    confirmButtonText: 'OK'
                });
            }
        });
    });

    // Manipulador de clique para o botão de exclusão
    $('.delete-supplier').click(function() {
        var supplierId = $(this).data('id');
        var supplierCNPJ = $(this).data('cnpj');
        
        Swal.fire({
            title: 'Tem certeza?',
            text: `Deseja deletar este fornecedor (CNPJ: ${supplierCNPJ}) do catálogo?`,
            icon: 'warning',
            showCancelButton: true,
            confirmButtonColor: '#3085d6',
            cancelButtonColor: '#d33',
            confirmButtonText: 'Sim, deletar!',
            cancelButtonText: 'Cancelar'
        }).then((result) => {
            if (result.isConfirmed) {
                $.ajax({
                    url: '/api/v1/suppliers/' + supplierId,
                    type: 'DELETE',
                    success: function(result) {
                        Swal.fire(
                            'Deletado!',
                            'Fornecedor removido com sucesso.',
                            'success'
                        ).then(() => {
                            location.reload();
                        });
                    },
                    error: function(xhr, status, error) {
                        Swal.fire(
                            'Erro!',
                            'Erro ao deletar fornecedor: ' + error,
                            'error'
                        );
                    }
                });
            }
        });
    });

    // Função para atualizar a lista de categorias
    function updateCategoryList(categories) {
        var categorySelect = $('#category');
        categorySelect.empty();
        categorySelect.append($('<option>', {
            value: '',
            text: 'Selecione a categoria'
        }));
        $.each(categories, function(i, category) {
            categorySelect.append($('<option>', {
                value: category.ID,
                text: category.Name
            }));
        });
    }

    // Função para atualizar a lista de serviços
    function updateServiceList(services) {
        var serviceSelect = $('#service');
        serviceSelect.empty();
        serviceSelect.append($('<option>', {
            value: '',
            text: 'Selecione o serviço'
        }));
        $.each(services, function(i, service) {
            serviceSelect.append($('<option>', {
                value: service.ID,
                text: service.Name
            }));
        });
    }

    // Adicione event listeners para os formulários de criação de categoria e serviço
    $('#createCategoryForm').submit(function(e) {
        e.preventDefault();
        // ... código existente para enviar a requisição ...
        $.ajax({
            // ... configurações existentes ...
            success: function(response) {
                // ... código existente ...
                updateCategoryList(response.categories);
            }
        });
    });

    $('#createServiceForm').submit(function(e) {
        e.preventDefault();
        var serviceName = $('#serviceName').val();
        var categoryId = $('#serviceCategory').val();
        
        $.ajax({
            url: '/api/v1/services',
            type: 'POST',
            contentType: 'application/json',
            data: JSON.stringify({
                name: serviceName,
                category_id: parseInt(categoryId)
            }),
            success: function(response) {
                // ... (código existente para lidar com o sucesso)
                updateServiceList(response.services);
            },
            error: function(xhr, status, error) {
                // ... (código existente para lidar com erros)
            }
        });
    });

    $('#createProductForm').submit(function(e) {
        e.preventDefault();
        var productName = $('#productName').val();
        var serviceId = $('#productService').val();
        
        $.ajax({
            url: '/api/v1/products',
            type: 'POST',
            contentType: 'application/json',
            data: JSON.stringify({
                name: productName,
                service_id: parseInt(serviceId)
            }),
            success: function(response) {
                Swal.fire({
                    icon: 'success',
                    title: 'Sucesso!',
                    text: 'Produto criado com sucesso!',
                    confirmButtonText: 'OK'
                }).then(() => {
                    $('#productName').val('');
                    updateProductSelect(response.products);
                });
            },
            error: function(xhr, status, error) {
                Swal.fire({
                    icon: 'error',
                    title: 'Erro!',
                    text: 'Erro ao criar produto: ' + error,
                    confirmButtonText: 'OK'
                });
            }
        });
    });

    // Adicionar o handler para mudança no select de produtos
    $('#product').change(function() {
        $('#filterForm').submit();
    });

    // Adicionar um handler para o submit do formulário
    $('#filterForm').on('submit', function(e) {
        e.preventDefault();
        
        var categoryId = $('#category').val();
        var serviceId = $('#service').val();
        var productId = $('#product').val();
        var name = $('#name').val();

        // Construir a URL com os parâmetros
        var url = '/api/v1/catalogo?';
        var params = [];
        
        if (categoryId) params.push('category=' + categoryId);
        if (serviceId) params.push('service=' + serviceId);
        if (productId) params.push('product=' + productId);
        if (name) params.push('name=' + encodeURIComponent(name));
        
        url += params.join('&');
        
        // Redirecionar para a URL com os filtros
        window.location.href = url;
    });
});
