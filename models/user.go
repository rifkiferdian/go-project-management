package models

// User merepresentasikan data pada tabel users.
type User struct {
	ID               int
	Name             string
	Email            string
	RoleDisplay      string
	RoleNames        []string
	CreatedAt        string
	CreatedAtDisplay string
}

// UserCreateInput menampung data yang dikirimkan dari form create user.
type UserCreateInput struct {
	Name      string
	Email     string
	Password  string
	RoleNames []string
}

// UserUpdateInput menampung data yang dikirimkan dari form edit user.
type UserUpdateInput struct {
	ID        int
	Name      string
	Email     string
	Password  string
	RoleNames []string
}
