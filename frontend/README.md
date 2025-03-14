# Personal Notes Frontend

This is the frontend for the Personal Notes application. It's a Single Page Application (SPA) built with vanilla JavaScript, HTML, and CSS.

## Features

- Clean, modern UI with responsive design
- CRUD operations for notes and categories
- SPA navigation without page reloads
- Toast notifications for user feedback
- Modal dialogs for forms and confirmations
- Proper error handling

## Structure

The frontend follows a component-based architecture for clean code organization:

```
frontend/
├── css/
│   └── styles.css          # All styles for the application
├── js/
│   ├── components/         # UI components
│   │   ├── notes.js        # Notes component
│   │   └── categories.js   # Categories component
│   ├── services/           # Shared services
│   │   ├── api.js          # API communication service
│   │   └── toast.js        # Toast notification service
│   └── app.js              # Main application logic
└── index.html              # Main HTML file
```

## How to Use

1. Start the backend server:
   ```
   go run main.go
   ```

2. Access the application:
   - Open your browser and navigate to `http://localhost:8080`
   - The server will automatically redirect you to the frontend

3. Using the application:
   - Switch between Notes and Categories using the navigation menu
   - Add new notes/categories using the "Add" buttons
   - Edit or delete existing items using the buttons on each card
   - All changes are saved to the backend automatically

## Development

If you want to modify the frontend:

1. Edit the HTML, CSS, or JavaScript files as needed
2. Refresh your browser to see the changes
3. No build step is required as this is a vanilla JavaScript application

## Clean Code Principles

This frontend implementation follows these clean code principles:

- **Single Responsibility Principle**: Each component and service has a single responsibility
- **DRY (Don't Repeat Yourself)**: Common functionality is extracted into reusable services
- **Separation of Concerns**: UI, business logic, and data access are separated
- **Consistent Error Handling**: All errors are caught and displayed to the user
- **Defensive Programming**: Input validation and error checking throughout
- **Meaningful Names**: Variables, functions, and classes have clear, descriptive names
- **Comments and Documentation**: Code is well-commented with JSDoc style comments 