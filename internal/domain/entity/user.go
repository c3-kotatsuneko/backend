package entity

type User struct {
	ID       string
	Name     string
	Password string
}

func (u *User) DeepCopy() *User {
	return &User{
		ID:       u.ID,
		Name:     u.Name,
		Password: u.Password,
	}
}
