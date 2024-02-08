package internal

import (
	//"github.com/leokuzmanovic/ai-bistro-chef/internal/models"
	"github.com/leokuzmanovic/ai-bistro-chef/internal/services"
)

func WireDependencies() {
	//models.Wire()
	services.Wire()
}
