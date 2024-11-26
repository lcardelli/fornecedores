// Adicione ao seu arquivo JavaScript existente

$(document).ready(function() {
    // Handler para abrir o modal de permissões
    $(document).on('click', '.manage-permissions', function() {
        const userId = $(this).data('user-id');
        const userName = $(this).data('username');
        const isAdmin = $(this).closest('tr').find('.admin-toggle').prop('checked');
        
        $('#userId').val(userId);
        $('#userName').text(userName);
        $('#adminAccess').prop('checked', isAdmin);
        
        // Carrega as permissões atuais
        $.get(`/api/v1/users/${userId}/permissions`, function(data) {
            $('#department').val(data.department || 'Geral');
            $('#viewSuppliers').prop('checked', data.view_suppliers);
            $('#viewLicenses').prop('checked', data.view_licenses);
        });
        
        $('#permissionsModal').modal('show');
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