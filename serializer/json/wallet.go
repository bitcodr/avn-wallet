//Package json ...
package json

import (
	"encoding/json"

	"github.com/bitcodr/avn-wallet/domain/model"
	"github.com/pkg/errors"
)

type Wallet struct{}

func (m *Wallet) Encode(input *model.Wallet) ([]byte, error) {
	rawWallet, err := json.Marshal(input)
	if err != nil {
		return nil, errors.Wrap(err, "serializer.Wallet.Encode")
	}
	return rawWallet, nil
}

func (m *Wallet) Decode(input []byte) (*model.Wallet, error) {
	walletModel := new(model.Wallet)
	if err := json.Unmarshal(input, walletModel); err != nil {
		return nil, errors.Wrap(err, "serializer.Wallet.Decode")
	}
	return walletModel, nil
}
