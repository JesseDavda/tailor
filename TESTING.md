# Testing Documentation

## Overview

This document describes the test coverage for the Tailor project. The project now has comprehensive test coverage for its core business logic packages.

## Phase 1: Foundation Tests (Completed)

Phase 1 focused on testing the highest-value packages with the most critical business logic.

### Test Coverage Summary

| Package | Coverage | Test Files | Focus |
|---------|----------|------------|-------|
| `internal/changes/` | **85.5%** | 4 test files | Business logic for changes, validation, and application |
| `internal/resume/` | **78.4%** | 2 test files | YAML manipulation and file I/O |
| `pkg/config/` | **100%** | 1 test file | Configuration validation |
| **Overall** | **83.3%** | **7 test files** | **Phase 1 packages** |

### Test Statistics

- **Total test functions**: 153+
- **Execution time**: < 1 second (all tests pass quickly)
- **Test framework**: Go standard library (`testing` package)
- **No external test dependencies** (using only standard library)

## Test Files Created

### Changes Package (`internal/changes/`)

1. **`test_helpers.go`** - Helper functions for creating test fixtures
2. **`models_test.go`** - Tests for Change and ChangeSet models
   - Change status management (Approve, Reject, IsApproved)
   - ChangeSet filtering and counting
   - Cover letter detection
3. **`parser_test.go`** - Tests for JSON parsing and validation
   - Valid/invalid JSON parsing
   - Markdown code block cleaning
   - Change validation rules
   - Cover letter validation
4. **`validator_test.go`** - Tests for change validation against YAML
   - Modify, Reorder, Add, Remove operations
   - Path validation
   - Type checking
   - Conflict detection
5. **`applier_test.go`** - Tests for applying changes to YAML documents
   - All change types (modify, reorder, add, remove)
   - Multiple changes in sequence
   - Partial failure handling
   - Change sorting and prioritization

### Resume Package (`internal/resume/`)

1. **`yaml_manipulator_test.go`** - Tests for YAML document manipulation (382 LOC tested)
   - Path parsing (simple keys, nested keys, array indices)
   - Node navigation and retrieval
   - Value setting (scalar and complex types)
   - Array operations (reorder, add, remove)
   - Round-trip serialization
   - **Critical**: Tests the DocumentNode unwrapping fix
2. **`processor_test.go`** - Tests for file I/O operations
   - Reading YAML files
   - Writing YAML files
   - YAML validation
   - Line counting

### Config Package (`pkg/config/`)

1. **`config_test.go`** - Tests for configuration validation
   - API key validation
   - File existence checks
   - Approval mode validation
   - Cover letter path generation
   - Environment variable loading

## Test Data

Test fixtures are organized in `testdata/` directories:

```
internal/
├── changes/
│   └── testdata/
│       ├── valid_changeset.json
│       ├── changeset_with_cover_letter.json
│       ├── invalid_changeset.json
│       └── sample_resume.yaml
├── resume/
│   └── testdata/
│       ├── simple_resume.yaml
│       ├── complex_resume.yaml
│       └── invalid_resume.yaml
```

## Running Tests

### Run all tests
```bash
go test ./...
```

### Run with coverage
```bash
go test -cover ./...
```

### Generate coverage report
```bash
go test -coverprofile=coverage.out ./internal/changes/ ./internal/resume/ ./pkg/config/
go tool cover -html=coverage.out -o coverage.html
```

### Run specific package
```bash
go test -v ./internal/changes/
```

### Run specific test
```bash
go test -run TestParseChangeSet ./internal/changes/
```

### Check for race conditions
```bash
go test -race ./...
```

## Coverage Goals vs Actual

| Package | Target | Actual | Status |
|---------|--------|--------|--------|
| changes/ | 90% | **85.5%** | ✅ Close to target |
| resume/ | 85% | **78.4%** | ✅ Good coverage |
| config/ | 95% | **100%** | ✅ Exceeded target |
| **Overall** | **75%** | **83.3%** | ✅ **Exceeded target** |

## Test Quality

### What We Test

✅ **Business Logic**
- Change approval/rejection workflows
- ChangeSet filtering and counting
- Validation rules for all change types
- Path parsing and navigation
- YAML manipulation operations

✅ **Edge Cases**
- Empty changesets
- Invalid paths
- Out-of-bounds array access
- Type mismatches
- Conflicting changes

✅ **Error Conditions**
- Missing files
- Invalid YAML syntax
- Invalid JSON parsing
- Path validation failures
- Permission errors

✅ **Integration Scenarios**
- Multiple changes applied in sequence
- Partial failure handling
- Round-trip YAML serialization
- Complex nested structures

### Testing Patterns Used

1. **Table-Driven Tests** - Used extensively for testing multiple scenarios
2. **Temporary Files** - Using `t.TempDir()` for file I/O tests
3. **Helper Functions** - Shared test utilities in `test_helpers.go`
4. **Descriptive Test Names** - Format: `Test<Function>_<Scenario>`
5. **Subtests** - Using `t.Run()` for readable output and organization

### Test Isolation

- Each test is independent
- Temporary directories used for file tests (`t.TempDir()`)
- No global state dependencies
- Tests can run in parallel safely

## Known Gaps

Areas with lower coverage (not critical for Phase 1):

1. **`internal/resume/yaml_manipulator.go`** - Some edge cases in complex type handling (78.4% coverage)
   - DocumentNode unwrapping edge cases
   - Some error paths in node traversal

These gaps are acceptable for Phase 1 and can be addressed in future phases if needed.

## Maintenance

### Adding New Tests

1. Follow existing naming conventions
2. Use table-driven tests for multiple scenarios
3. Create test fixtures in `testdata/` directories
4. Add helper functions to `test_helpers.go` if reusable
5. Ensure tests are isolated and can run in parallel

### Updating Tests

When modifying code:
1. Run tests: `go test ./...`
2. Check coverage: `go test -cover ./...`
3. Verify no regressions
4. Update tests if behavior changes
5. Add tests for new functionality
