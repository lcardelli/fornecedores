$(document).ready(function() {
    // Função para carregar as licenças
    function loadLicenses(filters = {}) {
        console.log('Enviando request com filtros:', filters); // Debug

        $.ajax({
            url: '/api/v1/licenses/list',
            method: 'GET',
            data: filters,
            success: function(response) {
                console.log('Resposta recebida:', response); // Debug
                updateLicensesTable(response.licenses);
            },
            error: function(xhr) {
                console.error('Erro na requisição:', xhr); // Debug
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
        const tbody = $('#licensesList');
        tbody.empty();

        if (licenses.length === 0) {
            tbody.append(`
                <tr>
                    <td colspan="7" class="text-center">
                        <div class="empty-state">
                            <i class="fas fa-key"></i>
                            <p>Nenhuma licença encontrada</p>
                        </div>
                    </td>
                </tr>
            `);
            return;
        }

        licenses.forEach(license => {
            const statusClass = getStatusClass(license.status.name);
            const periodRenew = getPeriodRenewText(license.period_renew);
            
            tbody.append(`
                <tr>
                    <td class="text-center align-middle">${license.software ? license.software.name : '-'}</td>
                    <td class="text-center align-middle">${license.type}</td>
                    <td class="text-center align-middle">${license.quantity || '-'}</td>
                    <td class="text-center align-middle">${periodRenew}</td>
                    <td class="text-center align-middle">${formatDate(license.purchase_date)}</td>
                    <td class="text-center align-middle">${formatDate(license.expiry_date)}</td>
                    <td class="text-center align-middle">
                        <span class="badge ${statusClass}" data-status-id="${license.status.id}">
                            ${license.status.name}
                        </span>
                    </td>
                </tr>
            `);
        });
    }

    // Função para formatar a data
    function formatDate(dateString) {
        if (!dateString) return '-';
        const date = new Date(dateString);
        return date.toLocaleDateString('pt-BR');
    }

    // Função para obter o texto do período de renovação
    function getPeriodRenewText(period) {
        switch (period) {
            case 1: return 'Mensal';
            case 3: return 'Trimestral';
            case 12: return 'Anual';
            default: return '-';
        }
    }

    // Função para definir a classe do status
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

    // Função para aplicar filtros
    function applyFilters() {
        const filters = {
            search: $('#licenseSearch').val() || '',
            status: $('#statusFilter').val() || '',
            date: $('#dateFilter').val() || ''
        };

        console.log('Aplicando filtros:', filters); // Debug

        loadLicenses(filters);
    }

    // Event listeners para filtros
    $('#licenseSearch').on('input', debounce(function() {
        applyFilters();
    }, 300));

    $('#statusFilter, #dateFilter').on('change', function() {
        console.log('Filtro alterado:', $(this).attr('id')); // Debug
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

    // Carrega as licenças inicialmente
    loadLicenses();
}); 