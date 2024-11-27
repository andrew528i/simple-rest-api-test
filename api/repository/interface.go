package repository

type CustomerRepository interface {
	DeleteByPrefix(prefix []string) (DeleteInfo, error)
	GetByPrefix(prefix []string) ([]CustomerInfo, error)
}
