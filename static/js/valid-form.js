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
                    $.each(data, function(index, service) {
                        serviceSelect.append(new Option(service.name, service.ID));
                    });
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

        var formData = {
            supplier_cnpj: $('#supplier_cnpj').val(),
            category_id: $('#category_id').val(),
            service_ids: $('#service_ids').val()
        };

        $.ajax({
            url: '/api/v1/suppliers',
            type: 'POST',
            contentType: 'application/json',
            data: JSON.stringify(formData),
            success: function(response) {
                Swal.fire({
                    icon: 'success',
                    title: 'Sucesso!',
                    text: 'Fornecedor cadastrado com sucesso.',
                    confirmButtonText: 'OK'
                }).then((result) => {
                    if (result.isConfirmed) {
                        // Limpa o formulário ou redireciona para outra página
                        $('#supplierForm')[0].reset();
                        $('.select2').val(null).trigger('change');
                    }
                });
            },
            error: function(xhr, status, error) {
                console.error('Erro ao cadastrar fornecedor:', error);
                Swal.fire({
                    icon: 'error',
                    title: 'Erro',
                    text: 'Não foi possível cadastrar o fornecedor. Por favor, tente novamente.'
                });
            }
        });
    });
});
