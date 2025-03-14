/**
 * Key Generator Component
 * Handles key generation functionality
 */
class KeyGeneratorComponent {
    constructor() {
        // DOM Elements
        this.keyGeneratorBtn = document.getElementById('key-generator-btn');
        this.keyGeneratorModal = document.getElementById('key-generator-modal');
        this.keyGeneratorForm = document.getElementById('key-generator-form');
        this.keyInput = document.getElementById('key-input');
        this.generatedKeyInput = document.getElementById('generated-key');
        this.copyKeyBtn = document.getElementById('copy-key-btn');

        // Initialize
        this.init();
    }

    /**
     * Initialize the component
     */
    init() {
        // Add event listeners
        this.keyGeneratorBtn.addEventListener('click', () => this.openModal());
        this.keyGeneratorForm.addEventListener('submit', (e) => this.handleGenerateKey(e));
        this.copyKeyBtn.addEventListener('click', () => this.copyGeneratedKey());

        // Close modal buttons
        const closeButtons = this.keyGeneratorModal.querySelectorAll('.close-modal');
        closeButtons.forEach(button => {
            button.addEventListener('click', () => this.closeModal());
        });
    }

    /**
     * Open the key generator modal
     */
    openModal() {
        this.keyGeneratorModal.classList.add('active');
        this.keyInput.value = '';
        this.generatedKeyInput.value = '';
    }

    /**
     * Close the key generator modal
     */
    closeModal() {
        this.keyGeneratorModal.classList.remove('active');
    }

    /**
     * Handle key generation form submission
     * @param {Event} event - The form submit event
     */
    async handleGenerateKey(event) {
        event.preventDefault();

        try {
            const inputText = this.keyInput.value;
            const response = await apiService.generateKey(inputText);
            this.generatedKeyInput.value = response.key;
            toastService.success('Key generated successfully');
        } catch (error) {
            toastService.error('Failed to generate key: ' + error.message);
        }
    }

    /**
     * Copy the generated key to clipboard
     */
    copyGeneratedKey() {
        const key = this.generatedKeyInput.value;
        if (!key) {
            toastService.error('No key to copy');
            return;
        }

        navigator.clipboard.writeText(key)
            .then(() => {
                toastService.success('Key copied to clipboard');
            })
            .catch(error => {
                toastService.error('Failed to copy key: ' + error.message);
            });
    }
}

// Create and export a singleton instance
const keyGeneratorComponent = new KeyGeneratorComponent(); 