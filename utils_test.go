package hyperacc

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"errors"
	"testing"

	"github.com/hyperledger/fabric-chaincode-go/v2/shim"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/IceTweak/hyperacc/mocks"
)

var _ shim.ChaincodeStubInterface = (*mockStub)(nil)

func TestGetCallerInfo(t *testing.T) {
	tests := []struct {
		name         string
		setup        func() (*mocks.MockTransactionContextInterface, *mocks.MockClientIdentity)
		expectedInfo *CallerInfo
		expectError  bool
	}{
		{
			name: "successfully getting caller information",
			setup: func() (*mocks.MockTransactionContextInterface, *mocks.MockClientIdentity) {
				ctxMock, identityMock := setupMocks()
				identityMock.On("GetID").Return("user123", nil)
				identityMock.On("GetAttributeValue", "role").Return("admin", true, nil)
				cert := &x509.Certificate{
					Subject: pkix.Name{
						OrganizationalUnit: []string{"admin", "user"},
					},
				}
				identityMock.On("GetX509Certificate").Return(cert, nil)
				identityMock.On("GetMSPID").Return("Org1MSP", nil)
				return ctxMock, identityMock
			},
			expectedInfo: &CallerInfo{
				ID:    "user123",
				Role:  "admin",
				OUs:   []string{"admin", "user"},
				MSPID: "Org1MSP",
			},
			expectError: false,
		},
		{
			name: "role missing",
			setup: func() (*mocks.MockTransactionContextInterface, *mocks.MockClientIdentity) {
				ctxMock, identityMock := setupMocks()
				identityMock.On("GetID").Return("user123", nil)
				identityMock.On("GetAttributeValue", "role").Return("", false, nil)
				cert := &x509.Certificate{
					Subject: pkix.Name{
						OrganizationalUnit: []string{"user"},
					},
				}
				identityMock.On("GetX509Certificate").Return(cert, nil)
				identityMock.On("GetMSPID").Return("Org1MSP", nil)
				return ctxMock, identityMock
			},
			expectedInfo: &CallerInfo{
				ID:    "user123",
				Role:  "",
				OUs:   []string{"user"},
				MSPID: "Org1MSP",
			},
			expectError: false,
		},
		{
			name: "error getting ID",
			setup: func() (*mocks.MockTransactionContextInterface, *mocks.MockClientIdentity) {
				ctxMock, identityMock := setupMocks()
				identityMock.On("GetID").Return("", errors.New("failed to get ID"))
				return ctxMock, identityMock
			},
			expectError: true,
		},
		{
			name: "error getting role",
			setup: func() (*mocks.MockTransactionContextInterface, *mocks.MockClientIdentity) {
				ctxMock, identityMock := setupMocks()
				identityMock.On("GetID").Return("user123", nil)
				identityMock.On("GetAttributeValue", "role").Return("", false, errors.New("failed to get role"))
				return ctxMock, identityMock
			},
			expectError: true,
		},
		{
			name: "error getting MSPID",
			setup: func() (*mocks.MockTransactionContextInterface, *mocks.MockClientIdentity) {
				ctxMock, identityMock := setupMocks()
				identityMock.On("GetID").Return("user123", nil)
				identityMock.On("GetAttributeValue", "role").Return("admin", true, nil)
				identityMock.On("GetMSPID").Return("", errors.New("failed to get MSPID"))
				return ctxMock, identityMock
			},
			expectError: true,
		},
		{
			name: "error getting certificate",
			setup: func() (*mocks.MockTransactionContextInterface, *mocks.MockClientIdentity) {
				ctxMock, identityMock := setupMocks()
				identityMock.On("GetID").Return("user123", nil)
				identityMock.On("GetAttributeValue", "role").Return("admin", true, nil)
				identityMock.On("GetMSPID").Return("Org1MSP", nil)
				identityMock.On("GetX509Certificate").Return(nil, errors.New("failed to get certificate"))
				return ctxMock, identityMock
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctxMock, identityMock := tt.setup()

			info, err := GetCallerInfo(ctxMock)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, info)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, info)
				if tt.expectedInfo != nil {
					assert.Equal(t, tt.expectedInfo.ID, info.ID)
					assert.Equal(t, tt.expectedInfo.Role, info.Role)
					assert.Equal(t, tt.expectedInfo.OUs, info.OUs)
					assert.Equal(t, tt.expectedInfo.MSPID, info.MSPID)
				}
			}

			ctxMock.AssertExpectations(t)
			identityMock.AssertExpectations(t)
		})
	}
}

func TestHasAttribute(t *testing.T) {
	tests := []struct {
		name        string
		setup       func() (*mocks.MockTransactionContextInterface, *mocks.MockClientIdentity)
		attrName    string
		expected    bool
		expectError bool
	}{
		{
			name: "attribute found",
			setup: func() (*mocks.MockTransactionContextInterface, *mocks.MockClientIdentity) {
				ctxMock, identityMock := setupMocks()
				identityMock.On("GetAttributeValue", "role").Return("admin", true, nil)
				return ctxMock, identityMock
			},
			attrName:    "role",
			expected:    true,
			expectError: false,
		},
		{
			name: "attribute not found",
			setup: func() (*mocks.MockTransactionContextInterface, *mocks.MockClientIdentity) {
				ctxMock, identityMock := setupMocks()
				identityMock.On("GetAttributeValue", "role").Return("", false, nil)
				return ctxMock, identityMock
			},
			attrName:    "role",
			expected:    false,
			expectError: false,
		},
		{
			name: "error checking attribute",
			setup: func() (*mocks.MockTransactionContextInterface, *mocks.MockClientIdentity) {
				ctxMock, identityMock := setupMocks()
				identityMock.On("GetAttributeValue", "role").Return("", false, errors.New("failed to check attribute"))
				return ctxMock, identityMock
			},
			attrName:    "role",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctxMock, identityMock := tt.setup()

			hasAttr, err := HasAttribute(ctxMock, tt.attrName)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, hasAttr)
			}

			ctxMock.AssertExpectations(t)
			identityMock.AssertExpectations(t)
		})
	}
}

func TestAttributeRule_Check(t *testing.T) {
	tests := []struct {
		name        string
		setup       func() (*mocks.MockTransactionContextInterface, *mocks.MockClientIdentity)
		attrName    string
		attrValue   string
		expectedErr error
		expectError bool
	}{
		{
			name: "attribute found and value matches",
			setup: func() (*mocks.MockTransactionContextInterface, *mocks.MockClientIdentity) {
				ctxMock, identityMock := setupMocks()
				identityMock.On("GetAttributeValue", "department").Return("IT", true, nil)
				return ctxMock, identityMock
			},
			attrName:    "department",
			attrValue:   "IT",
			expectedErr: nil,
			expectError: false,
		},
		{
			name: "attribute found but value does not match",
			setup: func() (*mocks.MockTransactionContextInterface, *mocks.MockClientIdentity) {
				ctxMock, identityMock := setupMocks()
				identityMock.On("GetAttributeValue", "department").Return("HR", true, nil)
				return ctxMock, identityMock
			},
			attrName:    "department",
			attrValue:   "IT",
			expectError: true,
		},
		{
			name: "attribute not found",
			setup: func() (*mocks.MockTransactionContextInterface, *mocks.MockClientIdentity) {
				ctxMock, identityMock := setupMocks()
				identityMock.On("GetAttributeValue", "department").Return("", false, nil)
				return ctxMock, identityMock
			},
			attrName:    "department",
			attrValue:   "IT",
			expectError: true,
		},
		{
			name: "value not specified - only presence is checked",
			setup: func() (*mocks.MockTransactionContextInterface, *mocks.MockClientIdentity) {
				ctxMock, identityMock := setupMocks()
				identityMock.On("GetAttributeValue", "department").Return("IT", true, nil)
				return ctxMock, identityMock
			},
			attrName:    "department",
			attrValue:   "",
			expectedErr: nil,
			expectError: false,
		},
		{
			name: "error getting attribute",
			setup: func() (*mocks.MockTransactionContextInterface, *mocks.MockClientIdentity) {
				ctxMock, identityMock := setupMocks()
				identityMock.On("GetAttributeValue", "department").Return("", false, errors.New("failed to get attribute"))
				return ctxMock, identityMock
			},
			attrName:    "department",
			attrValue:   "IT",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctxMock, identityMock := tt.setup()

			rule := RequireAttribute(tt.attrName, tt.attrValue)
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

func TestHasAttributeRule_Check(t *testing.T) {
	tests := []struct {
		name        string
		setup       func() (*mocks.MockTransactionContextInterface, *mocks.MockClientIdentity)
		attrName    string
		expectedErr error
		expectError bool
	}{
		{
			name: "attribute found - access granted",
			setup: func() (*mocks.MockTransactionContextInterface, *mocks.MockClientIdentity) {
				ctxMock, identityMock := setupMocks()
				identityMock.On("GetAttributeValue", "department").Return("IT", true, nil)
				return ctxMock, identityMock
			},
			attrName:    "department",
			expectedErr: nil,
			expectError: false,
		},
		{
			name: "attribute not found - access denied",
			setup: func() (*mocks.MockTransactionContextInterface, *mocks.MockClientIdentity) {
				ctxMock, identityMock := setupMocks()
				identityMock.On("GetAttributeValue", "department").Return("", false, nil)
				return ctxMock, identityMock
			},
			attrName:    "department",
			expectError: true,
		},
		{
			name: "error checking attribute",
			setup: func() (*mocks.MockTransactionContextInterface, *mocks.MockClientIdentity) {
				ctxMock, identityMock := setupMocks()
				identityMock.On("GetAttributeValue", "department").Return("", false, errors.New("failed to check attribute"))
				return ctxMock, identityMock
			},
			attrName:    "department",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctxMock, identityMock := tt.setup()

			rule := RequireHasAttribute(tt.attrName)
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

func TestIsHLFAdmintRule_Check(t *testing.T) {
	tests := []struct {
		name        string
		setup       func() (*mocks.MockTransactionContextInterface, *mocks.MockClientIdentity)
		expectedErr error
		expectError bool
	}{
		{
			name: "is administrator - access granted",
			setup: func() (*mocks.MockTransactionContextInterface, *mocks.MockClientIdentity) {
				ctxMock, identityMock := setupMocks()
				identityMock.On("AssertAttributeValue", "hf.Type", "admin").Return(nil)
				return ctxMock, identityMock
			},
			expectedErr: nil,
			expectError: false,
		},
		{
			name: "is not administrator - access denied",
			setup: func() (*mocks.MockTransactionContextInterface, *mocks.MockClientIdentity) {
				ctxMock, identityMock := setupMocks()
				identityMock.On("AssertAttributeValue", "hf.Type", "admin").Return(errors.New("not an admin"))
				return ctxMock, identityMock
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctxMock, identityMock := tt.setup()

			rule := RequireHLFAdmin()
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

func TestIsHLFClientRule_Check(t *testing.T) {
	tests := []struct {
		name        string
		setup       func() (*mocks.MockTransactionContextInterface, *mocks.MockClientIdentity)
		expectedErr error
		expectError bool
	}{
		{
			name: "is client - access granted",
			setup: func() (*mocks.MockTransactionContextInterface, *mocks.MockClientIdentity) {
				ctxMock, identityMock := setupMocks()
				identityMock.On("AssertAttributeValue", "hf.Type", "client").Return(nil)
				return ctxMock, identityMock
			},
			expectedErr: nil,
			expectError: false,
		},
		{
			name: "is not client - access denied",
			setup: func() (*mocks.MockTransactionContextInterface, *mocks.MockClientIdentity) {
				ctxMock, identityMock := setupMocks()
				identityMock.On("AssertAttributeValue", "hf.Type", "client").Return(errors.New("not a client"))
				return ctxMock, identityMock
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctxMock, identityMock := tt.setup()

			rule := RequireHLFClient()
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

func TestCreateMiddleware(t *testing.T) {
	tests := []struct {
		name        string
		setup       func() (*mocks.MockTransactionContextInterface, []Rule)
		expectedErr error
		expectError bool
	}{
		{
			name: "middleware successfully created and works",
			setup: func() (*mocks.MockTransactionContextInterface, []Rule) {
				ctxMock, identityMock := setupMocks()
				identityMock.On("GetMSPID").Return("Org1MSP", nil)
				rule := RequireMSPID("Org1MSP")
				return ctxMock, []Rule{rule}
			},
			expectedErr: nil,
			expectError: false,
		},
		{
			name: "middleware returns error on check",
			setup: func() (*mocks.MockTransactionContextInterface, []Rule) {
				ctxMock, identityMock := setupMocks()
				identityMock.On("GetMSPID").Return("Org2MSP", nil)
				rule := RequireMSPID("Org1MSP")
				return ctxMock, []Rule{rule}
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctxMock, rules := tt.setup()

			middleware := CreateMiddleware(rules...)
			err := middleware(ctxMock)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			ctxMock.AssertExpectations(t)
		})
	}
}

func TestLogAccessDenied(t *testing.T) {
	tests := []struct {
		name        string
		setup       func() (*mocks.MockTransactionContextInterface, *mocks.MockClientIdentity, *mockStub)
		err         error
		expectEvent bool
		eventName   string
	}{
		{
			name: "logging AccessError with full caller information",
			setup: func() (*mocks.MockTransactionContextInterface, *mocks.MockClientIdentity, *mockStub) {
				ctxMock, identityMock := setupMocks()
				stubMock := new(mockStub)
				ctxMock.On("GetStub").Return(stubMock)
				identityMock.On("GetID").Return("user123", nil)
				identityMock.On("GetAttributeValue", "role").Return("admin", true, nil)
				cert := &x509.Certificate{
					Subject: pkix.Name{
						OrganizationalUnit: []string{"admin"},
					},
				}
				identityMock.On("GetX509Certificate").Return(cert, nil)
				identityMock.On("GetMSPID").Return("Org1MSP", nil)
				stubMock.On("SetEvent", "AccessDenied", mock.AnythingOfType("[]uint8")).Return(nil)
				return ctxMock, identityMock, stubMock
			},
			err:         NewAccessError("access denied"),
			expectEvent: true,
			eventName:   "AccessDenied",
		},
		{
			name: "logging regular error with full information",
			setup: func() (*mocks.MockTransactionContextInterface, *mocks.MockClientIdentity, *mockStub) {
				ctxMock, identityMock := setupMocks()
				stubMock := new(mockStub)
				ctxMock.On("GetStub").Return(stubMock)
				identityMock.On("GetID").Return("user123", nil)
				identityMock.On("GetAttributeValue", "role").Return("admin", true, nil)
				cert := &x509.Certificate{
					Subject: pkix.Name{
						OrganizationalUnit: []string{"admin"},
					},
				}
				identityMock.On("GetX509Certificate").Return(cert, nil)
				identityMock.On("GetMSPID").Return("Org1MSP", nil)
				stubMock.On("SetEvent", "AccessDenied", mock.AnythingOfType("[]uint8")).Return(nil)
				return ctxMock, identityMock, stubMock
			},
			err:         errors.New("some error"),
			expectEvent: true,
			eventName:   "AccessDenied",
		},
		{
			name: "error getting caller information",
			setup: func() (*mocks.MockTransactionContextInterface, *mocks.MockClientIdentity, *mockStub) {
				ctxMock, identityMock := setupMocks()
				stubMock := new(mockStub)
				ctxMock.On("GetStub").Return(stubMock)
				identityMock.On("GetID").Return("", errors.New("failed to get ID"))
				stubMock.On("SetEvent", "AccessDenied", mock.AnythingOfType("[]uint8")).Return(nil)
				return ctxMock, identityMock, stubMock
			},
			err:         NewAccessError("access denied"),
			expectEvent: true,
			eventName:   "AccessDenied",
		},
		{
			name: "logging without role",
			setup: func() (*mocks.MockTransactionContextInterface, *mocks.MockClientIdentity, *mockStub) {
				ctxMock, identityMock := setupMocks()
				stubMock := new(mockStub)
				ctxMock.On("GetStub").Return(stubMock)
				identityMock.On("GetID").Return("user123", nil)
				identityMock.On("GetAttributeValue", "role").Return("", false, nil)
				cert := &x509.Certificate{
					Subject: pkix.Name{
						OrganizationalUnit: []string{"user"},
					},
				}
				identityMock.On("GetX509Certificate").Return(cert, nil)
				identityMock.On("GetMSPID").Return("Org1MSP", nil)
				stubMock.On("SetEvent", "AccessDenied", mock.AnythingOfType("[]uint8")).Return(nil)
				return ctxMock, identityMock, stubMock
			},
			err:         NewAccessError("access denied"),
			expectEvent: true,
			eventName:   "AccessDenied",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctxMock, identityMock, stubMock := tt.setup()

			LogAccessDenied(ctxMock, tt.err)

			if tt.expectEvent {
				stubMock.AssertCalled(t, "SetEvent", tt.eventName, mock.AnythingOfType("[]uint8"))
			}

			ctxMock.AssertExpectations(t)
			identityMock.AssertExpectations(t)
			stubMock.AssertExpectations(t)
		})
	}
}
