package hyperacc

import (
	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

// Rule interface for access rules
type Rule interface {
	Check(ctx contractapi.TransactionContextInterface) error
}

// controller main structure for access control
type controller struct {
	rules []Rule
}

// New creates a new instance of controller
func New(rules ...Rule) *controller {
	return &controller{
		rules: rules,
	}
}

// Check checks all access rules
func (c *controller) Check(ctx contractapi.TransactionContextInterface) error {
	if len(c.rules) == 0 {
		return nil // No restrictions
	}

	for _, rule := range c.rules {
		if err := rule.Check(ctx); err != nil {
			return err
		}
	}

	return nil
}

// CheckAccess helper function for access checking
func CheckAccess(ctx contractapi.TransactionContextInterface, rules ...Rule) error {
	c := New(rules...)
	return c.Check(ctx)
}
