$(document).ready(function() {
    // Inicialização
    initializeComponents();
    setupEventHandlers();
    setupFilters();
    setupCheckboxes();

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
        $('#filterSoftware, #filterType, #filterStatus').select2({
            width: '100%',
            placeholder: 'Selecione...',
            allowClear: true
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
    }

    function setupClearFiltersHandler() {
        $('#clearFilters').click(function() {
            $('#filterSoftware').val('').trigger('change');
            $('#filterType').val('').trigger('change');
            $('#filterStatus').val('').trigger('change');
            $('#filterDepartment').val('');
            $('#licensesTable tr').show();
            $('#noResultsMessage').remove();
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

            const formData = new FormData(form[0]);
            const data = Object.fromEntries(formData.entries());
            
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
                    Swal.fire({
                        icon: 'error',
                        title: 'Erro!',
                        text: `Erro ao ${isEdit ? 'atualizar' : 'cadastrar'} licença: ` + xhr.responseText
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
                            row.remove();
                            Swal.fire(
                                'Deletado!',
                                'A licença foi deletada com sucesso.',
                                'success'
                            );
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

    // ==================== Checkbox Handlers ====================
    function setupCheckboxes() {
        $("#selectAll").change(function() {
            $(".license-checkbox").prop('checked', $(this).prop("checked"));
            updateDeleteButtonVisibility();
        });

        $(".license-checkbox").change(function() {
            updateSelectAllCheckbox();
            updateDeleteButtonVisibility();
        });
    }

    function updateSelectAllCheckbox() {
        var totalCheckboxes = $(".license-checkbox").length;
        var checkedCheckboxes = $(".license-checkbox:checked").length;
        $("#selectAll").prop('checked', totalCheckboxes === checkedCheckboxes);
    }

    function updateDeleteButtonVisibility() {
        var checkedCheckboxes = $(".license-checkbox:checked").length;
        if (checkedCheckboxes > 0) {
            $("#deleteSelected").show();
        } else {
            $("#deleteSelected").hide();
        }
    }

    // ==================== Filter Handlers ====================
    function setupFilters() {
        $('#filterSoftware, #filterType, #filterStatus').on('change', applyFilters);
        $('#filterDepartment').on('input', applyFilters);
    }

    function applyFilters() {
        const softwareFilter = $('#filterSoftware').val().toLowerCase();
        const typeFilter = $('#filterType').val();
        const statusFilter = $('#filterStatus').val();
        const departmentFilter = $('#filterDepartment').val().toLowerCase();

        $('#licensesTable tr').each(function() {
            const row = $(this);
            if (row.find('td').length > 0) {
                const software = row.find('td:eq(1)').text().toLowerCase();
                const type = row.find('td:eq(3)').text();
                const status = row.find('.badge').text().trim();
                const department = row.find('td').text().toLowerCase();

                const matchesSoftware = !softwareFilter || software.includes(softwareFilter);
                const matchesType = !typeFilter || type === typeFilter;
                const matchesStatus = !statusFilter || status === statusFilter;
                const matchesDepartment = !departmentFilter || department.includes(departmentFilter);

                if (matchesSoftware && matchesType && matchesStatus && matchesDepartment) {
                    row.show();
                } else {
                    row.hide();
                }
            }
        });

        updateTableStatus();
    }

    // ==================== Utility Functions ====================
    function fillEditForm(license) {
        const form = $('#licenseForm');
        
        form.find('[name="software_id"]').val(license.software_id).trigger('change');
        form.find('[name="license_key"]').val(license.license_key);
        form.find('[name="username"]').val(license.username);
        form.find('[name="password"]').val(license.password);
        form.find('[name="type"]').val(license.type);
        form.find('[name="department"]').val(license.department);
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
        
        const cost = license.cost ? license.cost * 100 : 0;
        form.find('[name="cost"]').val(formatMoney(String(cost)));
        
        if (license.assigned_users && license.assigned_users.length > 0) {
            const userIds = license.assigned_users.map(user => user.ID);
            form.find('[name="assigned_users[]"]').val(userIds).trigger('change');
        }
        
        const periodRenew = license.period_renew || '';
        form.find('[name="period_renew"]').val(periodRenew);

        // Atualiza o status na tabela
        const row = $(`tr:has(button[data-id="${license.id}"])`);
        if (row.length && license.status) {
            const statusBadge = row.find('.badge');
            const statusClass = getStatusClass(license.status.name);
            statusBadge.removeClass().addClass(`badge ${statusClass}`).text(license.status.name);
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
        data.period_renew = parseInt(data.period_renew) || 0;
        return data;
    }

    function formatDateToISO(dateStr) {
        if (!dateStr) return null;
        const parts = dateStr.split('/');
        return parts.length === 3 ? `${parts[2]}-${parts[1]}-${parts[0]}T00:00:00Z` : null;
    }

    function formatMoney(value) {
        value = value.replace(/\D/g, '');
        value = (Number(value) / 100).toFixed(2);
        value = value.replace('.', ',');
        value = value.replace(/(\d)(?=(\d{3})+(?!\d))/g, '$1.');
        return 'R$ ' + value;
    }

    function unformatMoney(value) {
        return value.replace(/[^\d,]/g, '').replace(',', '.');
    }

    function handleCostInput(input) {
        let value = input.value.replace(/\D/g, '');
        value = value.replace(/^0+/, '');
        
        if (value === '') {
            input.value = 'R$ 0,00';
            return;
        }
        
        while (value.length < 3) {
            value = value + '0';
        }
        
        const cents = value.slice(-2);
        const integers = value.slice(0, -2);
        let formattedValue = integers.replace(/(\d)(?=(\d{3})+(?!\d))/g, '$1.');
        const finalValue = `R$ ${formattedValue},${cents}`;
        input.value = finalValue;
        
        const numberEnd = finalValue.length - 3;
        input.setSelectionRange(numberEnd, numberEnd);
    }

    function updateTableStatus() {
        const visibleRows = $('#licensesTable tr:visible').filter(function() {
            return $(this).find('td').length > 0;
        }).length;

        if (visibleRows === 0) {
            if ($('#noResultsMessage').length === 0) {
                $('#licensesTable').after(
                    '<div id="noResultsMessage" class="alert alert-info text-center">' +
                    'Nenhuma licença encontrada com os filtros selecionados.' +
                    '</div>'
                );
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

    // Atualizar os filtros de status para usar o nome do status
    function applyFilters() {
        const statusFilter = $('#filterStatus').val();
        
        $('#licensesTable tr').each(function() {
            const row = $(this);
            if (row.find('td').length > 0) {
                const status = row.find('.badge').text().trim();
                const matchesStatus = !statusFilter || status === statusFilter;
                
                // ... resto do código de filtro ...
                
                if (matchesStatus /* && outros filtros */) {
                    row.show();
                } else {
                    row.hide();
                }
            }
        });
        
        updateTableStatus();
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
});
