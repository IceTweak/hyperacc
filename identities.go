package hyperacc

import (
	"fmt"
	"slices"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

// IDRule checks the caller's enrollment ID
type IDRule struct {
	id string
}

// RequireID creates a rule to check for a specific enrollment ID
func RequireID(id string) *IDRule {
	return &IDRule{id: id}
}

// Check checks the caller's enrollment ID
func (r *IDRule) Check(ctx contractapi.TransactionContextInterface) error {
	identity := ctx.GetClientIdentity()
	id, err := identity.GetID()
	if err != nil {
		return fmt.Errorf("failed to get ID: %w", err)
	}

	if id != r.id {
		return NewAccessError(fmt.Sprintf("required ID '%s', got '%s'", r.id, id))
	}

	return nil
}

// AnyIDRule checks for one of the specified enrollment IDs
type AnyIDRule struct {
	ids []string
}

// RequireAnyID creates a rule to check for one of the enrollment IDs
func RequireAnyID(ids ...string) *AnyIDRule {
	return &AnyIDRule{ids: ids}
}

// Check checks if the caller has one of the enrollment IDs
func (r *AnyIDRule) Check(ctx contractapi.TransactionContextInterface) error {
	identity := ctx.GetClientIdentity()
	id, err := identity.GetID()
	if err != nil {
		return fmt.Errorf("failed to get ID: %w", err)
	}

	if slices.Contains(r.ids, id) {
		return nil
	}

	return NewAccessError(fmt.Sprintf("required one of IDs %v, got '%s'", r.ids, id))
}

// AffiliationRule checks the caller's affiliation (hf.Affiliation attribute)
type AffiliationRule struct {
	affiliation string
}

// RequireAffiliation creates a rule to check for a specific affiliation
func RequireAffiliation(affiliation string) *AffiliationRule {
	return &AffiliationRule{affiliation: affiliation}
}

// Check checks the caller's affiliation
func (r *AffiliationRule) Check(ctx contractapi.TransactionContextInterface) error {
	identity := ctx.GetClientIdentity()
	err := identity.AssertAttributeValue("hf.Affiliation", r.affiliation)
	if err != nil {
		return WrapAccessError(
			fmt.Sprintf("required affiliation '%s'", r.affiliation), err,
		)
	}
	return nil
}
