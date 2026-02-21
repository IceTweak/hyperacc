package hyperacc

import (
	"fmt"
	"strings"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

// AndRule combines multiple rules with AND logic (all must pass)
type AndRule struct {
	rules []Rule
}

// And creates a rule that requires all nested rules to pass
func And(rules ...Rule) *AndRule {
	return &AndRule{rules: rules}
}

// Check checks all rules (all must pass)
func (r *AndRule) Check(ctx contractapi.TransactionContextInterface) error {
	var errors []string

	for i, rule := range r.rules {
		if err := rule.Check(ctx); err != nil {
			errors = append(errors, fmt.Sprintf("rule %d: %v", i+1, err))
		}
	}

	if len(errors) > 0 {
		return NewAccessError(fmt.Sprintf("AND rule failed: %s", strings.Join(errors, "; ")))
	}

	return nil
}

// OrRule combines multiple rules with OR logic (at least one must pass)
type OrRule struct {
	rules []Rule
}

// Or creates a rule that requires at least one nested rule to pass
func Or(rules ...Rule) *OrRule {
	return &OrRule{rules: rules}
}

// Check checks rules (at least one must pass)
func (r *OrRule) Check(ctx contractapi.TransactionContextInterface) error {
	if len(r.rules) == 0 {
		return NewAccessError("OR rule: no rules defined")
	}

	var errors []string

	for i, rule := range r.rules {
		if err := rule.Check(ctx); err == nil {
			return nil // At least one rule passed
		} else {
			errors = append(errors, fmt.Sprintf("rule %d: %v", i+1, err))
		}
	}

	return NewAccessError(fmt.Sprintf("OR rule failed: none of the rules passed: %s", strings.Join(errors, "; ")))
}

// NotRule inverts a rule
type NotRule struct {
	rule Rule
}

// Not creates a rule that inverts the result of a nested rule
func Not(rule Rule) *NotRule {
	return &NotRule{rule: rule}
}

// Check inverts the result of rule checking
func (r *NotRule) Check(ctx contractapi.TransactionContextInterface) error {
	err := r.rule.Check(ctx)
	if err == nil {
		return NewAccessError("NOT rule: rule should not pass")
	}
	return nil
}
