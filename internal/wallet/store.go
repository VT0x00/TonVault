package wallet

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/VT0x00/tonvault/internal/models"
)

type Store struct {
	mu       sync.RWMutex
	dir      string
	wallets  []models.Wallet
}

func NewStore() (*Store, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, errors.New("cannot determine home directory")
	}

	dir := filepath.Join(home, ".config", "tonvault", "wallets")
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, err
	}

	s := &Store{
		dir: dir,
	}

	if err := s.load(); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Store) load() error {
	data, err := os.ReadFile(filepath.Join(s.dir, "wallets.json"))
	if err != nil {
		if os.IsNotExist(err) {
			s.wallets = []models.Wallet{}
			return nil
		}
		return err
	}
	return json.Unmarshal(data, &s.wallets)
}

func (s *Store) save() error {
	data, err := json.MarshalIndent(s.wallets, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(s.dir, "wallets.json"), data, 0600)
}

func (s *Store) List() []models.Wallet {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]models.Wallet, len(s.wallets))
	copy(result, s.wallets)
	return result
}

func (s *Store) Get(idOrName string) (*models.Wallet, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for i := range s.wallets {
		if s.wallets[i].ID == idOrName || s.wallets[i].Name == idOrName || strings.HasPrefix(s.wallets[i].ID, idOrName) {
			return &s.wallets[i], nil
		}
	}
	return nil, errors.New("wallet not found")
}

func (s *Store) Add(w *models.Wallet) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.wallets {
		if s.wallets[i].ID == w.ID {
			return errors.New("wallet with this ID already exists")
		}
	}

	s.wallets = append(s.wallets, *w)
	return s.save()
}

func (s *Store) Delete(idOrName string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.wallets {
		if s.wallets[i].ID == idOrName || s.wallets[i].Name == idOrName || strings.HasPrefix(s.wallets[i].ID, idOrName) {
			s.wallets = append(s.wallets[:i], s.wallets[i+1:]...)
			return s.save()
		}
	}
	return errors.New("wallet not found")
}

func (s *Store) SetDefault(idOrName string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	found := false
	for i := range s.wallets {
		match := s.wallets[i].ID == idOrName || s.wallets[i].Name == idOrName || strings.HasPrefix(s.wallets[i].ID, idOrName)
		if match {
			found = true
		}
		s.wallets[i].IsDefault = match
	}
	if !found {
		return errors.New("wallet not found")
	}
	return s.save()
}

func (s *Store) GetDefault() *models.Wallet {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for i := range s.wallets {
		if s.wallets[i].IsDefault {
			return &s.wallets[i]
		}
	}
	if len(s.wallets) > 0 {
		return &s.wallets[0]
	}
	return nil
}

func (s *Store) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.wallets)
}
