//Package msgpack ...
package msgpack

import (
	"github.com/amiraliio/avn-wallet/domain/model"
	"github.com/pkg/errors"
	"github.com/vmihailenco/msgpack/v4"
)

type Wallet struct{}

func (m *Wallet) Encode(input *model.Wallet) ([]byte, error) {
	rawWallet, err := msgpack.Marshal(input)
	if err != nil {
		return nil, errors.Wrap(err, "serializer.Wallet.Encode")
	}
	return rawWallet, nil
}

func (m *Wallet) Decode(input []byte) (*model.Wallet, error) {
	walletModel := new(model.Wallet)
	if err := msgpack.Unmarshal(input, walletModel); err != nil {
		return nil, errors.Wrap(err, "serializer.Wallet.Decode")
	}
	return walletModel, nil
}
