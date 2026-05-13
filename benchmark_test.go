package hyperacc

import (
	"testing"
)

func BenchmarkCheckAccess_SingleRule(b *testing.B) {
	ctxMock, identityMock := setupMocks()
	identityMock.On("GetMSPID").Return("Org1MSP", nil)

	rule := RequireMSPID("Org1MSP")
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = rule.Check(ctxMock)
	}

	ctxMock.AssertExpectations(b)
	identityMock.AssertExpectations(b)
}

func BenchmarkCheckAccess_AndCombinator(b *testing.B) {
	ctxMock, identityMock := setupMocks()
	identityMock.On("GetAttributeValue", "role").Return("admin", true, nil)
	identityMock.On("GetMSPID").Return("Org1MSP", nil)

	rule := All(RequireRole("admin"), RequireMSPID("Org1MSP"))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = rule.Check(ctxMock)
	}

	ctxMock.AssertExpectations(b)
	identityMock.AssertExpectations(b)
}

func BenchmarkCheckAccess_OrCombinator(b *testing.B) {
	ctxMock, identityMock := setupMocks()
	identityMock.On("GetMSPID").Return("Org2MSP", nil)

	rule := Any(RequireMSPID("Org1MSP"), RequireMSPID("Org2MSP"))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = rule.Check(ctxMock)
	}

	ctxMock.AssertExpectations(b)
	identityMock.AssertExpectations(b)
}

func BenchmarkCheckAccess_ComplexNested(b *testing.B) {
	ctxMock, identityMock := setupMocks()
	identityMock.On("GetAttributeValue", "role").Return("admin", true, nil)
	identityMock.On("GetMSPID").Return("Org1MSP", nil)

	rule := All(
		Any(RequireRole("admin"), RequireRole("superuser")),
		RequireMSPID("Org1MSP"),
	)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = rule.Check(ctxMock)
	}

	ctxMock.AssertExpectations(b)
	identityMock.AssertExpectations(b)
}

func BenchmarkCheckAccess_Controller(b *testing.B) {
	ctxMock, identityMock := setupMocks()
	identityMock.On("GetMSPID").Return("Org1MSP", nil)

	controller := New(RequireMSPID("Org1MSP"))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = controller.Check(ctxMock)
	}

	ctxMock.AssertExpectations(b)
	identityMock.AssertExpectations(b)
}

func BenchmarkCheckAccess_Helper(b *testing.B) {
	ctxMock, identityMock := setupMocks()
	identityMock.On("GetMSPID").Return("Org1MSP", nil)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = CheckAccess(ctxMock, RequireMSPID("Org1MSP"))
	}

	ctxMock.AssertExpectations(b)
	identityMock.AssertExpectations(b)
}

func BenchmarkCheckAccess_PolicySet(b *testing.B) {
	ctxMock, identityMock := setupMocks()
	identityMock.On("GetAttributeValue", "role").Return("admin", true, nil)

	ps := NewPolicySet()
	ps.AddPolicy("CreateAsset", RequireRole("admin"))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = ps.Check(ctxMock, "CreateAsset")
	}

	ctxMock.AssertExpectations(b)
	identityMock.AssertExpectations(b)
}
