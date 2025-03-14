/**
 * Notes Component
 * Handles all note-related functionality
 */
class NotesComponent {
    constructor() {
        // DOM Elements
        this.notesContainer = document.getElementById('notes-list');
        this.addNoteBtn = document.getElementById('add-note-btn');
        this.noteModal = document.getElementById('note-modal');
        this.noteForm = document.getElementById('note-form');
        this.noteModalTitle = document.getElementById('note-modal-title');
        this.noteIdInput = document.getElementById('note-id');
        this.noteSubjectInput = document.getElementById('note-subject');
        this.noteContentInput = document.getElementById('note-content');
        this.notePriorityInput = document.getElementById('note-priority');
        this.noteTagsInput = document.getElementById('note-tags');
        this.noteCategoryInput = document.getElementById('note-category');
        this.confirmModal = document.getElementById('confirm-modal');
        this.confirmMessage = document.getElementById('confirm-message');
        this.confirmYesBtn = document.getElementById('confirm-yes');
        
        // Search elements
        this.searchInput = document.getElementById('search-notes');
        this.searchClearBtn = document.getElementById('search-clear-btn');
        this.showAllBtn = document.getElementById('show-all-notes-btn');
        
        // State
        this.notes = [];
        this.categories = [];
        this.currentNoteId = null;
        this.searchQuery = '';
        this.searchTimeout = null;
        this.showAllNotes = false;
        
        // Initialize
        this.init();
    }
    
    /**
     * Initialize the component
     */
    async init() {
        // Add event listeners
        this.addNoteBtn.addEventListener('click', () => this.openAddNoteModal());
        this.noteForm.addEventListener('submit', (e) => this.handleNoteSubmit(e));
        
        // Search event listeners
        this.searchInput.addEventListener('input', () => this.handleSearchInput());
        this.searchClearBtn.addEventListener('click', () => this.clearSearch());
        this.showAllBtn.addEventListener('click', () => this.toggleShowAll());
        
        // Close modal buttons
        const closeButtons = this.noteModal.querySelectorAll('.close-modal');
        closeButtons.forEach(button => {
            button.addEventListener('click', () => this.closeNoteModal());
        });
        
        // Load data
        await this.loadCategories();
        await this.loadNotes();
    }
    
    /**
     * Toggle between showing all notes and limited notes
     */
    toggleShowAll() {
        this.showAllNotes = !this.showAllNotes;
        
        // Update button appearance
        if (this.showAllNotes) {
            this.showAllBtn.classList.add('active');
            this.showAllBtn.innerHTML = '<i class="fas fa-list-alt"></i> Showing All';
        } else {
            this.showAllBtn.classList.remove('active');
            this.showAllBtn.innerHTML = '<i class="fas fa-list"></i> Show All';
        }
        
        // Reload notes with new setting
        this.loadNotes();
    }
    
    /**
     * Handle search input with debounce
     */
    handleSearchInput() {
        const query = this.searchInput.value.trim();
        
        // Show/hide clear button
        if (query) {
            this.searchClearBtn.classList.add('visible');
        } else {
            this.searchClearBtn.classList.remove('visible');
        }
        
        // Debounce search
        clearTimeout(this.searchTimeout);
        this.searchTimeout = setTimeout(() => {
            this.searchQuery = query;
            this.loadNotes();
        }, 300);
    }
    
    /**
     * Clear search input and reset results
     */
    clearSearch() {
        this.searchInput.value = '';
        this.searchClearBtn.classList.remove('visible');
        this.searchQuery = '';
        this.loadNotes();
    }
    
    /**
     * Load all notes from the API
     */
    async loadNotes() {
        try {
            this.notesContainer.innerHTML = '<div class="loading">Loading notes...</div>';
            
            // Build URL with parameters
            let url = '/notes';
            const params = new URLSearchParams();
            
            if (this.searchQuery) {
                params.append('q', this.searchQuery);
            }
            
            if (this.showAllNotes) {
                params.append('all', 'true');
            }
            
            const queryString = params.toString();
            if (queryString) {
                url += `?${queryString}`;
            }
            
            const notes = await apiService.request(url);
            
            // Ensure notes is always an array, even if API returns null or undefined
            this.notes = Array.isArray(notes) ? notes : [];
            
            this.renderNotes();
        } catch (error) {
            this.notesContainer.innerHTML = `
                <div class="empty-state">
                    <i class="fas fa-exclamation-circle"></i>
                    <p>We were unable to retrieve your notes at this time. ${error.message}</p>
                    <button class="btn btn-secondary" onclick="notesComponent.loadNotes()">
                        Try Again
                    </button>
                </div>
            `;
            console.error('Error loading notes:', error);
        }
    }
    
    /**
     * Load all categories from the API
     */
    async loadCategories() {
        try {
            const categories = await apiService.getCategories();
            // Ensure categories is always an array, even if API returns null or undefined
            this.categories = Array.isArray(categories) ? categories : [];
            this.populateCategoryDropdown();
        } catch (error) {
            toastService.error('We were unable to load categories. Please try again later.');
            console.error('Error loading categories:', error);
        }
    }
    
    /**
     * Populate the category dropdown in the note form
     */
    populateCategoryDropdown() {
        this.noteCategoryInput.innerHTML = '<option value="">Select a category</option>';
        
        this.categories.forEach(category => {
            const option = document.createElement('option');
            option.value = category.id;
            option.textContent = category.name;
            this.noteCategoryInput.appendChild(option);
        });
    }
    
    /**
     * Render all notes in the container
     */
    renderNotes() {
        // Ensure notes is always an array
        if (!Array.isArray(this.notes)) {
            this.notes = [];
        }
        
        if (this.notes.length === 0) {
            if (this.searchQuery) {
                this.notesContainer.innerHTML = `
                    <div class="empty-state">
                        <i class="fas fa-search"></i>
                        <p>No notes were found matching "${this.escapeHtml(this.searchQuery)}"</p>
                        <button class="btn btn-secondary" onclick="notesComponent.clearSearch()">
                            Clear Search
                        </button>
                    </div>
                `;
            } else {
                this.notesContainer.innerHTML = `
                    <div class="empty-state">
                        <i class="fas fa-sticky-note"></i>
                        <p>You don't have any notes yet. Click "Add Note" to create your first note.</p>
                    </div>
                `;
            }
            return;
        }
        
        // Add note count info
        let noteCountInfo = '';
        if (!this.showAllNotes && this.notes.length >= 1) {
            // If we're showing limited notes (not showing all)
            noteCountInfo = `
                <div class="note-count-info">
                    <p>Showing ${this.notes.length} notes. <a href="#" onclick="event.preventDefault(); notesComponent.toggleShowAll();">Show all notes</a>.</p>
                </div>
            `;
        } else if (this.showAllNotes) {
            noteCountInfo = `
                <div class="note-count-info">
                    <p>Showing all ${this.notes.length} notes.</p>
                </div>
            `;
        }
        
        let html = noteCountInfo;
        
        this.notes.forEach(note => {
            // Find category name if available
            let categoryName = '';
            if (note.category_id) {
                const category = this.categories.find(c => c.id === note.category_id);
                if (category) {
                    categoryName = category.name;
                }
            }
            
            // Highlight search matches if search query exists
            let subject = note.subject;
            let content = note.content;
            
            if (this.searchQuery) {
                subject = this.highlightText(subject, this.searchQuery);
                content = this.highlightText(content, this.searchQuery);
            }
            
            html += `
                <div class="card note-card">
                    <div class="card-header">
                        <h3 class="card-title">${subject}</h3>
                        <div class="card-actions">
                            <button class="btn btn-secondary btn-sm" onclick="notesComponent.openEditNoteModal('${note.id}')">
                                <i class="fas fa-edit"></i>
                            </button>
                            <button class="btn btn-danger btn-sm" onclick="notesComponent.confirmDeleteNote('${note.id}')">
                                <i class="fas fa-trash"></i>
                            </button>
                        </div>
                    </div>
                    <div class="card-content">
                        ${content}
                    </div>
                    <div class="card-footer">
                        <div>
                            <span class="priority priority-${note.priority || 'medium'}">${note.priority || 'Medium'}</span>
                            ${categoryName ? `<span class="category-badge">${this.escapeHtml(categoryName)}</span>` : ''}
                        </div>
                        <div class="tags">
                            ${this.renderTags(note.tags)}
                        </div>
                    </div>
                </div>
            `;
        });
        
        this.notesContainer.innerHTML = html;
    }
    
    /**
     * Highlight search text in content
     * @param {string} text - The text to search in
     * @param {string} query - The search query
     * @returns {string} - Text with highlighted search matches
     */
    highlightText(text, query) {
        if (!query) return this.escapeHtml(text);
        
        const escapedText = this.escapeHtml(text);
        const escapedQuery = this.escapeHtml(query);
        
        // Case insensitive search
        const regex = new RegExp(escapedQuery, 'gi');
        return escapedText.replace(regex, match => `<span class="highlight">${match}</span>`);
    }
    
    /**
     * Render tags as badges
     * @param {string} tagsString - Comma separated tags
     * @returns {string} - HTML for tags
     */
    renderTags(tagsString) {
        if (!tagsString) return '';
        
        const tags = tagsString.split(',').map(tag => tag.trim()).filter(tag => tag);
        return tags.map(tag => `<span class="tag">${this.escapeHtml(tag)}</span>`).join('');
    }
    
    /**
     * Escape HTML to prevent XSS
     * @param {string} unsafe - Unsafe string
     * @returns {string} - Safe string
     */
    escapeHtml(unsafe) {
        if (!unsafe) return '';
        return unsafe
            .replace(/&/g, "&amp;")
            .replace(/</g, "&lt;")
            .replace(/>/g, "&gt;")
            .replace(/"/g, "&quot;")
            .replace(/'/g, "&#039;");
    }
    
    /**
     * Open the modal for adding a new note
     */
    openAddNoteModal() {
        this.noteModalTitle.textContent = 'Add New Note';
        this.noteForm.reset();
        this.noteIdInput.value = '';
        this.currentNoteId = null;
        this.noteModal.classList.add('active');
    }
    
    /**
     * Open the modal for editing an existing note
     * @param {string} noteId - The ID of the note to edit
     */
    async openEditNoteModal(noteId) {
        try {
            this.noteModalTitle.textContent = 'Edit Note';
            this.noteForm.reset();
            
            // Get the note data
            const note = this.notes.find(n => n.id === noteId);
            if (!note) {
                throw new Error('Note not found');
            }
            
            // Populate the form
            this.noteIdInput.value = note.id;
            this.noteSubjectInput.value = note.subject || '';
            this.noteContentInput.value = note.content || '';
            this.notePriorityInput.value = note.priority || 'medium';
            this.noteTagsInput.value = note.tags || '';
            this.noteCategoryInput.value = note.category_id || '';
            
            this.currentNoteId = noteId;
            this.noteModal.classList.add('active');
        } catch (error) {
            toastService.error('We were unable to load the note for editing. Please try again later.');
            console.error('Error loading note for edit:', error);
        }
    }
    
    /**
     * Close the note modal
     */
    closeNoteModal() {
        this.noteModal.classList.remove('active');
    }
    
    /**
     * Handle note form submission
     * @param {Event} event - The form submit event
     */
    async handleNoteSubmit(event) {
        event.preventDefault();
        
        try {
            const noteData = {
                subject: this.noteSubjectInput.value,
                content: this.noteContentInput.value,
                priority: this.notePriorityInput.value,
                tags: this.noteTagsInput.value,
                category_id: this.noteCategoryInput.value || null
            };
            
            let result;
            
            if (this.currentNoteId) {
                // Update existing note - include ID in the URL, not in the body
                result = await apiService.updateNote(this.currentNoteId, noteData);
                toastService.success('Your note has been updated successfully.');
            } else {
                // Create new note - ID will be generated by the server
                result = await apiService.createNote(noteData);
                toastService.success('Your note has been created successfully.');
            }
            
            // Reload notes and close modal
            await this.loadNotes();
            this.closeNoteModal();
        } catch (error) {
            toastService.error('We were unable to save your note. Please try again later.');
            console.error('Error saving note:', error);
        }
    }
    
    /**
     * Show confirmation dialog for deleting a note
     * @param {string} noteId - The ID of the note to delete
     */
    confirmDeleteNote(noteId) {
        const note = this.notes.find(n => n.id === noteId);
        if (!note) return;
        
        this.confirmMessage.textContent = `Are you sure you want to delete the note "${note.subject}"?`;
        
        // Remove previous event listeners
        const newConfirmYesBtn = this.confirmYesBtn.cloneNode(true);
        this.confirmYesBtn.parentNode.replaceChild(newConfirmYesBtn, this.confirmYesBtn);
        this.confirmYesBtn = newConfirmYesBtn;
        
        // Add new event listener
        this.confirmYesBtn.addEventListener('click', async () => {
            await this.deleteNote(noteId);
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
     * Delete a note
     * @param {string} noteId - The ID of the note to delete
     */
    async deleteNote(noteId) {
        try {
            await apiService.deleteNote(noteId);
            toastService.success('Your note has been deleted successfully.');
            await this.loadNotes();
        } catch (error) {
            toastService.error('We were unable to delete your note. Please try again later.');
            console.error('Error deleting note:', error);
        }
    }
}

// Create and export a singleton instance
const notesComponent = new NotesComponent(); 