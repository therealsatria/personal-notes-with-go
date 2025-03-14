/**
 * Encryption Status Service
 * Handles checking and managing encryption status
 */
class EncryptionStatusService {
    constructor() {
        this.isValid = false; // Default to false until checked
        this.statusMessage = "";
        this.listeners = [];
        console.log("EncryptionStatusService initialized");
    }

    /**
     * Check the encryption status from the server
     */
    async checkStatus() {
        try {
            console.log("Checking encryption status from server...");
            const response = await apiService.getEncryptionStatus();
            console.log("Encryption status response:", response);
            
            this.isValid = response.encryption_valid;
            this.statusMessage = response.message;
            
            // Notify all listeners of the status change
            this.notifyListeners();
            
            return this.isValid;
        } catch (error) {
            console.error('Failed to check encryption status:', error);
            this.isValid = false;
            this.statusMessage = "Unable to verify encryption status. Data modification is disabled for security reasons.";
            
            // Notify all listeners of the status change
            this.notifyListeners();
            
            return false;
        }
    }

    /**
     * Get the current encryption status
     */
    isEncryptionValid() {
        return this.isValid;
    }

    /**
     * Get the current status message
     */
    getStatusMessage() {
        return this.statusMessage;
    }

    /**
     * Add a listener to be notified when encryption status changes
     * @param {Function} listener - Function to call when status changes
     */
    addListener(listener) {
        if (typeof listener === 'function' && !this.listeners.includes(listener)) {
            this.listeners.push(listener);
            console.log("Listener added, total listeners:", this.listeners.length);
        }
    }

    /**
     * Remove a listener
     * @param {Function} listener - Function to remove
     */
    removeListener(listener) {
        const index = this.listeners.indexOf(listener);
        if (index !== -1) {
            this.listeners.splice(index, 1);
            console.log("Listener removed, total listeners:", this.listeners.length);
        }
    }

    /**
     * Notify all listeners of a status change
     */
    notifyListeners() {
        console.log("Notifying listeners of encryption status change:", 
                   this.isValid, this.statusMessage, "Listeners:", this.listeners.length);
        
        this.listeners.forEach(listener => {
            try {
                listener(this.isValid, this.statusMessage);
            } catch (error) {
                console.error('Error in encryption status listener:', error);
            }
        });
    }
}

// Create and export a singleton instance
const encryptionStatusService = new EncryptionStatusService();

// Check status immediately after page load
document.addEventListener('DOMContentLoaded', () => {
    setTimeout(() => {
        encryptionStatusService.checkStatus();
    }, 500); // Small delay to ensure DOM is fully loaded
}); 