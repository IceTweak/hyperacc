package hyperacc

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/IceTweak/hyperacc/mocks"
)

func TestIDRule_Check(t *testing.T) {
	tests := []struct {
		name        string
		setup       func() (*mocks.MockTransactionContextInterface, *mocks.MockClientIdentity)
		requiredID  string
		expectError bool
	}{
		{
			name: "ID matches - access granted",
			setup: func() (*mocks.MockTransactionContextInterface, *mocks.MockClientIdentity) {
				ctxMock, identityMock := setupMocks()
				identityMock.On("GetID").Return("user-123", nil)
				return ctxMock, identityMock
			},
			requiredID:  "user-123",
			expectError: false,
		},
		{
			name: "ID does not match - access denied",
			setup: func() (*mocks.MockTransactionContextInterface, *mocks.MockClientIdentity) {
				ctxMock, identityMock := setupMocks()
				identityMock.On("GetID").Return("user-456", nil)
				return ctxMock, identityMock
			},
			requiredID:  "user-123",
			expectError: true,
		},
		{
			name: "error getting ID",
			setup: func() (*mocks.MockTransactionContextInterface, *mocks.MockClientIdentity) {
				ctxMock, identityMock := setupMocks()
				identityMock.On("GetID").Return("", errors.New("failed to get ID"))
				return ctxMock, identityMock
			},
			requiredID:  "user-123",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctxMock, identityMock := tt.setup()

			rule := RequireID(tt.requiredID)
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

func TestAnyIDRule_Check(t *testing.T) {
	tests := []struct {
		name        string
		setup       func() (*mocks.MockTransactionContextInterface, *mocks.MockClientIdentity)
		allowedIDs  []string
		expectError bool
	}{
		{
			name: "ID in allowed list - access granted",
			setup: func() (*mocks.MockTransactionContextInterface, *mocks.MockClientIdentity) {
				ctxMock, identityMock := setupMocks()
				identityMock.On("GetID").Return("user-123", nil)
				return ctxMock, identityMock
			},
			allowedIDs:  []string{"user-123", "admin-001"},
			expectError: false,
		},
		{
			name: "ID not in allowed list - access denied",
			setup: func() (*mocks.MockTransactionContextInterface, *mocks.MockClientIdentity) {
				ctxMock, identityMock := setupMocks()
				identityMock.On("GetID").Return("user-999", nil)
				return ctxMock, identityMock
			},
			allowedIDs:  []string{"user-123", "admin-001"},
			expectError: true,
		},
		{
			name: "single allowed ID matches",
			setup: func() (*mocks.MockTransactionContextInterface, *mocks.MockClientIdentity) {
				ctxMock, identityMock := setupMocks()
				identityMock.On("GetID").Return("admin-001", nil)
				return ctxMock, identityMock
			},
			allowedIDs:  []string{"admin-001"},
			expectError: false,
		},
		{
			name: "error getting ID",
			setup: func() (*mocks.MockTransactionContextInterface, *mocks.MockClientIdentity) {
				ctxMock, identityMock := setupMocks()
				identityMock.On("GetID").Return("", errors.New("failed to get ID"))
				return ctxMock, identityMock
			},
			allowedIDs:  []string{"user-123"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctxMock, identityMock := tt.setup()

			rule := RequireAnyID(tt.allowedIDs...)
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
