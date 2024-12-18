$(document).ready(function() {
    $('#category').change(function() {
        var categoryId = $(this).val();
        $('#service').val('');
        $('#product').empty().append($('<option>', {
            value: '',
            text: 'Selecione o produto'
        }));

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
            // Se nenhuma categoria for selecionada, carregue todos os serviços e produtos
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

            // Carregar todos os produtos
            $.ajax({
                url: '/api/v1/products',
                type: 'GET',
                success: function(products) {
                    updateProductSelect(products);
                },
                error: function() {
                    console.error('Erro ao carregar produtos');
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
        if (Array.isArray(products) && products.length > 0) {
            products.forEach(function(product) {
                productSelect.append($('<option>', {
                    value: product.ID,
                    text: product.name
                }));
            });
        }
    }

    // Adicionar uma variável global para armazenar o estado atual do fornecedor
    let currentSupplierData = null;

    // Manipulador de clique para o botão de edição
    $('.edit-supplier').click(function() {
        var supplierId = $(this).data('id');
        var supplierCNPJ = $(this).data('cnpj');
        
        // Carregar dados do fornecedor
        $.ajax({
            url: '/api/v1/suppliers-by-id?id=' + supplierId,
            type: 'GET',
            success: function(supplier) {
                // Armazenar os dados do fornecedor globalmente
                currentSupplierData = supplier;
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
                        // Carregar serviços e produtos da categoria selecionada
                        loadServicesForCategory(supplier.Category.ID, supplier.Services);
                    },
                    error: function(xhr) {
                        console.error('Erro ao carregar categorias:', xhr);
                    }
                });
                
                $('#editSupplierModal').modal('show');
            },
            error: function(xhr) {
                console.error('Erro ao carregar dados do fornecedor:', xhr);
            }
        });
    });

    // Função para carregar serviços de uma categoria
    function loadServicesForCategory(categoryId, supplierServices) {
        $.ajax({
            url: '/api/v1/services-by-category/' + categoryId,
            type: 'GET',
            success: function(services) {
                const servicesDiv = document.getElementById('editServices');
                servicesDiv.innerHTML = ''; // Limpar serviços existentes

                services.forEach(function(service) {
                    var serviceDiv = document.createElement('div');
                    serviceDiv.className = 'form-check';

                    var input = document.createElement('input');
                    input.className = 'form-check-input';
                    input.type = 'checkbox';
                    input.name = 'services[]';
                    input.value = service.ID;
                    input.id = 'editService' + service.ID;
                    input.setAttribute('data-category', categoryId);

                    // Verifica se o serviço está nos serviços do fornecedor e não está deletado
                    var isChecked = supplierServices && supplierServices.some(function(supplierService) {
                        return (supplierService.ServiceID === service.ID || 
                               supplierService.ID === service.ID) && 
                               !supplierService.DeletedAt;
                    });
                    
                    input.checked = isChecked;

                    var label = document.createElement('label');
                    label.className = 'form-check-label';
                    label.htmlFor = 'editService' + service.ID;
                    label.textContent = service.name;

                    serviceDiv.appendChild(input);
                    serviceDiv.appendChild(label);
                    servicesDiv.appendChild(serviceDiv);

                    // Se o serviço estiver marcado, carrega seus produtos
                    if (isChecked) {
                        loadProductsForService(service.ID);
                    }
                });
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
        
        // Limpar os produtos existentes quando mudar de categoria
        $('#editProducts').empty();
        
        // Carregar os serviços da categoria selecionada
        $.ajax({
            url: '/api/v1/services-by-category/' + categoryId,
            type: 'GET',
            success: function(services) {
                const servicesDiv = document.getElementById('editServices');
                servicesDiv.innerHTML = ''; // Limpar serviços existentes

                services.forEach(function(service) {
                    var serviceDiv = document.createElement('div');
                    serviceDiv.className = 'form-check';

                    var input = document.createElement('input');
                    input.className = 'form-check-input';
                    input.type = 'checkbox';
                    input.name = 'services[]';
                    input.value = service.ID;
                    input.id = 'editService' + service.ID;
                    input.setAttribute('data-category', categoryId);

                    // Verifica se o serviço estava selecionado anteriormente e pertence à categoria atual
                    if (currentSupplierData && currentSupplierData.Services) {
                        var isChecked = currentSupplierData.Services.some(function(supplierService) {
                            return (supplierService.ServiceID === service.ID || 
                                   supplierService.ID === service.ID) && 
                                   !supplierService.DeletedAt &&
                                   supplierService.Service.CategoryID === parseInt(categoryId); // Verifica se o serviço pertence à categoria atual
                        });
                        input.checked = isChecked;
                    }

                    var label = document.createElement('label');
                    label.className = 'form-check-label';
                    label.htmlFor = 'editService' + service.ID;
                    label.textContent = service.name;

                    serviceDiv.appendChild(input);
                    serviceDiv.appendChild(label);
                    servicesDiv.appendChild(serviceDiv);

                    // Se o serviço estiver marcado e pertencer à categoria atual, carrega seus produtos
                    if (input.checked) {
                        loadProductsForService(service.ID);
                    }
                });
            },
            error: function(xhr) {
                console.error('Erro ao carregar serviços:', xhr);
            }
        });
    });

    // Quando um serviço é selecionado/deselecionado
    $(document).on('change', 'input[name="services[]"]', function() {
        var categoryId = $('#editCategory').val();
        var selectedServices = $('input[name="services[]"]:checked').map(function() {
            return {
                id: $(this).val(),
                category: $(this).attr('data-category') || categoryId
            };
        }).get();

        // Se o serviço foi desmarcado, remove apenas os produtos desse serviço
        if (!this.checked) {
            var serviceId = $(this).val();
            $(`#editProducts [data-service="${serviceId}"]`).remove();
            return;
        }

        // Se o serviço foi marcado, adiciona seus produtos sem limpar os existentes
        if (this.checked) {
            var serviceId = $(this).val();
            loadProductsForService(serviceId);
        }
    });

    // Função para carregar produtos de um serviço específico
    function loadProductsForService(serviceId) {
        $.ajax({
            url: '/api/v1/products-by-service/' + serviceId,
            type: 'GET',
            success: function(products) {
                var productsDiv = $('#editProducts');
                
                products.forEach(function(product) {
                    // Verifica se o produto já existe no div antes de adicionar
                    if (!$('#product_' + product.ID).length) {
                        // Verifica se o produto estava selecionado anteriormente
                        var isChecked = false;
                        if (currentSupplierData && currentSupplierData.Products) {
                            isChecked = currentSupplierData.Products.some(function(supplierProduct) {
                                return supplierProduct.ProductID === product.ID && 
                                       !supplierProduct.DeletedAt;
                            });
                        }

                        productsDiv.append(`
                            <div class="form-check" data-service="${serviceId}">
                                <input class="form-check-input" type="checkbox" name="products[]" 
                                       value="${product.ID}" id="product_${product.ID}"
                                       ${isChecked ? 'checked' : ''}>
                                <label class="form-check-label" for="product_${product.ID}">
                                    ${product.name}
                                </label>
                            </div>
                        `);
                    }
                });
            },
            error: function(xhr) {
                console.error('Erro ao carregar produtos:', xhr);
            }
        });
    }

    // Manipulador de clique para salvar alterações
    $('#saveSupplierChanges').click(function() {
        var supplierId = $('#editSupplierId').val();
        var categoryId = $('#editCategory').val();
        var selectedServices = $('input[name="services[]"]:checked').map(function() {
            return parseInt(this.value);
        }).get();
        var selectedProducts = $('input[name="products[]"]:checked').map(function() {
            return parseInt(this.value);
        }).get();

        var data = {
            category_id: parseInt(categoryId),
            service_ids: selectedServices,
            product_ids: selectedProducts
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

    // Função para carregar produtos dos serviços selecionados
    function loadProductsForServices(supplierServices, supplier) {
        if (!supplierServices || supplierServices.length === 0) return;

        // Guarda os dados do fornecedor em uma variável global
        window.supplierData = supplier;
        console.log('Dados do fornecedor carregados:', supplier);
        console.log('Produtos do fornecedor:', supplier.Products);

        $('#editProducts').empty();
        
        supplierServices.forEach(function(service) {
            var serviceId = service.ServiceID || service.ID;
            if (serviceId) {
                loadProductsForService(serviceId);
            }
        });
    }

    // Quando a categoria é alterada no modal de edição
    $('#editSupplierModal select[name="category"]').change(function() {
        var categoryId = $(this).val();
        
        // Limpa os checkboxes de produtos
        $('#editSupplierModal input[name="products[]"]').prop('checked', false);
        
        // Atualiza os serviços baseado na categoria selecionada
        $.ajax({
            url: '/api/v1/services/by-category/' + categoryId,
            type: 'GET',
            success: function(services) {
                // Atualiza a lista de serviços
                var servicesList = $('#editSupplierModal .services-list');
                servicesList.empty();
                
                services.forEach(function(service) {
                    servicesList.append(`
                        <div class="form-check">
                            <input class="form-check-input service-checkbox" type="checkbox" 
                                   name="services[]" value="${service.id}" 
                                   id="service_${service.id}">
                            <label class="form-check-label" for="service_${service.id}">
                                ${service.name}
                            </label>
                        </div>
                    `);
                });

                // Atualiza os produtos baseado nos serviços disponíveis
                updateProductsList(services.map(s => s.id));
            },
            error: function(xhr, status, error) {
                console.error('Erro ao carregar serviços:', error);
            }
        });
    });

    // Função para atualizar a lista de produtos baseado nos serviços
    function updateProductsList(serviceIds) {
        if (!serviceIds || serviceIds.length === 0) {
            $('#editSupplierModal .products-list').empty();
            return;
        }

        $.ajax({
            url: '/api/v1/products/by-services',
            type: 'POST',
            data: JSON.stringify({ service_ids: serviceIds }),
            contentType: 'application/json',
            success: function(products) {
                var productsList = $('#editSupplierModal .products-list');
                productsList.empty();
                
                products.forEach(function(product) {
                    productsList.append(`
                        <div class="form-check">
                            <input class="form-check-input product-checkbox" type="checkbox" 
                                   name="products[]" value="${product.id}" 
                                   id="product_${product.id}">
                            <label class="form-check-label" for="product_${product.id}">
                                ${product.name}
                            </label>
                        </div>
                    `);
                });
            },
            error: function(xhr, status, error) {
                console.error('Erro ao carregar produtos:', error);
            }
        });
    }
});
