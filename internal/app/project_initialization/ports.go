package project_initialization

import "Paarthurnax/internal/domain/translation"

type ProjectLoader interface {
	Load() (*translation.Project, error)
}

type Selector interface {
	Select(question string, choices []string) (string, error)
}
