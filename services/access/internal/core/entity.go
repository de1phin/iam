package core

type AccessBinding struct {
	UserID   string
	RoleName string
	Resource string
}

type Role struct {
	Name        string
	Permissions []string
}

func (r *Role) HavePermission(expectedPermission string) bool {
	for _, permission := range r.Permissions {
		if expectedPermission == permission {
			return true
		}
	}

	return false
}
