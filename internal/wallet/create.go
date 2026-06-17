package wallet

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/VT0x00/tonvault/internal/models"
)

func CreateNewWallet(store *Store, network string) (*models.Wallet, error) {
	words, err := GenerateMnemonic()
	if err != nil {
		return nil, fmt.Errorf("failed to generate mnemonic: %w", err)
	}

	fmt.Println("Your 24-word seed phrase (write this down safely):")
	fmt.Println("================================================")
	for i, word := range words {
		fmt.Printf("%3d. %s\n", i+1, word)
	}
	fmt.Println("================================================")
	fmt.Println("⚠ Store these words securely. Anyone with access to them can control your wallet.")
	fmt.Println()

	wallet, seedData, err := CreateFromMnemonic(words, network)
	if err != nil {
		return nil, err
	}

	if err := EncryptAndStore(store, wallet, seedData); err != nil {
		return nil, err
	}

	return wallet, nil
}

func ImportWalletFromSeed(store *Store, network string) (*models.Wallet, error) {
	fmt.Println("Enter your 24-word seed phrase (separated by spaces):")
	fmt.Print("> ")

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("failed to read input: %w", err)
	}

	input = strings.TrimRight(input, "\n\r ")
	words := ParseSeedInput(input)

	if !ValidateMnemonic(words) {
		return nil, fmt.Errorf("invalid seed phrase")
	}

	wallet, seedData, err := CreateFromMnemonic(words, network)
	if err != nil {
		return nil, err
	}

	if err := EncryptAndStore(store, wallet, seedData); err != nil {
		return nil, err
	}

	return wallet, nil
}

func ParseSeedInput(input string) []string {
	words := []string{}
	current := ""
	for _, c := range input {
		if c == ' ' || c == '\t' || c == '\n' || c == '\r' {
			if current != "" {
				words = append(words, current)
				current = ""
			}
		} else {
			current += string(c)
		}
	}
	if current != "" {
		words = append(words, current)
	}
	return words
}
