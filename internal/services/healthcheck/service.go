package healthcheck

type repository interface {
}

type service struct {
	repository repository
}

func NewService(repository repository) *service {
	return &service{
		repository: repository,
	}
}
