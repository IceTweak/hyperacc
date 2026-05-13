package hyperacc

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/IceTweak/hyperacc/mocks"
)

func TestRoleRule_Check(t *testing.T) {
	tests := []struct {
		name         string
		setup        func() (*mocks.MockTransactionContextInterface, *mocks.MockClientIdentity)
		requiredRole string
		expectedErr  error
		expectError  bool
	}{
		{
			name: "role matches - access granted",
			setup: func() (*mocks.MockTransactionContextInterface, *mocks.MockClientIdentity) {
				ctxMock, identityMock := setupMocks()
				identityMock.On("GetAttributeValue", "role").Return("admin", true, nil)
				return ctxMock, identityMock
			},
			requiredRole: "admin",
			expectedErr:  nil,
			expectError:  false,
		},
		{
			name: "role does not match - access denied",
			setup: func() (*mocks.MockTransactionContextInterface, *mocks.MockClientIdentity) {
				ctxMock, identityMock := setupMocks()
				identityMock.On("GetAttributeValue", "role").Return("user", true, nil)
				return ctxMock, identityMock
			},
			requiredRole: "admin",
			expectError:  true,
		},
		{
			name: "role attribute not found",
			setup: func() (*mocks.MockTransactionContextInterface, *mocks.MockClientIdentity) {
				ctxMock, identityMock := setupMocks()
				identityMock.On("GetAttributeValue", "role").Return("", false, nil)
				return ctxMock, identityMock
			},
			requiredRole: "admin",
			expectError:  true,
		},
		{
			name: "error getting role attribute",
			setup: func() (*mocks.MockTransactionContextInterface, *mocks.MockClientIdentity) {
				ctxMock, identityMock := setupMocks()
				identityMock.On("GetAttributeValue", "role").Return("", false, errors.New("failed to get attribute"))
				return ctxMock, identityMock
			},
			requiredRole: "admin",
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctxMock, identityMock := tt.setup()

			rule := RequireRole(tt.requiredRole)
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

func TestAnyRoleRule_Check(t *testing.T) {
	tests := []struct {
		name         string
		setup        func() (*mocks.MockTransactionContextInterface, *mocks.MockClientIdentity)
		allowedRoles []string
		expectedErr  error
		expectError  bool
	}{
		{
			name: "role in allowed list - access granted",
			setup: func() (*mocks.MockTransactionContextInterface, *mocks.MockClientIdentity) {
				ctxMock, identityMock := setupMocks()
				identityMock.On("GetAttributeValue", "role").Return("admin", true, nil)
				return ctxMock, identityMock
			},
			allowedRoles: []string{"admin", "manager", "user"},
			expectedErr:  nil,
			expectError:  false,
		},
		{
			name: "role not in allowed list - access denied",
			setup: func() (*mocks.MockTransactionContextInterface, *mocks.MockClientIdentity) {
				ctxMock, identityMock := setupMocks()
				identityMock.On("GetAttributeValue", "role").Return("guest", true, nil)
				return ctxMock, identityMock
			},
			allowedRoles: []string{"admin", "manager", "user"},
			expectError:  true,
		},
		{
			name: "second role in list - access granted",
			setup: func() (*mocks.MockTransactionContextInterface, *mocks.MockClientIdentity) {
				ctxMock, identityMock := setupMocks()
				identityMock.On("GetAttributeValue", "role").Return("manager", true, nil)
				return ctxMock, identityMock
			},
			allowedRoles: []string{"admin", "manager", "user"},
			expectedErr:  nil,
			expectError:  false,
		},
		{
			name: "role attribute not found",
			setup: func() (*mocks.MockTransactionContextInterface, *mocks.MockClientIdentity) {
				ctxMock, identityMock := setupMocks()
				identityMock.On("GetAttributeValue", "role").Return("", false, nil)
				return ctxMock, identityMock
			},
			allowedRoles: []string{"admin", "manager"},
			expectError:  true,
		},
		{
			name: "error getting role attribute",
			setup: func() (*mocks.MockTransactionContextInterface, *mocks.MockClientIdentity) {
				ctxMock, identityMock := setupMocks()
				identityMock.On("GetAttributeValue", "role").Return("", false, errors.New("failed to get attribute"))
				return ctxMock, identityMock
			},
			allowedRoles: []string{"admin", "manager"},
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctxMock, identityMock := tt.setup()

			rule := RequireAnyRole(tt.allowedRoles...)
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
