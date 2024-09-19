package service

import "github.com/username/myapp/repository"

type UserService interface {
    DeleteUser(id int) error
    // Metode lain seperti GetUser, CreateUser, dsb.
}

type userServiceImpl struct {
    repo repository.UserRepository
}

// NewUserService adalah konstruktor untuk UserService
func NewUserService(repo repository.UserRepository) UserService {
    return &userServiceImpl{
        repo: repo,
    }
}

// DeleteUser menghapus pengguna berdasarkan ID dengan logika bisnis tambahan jika diperlukan
func (s *userServiceImpl) DeleteUser(id int) error {
    return s.repo.DeleteUser(id)
}
