package web

import "github.com/Shalqarov/forum/pkg/models"

type templateData struct {
	User  *models.User
	Users []*models.User
}
