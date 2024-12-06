document.addEventListener('DOMContentLoaded', function() {
    // Inicialização
    initializeComponents();
    setupEventHandlers();
    setupFilters();
    setupCheckboxes();
    setupTableScroll();

    // ==================== Inicialização de Componentes ====================
    function initializeComponents() {
        // Configuração do DateRangePicker
        $('#dateRange').daterangepicker({
            autoUpdateInput: false,
            locale: {
                format: 'DD/MM/YYYY',
                applyLabel: 'Aplicar',
                cancelLabel: 'Limpar',
                fromLabel: 'De',
                toLabel: 'Até',
                customRangeLabel: 'Personalizado',
                daysOfWeek: ['Dom', 'Seg', 'Ter', 'Qua', 'Qui', 'Sex', 'Sáb'],
                monthNames: ['Janeiro', 'Fevereiro', 'Março', 'Abril', 'Maio', 'Junho',
                           'Julho', 'Agosto', 'Setembro', 'Outubro', 'Novembro', 'Dezembro']
            }
        });

        // Handlers do DateRangePicker
        $('#dateRange').on('apply.daterangepicker', function(ev, picker) {
            $(this).val(picker.startDate.format('DD/MM/YYYY') + ' - ' + picker.endDate.format('DD/MM/YYYY'));
            applyFilters();
        });

        $('#dateRange').on('cancel.daterangepicker', function(ev, picker) {
            $(this).val('');
            applyFilters();
        });

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
                }
            },
            allowInput: true,
            wrap: true
        };

        // Inicializa o Flatpickr
        document.querySelectorAll('.datepicker-wrap').forEach(wrapper => {
            flatpickr(wrapper, flatpickrConfig);
        });

        // Select2
        $('.select2').select2({
            width: '100%',
            placeholder: "Selecione..."
        });

        // Inicializa o campo de valor
        const valueInput = document.querySelector('input[name="value"]');
        valueInput.value = 'R$ 0,00';
        
        // Estilização dos selects de filtro com Select2
        $('#filterStatus, #filterDepartment, #filterBranch').select2({
            width: '100%',
            placeholder: 'Selecione...',
            allowClear: true,
            dropdownParent: $('body'),
            containerCssClass: 'manage-contracts-select-container',
            dropdownCssClass: 'manage-contracts-dropdown'
        });

        // Inicialização dos selects do modal
        $('.manage-contracts-modal-select').select2({
            width: '100%',
            placeholder: "Selecione...",
            allowClear: true,
            dropdownParent: $('#addContractModal'),
            containerCssClass: 'manage-contracts-modal-container',
            dropdownCssClass: 'manage-contracts-modal-dropdown'
        });
    }

    // ==================== Event Handlers ====================
    function setupEventHandlers() {
        setupClearFiltersHandler();
        setupEditContractHandler();
        setupSaveContractHandler();
        setupValueInputHandler();
        setupDeleteContractHandler();
        setupNewContractHandler();
    }

    function setupClearFiltersHandler() {
        $('#clearFilters').click(function() {
            $('#filterStatus').val('').trigger('change');
            $('#filterDepartment').val('').trigger('change');
            $('#filterBranch').val('').trigger('change');
            $('#dateRange').val('');
            applyFilters();
        });
    }

    function setupEditContractHandler() {
        $('.edit-contract').click(function() {
            const contractId = $(this).data('id');
            $('#contractForm')[0].reset();
            $('#contractForm').removeData('contract-id');
            
            $.ajax({
                url: `/api/v1/contracts/${contractId}`,
                type: 'GET',
                success: function(response) {
                    fillEditForm(response);
                    $('#modalTitle').text('Editar Contrato');
                    $('#contractForm').data('contract-id', contractId);
                    $('#addContractModal').modal('show');
                },
                error: function(xhr) {
                    let errorMessage = 'Erro ao carregar dados do contrato';
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

    function setupSaveContractHandler() {
        $('#saveContract').click(function() {
            const form = $('#contractForm');
            const contractId = form.data('contract-id');
            const isEdit = !!contractId;
            
            if (!form[0].checkValidity()) {
                form[0].reportValidity();
                return;
            }

            const formData = new FormData(form[0]);
            const data = {};
            
            for (let [key, value] of formData.entries()) {
                if (key === 'department_id' || key === 'branch_id') {
                    const id = parseInt(value);
                    if (!id) {
                        Swal.fire({
                            icon: 'error',
                            title: 'Erro!',
                            text: `Por favor, selecione um ${key === 'department_id' ? 'departamento' : 'filial'}`
                        });
                        return;
                    }
                    data[key] = id;
                } else if (key === 'value') {
                    data[key] = unformatMoney(value);
                } else if (key === 'initial_date' || key === 'final_date') {
                    const [day, month, year] = value.split('/');
                    data[key] = `${year}-${month}-${day}T00:00:00Z`;
                } else {
                    data[key] = value;
                }
            }

            const url = isEdit ? `/api/v1/contracts/${contractId}` : '/api/v1/contracts';
            const method = isEdit ? 'PUT' : 'POST';

            $.ajax({
                url: url,
                type: method,
                contentType: 'application/json',
                data: JSON.stringify(data),
                success: function(response) {
                    $('#addContractModal').modal('hide');
                    Swal.fire({
                        icon: 'success',
                        title: 'Sucesso!',
                        text: `Contrato ${isEdit ? 'atualizado' : 'cadastrado'} com sucesso!`
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
                        text: `Erro ao ${isEdit ? 'atualizar' : 'cadastrar'} contrato: ${errorMessage}`
                    });
                }
            });
        });
    }

    function setupValueInputHandler() {
        const valueInput = document.querySelector('input[name="value"]');
        valueInput.addEventListener('input', function(e) {
            handleValueInput(this);
        });
        
        valueInput.addEventListener('focus', function() {
            this.setSelectionRange(3, this.value.length - 3);
        });
    }

    function setupDeleteContractHandler() {
        $('.delete-contract').click(function() {
            const contractId = $(this).data('id');
            
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
                        url: `/api/v1/contracts/${contractId}`,
                        type: 'DELETE',
                        success: function() {
                            Swal.fire({
                                icon: 'success',
                                title: 'Deletado!',
                                text: 'O contrato foi deletado com sucesso.'
                            }).then(() => {
                                window.location.reload();
                            });
                        },
                        error: function(xhr) {
                            Swal.fire({
                                icon: 'error',
                                title: 'Erro!',
                                text: 'Erro ao deletar contrato: ' + xhr.responseText
                            });
                        }
                    });
                }
            });
        });
    }

    // ==================== Checkbox Handlers ====================
    function setupCheckboxes() {
        $('.select-all-checkbox').change(function() {
            const isChecked = $(this).prop('checked');
            $('.contract-checkbox').prop('checked', isChecked);
            updateDeleteSelectedButton();
        });

        $('.contract-checkbox').change(function() {
            const allChecked = $('.contract-checkbox:checked').length === $('.contract-checkbox').length;
            $('.select-all-checkbox').prop('checked', allChecked);
            updateDeleteSelectedButton();
        });
    }

    function updateDeleteSelectedButton() {
        const checkedCount = $('.contract-checkbox:checked').length;
        if (checkedCount > 0) {
            $('#deleteSelectedBtn').fadeIn(200).css('display', 'inline-flex');
        } else {
            $('#deleteSelectedBtn').fadeOut(200);
        }
    }

    // ==================== Filter Handlers ====================
    function setupFilters() {
        $('#filterStatus, #filterDepartment, #filterBranch, #dateRange').on('change', function() {
            applyFilters();
        });
    }

    function applyFilters() {
        const statusFilter = $('#filterStatus').val() || '';
        const departmentFilter = $('#filterDepartment').val() || '';
        const branchFilter = $('#filterBranch').val() || '';
        const dateRange = $('#dateRange').val();

        let totalValue = 0;
        let visibleRowCount = 0;

        let startDate, endDate;
        if (dateRange) {
            const [start, end] = dateRange.split(' - ');
            startDate = moment(start, 'DD/MM/YYYY');
            endDate = moment(end, 'DD/MM/YYYY');
        }

        $('#contractsTable tr').each(function() {
            const row = $(this);
            if (row.find('td').length > 0) {
                const statusId = row.find('.badge').data('status-id');
                const departmentId = row.find('td:eq(3)').data('department-id');
                const branchId = row.find('td:eq(4)').data('branch-id');
                
                // Pega a data inicial do contrato
                const initialDateText = row.find('td:eq(6)').text();
                const contractDate = moment(initialDateText, 'DD/MM/YYYY');
                
                const matchesStatus = !statusFilter || statusId === parseInt(statusFilter);
                const matchesDepartment = !departmentFilter || departmentId === parseInt(departmentFilter);
                const matchesBranch = !branchFilter || branchId === parseInt(branchFilter);
                const matchesDate = !dateRange || (contractDate.isBetween(startDate, endDate, 'day', '[]'));

                if (matchesStatus && matchesDepartment && matchesBranch && matchesDate) {
                    row.show();
                    row.css({
                        'animation': 'none',
                        'opacity': '0'
                    });
                    setTimeout(() => {
                        row.css({
                            'animation': 'fadeIn 0.3s ease-out forwards',
                            'animation-delay': `${visibleRowCount * 0.05}s`
                        });
                    }, 0);
                    visibleRowCount++;
                    
                    const valueText = row.find('td:eq(5)').text().trim();
                    const value = parseFloat(valueText.replace('R$ ', '').replace('.', '').replace(',', '.')) || 0;
                    totalValue += value;
                } else {
                    row.hide();
                }
            }
        });

        const formattedTotal = formatMoneyBR(totalValue);
        $('.table-footer .total-value').text(formattedTotal);

        updateTableStatus();
    }

    // ==================== Utility Functions ====================
    function fillEditForm(contract) {
        const form = $('#contractForm');
        
        form.find('[name="contract_number"]').val(contract.contract_number);
        form.find('[name="name"]').val(contract.name);
        form.find('[name="department_id"]').val(contract.department_id).trigger('change');
        form.find('[name="branch_id"]').val(contract.branch_id).trigger('change');
        form.find('[name="value"]').val(formatMoney(contract.value));
        form.find('[name="notes"]').val(contract.notes);
        
        if (contract.initial_date) {
            const initialDate = new Date(contract.initial_date);
            const formattedInitialDate = formatDateBR(initialDate);
            const initialPicker = form.find('[name="initial_date"]')[0]._flatpickr;
            initialPicker.setDate(formattedInitialDate, true);
        }
        
        if (contract.final_date) {
            const finalDate = new Date(contract.final_date);
            const formattedFinalDate = formatDateBR(finalDate);
            const finalPicker = form.find('[name="final_date"]')[0]._flatpickr;
            finalPicker.setDate(formattedFinalDate, true);
        }
        
        form.find('[name="status"]').val(contract.status_id).trigger('change');
    }

    function formatDateBR(date) {
        if (!date) return '';
        const d = new Date(date);
        const day = d.getDate().toString().padStart(2, '0');
        const month = (d.getMonth() + 1).toString().padStart(2, '0');
        const year = d.getFullYear();
        return `${day}/${month}/${year}`;
    }

    function formatMoney(value) {
        return new Intl.NumberFormat('pt-BR', {
            style: 'currency',
            currency: 'BRL'
        }).format(value);
    }

    function unformatMoney(value) {
        return parseFloat(value.replace(/[^\d,]/g, '').replace(',', '.'));
    }

    function handleValueInput(input) {
        let value = input.value.replace(/\D/g, '');
        
        if (value === '') {
            input.value = 'R$ 0,00';
            return;
        }
        
        value = (parseFloat(value) / 100).toFixed(2);
        input.value = formatMoney(value);
    }

    function updateTableStatus() {
        const visibleRows = $('#contractsTable tr:visible').length;
        if (visibleRows === 0) {
            if ($('#noResultsMessage').length === 0) {
                $('#contractsTable').after(`
                    <div id="noResultsMessage" class="text-center p-4">
                        <div class="empty-state">
                            <i class="fas fa-file-contract fa-3x mb-3"></i>
                            <p class="h5">Nenhum contrato encontrado</p>
                            <p class="text-muted">Tente ajustar seus filtros de busca</p>
                        </div>
                    </div>
                `);
            }
        } else {
            $('#noResultsMessage').remove();
        }
    }

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

    function setupNewContractHandler() {
        $('button[data-target="#addContractModal"]').click(function() {
            // Limpa o formulário
            $('#contractForm')[0].reset();
            $('#contractForm').removeData('contract-id');
            
            // Reseta os selects do Select2
            $('.select2').val(null).trigger('change');
            
            // Inicializa o campo de valor com R$ 0,00
            $('input[name="value"]').val(formatMoney(0));
            
            // Reseta as datas do Flatpickr
            const initialPicker = $('input[name="initial_date"]')[0]._flatpickr;
            const finalPicker = $('input[name="final_date"]')[0]._flatpickr;
            initialPicker.clear();
            finalPicker.clear();
            
            $('#modalTitle').text('Novo Contrato');
        });

        // Também inicializa quando o modal é mostrado
        $('#addContractModal').on('shown.bs.modal', function() {
            if (!$('input[name="value"]').val()) {
                $('input[name="value"]').val(formatMoney(0));
            }
        });
    }

    // Função auxiliar para formatar valor monetário
    function formatMoney(value) {
        return new Intl.NumberFormat('pt-BR', {
            style: 'currency',
            currency: 'BRL',
            minimumFractionDigits: 2,
            maximumFractionDigits: 2
        }).format(value);
    }
}); 