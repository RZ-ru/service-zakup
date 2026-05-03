package repository // Определяет интерфейс хранилища заявок.
import "permission-service/internal/models"

type PermissionRepository interface {
	Create(p *models.Permission) error
	Exists(userID, taskID string) (bool, error)
}
