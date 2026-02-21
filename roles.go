package hyperacc

import (
	"fmt"
	"slices"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

// RoleRule checks for a specific role
type RoleRule struct {
	role string
}

// RequireRole creates a rule to check for a specific role
func RequireRole(role string) *RoleRule {
	return &RoleRule{role: role}
}

// Check checks if the caller has the role
func (r *RoleRule) Check(ctx contractapi.TransactionContextInterface) error {
	identity := ctx.GetClientIdentity()
	role, found, err := identity.GetAttributeValue("role")
	if err != nil {
		return fmt.Errorf("failed to get role attribute: %w", err)
	}

	if !found {
		return fmt.Errorf("role attribute not found in identity")
	}

	if role != r.role {
		return NewAccessError(fmt.Sprintf("required role '%s', got '%s'", r.role, role))
	}

	return nil
}

// AnyRoleRule checks for one of the specified roles
type AnyRoleRule struct {
	roles []string
}

// RequireAnyRole creates a rule to check for one of the roles
func RequireAnyRole(roles ...string) *AnyRoleRule {
	return &AnyRoleRule{roles: roles}
}

// Check checks if the caller has one of the roles
func (r *AnyRoleRule) Check(ctx contractapi.TransactionContextInterface) error {
	identity := ctx.GetClientIdentity()
	role, found, err := identity.GetAttributeValue("role")
	if err != nil {
		return fmt.Errorf("failed to get role attribute: %w", err)
	}

	if !found {
		return fmt.Errorf("role attribute not found in identity")
	}

	if slices.Contains(r.roles, role) {
		return nil
	}

	return NewAccessError(fmt.Sprintf("required one of roles %v, got '%s'", r.roles, role))
}
