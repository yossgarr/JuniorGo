package service

import (
	"errors"
	"backend-go/internal/models"
	"backend-go/internal/repo"
	"backend-go/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo *repo.UserRepository
}

func NewUserService(r *repo.UserRepository) *UserService {
	return &UserService{repo: r}
}

func (s *UserService) Register(u *models.User) error {
	if u.Username == "" || u.Password == "" {
		return errors.New("username dan password wajib diisi")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)

	return s.repo.Simpan(u)
}

func (s *UserService) Login(username, password string) (*models.AuthResponse, error) {
	u, err := s.repo.CariBerdasarkanUsername(username)
	if err != nil {
		return nil, errors.New("username atau password keliru")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return nil, errors.New("username atau password keliru")
	}

	token, err := utils.GenerateToken(u.ID, u.Username)
	if err != nil {
		return nil, errors.New("gagal memproses pembuatan token")
	}

	return &models.AuthResponse{
		Token:    token,
		Username: u.Username,
	}, nil
}

func (s *UserService) DaftarAtauAmbilUserOAuth(email string) (*models.User, error) {
	return s.repo.SimpanUserOAuth(email)
}