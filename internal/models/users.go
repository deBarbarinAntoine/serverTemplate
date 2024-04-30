package models

type UserModel struct {
	id int
}

func (m *UserModel) Exists(id int) (bool, error) {
	var exists bool

	// TODO checking whether the user exists or not
	if id == 1 {
		exists = true
	}

	return exists, nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {

	// TODO authentication with email and password
	if email == "thorgdar@gmail.com" {
		return 1, nil
	}

	return 0, ErrInvalidCredentials

}
