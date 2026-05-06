package repository // Определяет интерфейс хранилища заявок.
import "permission-service/internal/models"

type PermissionRepository interface {
	Create(p *models.Permission) error
	GetRole(userID, taskID string) (string, error)
}
