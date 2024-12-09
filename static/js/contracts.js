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
        flatpickr.localize(flatpickr.l10ns.pt);
        const flatpickrConfig = {
            dateFormat: "d/m/Y",
            locale: "pt",
            allowInput: true,
            onChange: function(selectedDates, dateStr) {
                // Garante que a data selecionada seja mantida
                this.setDate(dateStr, false);
            }
        };

        // Inicializa o Flatpickr para cada campo de data
        document.querySelectorAll('.datepicker').forEach(input => {
            flatpickr(input, flatpickrConfig);
        });

        // Select2
        $('.select2').select2({
            width: '100%',
            placeholder: "Selecione..."
        });

        // Inicializa o campo de valor sem afetar o total
        const valueInput = document.querySelector('input[name="value"]');
        if (!valueInput.value) {
            valueInput.value = 'R$ 0,00';
        }
        
        // Estilização dos selects de filtro com Select2
        $('#filterStatus, #filterDepartment, #filterBranch, #filterTerminationCondition').select2({
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
        setupEditContractHandler();
        setupClearFiltersHandler();
        setupSaveContractHandler();
        setupValueInputHandler();
        setupDeleteContractHandler();
        setupNewContractHandler();
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

    function setupClearFiltersHandler() {
        $('#clearFilters').click(function() {
            $('#filterStatus').val('').trigger('change');
            $('#filterDepartment').val('').trigger('change');
            $('#filterBranch').val('').trigger('change');
            $('#dateRange').val('');
            applyFilters();
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
            
            // Extrair datas para cálculo do status
            let initialDate = '';
            let finalDate = '';
            
            for (let [key, value] of formData.entries()) {
                if (key === 'department_id' || key === 'branch_id' || key === 'cost_center_id' || key === 'termination_condition_id') {
                    const id = parseInt(value);
                    if (!id) {
                        Swal.fire({
                            icon: 'error',
                            title: 'Erro!',
                            text: `Por favor, selecione um ${
                                key === 'department_id' ? 'departamento' : 
                                key === 'branch_id' ? 'filial' : 
                                key === 'cost_center_id' ? 'centro de custo' :
                                'condição de rescisão'
                            }`
                        });
                        return;
                    }
                    data[key] = id;
                } else if (key === 'value') {
                    data[key] = unformatMoney(value);
                } else if (key === 'initial_date') {
                    initialDate = value;
                    const [day, month, year] = value.split('/');
                    data[key] = `${year}-${month}-${day}T00:00:00Z`;
                } else if (key === 'final_date') {
                    finalDate = value;
                    const [day, month, year] = value.split('/');
                    data[key] = `${year}-${month}-${day}T00:00:00Z`;
                } else if (key !== 'status_id') { // Ignorar o campo de status
                    data[key] = value;
                }
            }

            // Se for edição, verificar se tem aditivos
            if (isEdit) {
                $.ajax({
                    url: `/api/v1/contracts/${contractId}/aditivos`,
                    type: 'GET',
                    async: false,
                    success: function(response) {
                        data.status_id = calculateContractStatus(initialDate, finalDate, response.length > 0);
                    },
                    error: function() {
                        data.status_id = calculateContractStatus(initialDate, finalDate, false);
                    }
                });
            } else {
                data.status_id = calculateContractStatus(initialDate, finalDate, false);
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
        $('#filterStatus, #filterDepartment, #filterBranch, #filterTerminationCondition, #dateRange').on('change', function() {
            applyFilters();
        });
    }

    function applyFilters() {
        const statusFilter = $('#filterStatus').val();
        const departmentFilter = $('#filterDepartment').val();
        const branchFilter = $('#filterBranch').val();
        const terminationConditionFilter = $('#filterTerminationCondition').val();
        const dateRange = $('#dateRange').val();

        let totalValue = 0;

        $('table tbody tr').each(function() {
            const row = $(this);
            let show = true;

            // Status
            if (statusFilter) {
                const statusId = parseInt(row.find('td:eq(8) .badge').data('status-id'));
                if (statusId !== parseInt(statusFilter)) show = false;
            }

            // Departamento
            if (show && departmentFilter) {
                const departmentId = parseInt(row.find('td:eq(3)').data('department-id'));
                if (departmentId !== parseInt(departmentFilter)) show = false;
            }

            // Filial
            if (show && branchFilter) {
                const branchId = parseInt(row.find('td:eq(4)').data('branch-id'));
                if (branchId !== parseInt(branchFilter)) show = false;
            }

            // Condição de Rescisão
            if (show && terminationConditionFilter) {
                const terminationId = parseInt(row.find('td:eq(3)').data('termination-condition-id'));
                if (terminationId !== parseInt(terminationConditionFilter)) show = false;
            }

            // Data
            if (show && dateRange) {
                const [startStr, endStr] = dateRange.split(' - ');
                const start = moment(startStr, 'DD/MM/YYYY');
                const end = moment(endStr, 'DD/MM/YYYY');
                const contractDate = moment(row.find('td:eq(6)').text().trim(), 'DD/MM/YYYY');
                
                if (!contractDate.isBetween(start, end, 'day', '[]')) {
                    show = false;
                }
            }

            if (show) {
                row.show();
                // Calcula o total apenas para linhas visíveis
                const valueText = row.find('td:eq(5)').text().trim();
                const value = parseFloat(valueText.replace('R$ ', '').replace(/\./g, '').replace(',', '.')) || 0;
                totalValue += value;
            } else {
                row.hide();
            }
        });

        // Atualiza o valor total formatado
        $('.table-footer .total-value').text(formatMoney(totalValue));
    }

    // ==================== Utility Functions ====================
    function fillEditForm(contract) {
        const form = $('#contractForm');
        
        form.find('[name="contract_number"]').val(contract.contract_number);
        form.find('[name="name"]').val(contract.name);
        form.find('[name="department_id"]').val(contract.department_id).trigger('change');
        form.find('[name="branch_id"]').val(contract.branch_id).trigger('change');
        form.find('[name="cost_center_id"]').val(contract.cost_center_id).trigger('change');
        form.find('[name="termination_condition_id"]').val(contract.termination_condition_id).trigger('change');
        form.find('[name="value"]').val(formatMoney(contract.value));
        form.find('[name="notes"]').val(contract.notes);
        
        // Ajuste para as datas
        if (contract.initial_date) {
            const initialDate = moment(contract.initial_date).format('DD/MM/YYYY');
            const initialDateInput = form.find('[name="initial_date"]')[0];
            if (initialDateInput._flatpickr) {
                initialDateInput._flatpickr.setDate(initialDate);
            }
        }
        
        if (contract.final_date) {
            const finalDate = moment(contract.final_date).format('DD/MM/YYYY');
            const finalDateInput = form.find('[name="final_date"]')[0];
            if (finalDateInput._flatpickr) {
                finalDateInput._flatpickr.setDate(finalDate);
            }
        }
        
        form.find('[name="status"]').val(contract.status_id).trigger('change');
    }

    function formatDateBR(date) {
        if (!date) return '';
        const d = new Date(date);
        if (isNaN(d.getTime())) return ''; // Verifica se a data é válida
        const day = d.getDate().toString().padStart(2, '0');
        const month = (d.getMonth() + 1).toString().padStart(2, '0');
        const year = d.getFullYear();
        return `${day}/${month}/${year}`;
    }

    function formatMoney(value) {
        // Converte para string com 2 casas decimais
        const parts = value.toFixed(2).split('.');
        
        // Formata a parte inteira com pontos
        const intPart = parts[0].replace(/\B(?=(\d{3})+(?!\d))/g, ".");
        
        // Retorna no formato brasileiro
        return `R$ ${intPart},${parts[1]}`;
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
        // Handler para o botão de Novo Contrato
        $('button[data-target="#addContractModal"]').click(function() {
            // Limpa o formulário
            $('#contractForm')[0].reset();
            
            // Remove qualquer ID de contrato armazenado
            $('#contractForm').removeData('contract-id');
            
            // Reseta os selects do Select2
            $('.select2').val(null).trigger('change');
            
            // Reseta o campo de valor sem afetar o total
            $('input[name="value"]').val('R$ 0,00');
            
            // Reseta as datas do Flatpickr
            document.querySelectorAll('.datepicker').forEach(input => {
                if (input._flatpickr) {
                    input._flatpickr.clear();
                }
            });
            
            // Atualiza o título do modal
            $('#modalTitle').text('Novo Contrato');
        });

        // Também inicializa quando o modal é mostrado
        $('#addContractModal').on('shown.bs.modal', function() {
            if (!$('input[name="value"]').val()) {
                $('input[name="value"]').val('R$ 0,00');
            }
        });
    }

    // Modifique a função que atualiza o valor total para não atualizar quando o modal é aberto
    function updateTotalValue(newValue) {
        // Só atualiza o total quando um contrato é efetivamente salvo
        if (typeof newValue === 'number') {
            const totalValueElement = $('.table-footer .total-value');
            const currentTotal = unformatMoney(totalValueElement.text());
            const newTotal = currentTotal + newValue;
            totalValueElement.text(formatMoney(newTotal));
        }
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

    // Remover o select de status do modal e adicionar esta função
    function calculateContractStatus(initialDate, finalDate, hasAditivo = false) {
        if (hasAditivo) {
            return 4; // ID do status "Renovado por Aditivo"
        }

        const today = moment();
        const start = moment(initialDate, 'DD/MM/YYYY');
        const end = moment(finalDate, 'DD/MM/YYYY');
        
        if (!start.isValid() || !end.isValid()) {
            return null;
        }

        if (today.isAfter(end)) {
            return 3; // ID do status "Vencido"
        }

        // Verifica se está próximo ao vencimento (30 dias)
        const daysUntilExpiration = end.diff(today, 'days');
        if (daysUntilExpiration <= 30) {
            return 2; // ID do status "Próximo ao Vencimento"
        }

        return 1; // ID do status "Em Vigor"
    }

    // Adicione um handler para o fechamento do modal
    $('#addContractModal').on('hidden.bs.modal', function() {
        // Reseta o formulário sem afetar o total
        $('#contractForm')[0].reset();
        $('.select2').val(null).trigger('change');
        $('input[name="value"]').val('R$ 0,00');
        document.querySelectorAll('.datepicker').forEach(input => {
            if (input._flatpickr) {
                input._flatpickr.clear();
            }
        });
    });
}); 