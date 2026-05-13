package hyperacc

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/IceTweak/hyperacc/mocks"
)

func TestAndRule_Check(t *testing.T) {
	tests := []struct {
		name        string
		setup       func() (*mocks.MockTransactionContextInterface, []*mocks.MockRule)
		expectedErr error
		expectError bool
	}{
		{
			name: "all rules pass",
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
			name: "first rule fails",
			setup: func() (*mocks.MockTransactionContextInterface, []*mocks.MockRule) {
				ctxMock := new(mocks.MockTransactionContextInterface)
				rule1 := new(mocks.MockRule)
				rule2 := new(mocks.MockRule)
				rule1.On("Check", ctxMock).Return(NewAccessError("rule 1 failed"))
				rule2.On("Check", ctxMock).Return(nil)
				return ctxMock, []*mocks.MockRule{rule1, rule2}
			},
			expectError: true,
		},
		{
			name: "second rule fails",
			setup: func() (*mocks.MockTransactionContextInterface, []*mocks.MockRule) {
				ctxMock := new(mocks.MockTransactionContextInterface)
				rule1 := new(mocks.MockRule)
				rule2 := new(mocks.MockRule)
				rule1.On("Check", ctxMock).Return(nil)
				rule2.On("Check", ctxMock).Return(NewAccessError("rule 2 failed"))
				return ctxMock, []*mocks.MockRule{rule1, rule2}
			},
			expectError: true,
		},
		{
			name: "all rules fail",
			setup: func() (*mocks.MockTransactionContextInterface, []*mocks.MockRule) {
				ctxMock := new(mocks.MockTransactionContextInterface)
				rule1 := new(mocks.MockRule)
				rule2 := new(mocks.MockRule)
				rule3 := new(mocks.MockRule)
				rule1.On("Check", ctxMock).Return(NewAccessError("rule 1 failed"))
				rule2.On("Check", ctxMock).Return(NewAccessError("rule 2 failed"))
				rule3.On("Check", ctxMock).Return(NewAccessError("rule 3 failed"))
				return ctxMock, []*mocks.MockRule{rule1, rule2, rule3}
			},
			expectError: true,
		},
		{
			name: "empty rules list",
			setup: func() (*mocks.MockTransactionContextInterface, []*mocks.MockRule) {
				ctxMock := new(mocks.MockTransactionContextInterface)
				return ctxMock, []*mocks.MockRule{}
			},
			expectedErr: nil,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctxMock, mockRules := tt.setup()

			rules := make([]Rule, len(mockRules))
			for i, mr := range mockRules {
				rules[i] = mr
			}

			andRule := And(rules...)
			err := andRule.Check(ctxMock)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "AND rule failed")
			} else {
				assert.NoError(t, err)
			}

			for _, mr := range mockRules {
				mr.AssertExpectations(t)
			}
		})
	}
}

func TestOrRule_Check(t *testing.T) {
	tests := []struct {
		name        string
		setup       func() (*mocks.MockTransactionContextInterface, []*mocks.MockRule)
		expectedErr error
		expectError bool
	}{
		{
			name: "first rule passes",
			setup: func() (*mocks.MockTransactionContextInterface, []*mocks.MockRule) {
				ctxMock := new(mocks.MockTransactionContextInterface)
				rule1 := new(mocks.MockRule)
				rule2 := new(mocks.MockRule)
				rule1.On("Check", ctxMock).Return(nil)
				rule2.On("Check", ctxMock).Maybe()
				return ctxMock, []*mocks.MockRule{rule1, rule2}
			},
			expectedErr: nil,
			expectError: false,
		},
		{
			name: "second rule passes",
			setup: func() (*mocks.MockTransactionContextInterface, []*mocks.MockRule) {
				ctxMock := new(mocks.MockTransactionContextInterface)
				rule1 := new(mocks.MockRule)
				rule2 := new(mocks.MockRule)
				rule1.On("Check", ctxMock).Return(NewAccessError("rule 1 failed"))
				rule2.On("Check", ctxMock).Return(nil)
				return ctxMock, []*mocks.MockRule{rule1, rule2}
			},
			expectedErr: nil,
			expectError: false,
		},
		{
			name: "all rules fail",
			setup: func() (*mocks.MockTransactionContextInterface, []*mocks.MockRule) {
				ctxMock := new(mocks.MockTransactionContextInterface)
				rule1 := new(mocks.MockRule)
				rule2 := new(mocks.MockRule)
				rule1.On("Check", ctxMock).Return(NewAccessError("rule 1 failed"))
				rule2.On("Check", ctxMock).Return(NewAccessError("rule 2 failed"))
				return ctxMock, []*mocks.MockRule{rule1, rule2}
			},
			expectError: true,
		},
		{
			name: "empty rules list",
			setup: func() (*mocks.MockTransactionContextInterface, []*mocks.MockRule) {
				ctxMock := new(mocks.MockTransactionContextInterface)
				return ctxMock, []*mocks.MockRule{}
			},
			expectError: true,
		},
		{
			name: "third rule passes",
			setup: func() (*mocks.MockTransactionContextInterface, []*mocks.MockRule) {
				ctxMock := new(mocks.MockTransactionContextInterface)
				rule1 := new(mocks.MockRule)
				rule2 := new(mocks.MockRule)
				rule3 := new(mocks.MockRule)
				rule1.On("Check", ctxMock).Return(NewAccessError("rule 1 failed"))
				rule2.On("Check", ctxMock).Return(NewAccessError("rule 2 failed"))
				rule3.On("Check", ctxMock).Return(nil)
				return ctxMock, []*mocks.MockRule{rule1, rule2, rule3}
			},
			expectedErr: nil,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctxMock, mockRules := tt.setup()

			rules := make([]Rule, len(mockRules))
			for i, mr := range mockRules {
				rules[i] = mr
			}

			orRule := Or(rules...)
			err := orRule.Check(ctxMock)

			if tt.expectError {
				assert.Error(t, err)
				if len(mockRules) == 0 {
					assert.Contains(t, err.Error(), "no rules defined")
				} else {
					assert.Contains(t, err.Error(), "OR rule failed")
				}
			} else {
				assert.NoError(t, err)
			}

			for _, mr := range mockRules {
				mr.AssertExpectations(t)
			}
		})
	}
}

func TestNotRule_Check(t *testing.T) {
	tests := []struct {
		name        string
		setup       func() (*mocks.MockTransactionContextInterface, *mocks.MockRule)
		expectedErr error
		expectError bool
	}{
		{
			name: "rule fails - NotRule passes",
			setup: func() (*mocks.MockTransactionContextInterface, *mocks.MockRule) {
				ctxMock := new(mocks.MockTransactionContextInterface)
				rule := new(mocks.MockRule)
				rule.On("Check", ctxMock).Return(NewAccessError("rule failed"))
				return ctxMock, rule
			},
			expectedErr: nil,
			expectError: false,
		},
		{
			name: "rule passes - NotRule fails",
			setup: func() (*mocks.MockTransactionContextInterface, *mocks.MockRule) {
				ctxMock := new(mocks.MockTransactionContextInterface)
				rule := new(mocks.MockRule)
				rule.On("Check", ctxMock).Return(nil)
				return ctxMock, rule
			},
			expectError: true,
		},
		{
			name: "rule returns regular error - NotRule passes",
			setup: func() (*mocks.MockTransactionContextInterface, *mocks.MockRule) {
				ctxMock := new(mocks.MockTransactionContextInterface)
				rule := new(mocks.MockRule)
				rule.On("Check", ctxMock).Return(errors.New("some error"))
				return ctxMock, rule
			},
			expectedErr: nil,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctxMock, mockRule := tt.setup()

			notRule := Not(mockRule)
			err := notRule.Check(ctxMock)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "NOT rule: rule should not pass")
			} else {
				assert.NoError(t, err)
			}

			mockRule.AssertExpectations(t)
		})
	}
}

func TestAll_IsAliasForAnd(t *testing.T) {
	ctxMock := new(mocks.MockTransactionContextInterface)
	rule1 := new(mocks.MockRule)
	rule2 := new(mocks.MockRule)
	rule1.On("Check", ctxMock).Return(nil)
	rule2.On("Check", ctxMock).Return(nil)

	allRule := All(rule1, rule2)
	andRule := And(rule1, rule2)

	allErr := allRule.Check(ctxMock)
	andErr := andRule.Check(ctxMock)

	assert.NoError(t, allErr)
	assert.Equal(t, andErr, allErr)
}

func TestAny_IsAliasForOr(t *testing.T) {
	ctxMock := new(mocks.MockTransactionContextInterface)
	rule1 := new(mocks.MockRule)
	rule2 := new(mocks.MockRule)
	rule1.On("Check", ctxMock).Return(nil)
	rule2.On("Check", ctxMock).Maybe()

	anyRule := Any(rule1, rule2)
	orRule := Or(rule1, rule2)

	anyErr := anyRule.Check(ctxMock)
	orErr := orRule.Check(ctxMock)

	assert.NoError(t, anyErr)
	assert.Equal(t, orErr, anyErr)
}
