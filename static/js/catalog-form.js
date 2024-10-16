$(document).ready(function() {
    // Inicializa o Select2 para melhorar a aparência e funcionalidade dos selects
    $('.select2').select2({
        theme: 'bootstrap-5',
        width: '100%'
    });

    // Adiciona animação suave ao scroll quando um filtro é alterado
    $('#filterForm select, #filterForm input').change(function() {
        $('html, body').animate({
            scrollTop: $('#filterForm').offset().top - 20
        }, 500);
    });

    // Adiciona efeito de hover nos cards de fornecedores
    $(document).on('mouseenter', '.card', function() {
        $(this).addClass('shadow-lg').css('cursor', 'pointer');
    }).on('mouseleave', '.card', function() {
        $(this).removeClass('shadow-lg');
    });

    // Intercepta o envio do formulário de cadastro de fornecedor
    $('#supplierForm').submit(function(e) {
        e.preventDefault();
        
        // Mostra um loader enquanto processa
        Swal.fire({
            title: 'Cadastrando fornecedor...',
            text: 'Por favor, aguarde.',
            allowOutsideClick: false,
            showConfirmButton: false,
            willOpen: () => {
                Swal.showLoading();
            }
        });

        // Envia o formulário via AJAX
        $.ajax({
            url: '/api/v1/suppliers',
            type: 'POST',
            data: $(this).serialize(),
            success: function(response) {
                // Mostra uma mensagem de sucesso
                Swal.fire({
                    icon: 'success',
                    title: 'Fornecedor cadastrado com sucesso!',
                    text: 'O fornecedor foi vinculado às categorias e serviços selecionados.',
                    confirmButtonText: 'OK'
                }).then((result) => {
                    if (result.isConfirmed) {
                        // Limpa o formulário
                        $('#supplierForm')[0].reset();
                        $('.select2').val(null).trigger('change');
                    }
                });
            },
            error: function(xhr) {
                // Mostra uma mensagem de erro
                Swal.fire({
                    icon: 'error',
                    title: 'Oops...',
                    text: 'Ocorreu um erro ao cadastrar o fornecedor: ' + (xhr.responseJSON ? xhr.responseJSON.error : 'Erro desconhecido'),
                    confirmButtonText: 'OK'
                });
            }
        });
    });

    // Adiciona funcionalidade para mostrar mais detalhes do fornecedor
    $(document).on('click', '.card', function() {
        var fornecedor = $(this).data('fornecedor');
        Swal.fire({
            title: fornecedor.NOME,
            html: `
                <p><strong>CNPJ:</strong> ${fornecedor.CGCCFO}</p>
                <p><strong>Endereço:</strong> ${fornecedor.RUA}, ${fornecedor.NUMERO} - ${fornecedor.BAIRRO}, ${fornecedor.CIDADE} - ${fornecedor.UF}</p>
                <p><strong>CEP:</strong> ${fornecedor.CEP}</p>
                <p><strong>Telefone:</strong> ${fornecedor.TELEFONE}</p>
                <p><strong>Email:</strong> ${fornecedor.EMAIL}</p>
                <p><strong>Tipo:</strong> ${fornecedor.TIPO}</p>
            `,
            confirmButtonText: 'Fechar'
        });
    });

    // Código existente para o formulário de filtro...
});
