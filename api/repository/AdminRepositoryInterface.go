package repository

type AdminRepository interface {
	Create(email string, password string) (RepoResponse, Status)
}
