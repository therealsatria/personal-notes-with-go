/**
 * Main Application
 * Handles SPA navigation and initialization
 */
class App {
    constructor() {
        // DOM Elements
        this.navLinks = document.querySelectorAll('nav a');
        this.pages = document.querySelectorAll('.page');
        this.addNoteBtn = document.getElementById('add-note-btn');
        this.addCategoryBtn = document.getElementById('add-category-btn');
        this.encryptionStatusBanner = document.getElementById('encryption-status-banner');
        this.encryptionStatusMessage = document.getElementById('encryption-status-message');
        
        // Initialize
        this.init();
    }
    
    /**
     * Initialize the application
     */
    init() {
        console.log("App initializing...");
        
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
        
        // Add event listener for close button on encryption status banner
        if (this.encryptionStatusBanner) {
            const closeBtn = this.encryptionStatusBanner.querySelector('.btn-close-banner');
            if (closeBtn) {
                closeBtn.addEventListener('click', () => {
                    this.encryptionStatusBanner.classList.add('hidden');
                });
            }
        }
        
        // Log DOM elements for debugging
        console.log("Encryption status elements:", {
            banner: this.encryptionStatusBanner,
            message: this.encryptionStatusMessage
        });
        
        // Check encryption status
        setTimeout(() => {
            this.checkEncryptionStatus();
        }, 1000);
        
        // Add listener for encryption status changes
        if (typeof encryptionStatusService !== 'undefined') {
            encryptionStatusService.addListener((isValid, message) => {
                console.log("Encryption status listener called:", isValid, message);
                this.updateEncryptionUI(isValid, message);
            });
        } else {
            console.error("encryptionStatusService is not defined!");
        }
        
        // Add test button to floating navbar for debugging
        this.addEncryptionTestButton();
        
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
     * Add a test button to the floating navbar for debugging
     */
    addEncryptionTestButton() {
        const floatingNavContent = document.querySelector('.floating-nav-content');
        if (floatingNavContent) {
            const testButton = document.createElement('button');
            testButton.className = 'btn btn-primary btn-icon';
            testButton.title = 'Test Encryption Status';
            testButton.innerHTML = '<i class="fas fa-sync"></i>';
            testButton.addEventListener('click', () => {
                this.checkEncryptionStatus();
                toastService.info('Checking encryption status...');
            });
            floatingNavContent.appendChild(testButton);
        }
    }
    
    /**
     * Check the encryption status
     */
    async checkEncryptionStatus() {
        console.log("Checking encryption status...");
        if (typeof encryptionStatusService !== 'undefined') {
            const isValid = await encryptionStatusService.checkStatus();
            const message = encryptionStatusService.getStatusMessage();
            console.log("Encryption status checked:", isValid, message);
            this.updateEncryptionUI(isValid, message);
        } else {
            console.error("encryptionStatusService is not defined!");
        }
    }
    
    /**
     * Update UI based on encryption status
     * @param {boolean} isValid - Whether encryption is valid
     * @param {string} message - Status message
     */
    updateEncryptionUI(isValid, message) {
        console.log("Updating encryption UI:", isValid, message); // Debug log
        
        // Update buttons
        if (this.addNoteBtn) {
            this.addNoteBtn.disabled = !isValid;
            this.addNoteBtn.title = isValid ? 'Add Note' : 'Disabled: Encryption system not properly initialized';
        }
        
        if (this.addCategoryBtn) {
            this.addCategoryBtn.disabled = !isValid;
            this.addCategoryBtn.title = isValid ? 'Add Category' : 'Disabled: Encryption system not properly initialized';
        }
        
        // Update status banner
        if (this.encryptionStatusBanner && this.encryptionStatusMessage) {
            // Get the icon element
            const iconElement = this.encryptionStatusBanner.querySelector('i');
            
            if (!isValid) {
                // Error state
                this.encryptionStatusBanner.classList.remove('hidden');
                this.encryptionStatusBanner.classList.remove('success');
                this.encryptionStatusBanner.classList.add('error');
                this.encryptionStatusMessage.textContent = message;
                
                // Update icon
                if (iconElement) {
                    iconElement.className = 'fas fa-exclamation-triangle';
                }
            } else {
                // Success state
                this.encryptionStatusBanner.classList.remove('hidden');
                this.encryptionStatusBanner.classList.remove('error');
                this.encryptionStatusBanner.classList.add('success');
                this.encryptionStatusMessage.textContent = "Encryption system is working correctly.";
                
                // Update icon
                if (iconElement) {
                    iconElement.className = 'fas fa-shield-alt';
                }
                
                // Hide success message after 3 seconds
                setTimeout(() => {
                    this.encryptionStatusBanner.classList.add('hidden');
                }, 3000);
            }
        } else {
            console.error("Encryption status banner elements not found:", 
                         this.encryptionStatusBanner, 
                         this.encryptionStatusMessage);
        }
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