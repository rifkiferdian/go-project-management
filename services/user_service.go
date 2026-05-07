package services

import (
	"database/sql"
	"errors"
	"fmt"
	"gobase-app/config"
	"gobase-app/models"
	"gobase-app/repositories"
	"net/mail"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	Repo *repositories.UserRepository
}

func (s *UserService) GetUsers() ([]models.User, error) {
	return s.Repo.GetAll()
}

func (s *UserService) GetDivisions() ([]models.DivisionOption, error) {
	return s.Repo.GetDivisions()
}

// CreateUser memproses data dari form, validasi dasar, hashing password,
// lalu menyimpan user beserta role yang dipilih.
func (s *UserService) CreateUser(input models.UserCreateInput) error {
	name := strings.TrimSpace(input.Name)
	email := strings.TrimSpace(input.Email)

	if name == "" || input.Password == "" {
		return errors.New("nama dan password wajib diisi")
	}
	if email == "" {
		return errors.New("email wajib diisi")
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return errors.New("email tidak valid")
	}
	divisionIDs := uniqueInt64s(input.DivisionIDs)
	if len(divisionIDs) == 0 {
		return errors.New("minimal pilih 1 divisi")
	}

	exists, err := s.Repo.ExistsByEmail(email)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("email %s sudah digunakan", email)
	}

	existingDivisionMap, err := s.Repo.FindExistingDivisionIDs(divisionIDs)
	if err != nil {
		return err
	}
	var missingDivisions []string
	for _, divisionID := range divisionIDs {
		if !existingDivisionMap[divisionID] {
			missingDivisions = append(missingDivisions, fmt.Sprintf("%d", divisionID))
		}
	}
	if len(missingDivisions) > 0 {
		return fmt.Errorf("divisi tidak ditemukan: %s", strings.Join(missingDivisions, ", "))
	}

	roleNames := uniqueStrings(input.RoleNames)
	roleMap, err := s.Repo.GetRoleIDsByNames(roleNames)
	if err != nil {
		return err
	}

	var (
		roleIDs      []int64
		missingRoles []string
	)
	for _, roleName := range roleNames {
		if id, ok := roleMap[roleName]; ok {
			roleIDs = append(roleIDs, id)
		} else {
			missingRoles = append(missingRoles, roleName)
		}
	}
	if len(missingRoles) > 0 {
		return fmt.Errorf("role tidak ditemukan: %s", strings.Join(missingRoles, ", "))
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = s.Repo.CreateUserWithRoles(repositories.UserCreateParams{
		HashedPassword: string(hashedPassword),
		Name:           name,
		Email:          email,
		DivisionIDs:    divisionIDs,
	}, roleIDs)

	return err
}

// UpdateUser memperbarui data user yang sudah ada.
func (s *UserService) UpdateUser(input models.UserUpdateInput) error {
	name := strings.TrimSpace(input.Name)
	email := strings.TrimSpace(input.Email)

	if input.ID <= 0 {
		return errors.New("user tidak valid")
	}
	if name == "" {
		return errors.New("nama wajib diisi")
	}
	if email == "" {
		return errors.New("email wajib diisi")
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return errors.New("email tidak valid")
	}
	divisionIDs := uniqueInt64s(input.DivisionIDs)
	if len(divisionIDs) == 0 {
		return errors.New("minimal pilih 1 divisi")
	}

	exists, err := s.Repo.ExistsByEmailExceptID(email, input.ID)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("email %s sudah digunakan", email)
	}

	existingDivisionMap, err := s.Repo.FindExistingDivisionIDs(divisionIDs)
	if err != nil {
		return err
	}
	var missingDivisions []string
	for _, divisionID := range divisionIDs {
		if !existingDivisionMap[divisionID] {
			missingDivisions = append(missingDivisions, fmt.Sprintf("%d", divisionID))
		}
	}
	if len(missingDivisions) > 0 {
		return fmt.Errorf("divisi tidak ditemukan: %s", strings.Join(missingDivisions, ", "))
	}

	roleNames := uniqueStrings(input.RoleNames)
	roleMap, err := s.Repo.GetRoleIDsByNames(roleNames)
	if err != nil {
		return err
	}

	var (
		roleIDs      []int64
		missingRoles []string
	)
	for _, roleName := range roleNames {
		if id, ok := roleMap[roleName]; ok {
			roleIDs = append(roleIDs, id)
		} else {
			missingRoles = append(missingRoles, roleName)
		}
	}
	if len(missingRoles) > 0 {
		return fmt.Errorf("role tidak ditemukan: %s", strings.Join(missingRoles, ", "))
	}

	var hashedPassword string
	if strings.TrimSpace(input.Password) != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		hashedPassword = string(hashed)
	}

	return s.Repo.UpdateUserWithRoles(repositories.UserUpdateParams{
		ID:             input.ID,
		HashedPassword: hashedPassword,
		Name:           name,
		Email:          email,
		DivisionIDs:    divisionIDs,
	}, roleIDs)
}

// DeleteUser removes user data by ID.
func (s *UserService) DeleteUser(id int) error {
	if id <= 0 {
		return errors.New("user id tidak valid")
	}
	return s.Repo.DeleteUser(id)
}

func UserHasPermission(userID int, perm string) (bool, error) {
	var dummy int
	queryRole := `
		SELECT 1
		FROM model_has_roles mhr
		JOIN role_has_permissions rhp ON rhp.role_id = mhr.role_id
		JOIN permissions p ON p.id = rhp.permission_id
		WHERE mhr.model_id = ? AND mhr.model_type = ? AND p.name = ?
		LIMIT 1
	`
	err := config.DB.QueryRow(queryRole, userID, repositoriesUserModelType(), perm).Scan(&dummy)
	if err == nil {
		return true, nil
	}
	if err != sql.ErrNoRows {
		return false, err
	}

	queryDirect := `
		SELECT 1
		FROM model_has_permissions mhp
		JOIN permissions p ON p.id = mhp.permission_id
		WHERE mhp.model_id = ? AND mhp.model_type = ? AND p.name = ?
		LIMIT 1
	`

	err = config.DB.QueryRow(queryDirect, userID, repositoriesUserModelType(), perm).Scan(&dummy)
	if err == nil {
		return true, nil
	}
	if err == sql.ErrNoRows {
		return false, nil
	}

	return false, err
}

func GetUserPermissions(userID int) (map[string]bool, error) {
	perms := make(map[string]bool)

	rows, err := config.DB.Query(`
		SELECT DISTINCT p.name
		FROM permissions p
		JOIN role_has_permissions rhp ON rhp.permission_id = p.id
		JOIN model_has_roles mhr ON mhr.role_id = rhp.role_id
		WHERE mhr.model_id = ? AND mhr.model_type = ?

		UNION

		SELECT DISTINCT p2.name
		FROM permissions p2
		JOIN model_has_permissions mhp ON mhp.permission_id = p2.id
		WHERE mhp.model_id = ? AND mhp.model_type = ?
	`, userID, repositoriesUserModelType(), userID, repositoriesUserModelType())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		perms[name] = true
	}

	return perms, nil
}

func repositoriesUserModelType() string {
	return "App\\Models\\User"
}

func uniqueStrings(values []string) []string {
	seen := make(map[string]bool)
	var result []string
	for _, v := range values {
		v = strings.TrimSpace(v)
		if v == "" || seen[v] {
			continue
		}
		seen[v] = true
		result = append(result, v)
	}
	return result
}

func uniqueInt64s(values []int64) []int64 {
	seen := make(map[int64]bool)
	var result []int64
	for _, v := range values {
		if v <= 0 || seen[v] {
			continue
		}
		seen[v] = true
		result = append(result, v)
	}
	return result
}
