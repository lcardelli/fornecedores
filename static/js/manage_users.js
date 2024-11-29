// Adicione ao seu arquivo JavaScript existente

$(document).ready(function() {
    // Handler para deletar usuário
    $('.delete-user').click(function() {
        const userId = $(this).data('user-id');
        const row = $(this).closest('tr');
        
        Swal.fire({
            title: 'Tem certeza?',
            text: "Você não poderá reverter esta ação!",
            icon: 'warning',
            showCancelButton: true,
            confirmButtonColor: '#3085d6',
            cancelButtonColor: '#d33',
            confirmButtonText: 'Sim, deletar!',
            cancelButtonText: 'Cancelar'
        }).then((result) => {
            if (result.isConfirmed) {
                $.ajax({
                    url: `/api/v1/users/${userId}`,
                    type: 'DELETE',
                    success: function(response) {
                        row.remove();
                        Swal.fire(
                            'Deletado!',
                            'O usuário foi deletado com sucesso.',
                            'success'
                        );
                    },
                    error: function(xhr, status, error) {
                        Swal.fire({
                            icon: 'error',
                            title: 'Erro!',
                            text: 'Erro ao deletar usuário: ' + xhr.responseText
                        });
                    }
                });
            }
        });
    });

    // Função para carregar departamentos no select
    function loadDepartments() {
        $.get('/api/v1/departments', function(departments) {
            const select = $('#department');
            select.empty();
            departments.forEach(dept => {
                select.append(new Option(dept.name, dept.name));
            });
        });
    }

    // Carregar departamentos quando a página carregar
    loadDepartments();

    // Handler para abrir o modal de permissões
    $(document).on('click', '.manage-permissions', function() {
        const userId = $(this).data('user-id');
        const userName = $(this).data('username');
        
        $('#userId').val(userId);
        $('#userName').text(userName);
        
        // Carrega as permissões atuais
        loadUserPermissions(userId);
        
        $('#permissionsModal').modal('show');
    });

    // Atualiza visual quando um switch é alterado
    $('.custom-control-input').change(function() {
        const control = $(this).closest('.custom-control');
        const label = $(this).siblings('.custom-control-label');
        
        if ($(this).is(':checked')) {
            control.addClass('active-permission');
            label.addClass('text-success');
        } else {
            control.removeClass('active-permission');
            label.removeClass('text-success');
        }
    });

    // Handler para salvar as permissões
    $('#savePermissions').click(function() {
        const userId = $('#userId').val();
        const data = {
            user_id: parseInt(userId),
            department: $('#department').val(),
            view_suppliers: $('#viewSuppliers').is(':checked'),
            view_licenses: $('#viewLicenses').is(':checked'),
            admin_suppliers: $('#adminSuppliers').is(':checked'),
            admin_licenses: $('#adminLicenses').is(':checked')
        };
        
        $.ajax({
            url: '/api/v1/users/permissions',
            method: 'POST',
            contentType: 'application/json',
            data: JSON.stringify(data),
            success: function(response) {
                $('#permissionsModal').modal('hide');
                Swal.fire({
                    icon: 'success',
                    title: 'Sucesso!',
                    text: 'Permissões atualizadas com sucesso!'
                }).then(() => {
                    location.reload();
                });
            },
            error: function(xhr) {
                Swal.fire({
                    icon: 'error',
                    title: 'Erro!',
                    text: 'Erro ao atualizar permissões: ' + xhr.responseText
                });
            }
        });
    });

    // Handler para fechar o modal
    $('.btn-close, button[data-bs-dismiss="modal"]').click(function() {
        $('#permissionsModal').modal('hide');
    });

    // Atualizar o Bootstrap 5 modal
    var permissionsModal = new bootstrap.Modal(document.getElementById('permissionsModal'), {
        keyboard: true,
        backdrop: true
    });
});

// Função que carrega as permissões
function loadUserPermissions(userId) {
    $.get(`/api/v1/users/${userId}/permissions`, function(data) {
        console.log('Permissões carregadas:', data);
        $('#userId').val(userId);
        $('#isAdmin').prop('checked', data.is_admin);
        $('#department').val(data.department || 'Geral');
        $('#viewSuppliers').prop('checked', data.view_suppliers);
        $('#viewLicenses').prop('checked', data.view_licenses);
        $('#adminSuppliers').prop('checked', data.admin_suppliers);
        $('#adminLicenses').prop('checked', data.admin_licenses);
    });
} 