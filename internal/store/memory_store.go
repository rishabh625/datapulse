package store

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"datapulse/internal/model"
)

var ErrNotFound = errors.New("account not found")

type MemoryStore struct {
	mu       sync.RWMutex
	accounts map[string]model.Account
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		accounts: make(map[string]model.Account),
	}
}

func buildKey(tenantID string, cloud model.CloudProvider, accountIdentifier string) string {
	return fmt.Sprintf("%s:%s:%s", tenantID, cloud, accountIdentifier)
}

func (m *MemoryStore) Exists(ctx context.Context, tenantID string, cloud model.CloudProvider, accountIdentifier string) (bool, error) {
	_ = ctx
	m.mu.RLock()
	defer m.mu.RUnlock()
	_, ok := m.accounts[buildKey(tenantID, cloud, accountIdentifier)]
	return ok, nil
}

func (m *MemoryStore) Save(ctx context.Context, account model.Account) error {
	_ = ctx
	m.mu.Lock()
	defer m.mu.Unlock()
	m.accounts[buildKey(account.TenantID, account.Cloud, account.AccountIdentifier)] = account
	return nil
}

func (m *MemoryStore) Get(ctx context.Context, tenantID string, cloud model.CloudProvider, accountIdentifier string) (model.Account, error) {
	_ = ctx
	m.mu.RLock()
	defer m.mu.RUnlock()
	account, ok := m.accounts[buildKey(tenantID, cloud, accountIdentifier)]
	if !ok {
		return model.Account{}, ErrNotFound
	}
	return account, nil
}
