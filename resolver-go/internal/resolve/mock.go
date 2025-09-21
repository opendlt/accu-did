package resolve

import (
	"time"

	"github.com/opendlt/accu-did/resolver-go/internal/acc"
	"gitlab.com/accumulatenetwork/accumulate/pkg/url"
)

type mockClient struct {
	GetLatestDIDEntryFn   func(adi string) (acc.Envelope, error)
	GetEntryAtTimeFn      func(adi string, t time.Time) (acc.Envelope, error)
	GetKeyPageStateFn     func(url string) (acc.KeyPageState, error)
	GetDataAccountEntryFn func(dataAccountURL *url.URL) ([]byte, error)
}

var _ acc.Client = (*mockClient)(nil)

func (m *mockClient) GetLatestDIDEntry(adi string) (acc.Envelope, error) {
	if m.GetLatestDIDEntryFn != nil {
		return m.GetLatestDIDEntryFn(adi)
	}
	return acc.Envelope{}, nil
}

func (m *mockClient) GetEntryAtTime(adi string, t time.Time) (acc.Envelope, error) {
	if m.GetEntryAtTimeFn != nil {
		return m.GetEntryAtTimeFn(adi, t)
	}
	return acc.Envelope{}, nil
}

func (m *mockClient) GetKeyPageState(url string) (acc.KeyPageState, error) {
	if m.GetKeyPageStateFn != nil {
		return m.GetKeyPageStateFn(url)
	}
	return acc.KeyPageState{}, nil
}

func (m *mockClient) GetDataAccountEntry(dataAccountURL *url.URL) ([]byte, error) {
	if m.GetDataAccountEntryFn != nil {
		return m.GetDataAccountEntryFn(dataAccountURL)
	}
	// Return mock deactivated DID
	if dataAccountURL.Authority == "deactivated" {
		return []byte(`{"@context":["https://www.w3.org/ns/did/v1"],"id":"did:acc:deactivated","deactivated":true}`), nil
	}
	// Return mock DID document
	return []byte(`{"@context":["https://www.w3.org/ns/did/v1"],"id":"did:acc:mock"}`), nil
}

func NewMockClient() *mockClient {
	return &mockClient{}
}

func NewMockDeactivatedClient() *mockClient {
	return &mockClient{
		GetDataAccountEntryFn: func(dataAccountURL *url.URL) ([]byte, error) {
			return []byte(`{"@context":["https://www.w3.org/ns/did/v1"],"id":"did:acc:deactivated","deactivated":true}`), nil
		},
	}
}
