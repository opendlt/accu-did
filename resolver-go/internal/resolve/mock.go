package resolve

import (
	"time"

	"github.com/opendlt/accu-did/resolver-go/internal/acc"
	"gitlab.com/accumulatenetwork/accumulate/pkg/url"
)

// MockClient implements acc.Client and records inputs/outputs for tests.
type MockClient struct {
	// Optional function hooks to override behavior
	GetLatestDIDEntryFn   func(adi string) (acc.Envelope, error)
	GetEntryAtTimeFn      func(adi string, t time.Time) (acc.Envelope, error)
	GetKeyPageStateFn     func(url string) (acc.KeyPageState, error)
	GetDataAccountEntryFn func(dataAccountURL *url.URL) ([]byte, error)

	// Recorded values for assertions in tests
	LastADI            string
	LastAtTime         time.Time
	LastKeyPageURL     string
	LastDataAccountURL *url.URL

	LastEnvelope acc.Envelope
	LastBytes    []byte

	// Simple call counters (handy in tests)
	CallsGetLatestDIDEntry   int
	CallsGetEntryAtTime      int
	CallsGetKeyPageState     int
	CallsGetDataAccountEntry int
}

var _ acc.Client = (*MockClient)(nil)

func (m *MockClient) GetLatestDIDEntry(adi string) (acc.Envelope, error) {
	m.CallsGetLatestDIDEntry++
	m.LastADI = adi

	if m.GetLatestDIDEntryFn != nil {
		env, err := m.GetLatestDIDEntryFn(adi)
		m.LastEnvelope = env
		return env, err
	}
	// default: empty envelope
	m.LastEnvelope = acc.Envelope{}
	return acc.Envelope{}, nil
}

func (m *MockClient) GetEntryAtTime(adi string, t time.Time) (acc.Envelope, error) {
	m.CallsGetEntryAtTime++
	m.LastADI = adi
	m.LastAtTime = t

	if m.GetEntryAtTimeFn != nil {
		env, err := m.GetEntryAtTimeFn(adi, t)
		m.LastEnvelope = env
		return env, err
	}
	// default: empty envelope
	m.LastEnvelope = acc.Envelope{}
	return acc.Envelope{}, nil
}

func (m *MockClient) GetKeyPageState(u string) (acc.KeyPageState, error) {
	m.CallsGetKeyPageState++
	m.LastKeyPageURL = u

	if m.GetKeyPageStateFn != nil {
		return m.GetKeyPageStateFn(u)
	}
	// default: zero value
	return acc.KeyPageState{}, nil
}

func (m *MockClient) GetDataAccountEntry(dataAccountURL *url.URL) ([]byte, error) {
	m.CallsGetDataAccountEntry++
	m.LastDataAccountURL = dataAccountURL

	if m.GetDataAccountEntryFn != nil {
		b, err := m.GetDataAccountEntryFn(dataAccountURL)
		m.LastBytes = b
		return b, err
	}

	// Defaults for convenient tests:
	// If authority == "deactivated", return a deactivated DID; otherwise a basic DID doc
	if dataAccountURL != nil && dataAccountURL.Authority == "deactivated" {
		doc := []byte(`{"@context":["https://www.w3.org/ns/did/v1"],"id":"did:acc:deactivated","deactivated":true}`)
		m.LastBytes = doc
		return doc, nil
	}

	doc := []byte(`{"@context":["https://www.w3.org/ns/did/v1"],"id":"did:acc:mock"}`)
	m.LastBytes = doc
	return doc, nil
}

// Constructors

func NewMockClient() *MockClient {
	return &MockClient{}
}

func NewMockDeactivatedClient() *MockClient {
	return &MockClient{
		GetDataAccountEntryFn: func(dataAccountURL *url.URL) ([]byte, error) {
			doc := []byte(`{"@context":["https://www.w3.org/ns/did/v1"],"id":"did:acc:deactivated","deactivated":true}`)
			return doc, nil
		},
	}
}
