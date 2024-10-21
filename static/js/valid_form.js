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
                            serviceSelect.append(new Option(service.name, service.id)); // Usar service.id como valor
                        });
                    }
                    serviceSelect.trigger('change');
                },
                error: function(xhr, status, error) {
                    console.error('Erro ao carregar serviços:', error);
                    console.log('Resposta do servidor:', xhr.responseText);
                    Swal.fire({
                        icon: 'error',
                        title: 'Erro',
                        text: 'Não foi possível carregar os serviços. Por favor, tente novamente.'
                    });
                }
            });
        } else {
            $('#service_ids').empty().trigger('change');
        }
    });

    // Submissão do formulário
    $('#supplierForm').on('submit', function(e) {
        e.preventDefault();

        var serviceNames = $('#service_ids').find("option:selected").map(function() {
            return $(this).text();
        }).get();

        var formData = {
            supplier_cnpj: $('#supplier_cnpj').val(),
            category_id: $('#category_id').val(),
            service_ids: serviceNames
        };

        console.log('Dados do formulário:', formData);

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
