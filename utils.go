package hyperacc

import (
	"errors"
	"fmt"
	"log"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

// CallerInfo contains information about the caller
type CallerInfo struct {
	MSPID string
	ID    string
	Role  string
	OUs   []string
}

// GetCallerInfo gets information about the current caller
func GetCallerInfo(ctx contractapi.TransactionContextInterface) (*CallerInfo, error) {
	info := &CallerInfo{}
	identity := ctx.GetClientIdentity()

	// Get ID
	id, err := identity.GetID()
	if err != nil {
		return nil, fmt.Errorf("failed to get ID: %w", err)
	}
	info.ID = id

	// Get role (may not exist)
	role, found, err := identity.GetAttributeValue("role")
	if err != nil {
		return nil, fmt.Errorf("failed to get role: %w", err)
	}
	if found {
		info.Role = role
	}

	// Get MSPID
	mspid, err := identity.GetMSPID()
	if err != nil {
		return nil, fmt.Errorf("failed to get MSPID: %w", err)
	}
	info.MSPID = mspid

	// Get OUs
	cert, err := identity.GetX509Certificate()
	if err != nil {
		return nil, fmt.Errorf("failed to get certificate: %w", err)
	}
	info.OUs = cert.Subject.OrganizationalUnit

	return info, nil
}

// HasAttribute checks if the caller has an attribute
func HasAttribute(ctx contractapi.TransactionContextInterface, attrName string) (bool, error) {
	identity := ctx.GetClientIdentity()
	_, found, err := identity.GetAttributeValue(attrName)
	if err != nil {
		return false, fmt.Errorf("failed to check attribute: %w", err)
	}
	return found, nil
}

// AttributeRule checks for the presence and value of an arbitrary attribute
type AttributeRule struct {
	name  string
	value string // if empty, only presence is checked
}

// RequireAttribute creates a rule to check an attribute
func RequireAttribute(name, value string) *AttributeRule {
	return &AttributeRule{
		name:  name,
		value: value,
	}
}

// Check checks the attribute
func (r *AttributeRule) Check(ctx contractapi.TransactionContextInterface) error {
	identity := ctx.GetClientIdentity()
	attrValue, found, err := identity.GetAttributeValue(r.name)
	if err != nil {
		return fmt.Errorf("failed to get attribute '%s': %w", r.name, err)
	}

	if !found {
		return NewAccessError(fmt.Sprintf("attribute '%s' not found", r.name))
	}

	if r.value != "" && attrValue != r.value {
		return NewAccessError(fmt.Sprintf("attribute '%s' has value '%s', expected '%s'", r.name, attrValue, r.value))
	}

	return nil
}

// HasAttributeRule checks only for the presence of an attribute (value is not important)
type HasAttributeRule struct {
	name string
}

// RequireHasAttribute creates a rule to check for the presence of an attribute
func RequireHasAttribute(name string) *HasAttributeRule {
	return &HasAttributeRule{name: name}
}

// Check checks for the presence of an attribute
func (r *HasAttributeRule) Check(ctx contractapi.TransactionContextInterface) error {
	identity := ctx.GetClientIdentity()
	_, found, err := identity.GetAttributeValue(r.name)
	if err != nil {
		return fmt.Errorf("failed to check attribute '%s': %w", r.name, err)
	}

	if !found {
		return NewAccessError(fmt.Sprintf("attribute '%s' not found", r.name))
	}

	return nil
}

// IsHLFAdmintRule checks if the caller is
// an administrator in the HyperLedger Fabric network
type IsHLFAdmintRule struct{}

// RequireHLFAdmin creates a rule to check that the caller is an administrator
func RequireHLFAdmin() *IsHLFAdmintRule {
	return &IsHLFAdmintRule{}
}

// Check checks the caller type
func (r *IsHLFAdmintRule) Check(ctx contractapi.TransactionContextInterface) error {
	identity := ctx.GetClientIdentity()
	err := identity.AssertAttributeValue("hf.Type", "admin")
	if err != nil {
		return WrapAccessError("admin type required in hf.Type attribute", err)
	}

	return nil
}

// IsHLFClientRule checks if the caller is
// a client in the HyperLedger Fabric network
type IsHLFClientRule struct{}

// RequireHLFClient creates a rule to check that the caller is a client
func RequireHLFClient() *IsHLFClientRule {
	return &IsHLFClientRule{}
}

// Check checks the caller type
func (r *IsHLFClientRule) Check(ctx contractapi.TransactionContextInterface) error {
	identity := ctx.GetClientIdentity()
	err := identity.AssertAttributeValue("hf.Type", "client")
	if err != nil {
		return WrapAccessError("client type required in hf.Type attribute", err)
	}

	return nil
}

// Middleware function for integration with various frameworks
type Middleware func(ctx contractapi.TransactionContextInterface) error

// CreateMiddleware creates a middleware function from rules
func CreateMiddleware(rules ...Rule) Middleware {
	return func(ctx contractapi.TransactionContextInterface) error {
		return CheckAccess(ctx, rules...)
	}
}

// LogAccessDenied logs access denial
// Automatically determines if the error is an access denial error
func LogAccessDenied(ctx contractapi.TransactionContextInterface, err error) {
	stub := ctx.GetStub()

	// Extract caller information
	info, infoErr := GetCallerInfo(ctx)
	if infoErr != nil {
		err := stub.SetEvent("AccessDenied", fmt.Appendf(nil, "Error: %v", err))
		if err != nil {
			log.Printf("failed to create AccessDenied event: %v", err)
		}
		return
	}

	// Check if the error is an access error
	var accessErr *AccessError
	var reason string
	if errors.As(err, &accessErr) {
		reason = accessErr.Reason
	} else {
		reason = err.Error()
	}

	logMsg := fmt.Sprintf("Access denied for MSPID=%s, ID=%s, Role=%s, Reason: %s",
		info.MSPID, info.ID, info.Role, reason)
	err = stub.SetEvent("AccessDenied", []byte(logMsg))
	if err != nil {
		log.Printf("failed to create AccessDenied event: %v", err)
	}
}
