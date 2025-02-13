package main

type VideoRepository interface {
	FindAll() ([]Video, error)
}

type videoBusiness struct {
	repository VideoRepository
}

func NewVideoBusiness(repository VideoRepository) *videoBusiness {
	return &videoBusiness{repository: repository}
}

func (b *videoBusiness) FindAll() ([]Video, error) {
	return b.repository.FindAll()
}
