CREATE TABLE IF NOT EXISTS cloud_accounts (
  id TEXT PRIMARY KEY,
  tenant_id TEXT NOT NULL,
  cloud TEXT NOT NULL,
  account_identifier TEXT NOT NULL,
  trigger_mode TEXT NOT NULL,
  schedule_cron TEXT NOT NULL DEFAULT '',
  pipelines TEXT[] NOT NULL,
  checks_json JSONB NOT NULL,
  metadata_json JSONB NOT NULL,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL,
  UNIQUE (tenant_id, cloud, account_identifier)
);
