package models

// TicketListItem merepresentasikan ticket untuk halaman management tickets.
type TicketListItem struct {
	ID               int
	Code             string
	Name             string
	ProjectName      string
	StatusName       string
	StatusColor      string
	PriorityName     string
	PriorityColor    string
	TypeName         string
	TypeColor        string
	OwnerName        string
	ResponsibleName  string
	Estimation       float64
	EstimationText   string
	StartsAtDisplay  string
	EndsAtDisplay    string
	UpdatedAtDisplay string
}

// TicketDetailPage mewakili data lengkap halaman detail ticket.
type TicketDetailPage struct {
	Ticket      TicketDetail
	Comments    []TicketCommentItem
	Activities  []TicketActivityItem
	Hours       []TicketHourItem
	Subscribers []TicketSubscriberItem
	Attachments []TicketAttachmentItem
}

// TicketEditPage mewakili data halaman edit ticket.
type TicketEditPage struct {
	Form            TicketEditForm
	StatusOptions   []TicketFormOption
	PriorityOptions []TicketFormOption
	TypeOptions     []TicketFormOption
	UserOptions     []TicketUserOption
	EpicOptions     []TicketEpicOption
}

// TicketEditForm menampung data form edit ticket.
type TicketEditForm struct {
	ID            int
	ProjectID     int
	Code          string
	ProjectName   string
	Name          string
	Content       string
	StatusID      int
	PriorityID    int
	TypeID        int
	OwnerID       int
	ResponsibleID int
	EpicID        int
	Estimation    string
	StartsAt      string
	EndsAt        string
}

// TicketDetail memuat informasi utama ticket.
type TicketDetail struct {
	ID                  int
	Code                string
	Name                string
	ContentText         string
	ProjectName         string
	StatusName          string
	StatusColor         string
	PriorityName        string
	PriorityColor       string
	TypeName            string
	TypeColor           string
	OwnerName           string
	OwnerInitials       string
	ResponsibleName     string
	ResponsibleInitials string
	EpicName            string
	Estimation          float64
	EstimationText      string
	LoggedHours         float64
	LoggedHoursText     string
	LoggedPercent       int
	SubscribersCount    int
	StartsAtDisplay     string
	EndsAtDisplay       string
	CreatedAtDisplay    string
	CreatedAtRelative   string
	UpdatedAtDisplay    string
	UpdatedAtRelative   string
}

// TicketCommentItem merepresentasikan komentar pada ticket.
type TicketCommentItem struct {
	ID                int
	UserName          string
	UserInitials      string
	Content           string
	CreatedAtDisplay  string
	CreatedAtRelative string
}

// TicketActivityItem merepresentasikan riwayat perpindahan status ticket.
type TicketActivityItem struct {
	ID                int
	UserName          string
	UserInitials      string
	OldStatusName     string
	NewStatusName     string
	CreatedAtDisplay  string
	CreatedAtRelative string
}

// TicketHourItem merepresentasikan log waktu pada ticket.
type TicketHourItem struct {
	ID                int
	UserName          string
	UserInitials      string
	ActivityName      string
	Comment           string
	Value             float64
	ValueText         string
	CreatedAtDisplay  string
	CreatedAtRelative string
}

// TicketSubscriberItem merepresentasikan subscriber ticket.
type TicketSubscriberItem struct {
	ID       int
	Name     string
	Initials string
}

// TicketAttachmentItem merepresentasikan file attachment pada ticket.
type TicketAttachmentItem struct {
	ID                int
	OriginalName      string
	FileName          string
	FilePath          string
	FileSize          int64
	FileSizeText      string
	MimeType          string
	UploaderName      string
	CreatedAtDisplay  string
	CreatedAtRelative string
}

// TicketAttachmentCreateInput menampung metadata file upload ticket.
type TicketAttachmentCreateInput struct {
	TicketID     int
	UserID       int
	OriginalName string
	FileName     string
	FilePath     string
	FileSize     int64
	MimeType     string
}

// TicketFormOption dipakai untuk opsi select ticket.
type TicketFormOption struct {
	ID    int
	Name  string
	Color string
}

// TicketUserOption dipakai untuk opsi user pada form ticket.
type TicketUserOption struct {
	ID   int
	Name string
}

// TicketEpicOption dipakai untuk opsi epic pada form ticket.
type TicketEpicOption struct {
	ID   int
	Name string
}

// TicketUpdateInput menampung data update ticket dari form.
type TicketUpdateInput struct {
	ID            int
	Name          string
	Content       string
	StatusID      int
	PriorityID    int
	TypeID        int
	OwnerID       int
	ResponsibleID int
	EpicID        int
	Estimation    string
	StartsAt      string
	EndsAt        string
}

// BoardColumn merepresentasikan kolom status pada board.
type BoardColumn struct {
	ID          int
	Name        string
	Color       string
	ScopeLabel  string
	Order       int
	TicketCount int
	Tickets     []BoardTicket
}

// BoardTicket merepresentasikan kartu ticket pada board.
type BoardTicket struct {
	ID              int
	Code            string
	Name            string
	ProjectName     string
	PriorityName    string
	PriorityColor   string
	TypeName        string
	TypeColor       string
	ResponsibleName string
	EstimationText  string
	StatusID        int
}

// RoadmapEpic merepresentasikan ringkasan epic untuk halaman roadmap.
type RoadmapEpic struct {
	ID            int
	ProjectID     int
	Name          string
	ProjectName   string
	StartsAtISO   string
	EndsAtISO     string
	StartsAt      string
	EndsAt        string
	SprintCount   int
	TicketCount   int
	DoneCount     int
	Progress      int
	ProgressLabel string
}

// RoadmapSprint merepresentasikan ringkasan sprint untuk halaman roadmap.
type RoadmapSprint struct {
	ID            int
	Name          string
	ProjectName   string
	EpicID        int
	EpicName      string
	StartsAtISO   string
	EndsAtISO     string
	StartsAt      string
	EndsAt        string
	TicketCount   int
	DoneCount     int
	Progress      int
	ProgressLabel string
	StateLabel    string
}

// RoadmapTicket merepresentasikan ticket child pada roadmap.
type RoadmapTicket struct {
	ID           int
	EpicID       int
	ProjectID    int
	Name         string
	ProjectName  string
	ResourceName string
	Progress     int
	StartsAtISO  string
	EndsAtISO    string
	StartsAt     string
	EndsAt       string
}

// RoadmapEpicCreateInput menampung input pembuatan epic dari halaman roadmap.
type RoadmapEpicCreateInput struct {
	ProjectID int
	Name      string
	StartsAt  string
	EndsAt    string
}

// RoadmapTicketCreateInput menampung input pembuatan ticket dari halaman roadmap.
type RoadmapTicketCreateInput struct {
	ProjectID      int
	EpicID         *int
	Name           string
	ResourceUserID int
	Estimation     float64
	StartsAt       string
	EndsAt         string
}

// RoadmapEpicOption dipakai untuk select epic pada modal ticket.
type RoadmapEpicOption struct {
	ID          int
	Name        string
	ProjectID   int
	ProjectName string
}

// RoadmapWeek mewakili satu kolom minggu pada timeline roadmap.
type RoadmapWeek struct {
	YearLabel string
	DateLabel string
}

// RoadmapYearGroup mewakili grup header tahun pada timeline roadmap.
type RoadmapYearGroup struct {
	Label   string
	Count   int
	WidthPx int
}

// RoadmapTimelineRow mewakili satu baris gantt roadmap.
type RoadmapTimelineRow struct {
	Name           string
	Resource       string
	Progress       int
	ProgressLabel  string
	StartDateLabel string
	EndDateLabel   string
	BarLeftPx      int
	BarWidthPx     int
	BarColor       string
	BarAccentColor string
	BarProgressPct int
	ShowBar        bool
	IsChild        bool
	StyleClass     string
	ShowGroupMark  bool
	SearchText     string
	RowTone        string
}
