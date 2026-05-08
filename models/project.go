package models

// Project merepresentasikan ringkasan project untuk tampilan list.
type Project struct {
	ID                 int
	Name               string
	Description        string
	OwnerID            int
	OwnerName          string
	RequestDivision    string
	RequestDivisionIDs []int
	StatusID           int
	StatusName         string
	StatusColor        string
	PriorityID         int
	PriorityName       string
	PriorityColor      string
	TicketPrefix       string
	StatusType         string
	Type               string
	MemberCount        int
	TicketCount        int
	CreatedAt          string
	CreatedAtDisplay   string
}

// ProjectCreateInput menampung data form pembuatan project.
type ProjectCreateInput struct {
	Name         string
	Description  string
	OwnerID      int
	DivisionIDs  []int64
	StatusID     int
	PriorityID   int
	TicketPrefix string
	StatusType   string
	Type         string
}

// ProjectUpdateInput menampung data form perubahan project.
type ProjectUpdateInput struct {
	ID           int
	Name         string
	Description  string
	OwnerID      int
	DivisionIDs  []int64
	StatusID     int
	PriorityID   int
	TicketPrefix string
	StatusType   string
	Type         string
}

// ProjectStatusOption dipakai untuk opsi select status.
type ProjectStatusOption struct {
	ID    int
	Name  string
	Color string
}

// ProjectPriorityOption dipakai untuk opsi select prioritas.
type ProjectPriorityOption struct {
	ID    int
	Name  string
	Color string
}

// ProjectStatusChartItem merepresentasikan komposisi project per status.
type ProjectStatusChartItem struct {
	Name    string
	Color   string
	Count   int
	Percent int
}

// ProjectDivisionChartItem merepresentasikan jumlah project per divisi peminta.
type ProjectDivisionChartItem struct {
	Name         string
	Count        int
	WidthPercent int
}

// DashboardProjectListItem merepresentasikan item list sederhana di dashboard.
type DashboardProjectListItem struct {
	ID                      int
	Name                    string
	RequestDivision         string
	StatusName              string
	StatusColor             string
	HighPriorityTicketCount int
}
