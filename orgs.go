package hyperacc

import (
	"fmt"
	"slices"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

// MSPIDRule checks the caller's MSPID
type MSPIDRule struct {
	mspid string
}

// RequireMSPID creates a rule to check for a specific MSPID
func RequireMSPID(mspid string) *MSPIDRule {
	return &MSPIDRule{mspid: mspid}
}

// Check checks the caller's MSPID
func (r *MSPIDRule) Check(ctx contractapi.TransactionContextInterface) error {
	mspid, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("failed to get MSPID: %w", err)
	}

	if mspid != r.mspid {
		return NewAccessError(fmt.Sprintf("required MSPID '%s', got '%s'", r.mspid, mspid))
	}

	return nil
}

// AnyMSPIDRule checks for one of the specified MSPIDs
type AnyMSPIDRule struct {
	mspids []string
}

// RequireAnyMSPID creates a rule to check for one of the MSPIDs
func RequireAnyMSPID(mspids ...string) *AnyMSPIDRule {
	return &AnyMSPIDRule{mspids: mspids}
}

// Check checks if the caller has one of the MSPIDs
func (r *AnyMSPIDRule) Check(ctx contractapi.TransactionContextInterface) error {
	mspid, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("failed to get MSPID: %w", err)
	}

	if slices.Contains(r.mspids, mspid) {
		return nil
	}

	return NewAccessError(fmt.Sprintf("required one of MSPIDs %v, got '%s'", r.mspids, mspid))
}

// OURule checks the caller's Organizational Unit
type OURule struct {
	ou string
}

// RequireOU creates a rule to check for a specific OU
func RequireOU(ou string) *OURule {
	return &OURule{ou: ou}
}

// Check checks the caller's OU
func (r *OURule) Check(ctx contractapi.TransactionContextInterface) error {
	identity := ctx.GetClientIdentity()
	// Get MSPID
	mspid, err := identity.GetMSPID()
	if err != nil {
		return fmt.Errorf("failed to get MSPID: %w", err)
	}
	// Parse certificate to get OU
	cert, err := identity.GetX509Certificate()
	if err != nil {
		return fmt.Errorf("failed to get X509 certificate: %w", err)
	}

	found := slices.Contains(cert.Subject.OrganizationalUnit, r.ou)

	if !found {
		return NewAccessError(fmt.Sprintf("required OU '%s', MSP: %s", r.ou, mspid))
	}

	return nil
}

// AnyOURule checks for one of the specified OUs
type AnyOURule struct {
	ous []string
}

// RequireAnyOU creates a rule to check for one of the OUs
func RequireAnyOU(ous ...string) *AnyOURule {
	return &AnyOURule{ous: ous}
}

// Check checks if the caller has one of the OUs
func (r *AnyOURule) Check(ctx contractapi.TransactionContextInterface) error {
	identity := ctx.GetClientIdentity()
	// Get MSPID
	mspid, err := identity.GetMSPID()
	if err != nil {
		return fmt.Errorf("failed to get MSPID: %w", err)
	}
	// Parse certificate to get OU
	cert, err := identity.GetX509Certificate()
	if err != nil {
		return fmt.Errorf("failed to get X509 certificate: %w", err)
	}

	for _, certOU := range cert.Subject.OrganizationalUnit {
		if slices.Contains(r.ous, certOU) {
			return nil
		}
	}

	return NewAccessError(fmt.Sprintf("required one of OUs %v, MSP: %s", r.ous, mspid))
}
