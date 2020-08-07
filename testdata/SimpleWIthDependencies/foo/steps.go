package main

type User struct {
	ID string
}

type UserRepository interface {
	Find(string) (User, error)
}

type NotificationRepository interface {
	Send(string) error
}

func step1(repo UserRepository) func(id string) error {
	return func(id string) error {
		_, err := repo.Find(id)
		return err
	}
}

func step2(repo NotificationRepository) func(string) error {
	return func(userID string) error {
		return repo.Send(userID)
	}
}
