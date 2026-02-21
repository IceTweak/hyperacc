# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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

[0.1.0]: https://github.com/IceTweak/hyperacc/releases/tag/v0.1.0