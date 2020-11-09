package user

import (
	"errors"
	"github.com/manveru/faker"
	"testing"
)

func TestUser(t *testing.T) {
	//Get user repository
	repo, err := GetRepository()
	if err != nil {
		t.Error(err)
	}

	fake, err := faker.New("en")
	email := fake.SafeEmail()
	defer deleteUser(email, t, repo)
	if err != nil {
		t.Error(err)
	}

	//Create user with random email and password
	var user User
	user.Email = email
	user.Password = "testpassword"
	err = repo.Create(&user)
	if err != nil {
		t.Error(err)
	}

	//Try to login
	//reset password with non hashed one
	user.Password = "testpassword"
	token, err := repo.Login(&user)
	if err != nil {
		t.Error(err)
	}
	if len(token) == 0 {
		t.Error(errors.New("Token is empty"))
	}

	//Search by secret and token
	result, err := repo.FindUserByKeyAndSecret(user.AccessKey, user.Secret)
	if err != nil {
		t.Error(err)
	}
	if result.Email != user.Email {
		t.Error(errors.New("This is not the same user"))
	}
}

func deleteUser(email string, t *testing.T, repo UserRepository) {
	//Remove User and test if this is ok
	err := repo.Delete(email)
	if err != nil {
		t.Error(err)
	}
}
