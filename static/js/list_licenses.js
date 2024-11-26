$(document).ready(function() {
    // Função para carregar as licenças
    function loadLicenses(filters = {}) {
        $.ajax({
            url: '/api/v1/licenses/list',
            method: 'GET',
            data: filters,
            success: function(response) {
                updateLicensesTable(response.licenses);
            },
            error: function(xhr) {
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
                    <td colspan="6" class="text-center">
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
            const statusClass = getStatusClass(license.status);
            tbody.append(`
                <tr>
                    <td>${license.license_key}</td>
                    <td>${license.software ? license.software.name : '-'}</td>
                    <td>${license.type}</td>
                    <td>${formatDate(license.purchase_date)}</td>
                    <td>${formatDate(license.expiry_date)}</td>
                    <td><span class="badge ${statusClass}">${license.status}</span></td>
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

    // Função para definir a classe do status
    function getStatusClass(status) {
        const statusClasses = {
            'active': 'badge-success',
            'inactive': 'badge-secondary',
            'expired': 'badge-danger'
        };
        return statusClasses[status] || 'badge-secondary';
    }

    // Event listeners para filtros
    $('#licenseSearch').on('input', debounce(function() {
        applyFilters();
    }, 300));

    $('#productFilter, #statusFilter, #dateFilter').on('change', function() {
        applyFilters();
    });

    // Função para aplicar filtros
    function applyFilters() {
        const filters = {
            search: $('#licenseSearch').val(),
            product: $('#productFilter').val(),
            status: $('#statusFilter').val(),
            date: $('#dateFilter').val()
        };
        loadLicenses(filters);
    }

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