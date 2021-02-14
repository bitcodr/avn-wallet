package service

import (
	"errors"

	"github.com/bitcodr/avn-wallet/domain/model"
	"github.com/bitcodr/avn-wallet/helper"
)

var (
	ErrWalletNotFound = errors.New("Wallet Not Found")
	ErrWalletInvalid  = errors.New("Wallet Invalid")
)

type WalletService interface {
	Get(cellphone uint64) (*model.Wallet, error)
	Insert(wallet *model.Wallet) (*model.Wallet, error)
	Transactions(cellphone uint64) ([]*model.Transaction, error)
}

type WalletRepository interface {
	Get(cellphone uint64) (*model.Wallet, error)
	Insert(wallet *model.Wallet) (*model.Wallet, error)
	Transactions(cellphone uint64) ([]*model.Transaction, error)
}

type WalletSerializer interface {
	Encode(input *model.Wallet) ([]byte, error)
	Decode(input []byte) (*model.Wallet, error)
}

type walletService struct {
	walletRepo WalletRepository
}

func NewWalletService(walletRepo WalletRepository) WalletService {
	return &walletService{
		walletRepo,
	}
}

func (w *walletService) Get(cellphone uint64) (*model.Wallet, error) {
	return w.walletRepo.Get(cellphone)
}

func (w *walletService) Insert(wallet *model.Wallet) (*model.Wallet, error) {
	if err := helper.ValidateModel(wallet); err != nil {
		return nil, err
	}
	return w.walletRepo.Insert(wallet)
}

func (w *walletService) Transactions(cellphone uint64) ([]*model.Transaction, error) {
	return w.walletRepo.Transactions(cellphone)
}
