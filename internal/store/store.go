package store

import (
	"context"

	"datapulse/internal/model"
)

type AccountStore interface {
	Exists(ctx context.Context, tenantID string, cloud model.CloudProvider, accountIdentifier string) (bool, error)
	Save(ctx context.Context, account model.Account) error
	Get(ctx context.Context, tenantID string, cloud model.CloudProvider, accountIdentifier string) (model.Account, error)
}
