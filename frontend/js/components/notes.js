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
        
        // State
        this.notes = [];
        this.categories = [];
        this.currentNoteId = null;
        
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
     * Load all notes from the API
     */
    async loadNotes() {
        try {
            this.notesContainer.innerHTML = '<div class="loading">Loading notes...</div>';
            this.notes = await apiService.getNotes();
            this.renderNotes();
        } catch (error) {
            this.notesContainer.innerHTML = '<div class="empty-state"><i class="fas fa-exclamation-circle"></i><p>Failed to load notes</p></div>';
            toastService.error('Failed to load notes: ' + error.message);
        }
    }
    
    /**
     * Load all categories from the API
     */
    async loadCategories() {
        try {
            this.categories = await apiService.getCategories();
            this.populateCategoryDropdown();
        } catch (error) {
            toastService.error('Failed to load categories: ' + error.message);
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
        if (this.notes.length === 0) {
            this.notesContainer.innerHTML = `
                <div class="empty-state">
                    <i class="fas fa-sticky-note"></i>
                    <p>No notes found. Create your first note!</p>
                </div>
            `;
            return;
        }
        
        this.notesContainer.innerHTML = '';
        
        this.notes.forEach(note => {
            const category = this.categories.find(c => c.id === note.category_id);
            const categoryName = category ? category.name : 'Uncategorized';
            
            // Parse tags
            const tags = note.tags ? note.tags.split(',').map(tag => tag.trim()) : [];
            
            const noteElement = document.createElement('div');
            noteElement.className = 'card';
            noteElement.innerHTML = `
                <div class="card-header">
                    <h3 class="card-title">${this.escapeHtml(note.subject)}</h3>
                    <div class="card-actions">
                        <button class="btn btn-sm btn-secondary edit-note" data-id="${note.id}">
                            <i class="fas fa-edit"></i>
                        </button>
                        <button class="btn btn-sm btn-danger delete-note" data-id="${note.id}">
                            <i class="fas fa-trash"></i>
                        </button>
                    </div>
                </div>
                <div class="card-content">
                    ${this.formatContent(note.content)}
                </div>
                <div class="card-footer">
                    <div>
                        <span class="priority priority-${note.priority}">${note.priority}</span>
                        <span class="category">${this.escapeHtml(categoryName)}</span>
                    </div>
                    <div class="tags">
                        ${tags.map(tag => `<span class="tag">${this.escapeHtml(tag)}</span>`).join('')}
                    </div>
                </div>
            `;
            
            // Add event listeners
            const editBtn = noteElement.querySelector('.edit-note');
            const deleteBtn = noteElement.querySelector('.delete-note');
            
            editBtn.addEventListener('click', () => this.openEditNoteModal(note.id));
            deleteBtn.addEventListener('click', () => this.confirmDeleteNote(note.id));
            
            this.notesContainer.appendChild(noteElement);
        });
    }
    
    /**
     * Format note content for display
     * @param {string} content - The note content
     * @returns {string} - Formatted HTML
     */
    formatContent(content) {
        if (!content) return '';
        
        // Convert line breaks to <br>
        return this.escapeHtml(content).replace(/\n/g, '<br>');
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
            toastService.error('Failed to load note: ' + error.message);
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
                category_id: this.noteCategoryInput.value
            };
            
            let result;
            
            if (this.currentNoteId) {
                // Update existing note
                result = await apiService.updateNote(this.currentNoteId, noteData);
                toastService.success('Note updated successfully');
            } else {
                // Create new note
                result = await apiService.createNote(noteData);
                toastService.success('Note created successfully');
            }
            
            // Reload notes and close modal
            await this.loadNotes();
            this.closeNoteModal();
        } catch (error) {
            toastService.error('Failed to save note: ' + error.message);
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
            toastService.success('Note deleted successfully');
            await this.loadNotes();
        } catch (error) {
            toastService.error('Failed to delete note: ' + error.message);
        }
    }
}

// Create and export a singleton instance
const notesComponent = new NotesComponent(); 