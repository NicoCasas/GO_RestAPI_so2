package initializers

import "github.com/ICOMP-UNC/2023---soii---laboratorio-6-NicoCasas/users_service/model/modelService"

func CheckSSHGroup() {
	modelService.CreateSSHClientGroupIfNotExists()
}
