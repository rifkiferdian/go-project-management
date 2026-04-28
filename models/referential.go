package models

// Activity mewakili data master activity.
type Activity struct {
	ID               int
	Name             string
	Description      string
	CreatedAtDisplay string
}

// StatusReference dipakai untuk project status dan ticket priority.
type StatusReference struct {
	ID               int
	Name             string
	Color            string
	IsDefault        bool
	CreatedAtDisplay string
}

// TicketStatusReference mewakili ticket status beserta urutan dan project opsional.
type TicketStatusReference struct {
	ID               int
	Name             string
	Color            string
	IsDefault        bool
	Order            int
	ProjectID        int
	ProjectName      string
	CreatedAtDisplay string
}

// TicketTypeReference mewakili master ticket type.
type TicketTypeReference struct {
	ID               int
	Name             string
	Icon             string
	Color            string
	IsDefault        bool
	CreatedAtDisplay string
}

// ProjectOption dipakai untuk select project referential.
type ProjectOption struct {
	ID   int
	Name string
}
