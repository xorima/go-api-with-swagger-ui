package userservice

import "errors"

type User struct {
	ID   int
	Name string
}

type UserService struct {
	users []User
}

func NewUserService() *UserService {
	return &UserService{
		users: []User{},
	}
}

func (us *UserService) CreateUser(name string) {
	us.users = append(us.users, User{
		ID:   len(us.users) + 1,
		Name: name,
	})
	// Create a new user
}

func (us *UserService) GetUser(id int) (User, error) {
	// Get a user
	if len(us.users) < id {
		return User{}, errors.New("user not found")
	}
	return us.users[id], nil
}

func (us *UserService) UpdateUser(id int, name string) error {
	// Update a user
	if len(us.users) < id {
		return errors.New("user not found")
	}
	us.users[id].Name = name
	return nil
}

func (us *UserService) DeleteUser(id int) error {
	// Delete a user
	if len(us.users) < id {
		return errors.New("user not found")
	}

	us.users = append(us.users[:id], us.users[id+1:]...)
	return nil
}
