$(document).ready(function() {
    // Inicialização
    initializeComponents();
    setupEventHandlers();
    setupFilters();
    setupCheckboxes();
    setupTableScroll();

    // ==================== Inicialização de Componentes ====================
    function initializeComponents() {
        // Configuração do Flatpickr em português
        const flatpickrConfig = {
            dateFormat: "d/m/Y",
            locale: {
                firstDayOfWeek: 0,
                weekdays: {
                    shorthand: ["Dom", "Seg", "Ter", "Qua", "Qui", "Sex", "Sáb"],
                    longhand: ["Domingo", "Segunda", "Terça", "Quarta", "Quinta", "Sexta", "Sábado"]
                },
                months: {
                    shorthand: ["Jan", "Fev", "Mar", "Abr", "Mai", "Jun", "Jul", "Ago", "Set", "Out", "Nov", "Dez"],
                    longhand: ["Janeiro", "Fevereiro", "Março", "Abril", "Maio", "Junho", "Julho", "Agosto", "Setembro", "Outubro", "Novembro", "Dezembro"]
                },
                rangeSeparator: " até ",
                weekAbbreviation: "Sem",
                scrollTitle: "Role para aumentar",
                toggleTitle: "Clique para alternar",
                amPM: ["AM", "PM"],
                yearAriaLabel: "Ano",
                monthAriaLabel: "Mês",
                hourAriaLabel: "Hora",
                minuteAriaLabel: "Minuto"
            },
            allowInput: true,
            altInput: true,
            altFormat: "d/m/Y",
            time_24hr: true,
            defaultHour: 0
        };

        // Inicializa o Flatpickr com as configurações
        $('.datepicker').flatpickr(flatpickrConfig);

        // Select2
        $('.select2').select2({
            width: '100%',
            placeholder: "Selecione..."
        });

        // Inicializa o campo de custo
        const costInput = document.querySelector('input[name="cost"]');
        costInput.value = 'R$ 0,00';
        
        // Estilização dos selects de filtro com Select2
        $('#filterSoftware, #filterType, #filterStatus, #filterYear').select2({
            width: '100%',
            placeholder: 'Selecione...',
            allowClear: true,
            dropdownParent: $('body'),
            containerCssClass: 'manage-licenses-select-container',
            dropdownCssClass: 'manage-licenses-dropdown'
        });

        // Inicialização dos selects do modal
        $('.manage-licenses-modal-select').select2({
            width: '100%',
            placeholder: "Selecione...",
            allowClear: true,
            dropdownParent: $('#addLicenseModal'),
            containerCssClass: 'manage-licenses-modal-container',
            dropdownCssClass: 'manage-licenses-modal-dropdown'
        });
    }

    // ==================== Event Handlers ====================
    function setupEventHandlers() {
        setupClearFiltersHandler();
        setupEditLicenseHandler();
        setupTogglePasswordHandler();
        setupSaveLicenseHandler();
        setupCostInputHandler();
        setupDeleteLicenseHandler();
        setupNewLicenseHandler();
    }

    function setupClearFiltersHandler() {
        $('#clearFilters').click(function() {
            // Limpa todos os selects e dispara o evento change
            $('#filterSoftware').val('').trigger('change');
            $('#filterType').val('').trigger('change');
            $('#filterStatus').val('').trigger('change');
            $('#filterYear').val('').trigger('change');
            
            // Limpa o select2 do departamento
            $('#filterDepartment').val(null).trigger('change');
            
            // Aplica os filtros após limpar
            applyFilters();
        });
    }

    function setupEditLicenseHandler() {
        $('.edit-license').click(function() {
            const licenseId = $(this).data('id');
            $('#licenseForm')[0].reset();
            $('#licenseForm').removeData('license-id');
            
            $.ajax({
                url: `/api/v1/licenses/${licenseId}`,
                type: 'GET',
                success: function(response) {
                    fillEditForm(response);
                    $('#modalTitle').text('Editar Licença');
                    $('#licenseForm').data('license-id', licenseId);
                    $('#addLicenseModal').modal('show');
                },
                error: function(xhr) {
                    let errorMessage = 'Erro ao carregar dados da licença';
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
    }

    function setupTogglePasswordHandler() {
        $('.toggle-password').click(function() {
            const passwordField = $(this).closest('.input-group').find('input');
            const icon = $(this).find('i');
            
            if (passwordField.attr('type') === 'password') {
                passwordField.attr('type', 'text');
                icon.removeClass('fa-eye-slash').addClass('fa-eye');
            } else {
                passwordField.attr('type', 'password');
                icon.removeClass('fa-eye').addClass('fa-eye-slash');
            }
        });
    }

    function setupSaveLicenseHandler() {
        $('#saveLicense').click(function() {
            const form = $('#licenseForm');
            const licenseId = form.data('license-id');
            const isEdit = !!licenseId;
            
            if (!form[0].checkValidity()) {
                form[0].reportValidity();
                return;
            }

            // Coleta todos os dados do formulário
            const formData = new FormData(form[0]);
            const data = {};
            
            // Converte os dados do formulário para um objeto
            for (let [key, value] of formData.entries()) {
                if (key === 'department_id') {
                    const departmentId = parseInt(value);
                    if (!departmentId) {
                        Swal.fire({
                            icon: 'error',
                            title: 'Erro!',
                            text: 'Por favor, selecione um departamento'
                        });
                        return;
                    }
                    data[key] = departmentId;
                } else {
                    data[key] = value;
                }
            }

            // Prepara os dados para envio
            prepareFormData(data);

            const url = isEdit ? `/api/v1/licenses/${licenseId}` : '/api/v1/licenses';
            const method = isEdit ? 'PUT' : 'POST';

            $.ajax({
                url: url,
                type: method,
                contentType: 'application/json',
                data: JSON.stringify(data),
                success: function(response) {
                    $('#addLicenseModal').modal('hide');
                    Swal.fire({
                        icon: 'success',
                        title: 'Sucesso!',
                        text: `Licença ${isEdit ? 'atualizada' : 'cadastrada'} com sucesso!`
                    }).then(() => {
                        location.reload();
                    });
                },
                error: function(xhr) {
                    let errorMessage = xhr.responseText;
                    try {
                        const errorObj = JSON.parse(xhr.responseText);
                        errorMessage = errorObj.error || errorMessage;
                    } catch (e) {}
                    
                    Swal.fire({
                        icon: 'error',
                        title: 'Erro!',
                        text: `Erro ao ${isEdit ? 'atualizar' : 'cadastrar'} licença: ${errorMessage}`
                    });
                }
            });
        });
    }

    function setupCostInputHandler() {
        const costInput = document.querySelector('input[name="cost"]');
        costInput.addEventListener('input', function(e) {
            handleCostInput(this);
        });
        
        costInput.addEventListener('focus', function() {
            this.setSelectionRange(3, this.value.length - 3);
        });
    }

    function setupDeleteLicenseHandler() {
        $('.delete-license').click(function() {
            const licenseId = $(this).data('id');
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
                        url: `/api/v1/licenses/${licenseId}`,
                        type: 'DELETE',
                        success: function() {
                            Swal.fire({
                                icon: 'success',
                                title: 'Deletado!',
                                text: 'A licença foi deletada com sucesso.',
                                allowOutsideClick: false
                            }).then(() => {
                                // Recarrega a página após a exclusão
                                window.location.reload();
                            });
                        },
                        error: function(xhr) {
                            Swal.fire({
                                icon: 'error',
                                title: 'Erro!',
                                text: 'Erro ao deletar licença: ' + xhr.responseText
                            });
                        }
                    });
                }
            });
        });
    }

    function updateTotalCost(costChange) {
        const totalCell = $('tr.bg-light.font-weight-bold td:nth-child(9)');
        const currentTotal = parseFloat(totalCell.text().replace('R$ ', '').replace('.', '').replace(',', '.'));
        const newTotal = currentTotal + costChange;
        
        // Formata o novo total usando a mesma função de formatação de moeda
        const formattedTotal = formatMoneyBR(newTotal);
        totalCell.text(formattedTotal);
    }

    function formatMoneyBR(value) {
        return new Intl.NumberFormat('pt-BR', {
            style: 'currency',
            currency: 'BRL',
            minimumFractionDigits: 2
        }).format(value);
    }

    // ==================== Checkbox Handlers ====================
    function setupCheckboxes() {
        // Handler para o checkbox "Selecionar Todos"
        $('.select-all-checkbox').change(function() {
            const isChecked = $(this).prop('checked');
            $('.license-checkbox').prop('checked', isChecked);
            updateDeleteSelectedButton();
        });

        // Handler para os checkboxes individuais
        $('.license-checkbox').change(function() {
            const allChecked = $('.license-checkbox:checked').length === $('.license-checkbox').length;
            $('.select-all-checkbox').prop('checked', allChecked);
            updateDeleteSelectedButton();
        });
    }

    // Função para atualizar a visibilidade do botão de excluir selecionados
    function updateDeleteSelectedButton() {
        const checkedCount = $('.license-checkbox:checked').length;
        if (checkedCount > 0) {
            $('#deleteSelectedBtn').fadeIn(200).css('display', 'inline-flex');
        } else {
            $('#deleteSelectedBtn').fadeOut(200);
        }
    }

    // Handler para o botão de excluir selecionados
    $('#deleteSelectedBtn').click(function() {
        const selectedIds = $('.license-checkbox:checked').map(function() {
            return parseInt($(this).closest('tr').find('.license-checkbox').val());
        }).get();

        if (selectedIds.length === 0) return;

        Swal.fire({
            title: 'Tem certeza?',
            text: `Você está prestes a excluir ${selectedIds.length} licença(s). Esta ação não pode ser revertida!`,
            icon: 'warning',
            showCancelButton: true,
            confirmButtonColor: '#3085d6',
            cancelButtonColor: '#d33',
            confirmButtonText: 'Sim, excluir!',
            cancelButtonText: 'Cancelar'
        }).then((result) => {
            if (result.isConfirmed) {
                deleteMultipleLicenses(selectedIds);
            }
        });
    });

    // Função para deletar múltiplas licenças
    function deleteMultipleLicenses(ids) {
        $.ajax({
            url: '/api/v1/licenses/batch',
            type: 'DELETE',
            data: JSON.stringify({ ids: ids }),
            contentType: 'application/json',
            success: function(response) {
                Swal.fire({
                    icon: 'success',
                    title: 'Sucesso!',
                    text: 'As licenças selecionadas foram excluídas com sucesso.',
                }).then(() => {
                    location.reload();
                });
            },
            error: function(xhr, status, error) {
                Swal.fire({
                    icon: 'error',
                    title: 'Erro!',
                    text: 'Não foi possível excluir as licenças selecionadas: ' + 
                        (xhr.responseJSON ? xhr.responseJSON.error : error)
                });
            }
        });
    }

    // ==================== Filter Handlers ====================
    function setupFilters() {
        $('#filterSoftware, #filterType, #filterStatus, #filterYear, #filterDepartment').on('change', function() {
            applyFilters();
        });
    }

    function applyFilters() {
        const softwareFilter = $('#filterSoftware').val()?.toLowerCase() || '';
        const typeFilter = $('#filterType').val() || '';
        const statusFilter = $('#filterStatus').val() || '';
        const departmentFilter = $('#filterDepartment').val() || '';
        const yearFilter = $('#filterYear').val() || '';

        let totalCost = 0;
        let visibleRowCount = 0;

        $('#licensesTable tr').each(function() {
            const row = $(this);
            if (row.find('td').length > 0) {
                const software = row.find('td:eq(1)').text().toLowerCase();
                const type = row.find('td:eq(2)').text();
                const statusId = row.find('.badge').data('status-id');
                const departmentId = row.find('td:eq(7)').data('department-id');
                const expiryDateText = row.find('td:eq(6)').text();
                const year = expiryDateText !== '-' ? expiryDateText.split('/')[2] : '';
                
                const costText = row.find('.license-cost').text().trim();
                const cost = parseFloat(costText.replace('R$ ', '').replace('.', '').replace(',', '.')) || 0;

                const matchesSoftware = !softwareFilter || software.includes(softwareFilter);
                const matchesType = !typeFilter || type === typeFilter;
                const matchesStatus = !statusFilter || statusId === parseInt(statusFilter);
                const matchesDepartment = !departmentFilter || departmentId === parseInt(departmentFilter);
                const matchesYear = !yearFilter || year === yearFilter;

                if (matchesSoftware && matchesType && matchesStatus && matchesDepartment && matchesYear) {
                    row.show();
                    row.css({
                        'animation': 'none',
                        'opacity': '0'
                    });
                    setTimeout(() => {
                        row.css({
                            'animation': 'fadeIn 0.5s ease-out forwards',
                            'animation-delay': `${visibleRowCount * 0.1}s`
                        });
                    }, 0);
                    visibleRowCount++;
                    totalCost += cost;
                } else {
                    row.hide();
                }
            }
        });

        // Atualiza o valor total no novo formato do rodapé
        const formattedTotal = formatMoneyBR(totalCost);
        $('.table-footer .total-value').text(formattedTotal);

        updateTableStatus();
    }

    // ==================== Utility Functions ====================
    function fillEditForm(license) {
        const form = $('#licenseForm');
        
        form.find('[name="software_id"]').val(license.software_id).trigger('change');
        form.find('[name="license_key"]').val(license.license_key);
        form.find('[name="username"]').val(license.username);
        form.find('[name="password"]').val(license.password);
        
        // Atualiza o campo tipo
        const typeSelect = form.find('[name="type"]');
        typeSelect.val(license.type).trigger('change');
        
        // Atualiza o campo período de renovação
        const periodSelect = form.find('[name="period_renew_id"]');
        if (license.period_renew_id) {
            periodSelect.val(license.period_renew_id).trigger('change');
        } else {
            periodSelect.val('').trigger('change');
        }
        
        if (license.department_id) {
            form.find('[name="department_id"]').val(license.department_id).trigger('change');
        }
        form.find('[name="quantity"]').val(license.quantity);
        form.find('[name="seats"]').val(license.seats);
        form.find('[name="notes"]').val(license.notes);
        
        if (license.purchase_date) {
            const purchaseDate = new Date(license.purchase_date);
            const formattedPurchaseDate = formatDateBR(purchaseDate);
            const purchasePicker = form.find('[name="purchase_date"]')[0]._flatpickr;
            purchasePicker.setDate(formattedPurchaseDate, true);
        }
        
        if (license.expiry_date) {
            const expiryDate = new Date(license.expiry_date);
            const formattedExpiryDate = formatDateBR(expiryDate);
            const expiryPicker = form.find('[name="expiry_date"]')[0]._flatpickr;
            expiryPicker.setDate(formattedExpiryDate, true);
        }
        
        const cost = license.cost || 0;
        form.find('[name="cost"]').val(formatMoney(cost));
        
        // Atualiza o status na tabela
        const row = $(`tr:has(button[data-id="${license.id}"])`);
        if (row.length && license.status) {
            const statusBadge = row.find('.badge');
            const statusClass = getStatusClass(license.status.name);
            statusBadge.removeClass().addClass(`badge ${statusClass}`).text(license.status.name);
        }
        
        // Atualiza a visibilidade do campo período de renovação baseado no tipo
        const periodField = periodSelect.closest('.form-group');
        if (license.type === 'Subscrição') {
            periodField.show();
            periodSelect.prop('required', true);
        } else {
            periodField.hide();
            periodSelect.prop('required', false);
        }
    }

    function prepareFormData(data) {
        // Garante que o license_key seja enviado mesmo vazio
        data.license_key = data.license_key || '';  // Se for undefined ou null, usa string vazia
        
        data.cost = parseFloat(unformatMoney(data.cost || '0'));
        data.quantity = parseInt(data.quantity) || 0;
        data.seats = parseInt(data.seats) || 0;

        if (data.purchase_date) {
            try {
                const [dia, mes, ano] = data.purchase_date.split('/');
                const date = new Date(ano, mes - 1, dia, 0, 0, 0);
                data.purchase_date = date.toISOString();
            } catch (e) {
                console.error('Erro ao converter data de compra:', e);
                data.purchase_date = null;
            }
        }

        if (data.expiry_date) {
            try {
                const [dia, mes, ano] = data.expiry_date.split('/');
                const date = new Date(ano, mes - 1, dia, 0, 0, 0);
                data.expiry_date = date.toISOString();
            } catch (e) {
                console.error('Erro ao converter data de expiração:', e);
                data.expiry_date = null;
            }
        }

        data.software_id = parseInt(data.software_id);
        data.period_renew_id = data.period_renew_id ? parseInt(data.period_renew_id) : null;
        return data;
    }

    function formatDateToISO(dateStr) {
        if (!dateStr) return null;
        const parts = dateStr.split('/');
        return parts.length === 3 ? `${parts[2]}-${parts[1]}-${parts[0]}T00:00:00Z` : null;
    }

    function formatMoney(value) {
        // Garante que o valor é um número
        const numValue = parseFloat(value);
        
        return new Intl.NumberFormat('pt-BR', {
            style: 'currency',
            currency: 'BRL',
            minimumFractionDigits: 2,
            maximumFractionDigits: 2
        }).format(numValue);
    }

    function unformatMoney(value) {
        // Remove todos os caracteres não numéricos exceto vírgula e ponto
        const cleanValue = value.replace(/[^\d,\.]/g, '');
        // Substitui vírgula por ponto e converte para float
        return cleanValue.replace(/\./g, '').replace(',', '.');
    }

    function handleCostInput(input) {
        let value = input.value.replace(/\D/g, '');
        
        if (value === '') {
            input.value = 'R$ 0,00';
            return;
        }
        
        // Converte para centavos
        value = (parseFloat(value) / 100).toFixed(2);
        
        // Formata usando Intl.NumberFormat
        const formattedValue = new Intl.NumberFormat('pt-BR', {
            style: 'currency',
            currency: 'BRL',
            minimumFractionDigits: 2,
            maximumFractionDigits: 2
        }).format(value);
        
        input.value = formattedValue;
    }

    function updateTableStatus() {
        const visibleRows = $('#licensesTable tr:visible').filter(function() {
            return $(this).find('td').length > 0;
        }).length;

        if (visibleRows === 0) {
            if ($('#noResultsMessage').length === 0) {
                $('#licensesTable').after(`
                    <tr id="noResultsMessage">
                        <td colspan="10" class="text-center">
                            <div class="empty-state">
                                <i class="fas fa-key fa-3x mb-3"></i>
                                <p class="h5">Nenhuma licença encontrada</p>
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

    // Função auxiliar para formatar data no padrão brasileiro
    function formatDateBR(date) {
        if (!date) return '';
        const d = new Date(date);
        const day = d.getDate().toString().padStart(2, '0');
        const month = (d.getMonth() + 1).toString().padStart(2, '0');
        const year = d.getFullYear();
        return `${day}/${month}/${year}`;
    }

    // Adicione um handler para mostrar/esconder o campo de período baseado no tipo
    $('select[name="type"]').on('change', function() {
        const periodField = $('select[name="period_renew"]').closest('.form-group');
        if ($(this).val() === 'Subscrição') {
            periodField.show();
            $('select[name="period_renew"]').prop('required', true);
        } else {
            periodField.hide();
            $('select[name="period_renew"]').prop('required', false);
            $('select[name="period_renew"]').val('');
        }
    });

    // Atualizar a função que exibe o status na interface
    function displayStatus(license) {
        return license.status ? license.status.name : 'N/A';
    }

    // Função auxiliar para definir a classe CSS do status
    function getStatusClass(statusName) {
        switch (statusName) {
            case 'Ativa':
                return 'badge-success';
            case 'Próxima ao vencimento':
                return 'badge-warning';
            case 'Vencida':
                return 'badge-danger';
            default:
                return 'badge-secondary';
        }
    }

    function updateLicenseRow(license) {
        const row = $(`tr:has(button[data-id="${license.id}"])`);
        if (row.length) {
            const cells = row.find('td');
            
            // Atualiza cada célula na ordem correta
            $(cells[1]).text(license.software ? license.software.name : '-');
            $(cells[2]).text(license.type || '-');
            $(cells[3]).text(license.quantity || '-');
            $(cells[4]).text(license.period_renew ? license.period_renew.name : '-');
            $(cells[5]).text(formatDateBR(license.purchase_date));
            $(cells[6]).text(formatDateBR(license.expiry_date));
            $(cells[7]).text(license.department || '-');
            $(cells[8]).text(formatMoney(String(license.cost * 100)));

            // Atualiza o status
            const statusBadge = $(cells[9]).find('.badge');
            const statusClass = getStatusClass(license.status.name);
            statusBadge
                .removeClass('badge-success badge-warning badge-danger badge-secondary')
                .addClass(statusClass)
                .text(license.status.name);
        }
    }

    // Adicione estas funções auxiliares se ainda não existirem
    function formatMoneyBR(value) {
        return new Intl.NumberFormat('pt-BR', {
            style: 'currency',
            currency: 'BRL'
        }).format(value);
    }

    // Atualize o CSS para incluir a animação
    function addStyleToHead() {
        const style = `
            @keyframes fadeIn {
                from {
                    opacity: 0;
                    transform: translateY(20px);
                }
                to {
                    opacity: 1;
                    transform: translateY(0);
                }
            }
        `;
        
        if (!document.getElementById('animation-style')) {
            const styleElement = document.createElement('style');
            styleElement.id = 'animation-style';
            styleElement.textContent = style;
            document.head.appendChild(styleElement);
        }
    }

    // Chame esta função quando o documento estiver pronto
    addStyleToHead();

    // Função para buscar taxas de câmbio em tempo real
    async function getExchangeRates() {
        try {
            const apiKey = '5625ad94c54cb820a54838d0';
            const response = await fetch(`https://v6.exchangerate-api.com/v6/${apiKey}/latest/USD`);
            const data = await response.json();
            
            if (data.result === 'error') {
                throw new Error(data['error-type']);
            }
            
            return data.conversion_rates;
        } catch (error) {
            console.error('Erro ao buscar taxas de câmbio:', error);
            Swal.fire({
                icon: 'error',
                title: 'Erro',
                text: 'Não foi possível obter as taxas de câmbio atualizadas. Por favor, tente novamente mais tarde.',
                footer: 'Detalhes: ' + error.message
            });
            return null;
        }
    }

    // Função para calcular os valores
    async function calculateLicenseValues() {
        console.log('Calculando valores...'); // Debug

        const baseValue = parseFloat($('#baseValue').val()) || 0;
        const currency = $('#currency').val();
        const period = $('#period').val();
        const quantity = parseInt($('#quantity').val()) || 1;

        console.log('Valores de entrada:', { baseValue, currency, period, quantity }); // Debug

        // Busca taxas de câmbio atualizadas
        const rates = await getExchangeRates();
        if (!rates) {
            return;
        }

        // Converte para BRL
        let valueInBRL = baseValue;
        if (currency === 'USD') {
            valueInBRL = baseValue * rates.BRL;
        } else if (currency === 'EUR') {
            valueInBRL = baseValue * (rates.BRL / rates.EUR);
        }

        // Resto do código permanece igual...
        let monthlyValue = valueInBRL;
        switch (period) {
            case 'quarterly':
                monthlyValue = valueInBRL / 3;
                break;
            case 'semiannual':
                monthlyValue = valueInBRL / 6;
                break;
            case 'annual':
                monthlyValue = valueInBRL / 12;
                break;
        }

        monthlyValue *= quantity;

        const quarterlyValue = monthlyValue * 3;
        const semiannualValue = monthlyValue * 6;
        const annualValue = monthlyValue * 12;

        try {
            $('#monthlyResult').text(formatMoney(monthlyValue));
            $('#quarterlyResult').text(formatMoney(quarterlyValue));
            $('#semiannualResult').text(formatMoney(semiannualValue));
            $('#annualResult').text(formatMoney(annualValue));
        } catch (error) {
            console.error('Erro ao formatar valores:', error);
        }
    }

    // Event listeners para atualização automática
    $('#baseValue, #currency, #period, #quantity').on('change input', function() {
        calculateLicenseValues();
    });

    // Inicializa os valores quando o modal é aberto
    $('#calculatorModal').on('shown.bs.modal', function() {
        calculateLicenseValues();
    });

    // Event listener para limpar os campos quando o modal for fechado
    $('#calculatorModal').on('hidden.bs.modal', function() {
        // Limpa os inputs
        $('#baseValue').val('');
        $('#currency').val('BRL'); // Volta para a moeda padrão
        $('#period').val('monthly'); // Volta para o período padrão
        $('#quantity').val('1'); // Volta para quantidade padrão

        // Limpa os resultados
        $('#monthlyResult').text('R$ 0,00');
        $('#quarterlyResult').text('R$ 0,00');
        $('#semiannualResult').text('R$ 0,00');
        $('#annualResult').text('R$ 0,00');
    });

    // Adicionar esta nova função
    function setupNewLicenseHandler() {
        // Handler para o botão de Nova Licença
        $('button[data-target="#addLicenseModal"]').click(function() {
            // Limpa o formulário
            $('#licenseForm')[0].reset();
            
            // Remove qualquer ID de licença armazenado
            $('#licenseForm').removeData('license-id');
            
            // Reseta os selects do Select2
            $('.select2').val(null).trigger('change');
            
            // Reseta o campo de custo
            $('input[name="cost"]').val('R$ 0,00');
            
            // Reseta as datas do Flatpickr
            const purchasePicker = $('input[name="purchase_date"]')[0]._flatpickr;
            const expiryPicker = $('input[name="expiry_date"]')[0]._flatpickr;
            purchasePicker.clear();
            expiryPicker.clear();
            
            // Atualiza o título do modal
            $('#modalTitle').text('Nova Licença');
        });
    }

    // Adicione esta função no seu arquivo JS
    function setupTableScroll() {
        const tableWrapper = document.querySelector('.table-wrapper');
        const thead = document.querySelector('.table thead');

        if (tableWrapper && thead) {
            tableWrapper.addEventListener('scroll', function() {
                const scrollTop = this.scrollTop;
                thead.style.transform = `translateY(${scrollTop}px)`;
            });
        }
    }

    // Chame a função após o documento estar pronto
    setupTableScroll();

    // Certifique-se de que esta função é chamada após carregar as licenças
    function renderLicenses(licenses) {
        // ... código existente ...
        
        // Após renderizar a tabela, configure os checkboxes
        setupCheckboxes();
    }
});

