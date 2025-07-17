# AGENTS.md - Development Guidelines for ansicht

## Build/Test Commands
- **Build**: `go build -o ansicht main.go`
- **Run**: `go run main.go`
- **Test**: `go test ./...` (run all tests)
- **Test single package**: `go test ./internal/service`
- **Format**: `go fmt ./...`
- **Vet**: `go vet ./...`
- **Mod tidy**: `go mod tidy`

## Code Style Guidelines
- Use tabs for indentation (Go standard)
- Package names: lowercase, single word when possible
- Type names: PascalCase (e.g., `MessageID`, `SearchResultMsg`)
- Function/method names: PascalCase for exported, camelCase for unexported
- Variable names: camelCase (e.g., `selectedIndex`, `markedMessages`)
- Constants: PascalCase or ALL_CAPS for package-level constants

## Import Organization
- Standard library imports first
- Third-party imports second (with blank line separator)
- Local imports last (with blank line separator)
- Use import aliases for long package names (e.g., `tea "github.com/charmbracelet/bubbletea"`)

## Error Handling
- Return errors as the last return value
- Use `fmt.Errorf()` for error wrapping with context
- Check bounds before array/slice access and return descriptive errors
- Use early returns for error conditions

## Project Structure
- `internal/`: Private application code
- `internal/model/`: Data structures and types
- `internal/service/`: Business logic and data management
- `internal/ui/`: User interface components (Bubble Tea)
- `internal/runtime/`: Configuration and Lua runtime
- `internal/db/`: Database/notmuch integration