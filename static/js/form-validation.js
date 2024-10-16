document.addEventListener('DOMContentLoaded', function() {
    const form = document.getElementById('supplierForm');
    
    $('.select2').select2({
        theme: 'bootstrap-5',
        placeholder: "Selecione um ou mais serviços",
        allowClear: true
    });

    $('#phone').mask('(00) 00000-0000');
    $('#cnpj').mask('00.000.000/0000-00');

    $('#category_id').change(function() {
        var categoryId = $(this).val();
        if (categoryId) {
            $.ajax({
                url: '/api/services-by-category/' + categoryId,
                type: 'GET',
                success: function(data) {
                    var serviceSelect = $('#service_ids');
                    serviceSelect.empty();
                    $.each(data, function(index, service) {
                        serviceSelect.append(new Option(service.Name, service.ID));
                    });
                    serviceSelect.trigger('change');
                }
            });
        }
    });

    form.addEventListener('submit', function(event) {
        event.preventDefault();
        if (validateForm()) {
            const formData = new FormData(form);
            const jsonData = {
                name: formData.get('name'),
                email: formData.get('email'),
                phone: formData.get('phone'),
                address: formData.get('address'),
                cnpj: formData.get('cnpj'),
                category_id: parseInt(formData.get('category_id')),
                service_ids: Array.from(formData.getAll('service_ids[]')).map(id => parseInt(id)),
                description: formData.get('description')
            };

            fetch('/api/v1/suppliers', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(jsonData),
            })
            .then(response => response.json())
            .then(data => {
                if (data.error) {
                    Swal.fire({
                        icon: 'error',
                        title: 'Erro!',
                        text: 'Erro ao cadastrar fornecedor: ' + data.error,
                    });
                } else {
                    Swal.fire({
                        icon: 'success',
                        title: 'Sucesso!',
                        text: 'Fornecedor cadastrado com sucesso!',
                    }).then((result) => {
                        if (result.isConfirmed) {
                            form.reset();
                            $('.select2').val(null).trigger('change');
                        }
                    });
                }
            })
            .catch((error) => {
                console.error('Error:', error);
                Swal.fire({
                    icon: 'error',
                    title: 'Erro!',
                    text: 'Erro ao cadastrar fornecedor. Por favor, tente novamente.',
                });
            });
        }
    });

    // Adicione aqui as funções de validação existentes
    function validateForm() {
        // Implemente a lógica de validação aqui
        return true; // Retorne true se o formulário for válido, false caso contrário
    }
});
