package service

import "github.com/energimind/identity-service/core/domain"

type mockIDGenerator struct{}

func newMockIDGenerator() *mockIDGenerator {
	return &mockIDGenerator{}
}

// ensure mockIDGenerator implements domain.IDGenerator.
var _ domain.IDGenerator = (*mockIDGenerator)(nil)

func (m mockIDGenerator) GenerateID() string {
	return "1"
}
