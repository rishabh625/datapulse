package store

import (
	"context"
	"database/sql"
	"encoding/json"

	"datapulse/internal/model"
)

type CloudSQLStore struct {
	db *sql.DB
}

func NewCloudSQLStore(db *sql.DB) *CloudSQLStore {
	return &CloudSQLStore{db: db}
}

func (c *CloudSQLStore) Exists(ctx context.Context, tenantID string, cloud model.CloudProvider, accountIdentifier string) (bool, error) {
	const q = `
SELECT EXISTS(
  SELECT 1
  FROM cloud_accounts
  WHERE tenant_id = $1 AND cloud = $2 AND account_identifier = $3
)`
	var exists bool
	if err := c.db.QueryRowContext(ctx, q, tenantID, string(cloud), accountIdentifier).Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}

func (c *CloudSQLStore) Save(ctx context.Context, account model.Account) error {
	checksJSON, err := json.Marshal(account.Checks)
	if err != nil {
		return err
	}
	metadataJSON, err := json.Marshal(account.Metadata)
	if err != nil {
		return err
	}

	const q = `
INSERT INTO cloud_accounts (
  id, tenant_id, cloud, account_identifier, trigger_mode, schedule_cron, pipelines,
  checks_json, metadata_json, created_at, updated_at
) VALUES (
  $1, $2, $3, $4, $5, $6, $7,
  $8, $9, $10, $11
)`

	_, err = c.db.ExecContext(
		ctx, q,
		account.ID,
		account.TenantID,
		string(account.Cloud),
		account.AccountIdentifier,
		string(account.TriggerMode),
		account.ScheduleCron,
		pipelineTypesToStrings(account.Pipelines),
		checksJSON,
		metadataJSON,
		account.CreatedAt,
		account.UpdatedAt,
	)
	return err
}

func (c *CloudSQLStore) Get(ctx context.Context, tenantID string, cloud model.CloudProvider, accountIdentifier string) (model.Account, error) {
	const q = `
SELECT id, tenant_id, cloud, account_identifier, trigger_mode, schedule_cron, pipelines,
       checks_json, metadata_json, created_at, updated_at
FROM cloud_accounts
WHERE tenant_id = $1 AND cloud = $2 AND account_identifier = $3`

	var (
		account                  model.Account
		cloudValue               string
		triggerMode              string
		pipelineStrings          []string
		checksJSON, metadataJSON []byte
	)

	if err := c.db.QueryRowContext(ctx, q, tenantID, string(cloud), accountIdentifier).Scan(
		&account.ID,
		&account.TenantID,
		&cloudValue,
		&account.AccountIdentifier,
		&triggerMode,
		&account.ScheduleCron,
		&pipelineStrings,
		&checksJSON,
		&metadataJSON,
		&account.CreatedAt,
		&account.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return model.Account{}, ErrNotFound
		}
		return model.Account{}, err
	}

	account.Cloud = model.CloudProvider(cloudValue)
	account.TriggerMode = model.TriggerMode(triggerMode)
	account.Pipelines = stringsToPipelineTypes(pipelineStrings)
	if err := json.Unmarshal(checksJSON, &account.Checks); err != nil {
		return model.Account{}, err
	}
	if err := json.Unmarshal(metadataJSON, &account.Metadata); err != nil {
		return model.Account{}, err
	}

	return account, nil
}

func pipelineTypesToStrings(pipelines []model.PipelineType) []string {
	out := make([]string, 0, len(pipelines))
	for _, p := range pipelines {
		out = append(out, string(p))
	}
	return out
}

func stringsToPipelineTypes(values []string) []model.PipelineType {
	out := make([]model.PipelineType, 0, len(values))
	for _, v := range values {
		out = append(out, model.PipelineType(v))
	}
	return out
}
