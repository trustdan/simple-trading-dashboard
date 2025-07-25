### 16. Windows Development Best Practices

- **Node Modules Issues Prevention**:
  - Use pnpm instead of npm: `pnpm install`
  - Keep dependencies minimal
  - Delete node_modules and use `npm ci` for clean installs
  - Consider using Wails embedded assets without separate build
  - Alternative: Build frontend on WSL2 if issues persist
- Set max path length in Windows if needed
- Use `.npmrc` with `package-lock=false` if lock files cause issues
- Regular cleanup: `npx rimraf node_modules` before reinstalls# Options Trading Dashboard - Development Rules & Guidelines

## Core Development Principles

### 1. Progress Tracking
- **MANDATORY**: Maintain a `PROGRESS_LOG.md` file in the project root
- Update after every significant milestone or work session
- Use clear, timestamped entries for each feature
- Example format:
  ```markdown
  ## Progress Log
  
  ### Market Rating Tab
  - [2025-07-23 16:32] Completed slider integration with backend API
  - [2025-07-23 18:45] Added sector persistence to SQLite
  
  ### Options Trade Management
  - [2025-07-24 09:15] Implemented trade entry form validation
  - [2025-07-24 11:30] Color coding system for all 24 strategy types complete
  ```

### 2. Code Organization

#### Backend (Go)
```
backend/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── api/
│   │   ├── handlers/
│   │   ├── middleware/
│   │   └── routes/
│   ├── models/
│   ├── services/
│   └── database/
├── pkg/
│   └── utils/
└── go.mod
```

#### Frontend (Svelte)
```
frontend/
├── src/
│   ├── components/
│   │   ├── market/
│   │   ├── trades/
│   │   └── shared/
│   ├── stores/
│   ├── services/
│   ├── types/
│   └── utils/
├── public/
└── package.json
```

### 3. Coding Standards

#### Go Backend
- Use standard Go formatting (`gofmt`)
- Prefix interfaces with 'I' (e.g., `ITradeService`)
- Keep functions under 50 lines
- Return errors as the last value
- Use context for cancellation and timeouts
- Comment all exported functions

#### Svelte/JavaScript Frontend
- Use plain JavaScript (no TypeScript)
- Keep the build system simple
- Define data structures with JSDoc comments
- Keep components under 200 lines
- Use composition over inheritance
- Implement proper error boundaries
- Document component props with JSDoc

### 4. API Design Rules

- RESTful endpoints with clear naming
- Version the API from the start (`/api/v1/`)
- Use proper HTTP status codes
- Implement request/response logging
- Standard response format:
  ```json
  {
    "success": boolean,
    "data": any,
    "error": string | null,
    "timestamp": ISO8601
  }
  ```

### 5. State Management

#### Frontend
- Use Svelte stores for global state
- Keep component state local when possible
- Implement optimistic updates for better UX
- Cache API responses appropriately

#### Backend
- Stateless API design
- Use database for all persistent state
- Implement proper transaction management
- No in-memory caching initially

### 6. Error Handling

- **Never ignore errors**
- Log all errors with context
- User-friendly error messages in UI
- Implement global error handler in frontend
- Use error boundaries for component crashes

### 7. Testing Requirements

- Write tests for critical business logic
- API endpoint tests are mandatory
- Component tests for complex UI logic
- Maintain >70% code coverage for services
- Use table-driven tests in Go

### 8. Database Rules

- Use migrations for all schema changes
- Never delete columns, deprecate instead
- Index foreign keys and commonly queried fields
- Backup strategy from day one
- Use transactions for multi-table operations

### 9. Security Guidelines

- Validate all user inputs
- Sanitize data before storage
- Use prepared statements for SQL
- Implement CORS properly
- No sensitive data in logs

### 10. Documentation

- README.md must include setup instructions
- Document all API endpoints
- Component props must have JSDoc comments
- Keep an architecture decision record (ADR)
- Maintain a CHANGELOG.md

### 11. Git Workflow

- Commit message format: `type(scope): description`
  - Types: feat, fix, docs, style, refactor, test, chore
- Branch naming: `feature/description` or `fix/description`
- Commit frequently with meaningful messages
- Never commit sensitive data

### 12. Performance Guidelines

- Profile before optimizing
- Lazy load components where appropriate
- Implement pagination for large datasets
- Use indexes for frequent queries
- Monitor memory usage

### 13. UI/UX Consistency

- **Dial Design Standards**:
  - 200x200px for sector dials, 300x300px for overall market dial
  - Smooth SVG paths with gradient fills
  - Animated needle transitions (300ms ease-in-out)
  - Clear numerical display at dial center
  - Color gradients: Red (-3) → Yellow (0) → Green (+3)
  - Hover states with subtle glow effects
  - Click-and-drag or click-to-set interaction
- Strategy colors organized by category
- Consistent spacing (use 8px grid for visual breathing room)
- Loading states with skeleton dials
- Empty states for no data scenarios
- Keyboard navigation support (arrow keys for dial adjustment)
- Professional dark theme with subtle shadows and highlights

### 14. Development Workflow

1. Plan the feature/fix
2. Update PROGRESS_LOG.md with start time
3. Write tests (when applicable)
4. Implement the feature
5. Test manually
6. Update documentation
7. Commit with proper message
8. Update PROGRESS_LOG.md with completion

### 15. Code Review Checklist

Before considering any feature complete:
- [ ] Code follows project structure
- [ ] Error handling is comprehensive
- [ ] Tests are written and passing
- [ ] Documentation is updated
- [ ] No console.log or fmt.Println left
- [ ] Performance impact considered
- [ ] Security implications reviewed
- [ ] PROGRESS_LOG.md updated with completion timestamp
