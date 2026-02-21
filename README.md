[![Go Reference](https://pkg.go.dev/badge/github.com/IceTweak/hyperacc.svg)](https://pkg.go.dev/github.com/IceTweak/hyperacc)
[![Go Report Card](https://goreportcard.com/badge/github.com/IceTweak/hyperacc)](https://goreportcard.com/report/github.com/IceTweak/hyperacc)
![Go Version](https://img.shields.io/badge/go%20version-%3E=1.21-61CFDD.svg?style=flat-square)
[![CI](https://github.com/IceTweak/hyperacc/workflows/CI/badge.svg)](https://github.com/IceTweak/hyperacc/actions?query=workflow%3ACI)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Release](https://img.shields.io/github/release/IceTweak/hyperacc.svg)](https://github.com/IceTweak/hyperacc/releases/latest)

# 🔐 hyperacc
Flexible and convenient package for access control to chaincode methods in HyperLedger Fabric.

## Table of Contents

- [Installation](#installation)
- [Quick Start](#quick-start)
- [Main Features](#main-features)
- [Documentation](#documentation)
- [Examples](#examples)
- [Requirements](#requirements)
- [Contributing](#contributing)
- [License](#license)

---

## Installation

```bash
go get github.com/IceTweak/hyperacc
```
**Important:** The package uses HyperLedger Fabric v2 API:

- `github.com/hyperledger/fabric-chaincode-go/v2`
- `github.com/hyperledger/fabric-contract-api-go/v2`
- `github.com/hyperledger/fabric-protos-go-apiv2`

## Main Features
- [x] Role-based access check (role attribute)
- [x] MSPID check
- [x] Organizational Unit (OU) check
- [x] Support for administrators and clients based on `hf.Type`
- [x] Rule combinations (AND/OR/NOT)
- [x] Custom rules
- [x] Simple and intuitive API


## Quick Start

### Simple role check

```go
package main

import (
    "github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
    "github.com/IceTweak/hyperacc"
)

type MyContract struct {
    contractapi.Contract
}

func (c *MyContract) CreateAsset(ctx contractapi.TransactionContextInterface, id string, value string) error {
    // Access check: only "manager" role
    if err := hyperacc.CheckAccess(ctx, hyperacc.RequireRole("manager")); err != nil {
        return err
    }
    // Asset creation logic

    return ctx.GetStub().PutState(id, []byte(value))
}
```

### Administrator check

```go
func (c *MyContract) DeleteAsset(ctx contractapi.TransactionContextInterface, id string) error {
    // Only administrators can delete assets
    if err := hyperacc.CheckAccess(ctx, hyperacc.RequireHLFAdmin()); err != nil {
        return err
    }
    
    return ctx.GetStub().DelState(id)
}
```

### Multiple roles check (OR)

```go
func (c *MyContract) ReadAsset(ctx contractapi.TransactionContextInterface, id string) (string, error) {
    // Access for manager, viewer or auditor
    if err := hyperacc.CheckAccess(ctx, 
        hyperacc.RequireAnyRole("manager", "viewer", "auditor"),
    ); err != nil {
        return "", err
    }
    
    data, err := ctx.GetStub().GetState(id)
    return string(data), err
}
```

### MSPID check

```go
func (c *MyContract) OrgSpecificMethod(ctx contractapi.TransactionContextInterface) error {
    // Access only for specific organization
    if err := hyperacc.CheckAccess(ctx, 
        hyperacc.RequireMSPID("Org1MSP"),
    ); err != nil {
        return err
    }
    // Method logic

    return nil
}
```

### Organizational Unit check

```go
func (c *MyContract) DepartmentMethod(ctx contractapi.TransactionContextInterface) error {
    // Access only for users from OU "department1"
    if err := hyperacc.CheckAccess(ctx, 
        hyperacc.RequireOU("department1"),
    ); err != nil {
        return err
    }
    // Method logic

    return nil
}
```

## Rule combinations

### AND - all rules must pass

```go
func (c *MyContract) SensitiveOperation(ctx contractapi.TransactionContextInterface) error {
    // Must be manager AND from Org1MSP
    if err := hyperacc.CheckAccess(ctx, 
        hyperacc.And(
            hyperacc.RequireRole("manager"),
            hyperacc.RequireMSPID("Org1MSP"),
        ),
    ); err != nil {
        return err
    }
    // Operation logic

    return nil
}
```

### OR - at least one rule must pass

```go
func (c *MyContract) MultiOrgMethod(ctx contractapi.TransactionContextInterface) error {
    // Access for multiple organizations
    if err := hyperacc.CheckAccess(ctx, 
        hyperacc.Or(
            hyperacc.RequireMSPID("Org1MSP"),
            hyperacc.RequireMSPID("Org2MSP"),
            hyperacc.RequireMSPID("Org3MSP"),
        ),
    ); err != nil {
        return err
    }
    // Method logic

    return nil
}
```

### Complex combinations

```go
func (c *MyContract) ComplexAccess(ctx contractapi.TransactionContextInterface) error {
    // (admin OR (manager AND Org1MSP)) AND NOT department2
    if err := hyperacc.CheckAccess(ctx, 
        hyperacc.And(
            hyperacc.Or(
                hyperacc.RequireHLFAdmin(),
                hyperacc.And(
                    hyperacc.RequireRole("manager"),
                    hyperacc.RequireMSPID("Org1MSP"),
                ),
            ),
            hyperacc.Not(
                hyperacc.RequireOU("department2"),
            ),
        ),
    ); err != nil {
        return err
    }
    // Method logic

    return nil
}
```
## Custom rules


```go
func (c *MyContract) CustomAccessMethod(ctx contractapi.TransactionContextInterface) error {
    // Custom rule based on function
    customRule := hyperacc.Custom("business-hours", func(ctx contractapi.TransactionContextInterface) error {
        stub := ctx.GetStub()
        timestamp, err := stub.GetTxTimestamp()
        if err != nil {
            return err
        }
        
        hour := timestamp.AsTime().Hour()
        if hour < 9 || hour > 18 {
            return fmt.Errorf("operation only allowed during business hours (9-18)")
        }
        
        return nil
    })
    
    if err := hyperacc.CheckAccess(ctx, customRule); err != nil {
        return err
    }
    // Method logic

    return nil
}
```
## Usage with hyperacc object


```go
type MyContract struct {
    contractapi.Contract
}

func (c *MyContract) Init(ctx contractapi.TransactionContextInterface) error {
    return nil
}

func (c *MyContract) CreateAsset(ctx contractapi.TransactionContextInterface, id string, value string) error {
    // Create access control object
    ac := hyperacc.New(
        hyperacc.RequireRole("manager"),
        hyperacc.RequireMSPID("Org1MSP"),
    )
    // Access check

    if err := ac.Check(ctx); err != nil {
        return err
    }
    // Asset creation logic

    return ctx.GetStub().PutState(id, []byte(value))
}
```

## API Reference

### Role-based rules

- `RequireRole(role string)` - requires specific role
- `RequireAnyRole(roles ...string)` - requires one of the roles

### Organization-based rules

- `RequireMSPID(mspid string)` - requires specific MSPID
- `RequireAnyMSPID(mspids ...string)` - requires one of the MSPIDs
- `RequireOU(ou string)` - requires specific Organizational Unit
- `RequireAnyOU(ous ...string)` - requires one of the OUs

### Combinators

- `And(rules ...Rule)` - all rules must pass
- `Or(rules ...Rule)` - at least one rule must pass
- `Not(rule Rule)` - inverts the rule result

### Custom rules

- `Custom(name string, checkFunc func(...) error)` - creates custom rule
- `AlwaysDeny(message string)` - always denies access

### Main functions

- `New(rules ...Rule)` - creates new hyperacc object
- `CheckAccess(stub, rules...)` - checks access (helper function)
- `(ha *hyperacc).Check(stub)` - checks all rules in the object

## Error handling

```go
func (c *MyContract) SomeMethod(ctx contractapi.TransactionContextInterface) error {
    err := hyperacc.CheckAccess(ctx, hyperacc.RequireRole("manager"))
    
    if err != nil {
        // Check error type
        if accErr, ok := hyperacc.AsAccessError(err); ok {
            // Special handling for access denial
            return fmt.Errorf("unauthorized: %w", err)
        }
        return err
    }
    // Method is accessible

    return nil
}
```

## Real-world scenario examples

### Supply Chain contract

```go
func (c *SupplyChainContract) CreateShipment(ctx contractapi.TransactionContextInterface, 
    shipmentID string, data string) error {
    // Only manufacturer from Org1MSP can create shipments

    if err := hyperacc.CheckAccess(ctx,
        hyperacc.And(
            hyperacc.RequireRole("manufacturer"),
            hyperacc.RequireMSPID("Org1MSP"),
        ),
    ); err != nil {
        return err
    }
    
    return ctx.GetStub().PutState(shipmentID, []byte(data))
}

func (c *SupplyChainContract) ReceiveShipment(ctx contractapi.TransactionContextInterface, 
    shipmentID string) error {
    // Only distributor or retailer can receive shipments

    if err := hyperacc.CheckAccess(ctx,
        hyperacc.RequireAnyRole("distributor", "retailer"),
    ); err != nil {
        return err
    }
    // Shipment receiving logic

    return nil
}

func (c *SupplyChainContract) AuditShipment(ctx contractapi.TransactionContextInterface, 
    shipmentID string) (string, error) {
    // Auditors from any organization or administrators

    if err := hyperacc.CheckAccess(ctx,
        hyperacc.Or(
            hyperacc.RequireRole("auditor"),
            hyperacc.RequireAdmin(),
        ),
    ); err != nil {
        return "", err
    }
    
    data, err := ctx.GetStub().GetState(shipmentID)
    return string(data), err
}
```

### Financial contract

```go
func (c *FinanceContract) TransferFunds(ctx contractapi.TransactionContextInterface, 
    from, to string, amount int) error {
    // Only managers from financial departments of banks

    if err := hyperacc.CheckAccess(ctx,
        hyperacc.And(
            hyperacc.RequireRole("finance-manager"),
            hyperacc.RequireOU("finance"),
            hyperacc.RequireAnyMSPID("BankAMSP", "BankBMSP"),
        ),
    ); err != nil {
        return err
    }
    // Fund transfer logic

    return nil
}

func (c *FinanceContract) ApproveTransaction(ctx contractapi.TransactionContextInterface, 
    txID string) error {
    // Only compliance officers or administrators

    if err := hyperacc.CheckAccess(ctx,
        hyperacc.Or(
            hyperacc.RequireRole("compliance-officer"),
            hyperacc.RequireAdmin(),
        ),
    ); err != nil {
        return err
    }
    // Transaction approval logic

    return nil
}
```

## Requirements

- Go 1.21+
- HyperLedger Fabric 2.x+
- fabric-chaincode-go/v2
- fabric-contract-api-go/v2
- fabric-protos-go-apiv2

## Documentation

### 📖 Main documentation
- [README.md](README.md) - Main documentation and examples
- [API Reference](https://pkg.go.dev/github.com/IceTweak/hyperacc) - Complete API documentation

### 🔄 Migration and updates
- [CHANGELOG.md](CHANGELOG.md) - Change history

## Examples

Complete chaincode example using hyperacc is available in [examples/asset-contract](examples/asset-contract).

## Support
- 📝 [Open issue](https://github.com/IceTweak/hyperacc/issues)

- 💬 [Start discussion](https://github.com/IceTweak/hyperacc/discussions)

## Roadmap

- [ ] Time-based policy support
- [ ] Integration with external authentication systems
- [ ] Advanced audit capabilities
- [ ] Dynamic rule support from ledger

See [open issues](https://github.com/IceTweak/hyperacc/issues) for a complete list of proposed features and known issues.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Citation

If you use this project in research or publications, please cite:

```
@software{hyperacc,
  author = {Roman Savoschik [github.com/IceTweak]},
  title = {Access control for HyperLedger Fabric Go chaincodes},
  year = {2026},
  url = {https://github.com/IceTweak/hyperacc}
}
```

---

<div align="center">
Made with ❤️ for the HyperLedger Fabric community
</div>
