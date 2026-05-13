package hyperacc

import (
	"errors"
	"testing"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
	"github.com/stretchr/testify/assert"
)

func TestCustomRule_Check(t *testing.T) {
	tests := []struct {
		name        string
		ruleName    string
		checkFunc   func(contractapi.TransactionContextInterface) error
		expectError bool
	}{
		{
			name:     "custom check passes",
			ruleName: "my-custom-rule",
			checkFunc: func(ctx contractapi.TransactionContextInterface) error {
				return nil
			},
			expectError: false,
		},
		{
			name:     "custom check fails with access error",
			ruleName: "failing-rule",
			checkFunc: func(ctx contractapi.TransactionContextInterface) error {
				return NewAccessError("custom failure")
			},
			expectError: true,
		},
		{
			name:     "custom check fails with generic error",
			ruleName: "generic-error-rule",
			checkFunc: func(ctx contractapi.TransactionContextInterface) error {
				return errors.New("something went wrong")
			},
			expectError: true,
		},
		{
			name:     "empty rule name",
			ruleName: "",
			checkFunc: func(ctx contractapi.TransactionContextInterface) error {
				return nil
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctxMock, _ := setupMocks()

			rule := Custom(tt.ruleName, tt.checkFunc)

			assert.Equal(t, tt.ruleName, rule.name)

			err := rule.Check(ctxMock)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAlwaysDenyRule_Check(t *testing.T) {
	tests := []struct {
		name        string
		message     string
		expectedMsg string
	}{
		{
			name:        "deny with custom message",
			message:     "custom deny message",
			expectedMsg: "access error: custom deny message",
		},
		{
			name:        "deny with empty message uses default",
			message:     "",
			expectedMsg: "access error: access denied",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule := AlwaysDeny(tt.message)
			err := rule.Check(nil)

			assert.Error(t, err)
			assert.Equal(t, tt.expectedMsg, err.Error())
		})
	}
}
