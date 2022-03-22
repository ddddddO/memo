package repository

type HealthRepository interface {
	Check() error
}
