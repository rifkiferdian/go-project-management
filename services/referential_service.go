package services

import (
	"errors"
	"gobase-app/models"
	"gobase-app/repositories"
	"strings"
)

type ReferentialService struct {
	Repo *repositories.ReferentialRepository
}

func (s *ReferentialService) GetActivities() ([]models.Activity, error) {
	return s.Repo.GetActivities()
}

func (s *ReferentialService) CreateActivity(name, description string) error {
	name = strings.TrimSpace(name)
	description = strings.TrimSpace(description)
	if name == "" || description == "" {
		return errors.New("nama dan deskripsi activity wajib diisi")
	}
	return s.Repo.CreateActivity(name, description)
}

func (s *ReferentialService) UpdateActivity(id int, name, description string) error {
	if id <= 0 {
		return errors.New("activity tidak valid")
	}
	name = strings.TrimSpace(name)
	description = strings.TrimSpace(description)
	if name == "" || description == "" {
		return errors.New("nama dan deskripsi activity wajib diisi")
	}
	return s.Repo.UpdateActivity(id, name, description)
}

func (s *ReferentialService) DeleteActivity(id int) error {
	if id <= 0 {
		return errors.New("activity tidak valid")
	}
	return s.Repo.DeleteActivity(id)
}

func (s *ReferentialService) GetProjectStatuses() ([]models.StatusReference, error) {
	return s.Repo.GetProjectStatuses()
}

func (s *ReferentialService) CreateProjectStatus(name, color string, isDefault bool) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return errors.New("nama wajib diisi")
	}
	return s.Repo.CreateProjectStatus(name, defaultColor(color), isDefault)
}

func (s *ReferentialService) UpdateProjectStatus(id int, name, color string, isDefault bool) error {
	if id <= 0 {
		return errors.New("project status tidak valid")
	}
	name = strings.TrimSpace(name)
	if name == "" {
		return errors.New("nama wajib diisi")
	}
	return s.Repo.UpdateProjectStatus(id, name, defaultColor(color), isDefault)
}

func (s *ReferentialService) DeleteProjectStatus(id int) error {
	if id <= 0 {
		return errors.New("project status tidak valid")
	}
	return s.Repo.DeleteProjectStatus(id)
}

func (s *ReferentialService) GetTicketPriorities() ([]models.StatusReference, error) {
	return s.Repo.GetTicketPriorities()
}

func (s *ReferentialService) CreateTicketPriority(name, color string, isDefault bool) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return errors.New("nama wajib diisi")
	}
	return s.Repo.CreateTicketPriority(name, defaultColor(color), isDefault)
}

func (s *ReferentialService) UpdateTicketPriority(id int, name, color string, isDefault bool) error {
	if id <= 0 {
		return errors.New("ticket priority tidak valid")
	}
	name = strings.TrimSpace(name)
	if name == "" {
		return errors.New("nama wajib diisi")
	}
	return s.Repo.UpdateTicketPriority(id, name, defaultColor(color), isDefault)
}

func (s *ReferentialService) DeleteTicketPriority(id int) error {
	if id <= 0 {
		return errors.New("ticket priority tidak valid")
	}
	return s.Repo.DeleteTicketPriority(id)
}

func (s *ReferentialService) GetTicketStatuses() ([]models.TicketStatusReference, error) {
	return s.Repo.GetTicketStatuses()
}

func (s *ReferentialService) CreateTicketStatus(name, color string, isDefault bool, order int, projectID *int) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return errors.New("nama wajib diisi")
	}
	return s.Repo.CreateTicketStatus(name, defaultColor(color), isDefault, normalizeOrder(order), projectID)
}

func (s *ReferentialService) UpdateTicketStatus(id int, name, color string, isDefault bool, order int, projectID *int) error {
	if id <= 0 {
		return errors.New("ticket status tidak valid")
	}
	name = strings.TrimSpace(name)
	if name == "" {
		return errors.New("nama wajib diisi")
	}
	return s.Repo.UpdateTicketStatus(id, name, defaultColor(color), isDefault, normalizeOrder(order), projectID)
}

func (s *ReferentialService) DeleteTicketStatus(id int) error {
	if id <= 0 {
		return errors.New("ticket status tidak valid")
	}
	return s.Repo.DeleteTicketStatus(id)
}

func (s *ReferentialService) GetTicketTypes() ([]models.TicketTypeReference, error) {
	return s.Repo.GetTicketTypes()
}

func (s *ReferentialService) CreateTicketType(name, icon, color string, isDefault bool) error {
	icon = strings.TrimSpace(icon)
	name = strings.TrimSpace(name)
	if name == "" {
		return errors.New("nama wajib diisi")
	}
	if icon == "" {
		icon = "heroicon-o-check-circle"
	}
	return s.Repo.CreateTicketType(name, icon, defaultColor(color), isDefault)
}

func (s *ReferentialService) UpdateTicketType(id int, name, icon, color string, isDefault bool) error {
	if id <= 0 {
		return errors.New("ticket type tidak valid")
	}
	icon = strings.TrimSpace(icon)
	name = strings.TrimSpace(name)
	if name == "" {
		return errors.New("nama wajib diisi")
	}
	if icon == "" {
		icon = "heroicon-o-check-circle"
	}
	return s.Repo.UpdateTicketType(id, name, icon, defaultColor(color), isDefault)
}

func (s *ReferentialService) DeleteTicketType(id int) error {
	if id <= 0 {
		return errors.New("ticket type tidak valid")
	}
	return s.Repo.DeleteTicketType(id)
}

func (s *ReferentialService) GetProjectOptions() ([]models.ProjectOption, error) {
	return s.Repo.GetProjectOptions()
}

func defaultColor(color string) string {
	color = strings.TrimSpace(color)
	if color == "" {
		return "#cecece"
	}
	return color
}

func normalizeOrder(order int) int {
	if order <= 0 {
		return 1
	}
	return order
}
