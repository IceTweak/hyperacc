package hyperacc

import (
	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

// PolicySet maps chaincode function names to access rules.
// Provides a declarative way to apply different access policies
// to different chaincode functions.
type PolicySet struct {
	defaultRules []Rule
	policies     map[string][]Rule
}

// NewPolicySet creates a new PolicySet with optional default rules.
// Default rules apply when no specific policy exists for a function.
func NewPolicySet(defaultRules ...Rule) *PolicySet {
	return &PolicySet{
		defaultRules: defaultRules,
		policies:     make(map[string][]Rule),
	}
}

// AddPolicy adds a policy for a specific function name.
// When this function is called via Check, its rules are evaluated.
func (ps *PolicySet) AddPolicy(functionName string, rules ...Rule) *PolicySet {
	ps.policies[functionName] = rules
	return ps
}

// RemovePolicy removes the policy for a specific function.
func (ps *PolicySet) RemovePolicy(functionName string) *PolicySet {
	delete(ps.policies, functionName)
	return ps
}

// Check evaluates the access policy for the given function name.
// If a specific policy exists for the function, its rules are used.
// Otherwise, default rules are applied.
func (ps *PolicySet) Check(ctx contractapi.TransactionContextInterface, functionName string) error {
	rules, ok := ps.policies[functionName]
	if !ok {
		rules = ps.defaultRules
	}

	if len(rules) == 0 {
		return nil
	}

	for _, rule := range rules {
		if err := rule.Check(ctx); err != nil {
			return err
		}
	}

	return nil
}

// Middleware creates a middleware function for a specific function name.
// Useful for integration with chaincode dispatch logic.
func (ps *PolicySet) Middleware(functionName string) Middleware {
	return func(ctx contractapi.TransactionContextInterface) error {
		return ps.Check(ctx, functionName)
	}
}
