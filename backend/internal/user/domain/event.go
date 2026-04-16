package domain

// UserRegisteredEvent is emitted when a new user registers.
type UserRegisteredEvent struct {
	UserID uint64
	Email  string
}

func (UserRegisteredEvent) Topic() string {
	return "user.registered"
}
