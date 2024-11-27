$(document).ready(function() {
    // Função para verificar e abrir o submenu ativo
    function checkActiveSubmenu() {
        const activeLink = $('.nav-link.active');
        if (activeLink.length) {
            // Encontra o submenu pai do link ativo
            const parentSubmenu = activeLink.closest('.collapse.nav');
            if (parentSubmenu.length) {
                // Abre o submenu pai
                parentSubmenu.addClass('show');
                parentSubmenu.prev('.dropdown-toggle').attr('aria-expanded', 'true');
                
                // Salva o estado no localStorage
                const menuId = parentSubmenu.attr('id');
                localStorage.setItem('menu-#' + menuId, 'open');
            }
        }
    }

    // Função para verificar qual menu deve estar aberto baseado na URL atual
    function checkCurrentPath() {
        const currentPath = window.location.pathname;
        
        // Mapeamento de caminhos para seus respectivos submenus
        const menuMappings = {
            '/api/v1/catalogo': '#fornecedoresSubmenu',
            '/api/v1/lista-fornecedores': '#fornecedoresSubmenu',
            '/api/v1/cadastro-categoria': '#fornecedoresSubmenu',
            '/api/v1/services': '#fornecedoresSubmenu',
            '/api/v1/produtos': '#fornecedoresSubmenu',
            
            '/api/v1/licenses/active': '#licencasSubmenu',
            '/api/v1/licenses/expired': '#licencasSubmenu',
            '/api/v1/licenses/manage': '#licencasSubmenu',
            '/api/v1/licenses/software': '#licencasSubmenu',
            
            '/api/v1/manage-users': '#adminSubmenu'
        };

        const targetSubmenu = menuMappings[currentPath];
        if (targetSubmenu) {
            $(targetSubmenu).addClass('show');
            $(`a[href="${targetSubmenu}"]`).attr('aria-expanded', 'true');
            localStorage.setItem('menu-' + targetSubmenu, 'open');
        }
    }

    // Executa ao carregar a página
    checkActiveSubmenu();
    checkCurrentPath();

    // Handler para cliques nos toggles do dropdown
    $('.dropdown-toggle').click(function(e) {
        e.preventDefault();
        const target = $(this).attr('href');
        const submenu = $(target);
        
        // Toggle do submenu atual
        if (submenu.hasClass('show')) {
            submenu.removeClass('show');
            $(this).attr('aria-expanded', 'false');
            localStorage.removeItem('menu-' + target);
        } else {
            submenu.addClass('show');
            $(this).attr('aria-expanded', 'true');
            localStorage.setItem('menu-' + target, 'open');
        }
    });

    // Restaura o estado dos menus do localStorage
    $('.dropdown-toggle').each(function() {
        const menuId = $(this).attr('href');
        const savedState = localStorage.getItem('menu-' + menuId);
        
        if (savedState === 'open') {
            $(this).attr('aria-expanded', 'true');
            $(menuId).addClass('show');
        }
    });

    // Previne o fechamento do submenu quando clica em links dentro dele
    $('.collapse.nav .nav-link').click(function(e) {
        e.stopPropagation(); // Impede a propagação do evento
        const parentSubmenu = $(this).closest('.collapse.nav');
        const menuId = '#' + parentSubmenu.attr('id');
        localStorage.setItem('menu-' + menuId, 'open');
    });
});