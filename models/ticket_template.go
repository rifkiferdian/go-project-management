package models

// TicketTemplateSet merepresentasikan grup template ticket (contoh: new_project).
type TicketTemplateSet struct {
	ID               int
	Name             string
	Purpose          string
	Description      string
	IsActive         bool
	ItemCount        int
	EpicCount        int
	CreatedAtDisplay string
}

// TicketTemplateEpic merepresentasikan template epic dalam satu set.
type TicketTemplateEpic struct {
	ID               int
	SetID            int
	SetName          string
	SetPurpose       string
	Name             string
	Description      string
	StartOffsetDays  int
	DueOffsetDays    int
	SortOrder        int
	IsActive         bool
	CreatedAtDisplay string
}

// TicketTemplateItem merepresentasikan detail ticket default di dalam satu set.
type TicketTemplateItem struct {
	ID                     int
	SetID                  int
	SetName                string
	SetPurpose             string
	Title                  string
	Description            string
	TemplateEpicID         int
	TemplateEpicName       string
	DefaultTypeID          int
	DefaultTypeName        string
	DefaultPriorityID      int
	DefaultPriorityName    string
	DefaultStatusID        int
	DefaultStatusName      string
	DefaultOwnerID         int
	DefaultOwnerName       string
	DefaultResponsibleID   int
	DefaultResponsibleName string
	Estimation             float64
	EstimationText         string
	StartOffsetDays        int
	DueOffsetDays          int
	SortOrder              int
	IsActive               bool
	CreatedAtDisplay       string
}

// TicketTemplateOption dipakai untuk opsi select pada form template.
type TicketTemplateOption struct {
	ID   int
	Name string
}
