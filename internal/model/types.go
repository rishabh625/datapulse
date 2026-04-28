package model

import "time"

type CloudProvider string

const (
	CloudAWS CloudProvider = "aws"
	CloudGCP CloudProvider = "gcp"
)

type TriggerMode string

const (
	TriggerCompletion TriggerMode = "completion"
	TriggerFailure    TriggerMode = "failure"
	TriggerSchedule   TriggerMode = "schedule"
)

type PipelineType string

const (
	PipelineGlue    PipelineType = "glue"
	PipelineAirflow PipelineType = "airflow"
	PipelineDBT     PipelineType = "dbt"
)

type Account struct {
	ID                string
	TenantID          string
	Cloud             CloudProvider
	AccountIdentifier string
	TriggerMode       TriggerMode
	ScheduleCron      string
	Pipelines         []PipelineType
	Checks            []CheckResult
	Metadata          map[string]string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type AWSConfig struct {
	CURBucket     string `json:"cur_bucket"`
	CURPrefix     string `json:"cur_prefix"`
	CURFormat     string `json:"cur_format"`
	ReadOnlyIAM   string `json:"read_only_iam_role_arn"`
	Region        string `json:"region"`
	CloudWatchLog string `json:"cloudwatch_log_group,omitempty"`
}

type GCPConfig struct {
	ProjectID      string `json:"project_id"`
	BillingDataset string `json:"billing_dataset"`
	BillingTable   string `json:"billing_table"`
	ImpersonateSA  string `json:"impersonate_service_account"`
	Region         string `json:"region"`
}

type OnboardingRequest struct {
	TenantID          string         `json:"tenant_id"`
	Cloud             CloudProvider  `json:"cloud"`
	AccountIdentifier string         `json:"account_identifier"`
	TriggerMode       TriggerMode    `json:"trigger_mode"`
	ScheduleCron      string         `json:"schedule_cron,omitempty"`
	Pipelines         []PipelineType `json:"pipelines"`
	AWS               *AWSConfig     `json:"aws,omitempty"`
	GCP               *GCPConfig     `json:"gcp,omitempty"`
}

type CheckResult struct {
	Name    string `json:"name"`
	OK      bool   `json:"ok"`
	Details string `json:"details"`
}

type OnboardingResponse struct {
	AccountID string        `json:"account_id"`
	Checks    []CheckResult `json:"checks"`
	Warnings  []string      `json:"warnings,omitempty"`
}

type TriggerEvaluationRequest struct {
	TenantID          string `json:"tenant_id"`
	AccountIdentifier string `json:"account_identifier"`
	PipelineID        string `json:"pipeline_id"`
	EventStatus       string `json:"event_status"`
}

type TriggerEvaluationResponse struct {
	PipelineSummary map[string]any `json:"pipeline_summary"`
	CostSummary     map[string]any `json:"cost_summary"`
	Explanation     string         `json:"explanation"`
}
