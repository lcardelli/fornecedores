$(document).ready(function() {
    // Inicializa os selects com Select2
    $('.select2').select2({
        theme: 'bootstrap-5'
    });

    // Vincula os selects de nome e CNPJ
    $('#supplier_name, #supplier_cnpj').on('change', function() {
        var selectedValue = $(this).val();
        if (this.id === 'supplier_name') {
            $('#supplier_cnpj').val(selectedValue).trigger('change');
        } else {
            $('#supplier_name').val(selectedValue).trigger('change');
        }
    });

    // Carrega os serviços quando uma categoria é selecionada
    $('#category_id').on('change', function() {
        var categoryId = $(this).val();
        if (categoryId) {
            $.ajax({
                url: '/api/v1/services-by-category/' + categoryId,
                type: 'GET',
                success: function(data) {
                    var serviceSelect = $('#service_ids');
                    serviceSelect.empty();
                    if (data.length === 0) {
                        console.log('Nenhum serviço encontrado para esta categoria');
                        serviceSelect.append(new Option('Nenhum serviço disponível', ''));
                    } else {
                        $.each(data, function(index, service) {
                            var option = new Option(service.name, service.ID);
                            serviceSelect.append(option);
                        });
                    }
                    serviceSelect.trigger('change');
                    $('#product_ids').empty().trigger('change');
                },
                error: function(xhr, status, error) {
                    console.error('Erro ao carregar serviços:', error);
                    Swal.fire({
                        icon: 'error',
                        title: 'Erro',
                        text: 'Não foi possível carregar os serviços. Por favor, tente novamente.'
                    });
                }
            });
        } else {
            $('#service_ids').empty().trigger('change');
            $('#product_ids').empty().trigger('change');
        }
    });

    // Carrega os produtos quando serviços são selecionados
    $('#service_ids').on('change', function() {
        var serviceIds = $(this).val();
        console.log('Serviços selecionados:', serviceIds);

        if (serviceIds && serviceIds.length > 0) {
            var productSelect = $('#product_ids');
            productSelect.empty();

            serviceIds.forEach(function(serviceId) {
                $.ajax({
                    url: '/api/v1/products-by-service/' + serviceId,
                    type: 'GET',
                    success: function(products) {
                        console.log('Produtos recebidos para serviço ' + serviceId + ':', products);
                        
                        if (Array.isArray(products) && products.length > 0) {
                            products.forEach(function(product) {
                                if (productSelect.find("option[value='" + product.ID + "']").length === 0) {
                                    productSelect.append(new Option(product.name, product.ID));
                                }
                            });
                        }
                        productSelect.trigger('change');
                    },
                    error: function(xhr, status, error) {
                        console.error('Erro na requisição para serviço ' + serviceId + ':', error);
                        console.log('Status:', status);
                        console.log('Resposta:', xhr.responseText);
                        Swal.fire({
                            icon: 'error',
                            title: 'Erro',
                            text: 'Não foi possível carregar os produtos para o serviço selecionado.'
                        });
                    }
                });
            });
        } else {
            $('#product_ids').empty().trigger('change');
        }
    });

    // Submissão do formulário
    $('#supplierForm').on('submit', function(e) {
        e.preventDefault();

        var formData = {
            supplier_cnpj: $('#supplier_cnpj').val(),
            category_id: parseInt($('#category_id').val()),
            service_ids: $('#service_ids').val().map(Number).filter(id => !isNaN(id) && id > 0),
            product_ids: $('#product_ids').val() ? $('#product_ids').val().map(Number).filter(id => !isNaN(id) && id > 0) : []
        };

        console.log('Dados do formulário:', formData);

        // Validações adicionais
        if (!formData.service_ids.length) {
            Swal.fire({
                icon: 'error',
                title: 'Erro',
                text: 'Selecione pelo menos um serviço.'
            });
            return;
        }

        if (formData.product_ids.some(id => id <= 0)) {
            Swal.fire({
                icon: 'error',
                title: 'Erro',
                text: 'IDs de produtos inválidos detectados.'
            });
            return;
        }

        $.ajax({
            url: '/api/v1/suppliers',
            type: 'POST',
            data: JSON.stringify(formData),
            contentType: 'application/json',
            success: function(response) {
                console.log('Resposta de sucesso:', response);
                Swal.fire({
                    icon: 'success',
                    title: 'Sucesso!',
                    text: 'Fornecedor cadastrado com sucesso.',
                    confirmButtonText: 'OK'
                }).then((result) => {
                    if (result.isConfirmed) {
                        $('#supplierForm')[0].reset();
                        $('.select2').val(null).trigger('change');
                    }
                });
            },
            error: function(xhr, status, error) {
                console.error('Erro ao cadastrar fornecedor:', error);
                console.log('Status:', status);
                console.log('Resposta do servidor:', xhr.responseText);
                
                let errorMessage = 'Não foi possível cadastrar o fornecedor. ';
                if (xhr.responseJSON && xhr.responseJSON.error) {
                    errorMessage += xhr.responseJSON.error;
                } else {
                    errorMessage += 'Por favor, tente novamente.';
                }

                Swal.fire({
                    icon: 'error',
                    title: 'Erro',
                    text: errorMessage
                });
            }
        });
    });
});
