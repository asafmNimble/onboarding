package guessers_counters

type Manager struct {
	backend DBBackend
}

func NewManager(b DBBackend) *Manager {
	return &Manager{
		backend: b,
	}
}

func (m *Manager) CreateGuessersCounter(guesserID int64) error {
	err := m.backend.CreateGuessersCounter(guesserID)
	if err != nil {return err}
	return nil
}

func (m *Manager) IncreaseGuesserCounter(guesserID int64) error {
	err := m.backend.IncreaseGuesserCounter(guesserID)
	if err != nil {return err}
	return nil
}

func (m *Manager) GetGuesserCounter(guesserID int64) (int64, error) {
	count, err := m.backend.GetGuesserCounter(guesserID)
	if err != nil {return count, err}
	return count, nil
}
