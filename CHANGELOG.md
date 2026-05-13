# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.2.0] - 2026-05-13

### Added

#### Identity-Based Access Control
- **Enrollment ID Checks**: New `RequireID()` and `RequireAnyID()` rules to check the caller's enrollment ID from the Fabric certificate. Useful for per-user allow/deny and admin identity checks.
- **Affiliation Checks**: New `RequireAffiliation()` rule to check the `hf.Affiliation` attribute. Enables org-subdivision based access control (e.g., `org1.department1`).
- **Multi-Value Attribute Checks**: New `RequireAnyAttribute()` rule to check an attribute against multiple accepted values, completing the `RequireAny*` API symmetry.

#### Declarative Policy Management
- **PolicySet**: New `PolicySet` type for function-name based access control. Maps chaincode function names to access rules with support for default policies, per-function overrides, chainable API (`AddPolicy`/`RemovePolicy` return `*PolicySet`), and `Middleware()` for framework integration.

#### Combinator Aliases
- **All() / Any()**: Added `All()` and `Any()` aliases for `And()` and `Or()` combinators for more natural readability in policy declarations.

#### Performance
- **Benchmark Suite**: Added benchmarks covering single rule, And/Or combinators, complex nested rules, controller, helper function, and PolicySet evaluation paths.

### Fixed
- **AsAccessError clarity**: Refactored `AsAccessError()` return to use explicit two-step form (`ok := errors.As(...); return e, ok`) for unambiguous semantics.
- **IsHLFAdmintRule typo**: Renamed to `IsHLFAdminRule`. The old name `IsHLFAdmintRule` is preserved as a deprecated type alias for backward compatibility.

## [0.1.0] - 2026-02-21

### Added

#### Core Features
- **Access Control Framework**: Complete access control solution for Hyperledger Fabric Go chaincodes
- **Role-Based Access Control**: Support for checking specific roles via `RequireRole()` and multiple roles via `RequireAnyRole()`
- **Organization-Based Controls**: MSPID checks with `RequireMSPID()` and `RequireAnyMSPID()` functions
- **Organizational Unit Checks**: OU validation using `RequireOU()` and `RequireAnyOU()` functions
- **User Type Validation**: Support for checking Hyperledger Fabric administrators via `RequireHLFAdmin()`

#### Rule Combinators
- **AND Logic**: Combine multiple rules requiring all to pass using `And()` function
- **OR Logic**: Combine multiple rules requiring at least one to pass using `Or()` function
- **NOT Logic**: Invert rule results using `Not()` function
- **Complex Combinations**: Support for nested rule combinations for advanced access policies

#### Custom Rules
- **Custom Rule Creation**: Create custom access rules using `Custom()` function with callback functions
- **Always Deny Rule**: Force access denial using `AlwaysDeny()` function for special cases

#### API & Architecture
- **Controller Pattern**: Main `hyperacc` controller with `New()` and `Check()` methods
- **Helper Functions**: Convenient `CheckAccess()` function for quick access checks
- **Rule Interface**: Extensible `Rule` interface for creating custom rule implementations
- **Comprehensive Error Handling**: Detailed access error reporting with `AccessError` type

#### Documentation & Examples
- **Complete API Documentation**: Full documentation of all exported functions and types
- **Usage Examples**: Comprehensive examples for basic and advanced use cases
- **Real-world Scenarios**: Practical examples for supply chain and financial contracts
- **Integration Guide**: Clear instructions for integrating with existing chaincodes

#### Testing & Quality
- **Unit Tests**: Complete test coverage for all core functionality
- **Mock Generation**: Mock interfaces for testing with `mockery`
- **Quality Assurance**: Linting configuration with `.golangci.yml`

#### Dependencies & Requirements
- **Fabric Compatibility**: Built specifically for Hyperledger Fabric v2.x
- **Modern Go Support**: Requires Go 1.21+ with support for latest language features
- **Dependency Management**: Proper module management with go.mod

### Fixed
- N/A (Initial release)

### Security
- **Attribute Validation**: Secure validation of client identity attributes
- **Certificate Parsing**: Safe parsing of X.509 certificates for OU validation
- **Access Error Isolation**: Proper error wrapping to prevent information disclosure

### Performance
- **Optimized Rule Evaluation**: Efficient evaluation of rule combinations
- **Minimal Overhead**: Lightweight implementation with minimal impact on transaction performance

### Known Issues
- None reported for initial release

[0.2.0]: https://github.com/IceTweak/hyperacc/releases/tag/v0.2.0
[0.1.0]: https://github.com/IceTweak/hyperacc/releases/tag/v0.1.0