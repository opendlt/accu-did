package acc

import "github.com/opendlt/accu-did/registrar-go/internal/ops"

type MockClient struct {
	CreateIdentityFn    func(adiLabel string, keyPageURL string) (string, error)
	CreateDataAccountFn func(adiURL, dataAccountLabel string) (string, error)
	WriteDataEntryFn    func(dataAccountURL string, data []byte) (string, error)
	SubmitWriteDataFn   func(dataAccountURL string, payload *ops.Envelope) (string, error)
	UpdateKeyPageFn     func(keyPageURL string, operations []KeyPageOperation) (string, error)
	GetKeyPageStateFn   func(keyPageURL string) (*KeyPageState, error)

	// Recorded values for test inspection
	LastWriteData  []byte
	LastEnvelope   *ops.Envelope
	LastAccountURL string
}

var _ Submitter = (*MockClient)(nil)

func (m *MockClient) CreateIdentity(adiLabel string, keyPageURL string) (string, error) {
	if m.CreateIdentityFn != nil {
		return m.CreateIdentityFn(adiLabel, keyPageURL)
	}
	return "txid-create-identity-mock", nil
}

func (m *MockClient) CreateDataAccount(adiURL, dataAccountLabel string) (string, error) {
	if m.CreateDataAccountFn != nil {
		return m.CreateDataAccountFn(adiURL, dataAccountLabel)
	}
	return "txid-create-data-account-mock", nil
}

func (m *MockClient) WriteDataEntry(dataAccountURL string, data []byte) (string, error) {
	// record for tests
	m.LastAccountURL = dataAccountURL
	m.LastWriteData = data

	if m.WriteDataEntryFn != nil {
		return m.WriteDataEntryFn(dataAccountURL, data)
	}
	return "txid-write-data-mock", nil
}

func (m *MockClient) SubmitWriteData(dataAccountURL string, payload *ops.Envelope) (string, error) {
	// record for tests
	m.LastAccountURL = dataAccountURL
	m.LastEnvelope = payload

	if m.SubmitWriteDataFn != nil {
		return m.SubmitWriteDataFn(dataAccountURL, payload)
	}
	return "txid-submit-write-mock", nil
}

func (m *MockClient) UpdateKeyPage(keyPageURL string, operations []KeyPageOperation) (string, error) {
	if m.UpdateKeyPageFn != nil {
		return m.UpdateKeyPageFn(keyPageURL, operations)
	}
	return "txid-update-key-page-mock", nil
}

func (m *MockClient) GetKeyPageState(keyPageURL string) (*KeyPageState, error) {
	if m.GetKeyPageStateFn != nil {
		return m.GetKeyPageStateFn(keyPageURL)
	}
	return &KeyPageState{
		URL:       keyPageURL,
		Threshold: 1,
		Keys:      []KeyInfo{{PublicKey: "mock-key", KeyType: "ed25519"}},
		Height:    1,
	}, nil
}

func NewMockClient() *MockClient {
	return &MockClient{}
}
