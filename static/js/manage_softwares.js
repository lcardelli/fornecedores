$(document).ready(function() {
    // Handler para exclusão em massa
    $('#deleteSelected').click(function() {
        const selectedIds = $('.software-checkbox:checked').map(function() {
            return $(this).val();
        }).get();

        if (selectedIds.length === 0) return;

        Swal.fire({
            title: 'Tem certeza?',
            text: `Você está prestes a excluir ${selectedIds.length} software(s). Esta ação não poderá ser revertida!`,
            icon: 'warning',
            showCancelButton: true,
            confirmButtonColor: '#d33',
            cancelButtonColor: '#3085d6',
            confirmButtonText: 'Sim, deletar!',
            cancelButtonText: 'Cancelar'
        }).then((result) => {
            if (result.isConfirmed) {
                const deletePromises = selectedIds.map(id => 
                    $.ajax({
                        url: `/api/v1/licenses/software/${id}`,
                        type: 'DELETE'
                    })
                );

                Promise.all(deletePromises)
                    .then(() => {
                        Swal.fire(
                            'Deletados!',
                            'Os softwares foram deletados com sucesso.',
                            'success'
                        ).then(() => {
                            location.reload();
                        });
                    })
                    .catch((error) => {
                        Swal.fire({
                            icon: 'error',
                            title: 'Erro!',
                            text: 'Erro ao deletar softwares: ' + error.responseText
                        });
                    });
            }
        });
    });

    // Seleção de todos os checkboxes
    $("#selectAll").change(function() {
        $(".software-checkbox").prop('checked', $(this).prop("checked"));
        updateDeleteButtonVisibility();
    });

    // Atualiza o checkbox "Selecionar Todos" quando os individuais são clicados
    $(".software-checkbox").change(function() {
        updateSelectAllCheckbox();
        updateDeleteButtonVisibility();
    });

    function updateSelectAllCheckbox() {
        var totalCheckboxes = $(".software-checkbox").length;
        var checkedCheckboxes = $(".software-checkbox:checked").length;
        $("#selectAll").prop('checked', totalCheckboxes === checkedCheckboxes);
    }

    function updateDeleteButtonVisibility() {
        var checkedCheckboxes = $(".software-checkbox:checked").length;
        if (checkedCheckboxes > 0) {
            $("#deleteSelected").show();
        } else {
            $("#deleteSelected").hide();
        }
    }

    // Handler para salvar software
    $('#saveSoftware').click(function() {
        const formData = new FormData($('#softwareForm')[0]);
        const data = {};
        
        // Converte os dados do formulário para o formato correto
        formData.forEach((value, key) => {
            if (key === 'id' && value) {
                data[key] = parseInt(value); // Converte ID para número
            } else if (value) {
                data[key] = value;
            }
        });

        // Remove o ID se for um novo registro
        if (!data.id) {
            delete data.id;
        }
        
        $.ajax({
            url: `/api/v1/licenses/software${data.id ? '/' + data.id : ''}`,
            type: data.id ? 'PUT' : 'POST',
            contentType: 'application/json',
            data: JSON.stringify(data),
            success: function(response) {
                $('#addSoftwareModal').modal('hide');
                Swal.fire({
                    icon: 'success',
                    title: 'Sucesso!',
                    text: `Software ${data.id ? 'atualizado' : 'cadastrado'} com sucesso!`
                }).then(() => {
                    location.reload();
                });
            },
            error: function(xhr) {
                Swal.fire({
                    icon: 'error',
                    title: 'Erro!',
                    text: `Erro ao ${data.id ? 'atualizar' : 'cadastrar'} software: ` + xhr.responseText
                });
            }
        });
    });

    // Handler para editar software
    $('.edit-software').click(function() {
        const softwareId = $(this).data('id');
        
        $.get(`/api/v1/licenses/software/${softwareId}`, function(software) {
            $('#softwareId').val(software.ID);
            $('#softwareForm [name="name"]').val(software.Name);
            $('#softwareForm [name="publisher"]').val(software.Publisher);
            $('#softwareForm [name="description"]').val(software.Description);
            
            $('#modalTitle').text('Editar Software');
            $('#addSoftwareModal').modal('show');
        });
    });

    // Handler para deletar software
    $('.delete-software').click(function() {
        const softwareId = $(this).data('id');
        const row = $(this).closest('tr');

        Swal.fire({
            title: 'Tem certeza?',
            text: "Esta ação não poderá ser revertida!",
            icon: 'warning',
            showCancelButton: true,
            confirmButtonColor: '#3085d6',
            cancelButtonColor: '#d33',
            confirmButtonText: 'Sim, deletar!',
            cancelButtonText: 'Cancelar'
        }).then((result) => {
            if (result.isConfirmed) {
                $.ajax({
                    url: `/api/v1/licenses/software/${softwareId}`,
                    type: 'DELETE',
                    success: function() {
                        row.remove();
                        Swal.fire(
                            'Deletado!',
                            'O software foi deletado com sucesso.',
                            'success'
                        );
                    },
                    error: function(xhr) {
                        Swal.fire({
                            icon: 'error',
                            title: 'Erro!',
                            text: 'Erro ao deletar software: ' + xhr.responseText
                        });
                    }
                });
            }
        });
    });

    // Limpa o formulário quando o modal é fechado
    $('#addSoftwareModal').on('hidden.bs.modal', function() {
        $('#softwareForm')[0].reset();
        $('#softwareId').val('');
        $('#modalTitle').text('Novo Software');
    });

    // Função para aplicar os filtros
    function applyFilters() {
        const nameFilter = $('#filterName').val().toLowerCase();
        const publisherFilter = $('#filterPublisher').val().toLowerCase();
        const licensesFilter = $('#filterLicenses').val();

        $('#softwareTable tr').each(function() {
            const row = $(this);
            const name = row.find('td:eq(1)').text().toLowerCase();
            const publisher = row.find('td:eq(2)').text().toLowerCase();
            const licenses = parseInt(row.find('td:eq(4)').text()) || 0; // Número de licenças

            const matchesName = !nameFilter || name.includes(nameFilter);
            const matchesPublisher = !publisherFilter || publisher.includes(publisherFilter);
            const matchesLicenses = !licensesFilter || 
                (licensesFilter === 'com' && licenses > 0) || 
                (licensesFilter === 'sem' && licenses === 0);

            if (matchesName && matchesPublisher && matchesLicenses) {
                row.show();
            } else {
                row.hide();
            }
        });

        updateTableStatus();
    }

    // Função para atualizar status da tabela
    function updateTableStatus() {
        const visibleRows = $('#softwareTable tr:visible').filter(function() {
            return $(this).find('td').length > 0;
        }).length;

        if (visibleRows === 0) {
            if ($('#noResultsMessage').length === 0) {
                $('#softwareTable').after(`
                    <tr id="noResultsMessage">
                        <td colspan="6" class="text-center">
                            <div class="empty-state">
                                <i class="fas fa-cube fa-3x mb-3"></i>
                                <p class="h5">Nenhum software encontrado</p>
                                <p class="text-muted">Tente ajustar seus filtros de busca</p>
                            </div>
                        </td>
                    </tr>
                `);
            }
        } else {
            $('#noResultsMessage').remove();
        }
    }

    // Event listeners para os filtros
    $('#filterName, #filterPublisher').on('input', applyFilters);
    $('#filterLicenses').on('change', applyFilters);

    // Limpar filtros
    $('#clearFilters').click(function() {
        $('#filterName, #filterPublisher').val('');
        $('#filterLicenses').val('');
        applyFilters();
    });

    // Estilização do select com Select2
    $('#filterLicenses').select2({
        width: '100%',
        placeholder: 'Selecione...',
        allowClear: true
    });
});