package hyperacc

import (
	"github.com/hyperledger/fabric-chaincode-go/v2/shim"
	"github.com/hyperledger/fabric-protos-go-apiv2/peer"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type mockStub struct {
	mock.Mock
}

func (m *mockStub) GetArgs() [][]byte                                       { return nil }
func (m *mockStub) GetStringArgs() []string                                 { return nil }
func (m *mockStub) GetFunctionAndParameters() (string, []string)            { return "", nil }
func (m *mockStub) GetArgsSlice() ([]byte, error)                           { return nil, nil }
func (m *mockStub) GetTxID() string                                         { return "" }
func (m *mockStub) GetChannelID() string                                    { return "" }
func (m *mockStub) InvokeChaincode(string, [][]byte, string) *peer.Response { return nil }
func (m *mockStub) GetState(string) ([]byte, error)                         { return nil, nil }
func (m *mockStub) GetMultipleStates(...string) ([][]byte, error)           { return nil, nil }
func (m *mockStub) PutState(string, []byte) error                           { return nil }
func (m *mockStub) DelState(string) error                                   { return nil }
func (m *mockStub) SetStateValidationParameter(string, []byte) error        { return nil }
func (m *mockStub) GetStateValidationParameter(string) ([]byte, error)      { return nil, nil }
func (m *mockStub) GetStateByRange(string, string) (shim.StateQueryIteratorInterface, error) {
	return nil, nil
}

func (m *mockStub) GetStateByRangeWithPagination(string, string, int32, string) (shim.StateQueryIteratorInterface, *peer.QueryResponseMetadata, error) {
	return nil, nil, nil
}

func (m *mockStub) GetStateByPartialCompositeKey(string, []string) (shim.StateQueryIteratorInterface, error) {
	return nil, nil
}

func (m *mockStub) GetStateByPartialCompositeKeyWithPagination(string, []string, int32, string) (shim.StateQueryIteratorInterface, *peer.QueryResponseMetadata, error) {
	return nil, nil, nil
}

func (m *mockStub) GetAllStatesCompositeKeyWithPagination(int32, string) (shim.StateQueryIteratorInterface, *peer.QueryResponseMetadata, error) {
	return nil, nil, nil
}
func (m *mockStub) CreateCompositeKey(string, []string) (string, error) { return "", nil }
func (m *mockStub) SplitCompositeKey(string) (string, []string, error)  { return "", nil, nil }
func (m *mockStub) GetQueryResult(string) (shim.StateQueryIteratorInterface, error) {
	return nil, nil
}

func (m *mockStub) GetQueryResultWithPagination(string, int32, string) (shim.StateQueryIteratorInterface, *peer.QueryResponseMetadata, error) {
	return nil, nil, nil
}

func (m *mockStub) GetHistoryForKey(string) (shim.HistoryQueryIteratorInterface, error) {
	return nil, nil
}
func (m *mockStub) GetPrivateData(string, string) ([]byte, error)              { return nil, nil }
func (m *mockStub) GetMultiplePrivateData(string, ...string) ([][]byte, error) { return nil, nil }
func (m *mockStub) GetPrivateDataHash(string, string) ([]byte, error)          { return nil, nil }
func (m *mockStub) PutPrivateData(string, string, []byte) error                { return nil }
func (m *mockStub) DelPrivateData(string, string) error                        { return nil }
func (m *mockStub) PurgePrivateData(string, string) error                      { return nil }
func (m *mockStub) SetPrivateDataValidationParameter(string, string, []byte) error {
	return nil
}

func (m *mockStub) GetPrivateDataValidationParameter(string, string) ([]byte, error) {
	return nil, nil
}

func (m *mockStub) GetPrivateDataByRange(string, string, string) (shim.StateQueryIteratorInterface, error) {
	return nil, nil
}

func (m *mockStub) GetPrivateDataByPartialCompositeKey(string, string, []string) (shim.StateQueryIteratorInterface, error) {
	return nil, nil
}

func (m *mockStub) GetPrivateDataQueryResult(string, string) (shim.StateQueryIteratorInterface, error) {
	return nil, nil
}
func (m *mockStub) GetCreator() ([]byte, error)                      { return nil, nil }
func (m *mockStub) GetTransient() (map[string][]byte, error)         { return nil, nil }
func (m *mockStub) GetBinding() ([]byte, error)                      { return nil, nil }
func (m *mockStub) GetDecorations() map[string][]byte                { return nil }
func (m *mockStub) GetSignedProposal() (*peer.SignedProposal, error) { return nil, nil }
func (m *mockStub) GetTxTimestamp() (*timestamppb.Timestamp, error)  { return nil, nil }
func (m *mockStub) SetEvent(name string, payload []byte) error {
	args := m.Called(name, payload)
	return args.Error(0)
}
func (m *mockStub) StartWriteBatch()        {}
func (m *mockStub) FinishWriteBatch() error { return nil }
