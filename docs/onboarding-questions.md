# Onboarding Question Flow

The orchestrator should ask these in sequence:

1. Which account do you want to onboard (`aws` or `gcp`)?
2. Tenant ID and account identifier (AWS account ID or GCP project/billing identifier)?
3. Trigger mode (`completion`, `failure`, `schedule`)?
4. Which pipelines to monitor (`glue`, `airflow`, `dbt`)?

For AWS:
- CUR S3 bucket, prefix, and file format (`parquet` or `csv`)
- Read-only IAM role ARN
- Region

For GCP:
- Project ID
- Billing dataset/table
- Service account impersonation target
- Region

Validation expectations:
- Duplicate account prevention by `(tenant_id, cloud, account_identifier)`
- Data source presence checks
- Pipeline connectivity checks
- Trigger configuration checks
