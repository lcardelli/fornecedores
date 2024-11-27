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
});