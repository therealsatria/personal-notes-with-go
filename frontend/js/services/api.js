/**
 * API Service
 * Handles all communication with the backend API
 */
class ApiService {
    constructor() {
        this.baseUrl = 'http://localhost:8080';
    }

    /**
     * Generic method to make API requests
     * @param {string} endpoint - API endpoint
     * @param {string} method - HTTP method (GET, POST, PUT, DELETE)
     * @param {object} data - Request body data (optional)
     * @returns {Promise} - Promise with response data
     */
    async request(endpoint, method = 'GET', data = null) {
        const url = `${this.baseUrl}${endpoint}`;
        
        const options = {
            method,
            headers: {
                'Content-Type': 'application/json'
            }
        };

        if (data && (method === 'POST' || method === 'PUT')) {
            options.body = JSON.stringify(data);
        }

        try {
            const response = await fetch(url, options);
            
            // For DELETE requests with 204 No Content
            if (response.status === 204) {
                return { success: true };
            }
            
            // For other responses, parse JSON
            if (!response.ok) {
                const errorData = await response.json();
                throw new Error(errorData.error || 'API request failed');
            }
            
            return await response.json();
        } catch (error) {
            console.error('API request error:', error);
            throw error;
        }
    }

    /**
     * Shorthand method for GET requests
     * @param {string} endpoint - API endpoint
     * @returns {Promise} - Promise with response data
     */
    async get(endpoint) {
        return this.request(endpoint, 'GET');
    }

    /**
     * Shorthand method for POST requests
     * @param {string} endpoint - API endpoint
     * @param {object} data - Request body data
     * @returns {Promise} - Promise with response data
     */
    async post(endpoint, data) {
        return this.request(endpoint, 'POST', data);
    }

    /**
     * Shorthand method for PUT requests
     * @param {string} endpoint - API endpoint
     * @param {object} data - Request body data
     * @returns {Promise} - Promise with response data
     */
    async put(endpoint, data) {
        return this.request(endpoint, 'PUT', data);
    }

    /**
     * Shorthand method for DELETE requests
     * @param {string} endpoint - API endpoint
     * @returns {Promise} - Promise with response data
     */
    async delete(endpoint) {
        return this.request(endpoint, 'DELETE');
    }

    // Notes API methods
    async getNotes(searchQuery = '') {
        let endpoint = '/notes';
        if (searchQuery) {
            endpoint += `?q=${encodeURIComponent(searchQuery)}`;
        }
        return this.request(endpoint);
    }

    async getNoteById(id) {
        return this.request(`/notes/${id}`);
    }

    async createNote(noteData) {
        return this.request('/notes', 'POST', noteData);
    }

    async updateNote(id, noteData) {
        return this.request(`/notes/${id}`, 'PUT', noteData);
    }

    async deleteNote(id) {
        return this.request(`/notes/${id}`, 'DELETE');
    }

    // Categories API methods
    async getCategories() {
        return this.request('/categories');
    }

    async getCategoryById(id) {
        return this.request(`/categories/${id}`);
    }

    async createCategory(categoryData) {
        return this.request('/categories', 'POST', categoryData);
    }

    async updateCategory(id, categoryData) {
        return this.request(`/categories/${id}`, 'PUT', categoryData);
    }

    async deleteCategory(id) {
        return this.request(`/categories/${id}`, 'DELETE');
    }

    // Key Generator
    async generateKey(inputText) {
        return this.request('/generate-key', 'POST', { text: inputText });
    }

    // Encryption status
    async getEncryptionStatus() {
        return this.request('/encryption/status');
    }
}

// Create and export a singleton instance
const apiService = new ApiService(); 