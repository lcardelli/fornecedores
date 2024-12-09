document.addEventListener('DOMContentLoaded', function() {
    // Inicialização
    initializeComponents();
    setupEventHandlers();
    setupFilters();
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

        // Select2
        $('.select2').select2({
            width: '100%',
            placeholder: "Selecione..."
        });
    }

    // ==================== Event Handlers ====================
    function setupEventHandlers() {
        setupClearFiltersHandler();
    }

    function setupClearFiltersHandler() {
        $('#clearFilters').click(function() {
            // Limpa todos os selects e dispara o evento change
            $('#filterStatus').val('').trigger('change');
            $('#filterDepartment').val('').trigger('change');
            $('#filterBranch').val('').trigger('change');
            $('#filterCostCenter').val('').trigger('change');
            $('#yearStart').val('').trigger('change');
            $('#yearEnd').val('').trigger('change');
            
            // Aplica os filtros após limpar
            applyFilters();
        });
    }

    // ==================== Filtros ====================
    function setupFilters() {
        // Adiciona event listeners para todos os filtros
        $('#filterStatus, #filterDepartment, #filterBranch, #filterCostCenter, #yearStart, #yearEnd').on('change', function() {
            applyFilters();
        });
    }

    function applyFilters() {
        const rows = document.querySelectorAll('tbody tr');
        const statusFilter = $('#filterStatus').val();
        const departmentFilter = $('#filterDepartment').val();
        const branchFilter = $('#filterBranch').val();
        const costCenterFilter = $('#filterCostCenter').val();
        const yearStartFilter = $('#yearStart').val();
        const yearEndFilter = $('#yearEnd').val();

        let totalValue = 0;

        // Função auxiliar para extrair valor monetário
        function extractMoneyValue(text) {
            // Remove 'R$ ' e converte pontos de milhar para nada e vírgula para ponto
            const numStr = text
                .replace('R$ ', '')
                .replace(/\./g, '')
                .replace(',', '.');
            return parseFloat(numStr);
        }

        rows.forEach(row => {
            let show = true;

            // Filtro de Status
            if (statusFilter) {
                const statusBadge = row.querySelector('span[data-status-id]');
                const statusId = statusBadge ? statusBadge.dataset.statusId : '';
                if (statusId !== statusFilter) show = false;
            }

            // Filtro de Departamento
            if (show && departmentFilter) {
                const departmentCell = row.querySelector('td[data-department-id]');
                const departmentId = departmentCell ? departmentCell.dataset.departmentId : '';
                if (departmentId !== departmentFilter) show = false;
            }

            // Filtro de Filial
            if (show && branchFilter) {
                const branchCell = row.querySelector('td[data-branch-id]');
                const branchId = branchCell ? branchCell.dataset.branchId : '';
                if (branchId !== branchFilter) show = false;
            }

            // Filtro de Centro de Custo
            if (show && costCenterFilter) {
                const costCenterCell = row.querySelector('td[data-cost-center-id]');
                const costCenterId = costCenterCell ? costCenterCell.dataset.costCenterId : '';
                if (costCenterId !== costCenterFilter) show = false;
            }

            // Filtro de Ano
            if (show && (yearStartFilter || yearEndFilter)) {
                const initialDate = row.querySelector('td:nth-child(9)').textContent;
                const finalDate = row.querySelector('td:nth-child(10)').textContent;
                
                const initialYear = parseInt(initialDate.split('/')[2]);
                const finalYear = parseInt(finalDate.split('/')[2]);

                if (yearStartFilter && initialYear < parseInt(yearStartFilter)) show = false;
                if (yearEndFilter && finalYear > parseInt(yearEndFilter)) show = false;
            }

            // Mostra ou esconde a linha
            row.style.display = show ? '' : 'none';

            // Soma o valor se a linha estiver visível
            if (show) {
                const valueCell = row.querySelector('td:nth-child(7)'); // Coluna do valor
                if (valueCell) {
                    const value = extractMoneyValue(valueCell.textContent);
                    if (!isNaN(value)) {
                        totalValue += value;
                    }
                }
            }
        });

        // Atualiza o valor total
        const formattedValue = new Intl.NumberFormat('pt-BR', {
            style: 'currency',
            currency: 'BRL'
        }).format(totalValue);
        
        $('.total-value').text(formattedValue);
    }

    // ==================== Scroll da Tabela ====================
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

    // ==================== Calculadora ====================
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
                text: 'Não foi possível obter as taxas de câmbio atualizadas.',
                footer: 'Detalhes: ' + error.message
            });
            return null;
        }
    }

    // Função para calcular os valores
    async function calculateContractValues() {
        const baseValue = parseFloat($('#baseValue').val()) || 0;
        const currency = $('#currency').val();
        const period = $('#period').val();
        const duration = parseInt($('#duration').val()) || 12;

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

        // Calcula valores por período
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

        const quarterlyValue = monthlyValue * 3;
        const semiannualValue = monthlyValue * 6;
        const annualValue = monthlyValue * 12;
        const totalValue = monthlyValue * duration;

        // Atualiza os resultados na interface
        $('#monthlyResult').text(formatMoney(monthlyValue));
        $('#quarterlyResult').text(formatMoney(quarterlyValue));
        $('#semiannualResult').text(formatMoney(semiannualValue));
        $('#annualResult').text(formatMoney(annualValue));
        $('#totalResult').text(formatMoney(totalValue));
    }

    // Event listeners para atualização automática da calculadora
    $('#baseValue, #currency, #period, #duration').on('change input', function() {
        calculateContractValues();
    });

    // Inicializa os valores quando o modal é aberto
    $('#calculatorModal').on('shown.bs.modal', function() {
        calculateContractValues();
    });

    // Event listener para limpar os campos quando o modal for fechado
    $('#calculatorModal').on('hidden.bs.modal', function() {
        $('#baseValue').val('');
        $('#currency').val('BRL');
        $('#period').val('monthly');
        $('#duration').val('12');
        
        $('#monthlyResult').text('R$ 0,00');
        $('#quarterlyResult').text('R$ 0,00');
        $('#semiannualResult').text('R$ 0,00');
        $('#annualResult').text('R$ 0,00');
        $('#totalResult').text('R$ 0,00');
    });

    // Função auxiliar para formatar valores monetários
    function formatMoney(value) {
        return new Intl.NumberFormat('pt-BR', {
            style: 'currency',
            currency: 'BRL'
        }).format(value);
    }

    // ==================== Anexos ====================
    window.showAttachments = function(contractId, attachments) {
        console.log('showAttachments called:', { contractId, attachments });
        const attachmentsList = $('#attachmentsList');
        attachmentsList.empty();

        // Converte a string de template Go para objeto JavaScript
        const attachmentsData = JSON.parse(decodeURIComponent(attachments));
        console.log('parsed attachments:', attachmentsData);

        attachmentsData.forEach(attachment => {
            const listItem = $(`
                <div class="list-group-item">
                    <div>
                        <i class="fas fa-file-pdf"></i>
                        <span>${attachment.Name}</span>
                    </div>
                    <button class="download-btn" onclick="downloadAttachment(event, ${attachment.ID}, '${attachment.Name}')">
                        <i class="fas fa-download"></i> Baixar
                    </button>
                </div>
            `);
            attachmentsList.append(listItem);
        });

        $('#attachmentsModal').modal('show');
    };

    // Função para download de anexos
    window.downloadAttachment = function(event, id, filename) {
        event.preventDefault();
        event.stopPropagation();
        
        const downloadBtn = $(event.currentTarget);
        const originalContent = downloadBtn.html();
        downloadBtn.html('<i class="fas fa-spinner fa-spin"></i> Baixando...');
        downloadBtn.prop('disabled', true);
        
        fetch(`/api/v1/contracts/download/${id}`)
            .then(response => {
                if (!response.ok) {
                    throw new Error('Erro ao baixar arquivo');
                }
                return response.blob();
            })
            .then(blob => {
                const url = window.URL.createObjectURL(blob);
                const a = document.createElement('a');
                a.style.display = 'none';
                a.href = url;
                a.download = filename;
                document.body.appendChild(a);
                a.click();
                window.URL.revokeObjectURL(url);
                document.body.removeChild(a);
                downloadBtn.html(originalContent);
                downloadBtn.prop('disabled', false);
            })
            .catch(error => {
                console.error('Erro:', error);
                downloadBtn.html(originalContent);
                downloadBtn.prop('disabled', false);
                Swal.fire({
                    icon: 'error',
                    title: 'Erro!',
                    text: 'Erro ao baixar o arquivo'
                });
            });
    };
});
