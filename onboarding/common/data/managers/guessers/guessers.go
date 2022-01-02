package guessers

type Manager struct {
	backend DBBackend
}

func NewManager(b DBBackend) *Manager {
	return &Manager{
		backend: b,
	}
}

// TODO: AddGuesser, RemoveGuesser, QueryGuesser
