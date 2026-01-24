# Changelog

All notable changes to the IPTU API Go SDK will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [2.1.2] - 2026-01-24

### Fixed
- Updated all examples to use current API signatures (context.Context + typed params)
- Fixed `examples/basic/main.go` - was using deprecated string parameters
- Fixed `examples/valuation/main.go` - was using non-existent structs
- Fixed `examples/advanced/main.go` - missing context parameter
- Fixed `examples/multi-city/main.go` - was using non-existent functions

### Changed
- Version constant updated to 2.1.2

## [2.1.0] - 2025-12-15

### Added
- IPTU Tools endpoints for 2026 calendar data
  - `IPTUToolsCidades()` - List cities with IPTU calendar
  - `IPTUToolsCalendario()` - Get IPTU calendar for a city
  - `IPTUToolsSimulador()` - Simulate payment options
  - `IPTUToolsIsencao()` - Check exemption eligibility
  - `IPTUToolsProximoVencimento()` - Get next due date info
- Typed structs for all IPTU Tools responses
- Brasilia city support (`CidadeBrasilia`)

### Changed
- All methods use `context.Context` for cancellation support (already implemented)

## [2.0.0] - 2025-11-01

### Added
- Complete rewrite with idiomatic Go patterns
- Context support for all API methods
- Configurable retry with exponential backoff
- Rate limit tracking via `RateLimit` and `LastRequestID` fields
- Typed error types: `AuthenticationError`, `ForbiddenError`, `NotFoundError`, `RateLimitError`, `ValidationError`, `ServerError`
- Helper functions: `IsNotFound()`, `IsRateLimit()`, `IsAuthError()`, `IsForbidden()`, `IsServerError()`
- Functional options pattern: `WithBaseURL()`, `WithTimeout()`, `WithRetry()`, `WithLogger()`, `WithHTTPClient()`
- Valuation endpoints (Pro+): `ValuationEstimate()`, `ValuationBatch()`, `ValuationComparables()`
- Data endpoints: `DadosIPTUHistorico()`, `DadosCNPJ()`, `IPCACorrecao()`

### Changed
- Client initialization now uses `NewClient(apiKey, ...options)` pattern
- All methods return typed structs instead of `map[string]interface{}`

## [1.0.0] - 2025-09-01

### Added
- Initial release
- Basic consultation endpoints: `ConsultaEndereco()`, `ConsultaSQL()`, `ConsultaCEP()`
- Zoning query: `ConsultaZoneamento()`
- Support for multiple cities (SP, BH, Recife, POA, Fortaleza, Curitiba, RJ)
