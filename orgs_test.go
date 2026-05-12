package hyperacc

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/IceTweak/hyperacc/mocks"
)

func TestMSPIDRule_Check(t *testing.T) {
	tests := []struct {
		name        string
		setup       func() (*mocks.MockTransactionContextInterface, *mocks.MockClientIdentity)
		requiredMSP string
		expectedErr error
		expectError bool
	}{
		{
			name: "MSPID matches - access granted",
			setup: func() (*mocks.MockTransactionContextInterface, *mocks.MockClientIdentity) {
				ctxMock, identityMock := setupMocks()
				identityMock.On("GetMSPID").Return("Org1MSP", nil)
				return ctxMock, identityMock
			},
			requiredMSP: "Org1MSP",
			expectedErr: nil,
			expectError: false,
		},
		{
			name: "MSPID does not match - access denied",
			setup: func() (*mocks.MockTransactionContextInterface, *mocks.MockClientIdentity) {
				ctxMock, identityMock := setupMocks()
				identityMock.On("GetMSPID").Return("Org2MSP", nil)
				return ctxMock, identityMock
			},
			requiredMSP: "Org1MSP",
			expectError: true,
		},
		{
			name: "error getting MSPID",
			setup: func() (*mocks.MockTransactionContextInterface, *mocks.MockClientIdentity) {
				ctxMock, identityMock := setupMocks()
				identityMock.On("GetMSPID").Return("", errors.New("failed to get MSPID"))
				return ctxMock, identityMock
			},
			requiredMSP: "Org1MSP",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctxMock, identityMock := tt.setup()

			rule := RequireMSPID(tt.requiredMSP)
			err := rule.Check(ctxMock)

			if tt.expectError {
				assert.Error(t, err)
				if tt.requiredMSP != "" {
					// Check that it's either AccessError or MSPID retrieval error
					if accessErr, ok := AsAccessError(err); ok {
						assert.Contains(t, accessErr.Reason, tt.requiredMSP)
					} else {
						assert.Contains(t, err.Error(), "failed to get MSPID")
					}
				}
			} else {
				assert.NoError(t, err)
			}

			ctxMock.AssertExpectations(t)
			identityMock.AssertExpectations(t)
		})
	}
}

func TestAnyMSPIDRule_Check(t *testing.T) {
	tests := []struct {
		name        string
		setup       func() (*mocks.MockTransactionContextInterface, *mocks.MockClientIdentity)
		allowedMSPs []string
		expectedErr error
		expectError bool
	}{
		{
			name: "MSPID in allowed list - access granted",
			setup: func() (*mocks.MockTransactionContextInterface, *mocks.MockClientIdentity) {
				ctxMock, identityMock := setupMocks()
				identityMock.On("GetMSPID").Return("Org1MSP", nil)
				return ctxMock, identityMock
			},
			allowedMSPs: []string{"Org1MSP", "Org2MSP", "Org3MSP"},
			expectedErr: nil,
			expectError: false,
		},
		{
			name: "MSPID not in allowed list - access denied",
			setup: func() (*mocks.MockTransactionContextInterface, *mocks.MockClientIdentity) {
				ctxMock, identityMock := setupMocks()
				identityMock.On("GetMSPID").Return("Org4MSP", nil)
				return ctxMock, identityMock
			},
			allowedMSPs: []string{"Org1MSP", "Org2MSP", "Org3MSP"},
			expectError: true,
		},
		{
			name: "error getting MSPID",
			setup: func() (*mocks.MockTransactionContextInterface, *mocks.MockClientIdentity) {
				ctxMock, identityMock := setupMocks()
				identityMock.On("GetMSPID").Return("", errors.New("failed to get MSPID"))
				return ctxMock, identityMock
			},
			allowedMSPs: []string{"Org1MSP", "Org2MSP"},
			expectError: true,
		},
		{
			name: "second MSPID in list - access granted",
			setup: func() (*mocks.MockTransactionContextInterface, *mocks.MockClientIdentity) {
				ctxMock, identityMock := setupMocks()
				identityMock.On("GetMSPID").Return("Org2MSP", nil)
				return ctxMock, identityMock
			},
			allowedMSPs: []string{"Org1MSP", "Org2MSP", "Org3MSP"},
			expectedErr: nil,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctxMock, identityMock := tt.setup()

			rule := RequireAnyMSPID(tt.allowedMSPs...)
			err := rule.Check(ctxMock)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			ctxMock.AssertExpectations(t)
			identityMock.AssertExpectations(t)
		})
	}
}

func TestOURule_Check(t *testing.T) {
	tests := []struct {
		name        string
		setup       func() (*mocks.MockTransactionContextInterface, *mocks.MockClientIdentity)
		requiredOU  string
		expectedErr error
		expectError bool
	}{
		{
			name: "OU matches - access granted",
			setup: func() (*mocks.MockTransactionContextInterface, *mocks.MockClientIdentity) {
				ctxMock, identityMock := setupMocks()
				identityMock.On("GetMSPID").Return("Org1MSP", nil)
				cert := &x509.Certificate{
					Subject: pkix.Name{
						OrganizationalUnit: []string{"admin", "user"},
					},
				}
				identityMock.On("GetX509Certificate").Return(cert, nil)
				return ctxMock, identityMock
			},
			requiredOU:  "admin",
			expectedErr: nil,
			expectError: false,
		},
		{
			name: "OU does not match - access denied",
			setup: func() (*mocks.MockTransactionContextInterface, *mocks.MockClientIdentity) {
				ctxMock, identityMock := setupMocks()
				identityMock.On("GetMSPID").Return("Org1MSP", nil)
				cert := &x509.Certificate{
					Subject: pkix.Name{
						OrganizationalUnit: []string{"user", "guest"},
					},
				}
				identityMock.On("GetX509Certificate").Return(cert, nil)
				return ctxMock, identityMock
			},
			requiredOU:  "admin",
			expectError: true,
		},
		{
			name: "error getting MSPID",
			setup: func() (*mocks.MockTransactionContextInterface, *mocks.MockClientIdentity) {
				ctxMock, identityMock := setupMocks()
				identityMock.On("GetMSPID").Return("", errors.New("failed to get MSPID"))
				return ctxMock, identityMock
			},
			requiredOU:  "admin",
			expectError: true,
		},
		{
			name: "error getting certificate",
			setup: func() (*mocks.MockTransactionContextInterface, *mocks.MockClientIdentity) {
				ctxMock, identityMock := setupMocks()
				identityMock.On("GetMSPID").Return("Org1MSP", nil)
				identityMock.On("GetX509Certificate").Return(nil, errors.New("failed to get certificate"))
				return ctxMock, identityMock
			},
			requiredOU:  "admin",
			expectError: true,
		},
		{
			name: "empty OU list - access denied",
			setup: func() (*mocks.MockTransactionContextInterface, *mocks.MockClientIdentity) {
				ctxMock, identityMock := setupMocks()
				identityMock.On("GetMSPID").Return("Org1MSP", nil)
				cert := &x509.Certificate{
					Subject: pkix.Name{
						OrganizationalUnit: []string{},
					},
				}
				identityMock.On("GetX509Certificate").Return(cert, nil)
				return ctxMock, identityMock
			},
			requiredOU:  "admin",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctxMock, identityMock := tt.setup()

			rule := RequireOU(tt.requiredOU)
			err := rule.Check(ctxMock)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			ctxMock.AssertExpectations(t)
			identityMock.AssertExpectations(t)
		})
	}
}

func TestAnyOURule_Check(t *testing.T) {
	tests := []struct {
		name        string
		setup       func() (*mocks.MockTransactionContextInterface, *mocks.MockClientIdentity)
		allowedOUs  []string
		expectedErr error
		expectError bool
	}{
		{
			name: "OU in allowed list - access granted",
			setup: func() (*mocks.MockTransactionContextInterface, *mocks.MockClientIdentity) {
				ctxMock, identityMock := setupMocks()
				identityMock.On("GetMSPID").Return("Org1MSP", nil)
				cert := &x509.Certificate{
					Subject: pkix.Name{
						OrganizationalUnit: []string{"admin", "user"},
					},
				}
				identityMock.On("GetX509Certificate").Return(cert, nil)
				return ctxMock, identityMock
			},
			allowedOUs:  []string{"admin", "manager"},
			expectedErr: nil,
			expectError: false,
		},
		{
			name: "OU not in allowed list - access denied",
			setup: func() (*mocks.MockTransactionContextInterface, *mocks.MockClientIdentity) {
				ctxMock, identityMock := setupMocks()
				identityMock.On("GetMSPID").Return("Org1MSP", nil)
				cert := &x509.Certificate{
					Subject: pkix.Name{
						OrganizationalUnit: []string{"user", "guest"},
					},
				}
				identityMock.On("GetX509Certificate").Return(cert, nil)
				return ctxMock, identityMock
			},
			allowedOUs:  []string{"admin", "manager"},
			expectError: true,
		},
		{
			name: "second OU in list - access granted",
			setup: func() (*mocks.MockTransactionContextInterface, *mocks.MockClientIdentity) {
				ctxMock, identityMock := setupMocks()
				identityMock.On("GetMSPID").Return("Org1MSP", nil)
				cert := &x509.Certificate{
					Subject: pkix.Name{
						OrganizationalUnit: []string{"user", "manager"},
					},
				}
				identityMock.On("GetX509Certificate").Return(cert, nil)
				return ctxMock, identityMock
			},
			allowedOUs:  []string{"admin", "manager"},
			expectedErr: nil,
			expectError: false,
		},
		{
			name: "error getting MSPID",
			setup: func() (*mocks.MockTransactionContextInterface, *mocks.MockClientIdentity) {
				ctxMock, identityMock := setupMocks()
				identityMock.On("GetMSPID").Return("", errors.New("failed to get MSPID"))
				return ctxMock, identityMock
			},
			allowedOUs:  []string{"admin"},
			expectError: true,
		},
		{
			name: "error getting certificate",
			setup: func() (*mocks.MockTransactionContextInterface, *mocks.MockClientIdentity) {
				ctxMock, identityMock := setupMocks()
				identityMock.On("GetMSPID").Return("Org1MSP", nil)
				identityMock.On("GetX509Certificate").Return(nil, errors.New("failed to get certificate"))
				return ctxMock, identityMock
			},
			allowedOUs:  []string{"admin"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctxMock, identityMock := tt.setup()

			rule := RequireAnyOU(tt.allowedOUs...)
			err := rule.Check(ctxMock)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			ctxMock.AssertExpectations(t)
			identityMock.AssertExpectations(t)
		})
	}
}
