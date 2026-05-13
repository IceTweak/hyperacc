package hyperacc

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/IceTweak/hyperacc/mocks"
)

func TestPolicySet_DefaultRules(t *testing.T) {
	t.Run("no default rules and no specific policy - access granted", func(t *testing.T) {
		ctxMock := new(mocks.MockTransactionContextInterface)
		ps := NewPolicySet()

		err := ps.Check(ctxMock, "CreateAsset")
		assert.NoError(t, err)
		ctxMock.AssertExpectations(t)
	})

	t.Run("default rules applied when no specific policy exists", func(t *testing.T) {
		ctxMock, identityMock := setupMocks()
		identityMock.On("GetMSPID").Return("Org1MSP", nil)

		ps := NewPolicySet(RequireMSPID("Org1MSP"))
		err := ps.Check(ctxMock, "CreateAsset")

		assert.NoError(t, err)
		ctxMock.AssertExpectations(t)
		identityMock.AssertExpectations(t)
	})

	t.Run("default rules deny when not matched", func(t *testing.T) {
		ctxMock, identityMock := setupMocks()
		identityMock.On("GetMSPID").Return("Org2MSP", nil)

		ps := NewPolicySet(RequireMSPID("Org1MSP"))
		err := ps.Check(ctxMock, "CreateAsset")

		assert.Error(t, err)
		ctxMock.AssertExpectations(t)
		identityMock.AssertExpectations(t)
	})
}

func TestPolicySet_SpecificPolicy(t *testing.T) {
	t.Run("specific policy overrides defaults", func(t *testing.T) {
		ctxMock, identityMock := setupMocks()
		identityMock.On("GetAttributeValue", "role").Return("admin", true, nil)

		ps := NewPolicySet(RequireMSPID("Org1MSP"))
		ps.AddPolicy("CreateAsset", RequireRole("admin"))
		err := ps.Check(ctxMock, "CreateAsset")

		assert.NoError(t, err)
		ctxMock.AssertExpectations(t)
		identityMock.AssertExpectations(t)
	})

	t.Run("specific policy denies", func(t *testing.T) {
		ctxMock, identityMock := setupMocks()
		identityMock.On("GetAttributeValue", "role").Return("user", true, nil)

		ps := NewPolicySet()
		ps.AddPolicy("DeleteAsset", RequireRole("admin"))
		err := ps.Check(ctxMock, "DeleteAsset")

		assert.Error(t, err)
		ctxMock.AssertExpectations(t)
		identityMock.AssertExpectations(t)
	})

	t.Run("multiple functions with different policies", func(t *testing.T) {
		ctxMock, identityMock := setupMocks()
		identityMock.On("GetAttributeValue", "role").Return("admin", true, nil)

		ps := NewPolicySet(RequireMSPID("Org1MSP"))
		ps.AddPolicy("CreateAsset", RequireRole("admin"))
		ps.AddPolicy("ReadAsset", RequireRole("user"))

		// CreateAsset should pass (admin role)
		err := ps.Check(ctxMock, "CreateAsset")
		assert.NoError(t, err)

		// ReadAsset should fail (admin role, not user)
		err = ps.Check(ctxMock, "ReadAsset")
		assert.Error(t, err)
	})
}

func TestPolicySet_RemovePolicy(t *testing.T) {
	t.Run("removed policy falls back to defaults", func(t *testing.T) {
		ctxMock, identityMock := setupMocks()
		identityMock.On("GetMSPID").Return("Org1MSP", nil)

		ps := NewPolicySet(RequireMSPID("Org1MSP"))
		ps.AddPolicy("CreateAsset", RequireRole("admin"))
		ps.RemovePolicy("CreateAsset")

		err := ps.Check(ctxMock, "CreateAsset")
		assert.NoError(t, err)
		ctxMock.AssertExpectations(t)
		identityMock.AssertExpectations(t)
	})
}

func TestPolicySet_Middleware(t *testing.T) {
	t.Run("middleware wraps policy check for function", func(t *testing.T) {
		ctxMock, identityMock := setupMocks()
		identityMock.On("GetAttributeValue", "role").Return("admin", true, nil)

		ps := NewPolicySet()
		ps.AddPolicy("CreateAsset", RequireRole("admin"))
		mw := ps.Middleware("CreateAsset")

		err := mw(ctxMock)
		assert.NoError(t, err)
		ctxMock.AssertExpectations(t)
		identityMock.AssertExpectations(t)
	})

	t.Run("middleware returns error on denial", func(t *testing.T) {
		ctxMock, identityMock := setupMocks()
		identityMock.On("GetAttributeValue", "role").Return("user", true, nil)

		ps := NewPolicySet()
		ps.AddPolicy("CreateAsset", RequireRole("admin"))
		mw := ps.Middleware("CreateAsset")

		err := mw(ctxMock)
		assert.Error(t, err)
		ctxMock.AssertExpectations(t)
		identityMock.AssertExpectations(t)
	})
}

func TestPolicySet_ChainableAPI(t *testing.T) {
	t.Run("AddPolicy returns PolicySet for chaining", func(t *testing.T) {
		ps := NewPolicySet()
		ps2 := ps.AddPolicy("CreateAsset", RequireRole("admin"))
		assert.Same(t, ps, ps2)
	})

	t.Run("RemovePolicy returns PolicySet for chaining", func(t *testing.T) {
		ps := NewPolicySet()
		ps.AddPolicy("CreateAsset", RequireRole("admin"))
		ps2 := ps.RemovePolicy("CreateAsset")
		assert.Same(t, ps, ps2)
	})
}
