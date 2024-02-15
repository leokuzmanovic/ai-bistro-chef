package internal

import (
	"github.com/leokuzmanovic/ai-bistro-chef/internal/services"
)

func WireDependencies() {
	services.Wire()
}
