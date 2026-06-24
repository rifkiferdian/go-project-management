package models

// ApprovalFlow merepresentasikan master alur approval.
type ApprovalFlow struct {
	ID               int
	DivisionID       int
	DivisionName     string
	FlowCode         string
	FlowName         string
	IsActive         bool
	CreatedAtDisplay string
	UpdatedAtDisplay string
}

// ApprovalFlowOption dipakai untuk dropdown flow.
type ApprovalFlowOption struct {
	ID           int
	DivisionID   int
	DivisionName string
	FlowCode     string
	FlowName     string
	IsActive     bool
}

// ApprovalFlowStep merepresentasikan step alur approval.
type ApprovalFlowStep struct {
	ID               int
	ApprovalFlowID   int
	FlowCode         string
	FlowName         string
	StepOrder        int
	StepName         string
	ApprovalRule     string
	IsActive         bool
	CreatedAtDisplay string
	UpdatedAtDisplay string
}

// ApprovalFlowStepOption dipakai untuk dropdown step.
type ApprovalFlowStepOption struct {
	ID             int
	ApprovalFlowID int
	StepOrder      int
	StepName       string
	FlowCode       string
	FlowName       string
	IsActive       bool
}

// ApprovalFlowStepApprover merepresentasikan approver pada step tertentu.
type ApprovalFlowStepApprover struct {
	ID                 int
	ApprovalFlowStepID int
	StepOrder          int
	StepName           string
	FlowCode           string
	FlowName           string
	ApproverType       string
	ApproverUserID     int
	ApproverRoleID     int
	ApproverDivisionID int
	ApproverLabel      string
	IsActive           bool
	CreatedAtDisplay   string
	UpdatedAtDisplay   string
}

// LookupOption dipakai untuk dropdown entity user/role/division.
type LookupOption struct {
	ID   int
	Name string
}
