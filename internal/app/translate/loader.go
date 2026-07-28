package translate

import "Paarthurnax/internal/domain/translation"

type ProjectLoader interface {
	Load() (*translation.Project, error)
}
