package services

import (
	"errors"
	"fmt"
	"gobase-app/models"
	"gobase-app/repositories"
	"regexp"
	"strings"
)

var flowCodePattern = regexp.MustCompile(`^[A-Z0-9_-]+$`)

type ApprovalFlowService struct {
	Repo *repositories.ApprovalFlowRepository
}

func (s *ApprovalFlowService) GetApprovalFlows() ([]models.ApprovalFlow, error) {
	return s.Repo.GetAll()
}

func (s *ApprovalFlowService) CreateApprovalFlow(flowCode, flowName string, isActive bool) error {
	flowCode, flowName, err := validateApprovalFlowInput(flowCode, flowName)
	if err != nil {
		return err
	}

	exists, err := s.Repo.ExistsByCode(flowCode)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("flow code %s sudah digunakan", flowCode)
	}

	return s.Repo.Create(flowCode, flowName, isActive)
}

func (s *ApprovalFlowService) UpdateApprovalFlow(id int, flowCode, flowName string, isActive bool) error {
	if id <= 0 {
		return errors.New("approval flow tidak valid")
	}

	flowCode, flowName, err := validateApprovalFlowInput(flowCode, flowName)
	if err != nil {
		return err
	}

	exists, err := s.Repo.ExistsByCodeExceptID(flowCode, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("flow code %s sudah digunakan", flowCode)
	}

	return s.Repo.Update(id, flowCode, flowName, isActive)
}

func (s *ApprovalFlowService) DeleteApprovalFlow(id int) error {
	if id <= 0 {
		return errors.New("approval flow tidak valid")
	}
	return s.Repo.Delete(id)
}

func validateApprovalFlowInput(flowCode, flowName string) (string, string, error) {
	flowCode = strings.ToUpper(strings.TrimSpace(flowCode))
	flowName = strings.TrimSpace(flowName)

	if flowCode == "" {
		return "", "", errors.New("flow code wajib diisi")
	}
	if len(flowCode) > 100 {
		return "", "", errors.New("flow code maksimal 100 karakter")
	}
	if !flowCodePattern.MatchString(flowCode) {
		return "", "", errors.New("flow code hanya boleh huruf kapital, angka, underscore, dan dash")
	}
	if flowName == "" {
		return "", "", errors.New("nama flow wajib diisi")
	}
	if len(flowName) > 255 {
		return "", "", errors.New("nama flow maksimal 255 karakter")
	}

	return flowCode, flowName, nil
}
