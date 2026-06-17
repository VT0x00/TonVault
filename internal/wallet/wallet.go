package wallet

import (
	"crypto/ed25519"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/VT0x00/tonvault/internal/models"
	"github.com/xssnick/tonutils-go/ton/wallet"
	"golang.org/x/crypto/ssh/terminal"
)

const walletVersionV4R2 = "v4r2"

func GenerateMnemonic() ([]string, error) {
	words := wallet.NewSeed()
	if len(words) == 0 {
		return nil, errors.New("failed to generate seed phrase")
	}
	return words, nil
}

func ValidateMnemonic(words []string) bool {
	_, err := wallet.SeedToPrivateKeyWithOptions(words)
	return err == nil
}

func CreateFromMnemonic(words []string, network string) (*models.Wallet, []byte, error) {
	mnemonicWords := strings.Join(words, " ")

	w, err := wallet.FromSeed(nil, words, wallet.V4R2)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create wallet from seed: %w", err)
	}

	addr := w.Address()
	pubKey := w.PrivateKey().Public().(ed25519.PublicKey)

	id := fmt.Sprintf("%x", pubKey)[:16]

	walletModel := &models.Wallet{
		ID:          id,
		Name:        "wallet-" + id[:8],
		Address:     addr.String(),
		PublicKey:   hex.EncodeToString(pubKey),
		Version:     walletVersionV4R2,
		Network:     network,
		SubwalletID: wallet.DefaultSubwallet,
		CreatedAt:   time.Now(),
	}

	return walletModel, []byte(mnemonicWords), nil
}

func PromptPassword(confirm bool) ([]byte, error) {
	fmt.Print("Enter password: ")
	password, err := terminal.ReadPassword(0)
	fmt.Println()
	if err != nil {
		return nil, err
	}
	if len(password) == 0 {
		return nil, errors.New("password cannot be empty")
	}
	if confirm {
		fmt.Print("Confirm password: ")
		confirmPassword, err := terminal.ReadPassword(0)
		fmt.Println()
		if err != nil {
			return nil, err
		}
		if string(password) != string(confirmPassword) {
			return nil, errors.New("passwords do not match")
		}
	}
	return password, nil
}

func EncryptAndStore(store *Store, wallet *models.Wallet, seedData []byte) error {
	password, err := PromptPassword(true)
	if err != nil {
		return err
	}

	encrypted, err := Encrypt(seedData, password)
	if err != nil {
		return fmt.Errorf("failed to encrypt seed: %w", err)
	}
	wallet.EncryptedSeed = encrypted

	if err := store.Add(wallet); err != nil {
		return fmt.Errorf("failed to store wallet: %w", err)
	}

	fmt.Println("Wallet created and stored securely.")
	return nil
}

func RecoverSeedPhrase(store *Store, walletID string) ([]string, error) {
	w, err := store.Get(walletID)
	if err != nil {
		return nil, err
	}

	fmt.Print("Enter password to decrypt wallet: ")
	password, err := terminal.ReadPassword(0)
	fmt.Println()
	if err != nil {
		return nil, err
	}

	decrypted, err := Decrypt(w.EncryptedSeed, password)
	if err != nil {
		return nil, errors.New("incorrect password")
	}

	words := strings.Fields(string(decrypted))
	return words, nil
}

func RecoverPrivateKey(store *Store, walletID string) (ed25519.PrivateKey, error) {
	words, err := RecoverSeedPhrase(store, walletID)
	if err != nil {
		return nil, err
	}

	priv, err := wallet.SeedToPrivateKeyWithOptions(words)
	if err != nil {
		return nil, fmt.Errorf("failed to recover private key: %w", err)
	}

	return priv, nil
}
