package wallet

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"log"
	"os"
)

const walletsPath = "./tmp/wallets.data"

type Wallets struct {
	Wallets map[string]*Wallet
}

func (ws *Wallets) SaveFile() {
	var content bytes.Buffer

	gob.Register(elliptic.P256())

	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(ws)
	if err != nil {
		log.Panic(err)
	}

	err = os.WriteFile(walletsPath, content.Bytes(), 0644) // 0644 is the file permission for read/write
	if err != nil {
		log.Panic(err)
	}
}

func (ws *Wallets) LoadFile() error {
	if _, err := os.Stat(walletsPath); os.IsNotExist(err) {
		return err
	}

	var loadedWallets Wallets
	fileContent, err := os.ReadFile(walletsPath)
	if err != nil {
		return err
	}

	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decoder.Decode(&loadedWallets)
	if err != nil {
		return err
	}

	ws.Wallets = loadedWallets.Wallets

	return nil
}

func CreateWallets() (*Wallets, error) {
	wallets := Wallets{}
	wallets.Wallets = make(map[string]*Wallet)

	err := wallets.LoadFile()

	return &wallets, err
}

func (ws *Wallets) GetWallet(address string) Wallet {
	return *ws.Wallets[address]
}

func (ws *Wallets) GetAllAddresses() []string {
	var addresses []string

	for address := range ws.Wallets {
		addresses = append(addresses, address)
	}

	return addresses
}

func (ws *Wallets) AddWallet() string {
	newWallet := MakeWallet()
	address := string(newWallet.Address())

	ws.Wallets[address] = newWallet
	return address
}
