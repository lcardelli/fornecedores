$(document).ready(function() {
    var allServices = [];

    loadServices();

    $('#serviceForm').submit(function(e) {
        e.preventDefault();
        var serviceId = $('#serviceId').val();
        var serviceName = $('#serviceName').val();
        var categoryId = $('#serviceCategory').val();
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
                });
                resetForm();
                loadServices();
            },
            error: function(xhr, status, error) {
                console.error('Erro ao processar serviço:', xhr.responseText);
                Swal.fire({
                    icon: 'error',
                    title: 'Erro!',
                    text: 'Erro ao processar serviço: ' + error,
                });
            }
        });
    });

    $('#cancelBtn').click(function() {
        resetForm();
    });

    function loadServices() {
        $.ajax({
            url: '/api/v1/service-list',
            type: 'GET',
            success: function(services) {
                allServices = services;
                // Após carregar os serviços, aplica o filtro atual
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
            showEmptyState();
        } else {
            services.forEach(function(service) {
                list.append(
                    `<div class="service-item" data-category="${service.category_id}">
                        <div class="d-flex align-items-center">
                            <input type="checkbox" class="service-checkbox mr-2" data-id="${service.id}">
                            <span>${service.name}</span>
                        </div>
                        <div>
                            <button class="btn btn-sm btn-warning mr-2 edit-btn" data-id="${service.id}" data-name="${service.name}" data-category="${service.category_id}"><i class="fas fa-edit"></i></button>
                            <button class="btn btn-sm btn-danger delete-btn" data-id="${service.id}"><i class="fas fa-trash"></i></button>
                        </div>
                    </div>`
                );
            });
            setupEditButtons();
            setupDeleteButtons();
            setupCheckboxEvents();
        }
    }

    function showEmptyState() {
        var list = $('#servicesList');
        var categoryId = $('#categoryFilter').val();
        
        // Só mostra o estado vazio se nenhuma categoria estiver selecionada
        if (categoryId === '') {
            list.html(`
                <div class="empty-state">
                    <i class="fas fa-filter"></i>
                    <h4>Nenhuma categoria selecionada</h4>
                    <p>Selecione uma categoria ou digite na busca para ver os serviços.</p>
                </div>
            `);
        } else {
            list.html(`
                <div class="empty-state">
                    <i class="fas fa-info-circle"></i>
                    <h4>Nenhum serviço encontrado</h4>
                    <p>Não há serviços cadastrados nesta categoria.</p>
                </div>
            `);
        }
    }

    $('#serviceSearch, #categoryFilter').on('input change', filterServices);

    function filterServices() {
        var searchTerm = $('#serviceSearch').val().toLowerCase();
        var categoryId = $('#categoryFilter').val();

        // Controla a visibilidade do botão "Selecionar Todos"
        $('#selectAllBtn').toggle(categoryId !== '');

        // Filtra os serviços mesmo se não houver termo de busca
        var filteredServices = allServices.filter(function(service) {
            var matchesSearch = searchTerm === '' || service.name.toLowerCase().includes(searchTerm);
            var matchesCategory = categoryId === '' || service.category_id == categoryId;
            return matchesSearch && matchesCategory;
        });

        // Sempre renderiza os serviços se uma categoria estiver selecionada
        if (categoryId !== '') {
            renderServices(filteredServices);
        } else if (searchTerm !== '') {
            renderServices(filteredServices);
        } else {
            showEmptyState();
        }
    }

    function setupEditButtons() {
        $('.edit-btn').click(function() {
            var id = $(this).data('id');
            var name = $(this).data('name');
            var categoryId = $(this).data('category');
            $('#serviceId').val(id);
            $('#serviceName').val(name);
            $('#serviceCategory').val(categoryId);
            $('#submitBtn').html('<i class="fas fa-save mr-2"></i>Atualizar Serviço');
            $('#cancelBtn').show();
        });
    }

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
        $('#serviceCategory').val('');
        $('#submitBtn').html('<i class="fas fa-save mr-2"></i>Cadastrar Serviço');
        $('#cancelBtn').hide();
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
});