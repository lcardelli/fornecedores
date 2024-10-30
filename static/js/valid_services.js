$(document).ready(function() {
    var allServices = [];

    loadServices();

    $('#newServiceBtn').click(function() {
        $('#formSection').slideDown();
        resetForm();
    });

    $('#serviceForm').submit(function(e) {
        e.preventDefault();
        var serviceId = $('#serviceId').val();
        var serviceName = $('#serviceName').val().trim();
        var categoryId = $('#serviceCategory').val();

        if (!serviceName) {
            Swal.fire({
                icon: 'error',
                title: 'Erro!',
                text: 'O nome do serviço não pode estar vazio'
            });
            return;
        }

        var url = serviceId ? `/api/v1/services/${serviceId}` : '/api/v1/services';
        var method = serviceId ? 'PUT' : 'POST';

        var data = {
            id: serviceId ? parseInt(serviceId) : null,
            name: serviceName,
            category_id: parseInt(categoryId)
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
                    text: serviceId ? 'Serviço atualizado com sucesso!' : 'Serviço cadastrado com sucesso!',
                }).then(() => {
                    $('#formSection').slideUp();
                    resetForm();
                    loadServices();
                });
            },
            error: function(xhr, status, error) {
                let errorMessage = 'Erro ao processar serviço';
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

    function loadServices() {
        $.ajax({
            url: '/api/v1/service-list',
            type: 'GET',
            success: function(services) {
                allServices = services;
                filterServices();
            },
            error: function(xhr, status, error) {
                console.error('Erro ao carregar serviços:', xhr.responseText);
            }
        });
    }

    function renderServices(services) {
        var list = $('#servicesList');
        list.empty();
        
        if (services.length === 0) {
            list.html('<div class="empty-state"><i class="fas fa-info-circle"></i><h4>Nenhum serviço encontrado</h4><p>Não há serviços cadastrados com os filtros atuais.</p></div>');
            return;
        }

        var table = `
            <table class="table table-hover mb-0">
                <thead>
                    <tr>
                        <th width="40px"></th>
                        <th>Nome do Serviço</th>
                        <th>Categoria</th>
                        <th width="120px">Ações</th>
                    </tr>
                </thead>
                <tbody>
        `;
        
        services.forEach(function(service) {
            var category = $('#serviceCategory option[value="' + service.category_id + '"]').text();
            table += `
                <tr>
                    <td>
                        <input type="checkbox" class="service-checkbox" data-id="${service.id}">
                    </td>
                    <td>${service.name}</td>
                    <td>${category}</td>
                    <td>
                        <div class="btn-group-actions">
                            <button class="btn btn-sm btn-warning edit-btn" data-id="${service.id}" data-name="${service.name}" data-category="${service.category_id}">
                                <i class="fas fa-edit"></i>
                            </button>
                            <button class="btn btn-sm btn-danger delete-btn" data-id="${service.id}">
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
    }

    function setupEditButtons() {
        $('.edit-btn').click(function() {
            var id = $(this).data('id');
            var name = $(this).data('name');
            var categoryId = $(this).data('category');
            $('#serviceId').val(id);
            $('#serviceName').val(name);
            $('#serviceCategory').val(categoryId);
            $('#submitBtn').html('<i class="fas fa-save mr-2"></i>Atualizar');
            $('#formSection').slideDown();
            $('#cancelBtn').show();
        });
    }

    // Adicione este evento para atualizar a lista quando a categoria mudar
    $('#serviceCategory').change(function() {
        var categoryId = $(this).val();
        if (!categoryId) return;
        
        // Atualiza a lista de serviços baseado na categoria selecionada
        filterServices();
    });

    function setupDeleteButtons() {
        $('.delete-btn').click(function() {
            var id = $(this).data('id');
            console.log("ID do serviço a ser deletado:", id);
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
                    deleteService(id);
                }
            });
        });
    }

    function deleteService(id) {
        $.ajax({
            url: `/api/v1/services/${id}`,
            type: 'DELETE',
            success: function(response) {
                Swal.fire(
                    'Deletado!',
                    'O serviço foi deletado com sucesso.',
                    'success'
                );
                // Remover o serviço excluído da lista local
                allServices = allServices.filter(service => service.id !== id);
                // Atualizar a visualização
                filterServices();
            },
            error: function(xhr, status, error) {
                console.error('Erro ao deletar serviço:', xhr.responseText);
                Swal.fire(
                    'Erro!',
                    'Não foi possível deletar o serviço: ' + (xhr.responseJSON ? xhr.responseJSON.error : error),
                    'error'
                );
            }
        });
    }

    function resetForm() {
        $('#serviceId').val('');
        $('#serviceName').val('');
        $('#serviceCategory').val($('#serviceCategory option:first').val());
        $('#submitBtn').html('<i class="fas fa-save mr-2"></i>Salvar');
        $('#cancelBtn').show();
    }

    function setupCheckboxEvents() {
        $('.service-checkbox').change(function() {
            updateDeleteSelectedButton();
        });
    }

    function updateDeleteSelectedButton() {
        var checkedCount = $('.service-checkbox:checked').length;
        $('#deleteSelectedBtn').toggle(checkedCount > 0);
    }

    $('#selectAllBtn').click(function() {
        var isAllSelected = $('.service-checkbox:checked').length === $('.service-checkbox').length;
        $('.service-checkbox').prop('checked', !isAllSelected);
        // Atualiza o texto do botão baseado no estado
        $(this).html(
            isAllSelected ? 
            '<i class="fas fa-check-square mr-2"></i>Selecionar Todos' : 
            '<i class="fas fa-square mr-2"></i>Desmarcar Todos'
        );
        updateDeleteSelectedButton();
    });

    $('#deleteSelectedBtn').click(function() {
        var selectedIds = $('.service-checkbox:checked').map(function() {
            return $(this).data('id');
        }).get();

        if (selectedIds.length === 0) return;

        Swal.fire({
            title: 'Tem certeza?',
            text: `Você está prestes a excluir ${selectedIds.length} serviço(s). Esta ação não pode ser revertida!`,
            icon: 'warning',
            showCancelButton: true,
            confirmButtonColor: '#3085d6',
            cancelButtonColor: '#d33',
            confirmButtonText: 'Sim, excluir!',
            cancelButtonText: 'Cancelar'
        }).then((result) => {
            if (result.isConfirmed) {
                deleteMultipleServices(selectedIds);
            }
        });
    });

    function deleteMultipleServices(ids) {
        $.ajax({
            url: '/api/v1/services/batch',
            type: 'DELETE',
            data: JSON.stringify({ ids: ids }),
            contentType: 'application/json',
            success: function(response) {
                Swal.fire(
                    'Excluídos!',
                    'Os serviços selecionados foram excluídos com sucesso.',
                    'success'
                );
                // Remover os serviços excluídos da lista local
                allServices = allServices.filter(service => !ids.includes(service.id));
                
                // Resetar os botões
                $('#deleteSelectedBtn').hide();
                $('#selectAllBtn').hide();
                
                // Aplicar o filtro atual para atualizar a visualização
                filterServices();
            },
            error: function(xhr, status, error) {
                console.error('Erro ao excluir serviços:', xhr.responseText);
                Swal.fire(
                    'Erro!',
                    'Não foi possível excluir os serviços selecionados: ' + (xhr.responseJSON ? xhr.responseJSON.error : error),
                    'error'
                );
            }
        });
    }

    $('#serviceSearch, #categoryFilter').on('input change', filterServices);

    function filterServices() {
        var searchTerm = $('#serviceSearch').val().toLowerCase();
        var categoryId = $('#categoryFilter').val();

        var filteredServices = allServices.filter(function(service) {
            var matchesSearch = service.name.toLowerCase().includes(searchTerm);
            var matchesCategory = !categoryId || service.category_id == categoryId;
            return matchesSearch && matchesCategory;
        });

        renderServices(filteredServices);
    }
});