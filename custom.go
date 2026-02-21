package hyperacc

import (
	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

// CustomRule allows creating custom rules using a function
type CustomRule struct {
	checkFunc func(ctx contractapi.TransactionContextInterface) error
	name      string
}

// Custom creates a custom rule based on a function
func Custom(name string, checkFunc func(ctx contractapi.TransactionContextInterface) error) *CustomRule {
	return &CustomRule{
		checkFunc: checkFunc,
		name:      name,
	}
}

// Check executes the custom check function
func (r *CustomRule) Check(ctx contractapi.TransactionContextInterface) error {
	return r.checkFunc(ctx)
}

// AlwaysDenyRule always denies access
type AlwaysDenyRule struct {
	message string
}

// AlwaysDeny creates a rule that always denies access
func AlwaysDeny(message string) *AlwaysDenyRule {
	return &AlwaysDenyRule{message: message}
}

// Check always returns an error (access denied)
func (r *AlwaysDenyRule) Check(ctx contractapi.TransactionContextInterface) error {
	if r.message == "" {
		return NewAccessError("access denied")
	}
	return NewAccessError(r.message)
}
