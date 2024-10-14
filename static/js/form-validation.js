document.addEventListener('DOMContentLoaded', function() {
    const form = document.getElementById('supplierForm');
    
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
                    alert('Erro ao cadastrar fornecedor: ' + data.error);
                } else {
                    alert('Fornecedor cadastrado com sucesso!');
                    form.reset();
                    $('.select2').val(null).trigger('change');
                }
            })
            .catch((error) => {
                console.error('Error:', error);
                alert('Erro ao cadastrar fornecedor. Por favor, tente novamente.');
            });
        }
    });

    // Adicione aqui as funções de validação existentes
    function validateForm() {
        // Implemente a lógica de validação aqui
        return true; // Retorne true se o formulário for válido, false caso contrário
    }
});