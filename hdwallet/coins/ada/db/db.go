package db

import (
	"fmt"
	"git.mazdax.tech/blockchain/hdwallet/coins/ada/domain"
	"github.com/dgraph-io/badger/v3"
	"io/ioutil"
	"os"
	"time"
)

type DB interface {
	SaveWallet(*Wallet) error
	GetWallets() ([]*Wallet, error)
	GetWallet(name string) (*Wallet, error)
	DeleteWallet(string) error
	Backup(since uint64) (uint64, error)
	Load(backupFilePath string, maxPendingWrites int) error
	Close()
}

type badgerDB struct {
	db         *badger.DB
	DbPath     string
	BackupPath string
}

var WalletNotFoundErr = fmt.Errorf("no wallet found with this name")

func NewBadgerDB(config domain.ADAConfig) *badgerDB {
	if _, err := os.Stat(config.DBPath); err != nil {
		err = os.MkdirAll(config.DBPath, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}

	db, err := badger.Open(badger.DefaultOptions(config.DBPath).WithLoggingLevel(badger.ERROR))
	if err != nil {
		panic(err)
	}
	return &badgerDB{db, config.DBPath, config.BackupPath}
}

func (bdb *badgerDB) Close() {
	bdb.db.Close()
}

func (bdb *badgerDB) SaveWallet(w *Wallet) error {
	bytes, err := w.Marshal()
	if err != nil {
		return err
	}
	err = bdb.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(w.ID), bytes)
	})
	if err != nil {
		return err
	}
	go bdb.Backup(0)
	return nil
}

func (bdb *badgerDB) GetWallets() ([]*Wallet, error) {
	wallets := []*Wallet{}
	err := bdb.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			value, err := item.ValueCopy(nil)
			if err != nil {
				return err
			}
			wallet := &Wallet{}
			wallet.Unmarshal(value)
			wallets = append(wallets, wallet)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return wallets, nil
}

func (bdb *badgerDB) GetWallet(name string) (*Wallet, error) {
	wallets, err := bdb.GetWallets()
	if err != nil {
		return nil, err
	}
	for _, wall := range wallets {
		if wall.Name == name {
			return wall, nil
		}
	}
	return nil, WalletNotFoundErr
}

func (bdb *badgerDB) DeleteWallet(id string) error {
	err := bdb.db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(id))
	})
	return err
}

func (bdb *badgerDB) Backup(since uint64) (uint64, error) {
	today := time.Now().Format("20060102")
	path := bdb.BackupPath + "/" + today
	if _, err := os.Stat(path); err != nil {
		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return 0, err
		}
	}
	dir, err := ioutil.TempDir(path, "mazdax_")
	if err != nil {
		return 0, err
	}
	bak, err := ioutil.TempFile(dir, "backup_")

	return bdb.db.Backup(bak, since)
}

func (bdb *badgerDB) Load(backupFilePath string, maxPendingWrites int) error {
	bak, err := os.Open(backupFilePath)
	if err != nil {
		return err
	}

	return bdb.db.Load(bak, maxPendingWrites)
}
