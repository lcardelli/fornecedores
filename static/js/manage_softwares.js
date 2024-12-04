$(document).ready(function() {
    // Inicializa o Select2 para o filtro de nome
    $('.select2-software').select2({
        width: '100%',
        placeholder: 'Selecione um software...',
        allowClear: true,
        theme: 'default'
    });

    // Inicializa o Select2 para o filtro de fabricante
    $('.select2-publisher').select2({
        width: '100%',
        placeholder: 'Selecione um fabricante...',
        allowClear: true,
        theme: 'default'
    });

    // Atualiza a função de filtro
    function applyFilters() {
        const nameFilter = $('#filterName').val();
        const publisherFilter = $('#filterPublisher').val();
        const licensesFilter = $('#filterLicenses').val();

        let visibleCount = 0;
        $('#softwareTable tr').each(function() {
            const row = $(this);
            const name = row.find('td:eq(1)').text();
            const publisher = row.find('td:eq(2)').text();
            const licenses = parseInt(row.find('td:eq(4)').text()) || 0;

            const matchesName = !nameFilter || name === nameFilter;
            const matchesPublisher = !publisherFilter || publisher === publisherFilter;
            const matchesLicenses = !licensesFilter || 
                (licensesFilter === 'com' && licenses > 0) || 
                (licensesFilter === 'sem' && licenses === 0);

            if (matchesName && matchesPublisher && matchesLicenses) {
                // Remove e readiciona a linha para reiniciar a animação
                row.css('animation', 'none');
                row[0].offsetHeight; // Força um reflow
                row.css('animation', '');
                
                // Aplica o delay baseado na posição
                row.css('animation-delay', `${visibleCount * 0.02}s`);
                visibleCount++;
                row.show();
            } else {
                row.hide();
            }
        });

        updateTableStatus();
    }

    // Event listeners para os filtros
    $('#filterName').on('change', applyFilters);
    $('#filterPublisher').on('change', applyFilters);
    $('#filterLicenses').on('change', applyFilters);

    // Limpar filtros
    $('#clearFilters').click(function() {
        $('#filterName').val('').trigger('change');
        $('#filterPublisher').val('').trigger('change');
        $('#filterLicenses').val('').trigger('change');
        applyFilters();
    });

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
        const data = {
            name: formData.get('name'),
            publisher: formData.get('publisher'),
            description: formData.get('description')
        };

        // Corrigindo a lógica para identificar se é uma edição
        const softwareId = $('#softwareId').val();
        const isEdit = softwareId && softwareId !== '';
        const url = isEdit ? `/api/v1/licenses/software/${softwareId}` : '/api/v1/licenses/software';

        $.ajax({
            url: url,
            type: isEdit ? 'PUT' : 'POST',
            contentType: 'application/json',
            data: JSON.stringify(data),
            success: function(response) {
                $('#addSoftwareModal').modal('hide');
                Swal.fire({
                    icon: 'success',
                    title: 'Sucesso!',
                    text: `Software ${isEdit ? 'atualizado' : 'cadastrado'} com sucesso!`
                }).then(() => {
                    location.reload();
                });
            },
            error: function(xhr) {
                console.error('Erro:', xhr.responseText);
                Swal.fire({
                    icon: 'error',
                    title: 'Erro!',
                    text: `Erro ao ${isEdit ? 'atualizar' : 'cadastrar'} software: ` + xhr.responseText
                });
            }
        });
    });

    // Handler para editar software
    $('.edit-software').click(function() {
        const softwareId = $(this).data('id');
        
        $.ajax({
            url: `/api/v1/licenses/software/${softwareId}`,
            type: 'GET',
            success: function(software) {
                // Preenche os campos do formulário
                $('#softwareId').val(software.ID); // Ajustando para usar software.ID
                $('#softwareForm [name="name"]').val(software.name);
                $('#softwareForm [name="publisher"]').val(software.publisher);
                $('#softwareForm [name="description"]').val(software.description);
                
                // Atualiza o título do modal
                $('#modalTitle').text('Editar Software');
                
                // Abre o modal
                $('#addSoftwareModal').modal('show');
            },
            error: function(xhr) {
                Swal.fire({
                    icon: 'error',
                    title: 'Erro!',
                    text: 'Erro ao carregar dados do software: ' + xhr.responseText
                });
            }
        });
    });

    // Handler para deletar software
    $('.delete-software').click(function() {
        const softwareId = $(this).data('id');
        const row = $(this).closest('tr');
        const licenseCount = parseInt(row.find('td:eq(4)').text()) || 0; // Pega o número de licenças da coluna

        // Verifica se há licenças atribuídas
        if (licenseCount > 0) {
            Swal.fire({
                icon: 'error',
                title: 'Não é possível excluir',
                text: 'Este software possui licenças atribuídas. Remova todas as licenças antes de excluir o software.',
                confirmButtonColor: '#3085d6'
            });
            return;
        }

        // Se não houver licenças, prossegue com a confirmação de exclusão
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
        $('#softwareId').val('');  // Garante que o ID é limpo
        $('#modalTitle').text('Novo Software');
    });

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

    // Estilização do select com Select2
    $('#filterLicenses').select2({
        width: '100%',
        placeholder: 'Selecione...',
        allowClear: true
    });
});