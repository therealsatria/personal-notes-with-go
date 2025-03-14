/**
 * Categories Component
 * Handles all category-related functionality
 */
class CategoriesComponent {
    constructor() {
        // DOM Elements
        this.categoriesContainer = document.getElementById('categories-list');
        this.addCategoryBtn = document.getElementById('add-category-btn');
        this.categoryModal = document.getElementById('category-modal');
        this.categoryForm = document.getElementById('category-form');
        this.categoryModalTitle = document.getElementById('category-modal-title');
        this.categoryIdInput = document.getElementById('category-id');
        this.categoryNameInput = document.getElementById('category-name');
        this.confirmModal = document.getElementById('confirm-modal');
        this.confirmMessage = document.getElementById('confirm-message');
        this.confirmYesBtn = document.getElementById('confirm-yes');
        
        // State
        this.categories = [];
        this.currentCategoryId = null;
        
        // Initialize
        this.init();
    }
    
    /**
     * Initialize the component
     */
    async init() {
        // Add event listeners
        this.addCategoryBtn.addEventListener('click', () => this.openAddCategoryModal());
        this.categoryForm.addEventListener('submit', (e) => this.handleCategorySubmit(e));
        
        // Close modal buttons
        const closeButtons = this.categoryModal.querySelectorAll('.close-modal');
        closeButtons.forEach(button => {
            button.addEventListener('click', () => this.closeCategoryModal());
        });
        
        // Load data
        await this.loadCategories();
    }
    
    /**
     * Load all categories from the API
     */
    async loadCategories() {
        try {
            this.categoriesContainer.innerHTML = '<div class="loading">Loading categories...</div>';
            this.categories = await apiService.getCategories();
            this.renderCategories();
        } catch (error) {
            this.categoriesContainer.innerHTML = '<div class="empty-state"><i class="fas fa-exclamation-circle"></i><p>Failed to load categories</p></div>';
            toastService.error('Failed to load categories: ' + error.message);
        }
    }
    
    /**
     * Render all categories in the container
     */
    renderCategories() {
        if (this.categories.length === 0) {
            this.categoriesContainer.innerHTML = `
                <div class="empty-state">
                    <i class="fas fa-folder-open"></i>
                    <p>No categories found. Create your first category!</p>
                </div>
            `;
            return;
        }
        
        this.categoriesContainer.innerHTML = '';
        
        this.categories.forEach(category => {
            const categoryElement = document.createElement('div');
            categoryElement.className = 'category-card';
            categoryElement.innerHTML = `
                <h3 class="category-name">${this.escapeHtml(category.name)}</h3>
                <div class="card-actions">
                    <button class="btn btn-sm btn-secondary edit-category" data-id="${category.id}">
                        <i class="fas fa-edit"></i>
                    </button>
                    <button class="btn btn-sm btn-danger delete-category" data-id="${category.id}">
                        <i class="fas fa-trash"></i>
                    </button>
                </div>
            `;
            
            // Add event listeners
            const editBtn = categoryElement.querySelector('.edit-category');
            const deleteBtn = categoryElement.querySelector('.delete-category');
            
            editBtn.addEventListener('click', () => this.openEditCategoryModal(category.id));
            deleteBtn.addEventListener('click', () => this.confirmDeleteCategory(category.id));
            
            this.categoriesContainer.appendChild(categoryElement);
        });
    }
    
    /**
     * Escape HTML to prevent XSS
     * @param {string} unsafe - Unsafe string
     * @returns {string} - Safe string
     */
    escapeHtml(unsafe) {
        if (!unsafe) return '';
        return unsafe
            .replace(/&/g, '&amp;')
            .replace(/</g, '&lt;')
            .replace(/>/g, '&gt;')
            .replace(/"/g, '&quot;')
            .replace(/'/g, '&#039;');
    }
    
    /**
     * Open the modal for adding a new category
     */
    openAddCategoryModal() {
        this.categoryModalTitle.textContent = 'Add New Category';
        this.categoryForm.reset();
        this.categoryIdInput.value = '';
        this.currentCategoryId = null;
        this.categoryModal.classList.add('active');
    }
    
    /**
     * Open the modal for editing an existing category
     * @param {string} categoryId - The ID of the category to edit
     */
    async openEditCategoryModal(categoryId) {
        try {
            this.categoryModalTitle.textContent = 'Edit Category';
            this.categoryForm.reset();
            
            // Get the category data
            const category = this.categories.find(c => c.id === categoryId);
            if (!category) {
                throw new Error('Category not found');
            }
            
            // Populate the form
            this.categoryIdInput.value = category.id;
            this.categoryNameInput.value = category.name || '';
            
            this.currentCategoryId = categoryId;
            this.categoryModal.classList.add('active');
        } catch (error) {
            toastService.error('Failed to load category: ' + error.message);
        }
    }
    
    /**
     * Close the category modal
     */
    closeCategoryModal() {
        this.categoryModal.classList.remove('active');
    }
    
    /**
     * Handle category form submission
     * @param {Event} event - The form submit event
     */
    async handleCategorySubmit(event) {
        event.preventDefault();
        
        try {
            const categoryData = {
                name: this.categoryNameInput.value
            };
            
            let result;
            
            if (this.currentCategoryId) {
                // Update existing category
                result = await apiService.updateCategory(this.currentCategoryId, categoryData);
                toastService.success('Category updated successfully');
            } else {
                // Create new category
                result = await apiService.createCategory(categoryData);
                toastService.success('Category created successfully');
            }
            
            // Reload categories and close modal
            await this.loadCategories();
            this.closeCategoryModal();
            
            // If notes component exists, reload categories there too
            if (typeof notesComponent !== 'undefined') {
                await notesComponent.loadCategories();
            }
        } catch (error) {
            toastService.error('Failed to save category: ' + error.message);
        }
    }
    
    /**
     * Show confirmation dialog for deleting a category
     * @param {string} categoryId - The ID of the category to delete
     */
    confirmDeleteCategory(categoryId) {
        const category = this.categories.find(c => c.id === categoryId);
        if (!category) return;
        
        this.confirmMessage.textContent = `Are you sure you want to delete the category "${category.name}"?`;
        
        // Remove previous event listeners
        const newConfirmYesBtn = this.confirmYesBtn.cloneNode(true);
        this.confirmYesBtn.parentNode.replaceChild(newConfirmYesBtn, this.confirmYesBtn);
        this.confirmYesBtn = newConfirmYesBtn;
        
        // Add new event listener
        this.confirmYesBtn.addEventListener('click', async () => {
            await this.deleteCategory(categoryId);
            this.closeConfirmModal();
        });
        
        // Close buttons
        const closeButtons = this.confirmModal.querySelectorAll('.close-modal');
        closeButtons.forEach(button => {
            button.addEventListener('click', () => this.closeConfirmModal());
        });
        
        this.confirmModal.classList.add('active');
    }
    
    /**
     * Close the confirmation modal
     */
    closeConfirmModal() {
        this.confirmModal.classList.remove('active');
    }
    
    /**
     * Delete a category
     * @param {string} categoryId - The ID of the category to delete
     */
    async deleteCategory(categoryId) {
        try {
            await apiService.deleteCategory(categoryId);
            toastService.success('Category deleted successfully');
            await this.loadCategories();
            
            // If notes component exists, reload categories there too
            if (typeof notesComponent !== 'undefined') {
                await notesComponent.loadCategories();
                // Also reload notes as they might reference this category
                await notesComponent.loadNotes();
            }
        } catch (error) {
            toastService.error('Failed to delete category: ' + error.message);
        }
    }
}

// Create and export a singleton instance
const categoriesComponent = new CategoriesComponent(); 