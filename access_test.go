package hyperacc

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/IceTweak/hyperacc/mocks"
)

func setupMocks() (*mocks.MockTransactionContextInterface, *mocks.MockClientIdentity) {
	ctxMock := new(mocks.MockTransactionContextInterface)
	identityMock := new(mocks.MockClientIdentity)
	ctxMock.On("GetClientIdentity").Return(identityMock)
	return ctxMock, identityMock
}

func TestAccessControl_Check(t *testing.T) {
	tests := []struct {
		name        string
		setup       func() (*mocks.MockTransactionContextInterface, []*mocks.MockRule)
		expectedErr error
		expectError bool
	}{
		{
			name: "no rules - access granted",
			setup: func() (*mocks.MockTransactionContextInterface, []*mocks.MockRule) {
				ctxMock := new(mocks.MockTransactionContextInterface)
				return ctxMock, []*mocks.MockRule{}
			},
			expectedErr: nil,
			expectError: false,
		},
		{
			name: "one rule passes",
			setup: func() (*mocks.MockTransactionContextInterface, []*mocks.MockRule) {
				ctxMock := new(mocks.MockTransactionContextInterface)
				rule := new(mocks.MockRule)
				rule.On("Check", ctxMock).Return(nil)
				return ctxMock, []*mocks.MockRule{rule}
			},
			expectedErr: nil,
			expectError: false,
		},
		{
			name: "multiple rules pass",
			setup: func() (*mocks.MockTransactionContextInterface, []*mocks.MockRule) {
				ctxMock := new(mocks.MockTransactionContextInterface)
				rule1 := new(mocks.MockRule)
				rule2 := new(mocks.MockRule)
				rule3 := new(mocks.MockRule)
				rule1.On("Check", ctxMock).Return(nil)
				rule2.On("Check", ctxMock).Return(nil)
				rule3.On("Check", ctxMock).Return(nil)
				return ctxMock, []*mocks.MockRule{rule1, rule2, rule3}
			},
			expectedErr: nil,
			expectError: false,
		},
		{
			name: "first rule fails - access denied",
			setup: func() (*mocks.MockTransactionContextInterface, []*mocks.MockRule) {
				ctxMock := new(mocks.MockTransactionContextInterface)
				rule1 := new(mocks.MockRule)
				rule2 := new(mocks.MockRule)
				rule1.On("Check", ctxMock).Return(NewAccessError("access denied"))
				rule2.On("Check", ctxMock).Maybe()
				return ctxMock, []*mocks.MockRule{rule1, rule2}
			},
			expectError: true,
		},
		{
			name: "second rule fails - access denied",
			setup: func() (*mocks.MockTransactionContextInterface, []*mocks.MockRule) {
				ctxMock := new(mocks.MockTransactionContextInterface)
				rule1 := new(mocks.MockRule)
				rule2 := new(mocks.MockRule)
				rule1.On("Check", ctxMock).Return(nil)
				rule2.On("Check", ctxMock).Return(NewAccessError("access denied"))
				return ctxMock, []*mocks.MockRule{rule1, rule2}
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctxMock, mockRules := tt.setup()

			rules := make([]Rule, len(mockRules))
			for i, mr := range mockRules {
				rules[i] = mr
			}

			ac := New(rules...)
			err := ac.Check(ctxMock)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			for _, mr := range mockRules {
				mr.AssertExpectations(t)
			}
		})
	}
}

func TestCheckAccess(t *testing.T) {
	tests := []struct {
		name        string
		setup       func() (*mocks.MockTransactionContextInterface, []*mocks.MockRule)
		expectedErr error
		expectError bool
	}{
		{
			name: "helper function works correctly",
			setup: func() (*mocks.MockTransactionContextInterface, []*mocks.MockRule) {
				ctxMock := new(mocks.MockTransactionContextInterface)
				rule1 := new(mocks.MockRule)
				rule2 := new(mocks.MockRule)
				rule1.On("Check", ctxMock).Return(nil)
				rule2.On("Check", ctxMock).Return(nil)
				return ctxMock, []*mocks.MockRule{rule1, rule2}
			},
			expectedErr: nil,
			expectError: false,
		},
		{
			name: "helper function returns error",
			setup: func() (*mocks.MockTransactionContextInterface, []*mocks.MockRule) {
				ctxMock := new(mocks.MockTransactionContextInterface)
				rule := new(mocks.MockRule)
				rule.On("Check", ctxMock).Return(NewAccessError("access denied"))
				return ctxMock, []*mocks.MockRule{rule}
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctxMock, mockRules := tt.setup()

			rules := make([]Rule, len(mockRules))
			for i, mr := range mockRules {
				rules[i] = mr
			}

			err := CheckAccess(ctxMock, rules...)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			for _, mr := range mockRules {
				mr.AssertExpectations(t)
			}
		})
	}
}
