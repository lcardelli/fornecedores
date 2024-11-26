// Adicione ao seu arquivo JavaScript existente

$(document).ready(function() {
    // Handler para abrir o modal de permissões
    $(document).on('click', '.manage-permissions', function() {
        const userId = $(this).data('user-id');
        const userName = $(this).data('username');
        
        $('#userId').val(userId);
        $('#userName').text(userName);
        
        // Carrega as permissões atuais
        $.get(`/api/v1/users/${userId}/permissions`, function(data) {
            // Atualiza os switches com as permissões atuais
            $('#adminAccess').prop('checked', data.is_admin);
            $('#department').val(data.department || 'Geral');
            $('#viewSuppliers').prop('checked', data.view_suppliers);
            $('#viewLicenses').prop('checked', data.view_licenses);

            // Adiciona classe para switches ativos
            $('.custom-control-input:checked').each(function() {
                $(this).closest('.custom-control').addClass('active-permission');
                $(this).siblings('.custom-control-label').addClass('text-success');
            });
        });
        
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
        const isAdmin = $('#adminAccess').is(':checked');

        // Primeiro atualiza o status de admin
        $.ajax({
            url: `/api/v1/users/${userId}/toggle-admin`,
            type: 'PUT',
            contentType: 'application/json',
            data: JSON.stringify({ admin: isAdmin }),
            success: function() {
                // Depois atualiza as outras permissões
                const data = {
                    user_id: parseInt(userId),
                    department: $('#department').val(),
                    view_suppliers: $('#viewSuppliers').is(':checked'),
                    view_licenses: $('#viewLicenses').is(':checked')
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
            },
            error: function() {
                Swal.fire({
                    icon: 'error',
                    title: 'Erro!',
                    text: 'Erro ao atualizar privilégios de administrador'
                });
            }
        });
    });
}); 