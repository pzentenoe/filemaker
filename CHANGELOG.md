# Changelog

All notable changes to this project will be documented in this file.

## [2.0.0] - 2026-01-30

### Added
- **CI**: Add CI workflow for Go with testing, coverage, and linting steps (0718b3e)
- **Examples**: Add executable examples for FileMaker Go library with detailed README instructions (952ae53)
- **License**: Add MIT License (71fad04)
- **Auth**: Implement authentication strategies and enhance record service with pagination support (f04bacd)

### Changed
- **Refactor**: Refactor defer statements for improved error handling and readability in main and test files (063b772)
- **Docs**: Update FOSSA action version and enhance README with detailed example descriptions and setup instructions (f1cf976)
- **Docs**: Enhance README with detailed usage and features (71fad04)
- **Refactor**: Update defer statements for response body closure and error messages for consistency (9e89c9f)

### Tests
- Update authentication methods in tests to use SetBasicAuth and enhance test coverage for session creation (c06d882)
- Enhance tests for search service and record service, adding new functionalities and improving coverage (7ad2f40)

### Merged
- Merge pull request #3 from pzentenoe/feature/golangci (9abb74b)

## [1.0.0] - 2026-01-29

### Added
- **Builders**: Added `RecordBuilder`, `FindBuilder`, and `SessionBuilder` for fluent API design (677a09b)
- **Portals**: Added support for portal data in records (29fe6fa)
- **Groups**: Added GroupQuery functionalities (9ef4c8f)

### Fixed
- **Portals**: Remove omitempty from portal definitions to fix marshaling issues (a2628ef)
- **Typo**: Fix typo in searchData (9ef4c8f)