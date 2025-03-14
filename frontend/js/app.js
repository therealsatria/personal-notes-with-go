/**
 * Main Application
 * Handles SPA navigation and initialization
 */
class App {
    constructor() {
        // DOM Elements
        this.navLinks = document.querySelectorAll('nav a');
        this.pages = document.querySelectorAll('.page');
        
        // Initialize
        this.init();
    }
    
    /**
     * Initialize the application
     */
    init() {
        // Add event listeners for navigation
        this.navLinks.forEach(link => {
            link.addEventListener('click', (e) => {
                e.preventDefault();
                const targetPage = link.getAttribute('data-page');
                this.navigateTo(targetPage);
            });
        });
        
        // Initialize toast container if it doesn't exist
        if (!document.getElementById('toast-container')) {
            const toastContainer = document.createElement('div');
            toastContainer.id = 'toast-container';
            document.body.appendChild(toastContainer);
        }
        
        // Check for hash in URL for direct navigation
        const hash = window.location.hash.substring(1);
        if (hash) {
            this.navigateTo(hash);
        }
        
        // Add hash change listener for browser navigation
        window.addEventListener('hashchange', () => {
            const hash = window.location.hash.substring(1);
            if (hash) {
                this.navigateTo(hash);
            }
        });
    }
    
    /**
     * Navigate to a specific page
     * @param {string} pageName - The name of the page to navigate to
     */
    navigateTo(pageName) {
        // Update URL hash
        window.location.hash = pageName;
        
        // Update active nav link
        this.navLinks.forEach(link => {
            if (link.getAttribute('data-page') === pageName) {
                link.classList.add('active');
            } else {
                link.classList.remove('active');
            }
        });
        
        // Show active page
        this.pages.forEach(page => {
            if (page.id === `${pageName}-page`) {
                page.classList.add('active');
            } else {
                page.classList.remove('active');
            }
        });
    }
}

// Initialize the application when DOM is loaded
document.addEventListener('DOMContentLoaded', () => {
    // Initialize the app
    const app = new App();
    
    // Initialize toast service
    if (typeof toastService !== 'undefined') {
        // The container might not be in the DOM yet when the service is created
        toastService.container = document.getElementById('toast-container');
    }
    
    // Show welcome message
    if (typeof toastService !== 'undefined') {
        setTimeout(() => {
            toastService.info('Welcome to Personal Notes App!');
        }, 500);
    }
}); 