package keystore

import (
	"encoding/json"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/syndtr/goleveldb/leveldb"
)

type KeyStore struct {
	db *leveldb.DB
}

func NewKeyStore(filepath string) (types.KeyStore, error) {
	db, err := leveldb.OpenFile(filepath, nil)
	if err != nil {
		return nil, err
	}

	return &KeyStore{db}, nil
}

func (ks *KeyStore) List() ([]string, error) {
	wallets := []string{}

	iterator := ks.db.NewIterator(nil, nil)
	defer iterator.Release()
	for iterator.Next() {
		wallets = append(wallets, string(iterator.Key()))
	}
	if iterator.Error() != nil {
		return nil, iterator.Error()
	}
	return wallets, nil
}

func (ks *KeyStore) Get(k string) (types.KeyInfo, error) {
	value, err := ks.db.Get([]byte(k), nil)
	if err != nil {
		return types.KeyInfo{}, err
	}

	var keyInfo types.KeyInfo
	err = json.Unmarshal(value, &keyInfo)
	if err != nil {
		return types.KeyInfo{}, err
	}
	return keyInfo, nil
}

func (ks *KeyStore) Put(k string, ki types.KeyInfo) error {
	kiJson, err := json.Marshal(ki)
	if err != nil {
		return err
	}

	err = ks.db.Put([]byte(k), kiJson, nil)
	if err != nil {
		return err
	}
	return nil
}

func (ks *KeyStore) Delete(k string) error {
	err := ks.db.Delete([]byte(k), nil)
	if err != nil {
		return err
	}
	return nil
}
