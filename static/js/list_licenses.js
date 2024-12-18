$(document).ready(function() {
    // Função para carregar as licenças
    function loadLicenses(filters = {}) {
        console.log('Enviando request com filtros:', filters);

        // Converter IDs para números
        if (filters.software_id) {
            filters.software_id = parseInt(filters.software_id);
        }
        if (filters.status_id) {
            filters.status_id = parseInt(filters.status_id);
        }
        if (filters.department) {
            filters.department = parseInt(filters.department);
        }

        $.ajax({
            url: '/api/v1/licenses/list',
            method: 'GET',
            data: filters,
            success: function(response) {
                console.log('Resposta recebida:', response);
                updateLicensesTable(response.licenses);
            },
            error: function(xhr) {
                console.error('Erro na requisição:', xhr);
                Swal.fire({
                    icon: 'error',
                    title: 'Erro!',
                    text: 'Erro ao carregar licenças: ' + xhr.responseText
                });
            }
        });
    }

    // Função para atualizar a tabela de licenças
    function updateLicensesTable(licenses) {
        console.log('Licenses recebidas:', licenses);
        const tbody = $('#licensesList');
        tbody.empty();
        let totalCost = 0;
        let visibleRowCount = 0;

        if (licenses.length === 0) {
            tbody.append(`
                <tr>
                    <td colspan="10" class="text-center">
                        <div class="empty-state">
                            <i class="fas fa-key"></i>
                            <p>Nenhuma licença encontrada</p>
                        </div>
                    </td>
                </tr>
            `);
            $('#totalCost').text(formatMoney(0));
            return;
        }

        licenses.forEach((license, index) => {
            const statusClass = getStatusClass(license.status.name);
            const periodRenew = getPeriodRenewText(license.period_renew);
            const cost = license.cost || 0;
            totalCost += cost;
            
            const row = $(`
                <tr style="opacity: 0">
                    <td class="text-center align-middle">${license.software ? license.software.name : '-'}</td>
                    <td class="text-center align-middle">
                        ${license.site_url ? `<a href="${license.site_url}" target="_blank" class="btn btn-link" title="Acessar site">
                            <i class="fas fa-external-link-alt"></i>
                        </a>` : ''}
                    </td>
                    <td class="text-center align-middle">${license.type}</td>
                    <td class="text-center align-middle">${formatDate(license.purchase_date)}</td>
                    <td class="text-center align-middle">${formatDate(license.expiry_date)}</td>
                    <td class="text-center align-middle">${license.department ? license.department.name : '-'}</td>
                    <td class="text-center align-middle">${formatMoney(cost)}</td>
                    <td class="text-center align-middle">
                        <span class="badge ${statusClass}" data-status-id="${license.status.id}">
                            ${license.status.name}
                        </span>
                    </td>
                </tr>
            `);

            tbody.append(row);
            
            // Adiciona animação com delay progressivo
            setTimeout(() => {
                row.css({
                    'animation': 'fadeInUp 0.5s ease-out forwards',
                    'animation-delay': `${index * 0.05}s`
                });
            }, 50);
        });

        // Atualiza o total
        $('#totalCost').text(formatMoney(totalCost));
    }

    // Função para formatar a data
    function formatDate(dateString) {
        if (!dateString) return '-';
        const date = new Date(dateString);
        return date.toLocaleDateString('pt-BR');
    }

    // Função para obter o texto do período de renovação
    function getPeriodRenewText(periodRenew) {
        console.log('Processando period_renew:', periodRenew); // Debug do valor recebido
        
        if (!periodRenew) return '-';
        if (!periodRenew.ID && !periodRenew.id) return '-'; // Verifica tanto ID quanto id
        
        const id = periodRenew.ID || periodRenew.id; // Usa qualquer um que estiver disponível
        console.log('ID do period_renew:', id); // Debug do ID
        
        switch (id) {
            case 1:
                return 'Mensal';
            case 2:
                return 'Trimestral';
            case 3:
                return 'Anual';
            default:
                return '-';
        }
    }

    // Função para obter a classe do status
    function getStatusClass(statusName) {
        switch (statusName) {
            case 'Ativa':
                return 'badge-success';
            case 'Próxima do vencimento':
                return 'badge-warning';
            case 'Vencida':
                return 'badge-danger';
            case 'Cancelada':
                return 'badge-secondary';
            default:
                return 'badge-secondary';
        }
    }

    // Função para aplicar filtros
    function applyFilters() {
        const filters = {
            software_id: $('#licenseSearch').val() || '',
            status_id: $('#statusFilter').val() || '',
            date: $('#dateFilter').val() || '',
            department: $('#departmentFilter').val() || ''
        };

        console.log('Aplicando filtros:', filters);
        loadLicenses(filters);
    }

    // Event listeners para todos os filtros
    $('#licenseSearch, #statusFilter, #dateFilter, #departmentFilter').on('change', function() {
        applyFilters();
    });

    // Limpar filtros
    $('#clearFilters').click(function() {
        $('#licenseSearch').val('');
        $('#statusFilter').val('');
        $('#dateFilter').val('');
        $('#departmentFilter').val('');
        applyFilters();
    });

    // Função debounce para evitar múltiplas requisições
    function debounce(func, wait) {
        let timeout;
        return function executedFunction(...args) {
            const later = () => {
                clearTimeout(timeout);
                func(...args);
            };
            clearTimeout(timeout);
            timeout = setTimeout(later, wait);
        };
    }

    // Adicione esta função para formatar valores monetários
    function formatMoney(value) {
        return new Intl.NumberFormat('pt-BR', {
            style: 'currency',
            currency: 'BRL',
            minimumFractionDigits: 2
        }).format(value);
    }

    // Carrega as licenças inicialmente
    loadLicenses();
}); 