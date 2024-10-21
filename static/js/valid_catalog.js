$(document).ready(function() {
    $('#category').change(function() {
        var categoryId = $(this).val();
        $('#service').val('');
        if (categoryId) {
            $.ajax({
                url: '/api/v1/services-by-category/' + categoryId,
                type: 'GET',
                success: function(services) {
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
                },
                error: function() {
                    console.error('Erro ao carregar serviços');
                }
            });
        } else {
            // Se nenhuma categoria for selecionada, carregue todos os serviços
            $.ajax({
                url: '/api/v1/services',
                type: 'GET',
                success: function(services) {
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
                },
                error: function() {
                    console.error('Erro ao carregar serviços');
                }
            });
        }
        $('#filterForm').submit();
    });

    $('#service').change(function() {
        $('#filterForm').submit();
    });

    // Manipulador de clique para o botão de edição
    $('.edit-supplier').click(function() {
        var supplierId = $(this).data('id');
        var supplierCNPJ = $(this).data('cnpj');
        
        // Carregar dados do fornecedor
        $.ajax({
            url: '/api/v1/suppliers?id=' + supplierId,
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
                                selected: category.ID === supplier.CategoryID
                            }));
                        });
                    }
                });
                
                // Preencher serviços
                $.ajax({
                    url: '/api/v1/service-list',
                    type: 'GET',
                    success: function(services) {
                        const servicesDiv = document.getElementById('editServices')

                        
                        $.each(services, function(i, service) {
                            var serviceDiv = document.createElement('div');
                            serviceDiv.className = 'form-check';

                            var input = document.createElement('input');
                            input.className = 'form-check-input';
                            input.type = 'checkbox';
                            input.name = 'services[]';
                            input.value = service.id;
                            input.id = 'editService' + service.id;

                            var label = document.createElement('label');
                            label.className = 'form-check-label';
                            label.htmlFor = 'editService' + service.id;
                            label.textContent = service.name;

                            serviceDiv.appendChild(input);
                            serviceDiv.appendChild(label);

                            servicesDiv.appendChild(serviceDiv);
                        });
                    }
                });
                
                $('#editSupplierModal').modal('show');
            },
            error: function(xhr, status, error) {
                alert('Erro ao carregar dados do fornecedor: ' + error);
            }
        });
    });

    // Manipulador de clique para salvar alterações
    $('#saveSupplierChanges').click(function() {
        var supplierId = $('#editSupplierId').val();
        var categoryId = $('#editCategory').val();
        var selectedServices = $('input[name="services[]"]:checked').map(function() {
            return { service_id: parseInt(this.value) };
        }).get();

        var data = {
            category_id: parseInt(categoryId),
            services: selectedServices
        };

        console.log(data);

        $.ajax({
            url: '/api/v1/suppliers/' + supplierId,  // Note a mudança aqui
            type: 'PUT',
            contentType: 'application/json',
            data: JSON.stringify(data),
            success: function(result) {
                alert('Fornecedor atualizado com sucesso!');
                $('#editSupplierModal').modal('hide');
                location.reload();
            },
            error: function(xhr, status, error) {
                alert('Erro ao atualizar fornecedor: ' + error);
            }
        });
    });

    // Manipulador de clique para o botão de exclusão
    $('.delete-supplier').click(function() {
        var supplierId = $(this).data('id');
        var supplierCNPJ = $(this).data('cnpj');
        if (confirm('Tem certeza que deseja deletar este fornecedor (CNPJ: ' + supplierCNPJ + ') do catálogo?')) {
            $.ajax({
                url: '/api/v1/suppliers/' + supplierId,  // Note a mudança aqui
                type: 'DELETE',
                success: function(result) {
                    alert('Fornecedor removido com sucesso!');
                    location.reload();
                },
                error: function(xhr, status, error) {
                    alert('Erro ao deletar fornecedor: ' + error);
                }
            });
        }
    });
});