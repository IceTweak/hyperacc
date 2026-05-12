package hyperacc

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAccessError_Error(t *testing.T) {
	tests := []struct {
		name        string
		setup       func() *AccessError
		expectedMsg string
	}{
		{
			name: "error without cause",
			setup: func() *AccessError {
				return NewAccessError("access denied")
			},
			expectedMsg: "access error: access denied",
		},
		{
			name: "error with cause",
			setup: func() *AccessError {
				return WrapAccessError("access denied", errors.New("underlying error"))
			},
			expectedMsg: "access error: access denied: underlying error",
		},
		{
			name: "empty reason",
			setup: func() *AccessError {
				return NewAccessError("")
			},
			expectedMsg: "access error: ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.setup()
			assert.Contains(t, err.Error(), tt.expectedMsg)
		})
	}
}

func TestAccessError_Unwrap(t *testing.T) {
	tests := []struct {
		name        string
		setup       func() *AccessError
		expectedErr error
	}{
		{
			name: "error without cause - Unwrap returns nil",
			setup: func() *AccessError {
				return NewAccessError("access denied")
			},
			expectedErr: nil,
		},
		{
			name: "error with cause - Unwrap returns cause",
			setup: func() *AccessError {
				underlyingErr := errors.New("underlying error")
				return WrapAccessError("access denied", underlyingErr)
			},
			expectedErr: errors.New("underlying error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.setup()
			unwrapped := err.Unwrap()
			if tt.expectedErr == nil {
				assert.Nil(t, unwrapped)
			} else {
				assert.Equal(t, tt.expectedErr.Error(), unwrapped.Error())
			}
		})
	}
}

func TestNewAccessError(t *testing.T) {
	tests := []struct {
		name        string
		reason      string
		expectedErr *AccessError
	}{
		{
			name:   "creating new access error",
			reason: "access denied",
			expectedErr: &AccessError{
				Reason: "access denied",
				Cause:  nil,
			},
		},
		{
			name:   "empty reason",
			reason: "",
			expectedErr: &AccessError{
				Reason: "",
				Cause:  nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewAccessError(tt.reason)
			assert.NotNil(t, err)
			assert.Equal(t, tt.expectedErr.Reason, err.Reason)
			assert.Equal(t, tt.expectedErr.Cause, err.Cause)
		})
	}
}

func TestWrapAccessError(t *testing.T) {
	tests := []struct {
		name        string
		reason      string
		cause       error
		expectedErr *AccessError
	}{
		{
			name:   "wrapping error with cause",
			reason: "access denied",
			cause:  errors.New("underlying error"),
			expectedErr: &AccessError{
				Reason: "access denied",
				Cause:  errors.New("underlying error"),
			},
		},
		{
			name:   "wrapping nil error",
			reason: "access denied",
			cause:  nil,
			expectedErr: &AccessError{
				Reason: "access denied",
				Cause:  nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := WrapAccessError(tt.reason, tt.cause)
			assert.NotNil(t, err)
			assert.Equal(t, tt.expectedErr.Reason, err.Reason)
			if tt.cause != nil {
				assert.Equal(t, tt.expectedErr.Cause.Error(), err.Cause.Error())
			} else {
				assert.Nil(t, err.Cause)
			}
		})
	}
}

func TestAsAccessError(t *testing.T) {
	tests := []struct {
		name        string
		err         error
		expectedOk  bool
		expectedErr *AccessError
	}{
		{
			name:        "AccessError successfully identified",
			err:         NewAccessError("access denied"),
			expectedOk:  true,
			expectedErr: NewAccessError("access denied"),
		},
		{
			name:        "regular error not identified as AccessError",
			err:         errors.New("regular error"),
			expectedOk:  false,
			expectedErr: nil,
		},
		{
			name:        "wrapped AccessError identified",
			err:         WrapAccessError("access denied", errors.New("underlying")),
			expectedOk:  true,
			expectedErr: WrapAccessError("access denied", errors.New("underlying")),
		},
		{
			name:        "nil error",
			err:         nil,
			expectedOk:  false,
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accessErr, ok := AsAccessError(tt.err)
			assert.Equal(t, tt.expectedOk, ok)
			if tt.expectedOk {
				assert.NotNil(t, accessErr)
				if tt.expectedErr != nil {
					assert.Equal(t, tt.expectedErr.Reason, accessErr.Reason)
				}
			} else {
				assert.Nil(t, accessErr)
			}
		})
	}
}
