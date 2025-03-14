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
                throw new Error(errorData.message || 'API request failed');
            }
            
            return await response.json();
        } catch (error) {
            console.error('API request error:', error);
            throw error;
        }
    }

    // Notes API methods
    async getNotes() {
        return this.request('/notes');
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
}

// Create and export a singleton instance
const apiService = new ApiService(); 