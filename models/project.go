package models

// Project merepresentasikan ringkasan project untuk tampilan list.
type Project struct {
	ID               int
	Name             string
	Description      string
	OwnerID          int
	OwnerName        string
	StatusID         int
	StatusName       string
	StatusColor      string
	TicketPrefix     string
	StatusType       string
	Type             string
	MemberCount      int
	TicketCount      int
	CreatedAt        string
	CreatedAtDisplay string
}

// ProjectCreateInput menampung data form pembuatan project.
type ProjectCreateInput struct {
	Name         string
	Description  string
	OwnerID      int
	StatusID     int
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
	StatusID     int
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
