package repositories

import (
	"database/sql"
	"gobase-app/models"
	"strings"
)

type PermissionRepository struct {
	DB *sql.DB
}

// GetGrouped mengelompokkan permission berdasarkan resource yang dipakai.
func (r *PermissionRepository) GetGrouped() ([]models.PermissionGroup, error) {
	rows, err := r.DB.Query(`
		SELECT 
			id,
			name,
			guard_name
		FROM permissions
		ORDER BY name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	groupPermissions := make(map[string][]models.Permission)
	groupOrder := []string{}

	for rows.Next() {
		var perm models.Permission
		if err := rows.Scan(&perm.ID, &perm.Name, &perm.GuardName); err != nil {
			return nil, err
		}

		perm.GroupName = inferPermissionGroup(perm.Name)
		groupKey := perm.GroupName
		if _, exists := groupPermissions[groupKey]; !exists {
			groupOrder = append(groupOrder, groupKey)
		}

		groupPermissions[groupKey] = append(groupPermissions[groupKey], perm)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	groups := make([]models.PermissionGroup, 0, len(groupPermissions))
	for _, key := range groupOrder {
		groups = append(groups, models.PermissionGroup{
			Key:         key,
			Label:       formatGroupLabel(key),
			Permissions: groupPermissions[key],
		})
	}

	return groups, nil
}

func formatGroupLabel(groupKey string) string {
	if groupKey == "" {
		return "Others"
	}

	normalized := strings.ReplaceAll(groupKey, "_", " ")
	return strings.Title(normalized)
}

func inferPermissionGroup(permissionName string) string {
	parts := strings.Fields(strings.TrimSpace(permissionName))
	if len(parts) <= 1 {
		return "others"
	}

	group := strings.ToLower(strings.Join(parts[1:], " "))
	group = strings.ReplaceAll(group, " ", "_")
	return group
}
