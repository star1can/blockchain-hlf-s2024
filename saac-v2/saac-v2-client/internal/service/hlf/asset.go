package hlf

import (
	"encoding/json"
	"fmt"
	"github.com/hlf-mipt/saac-v2-core/pkg/model"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

func (svc *HLFService) CreateAsset(user string, item *model.Asset) error {
	buf, err := json.Marshal(*item)
	if err != nil {
		return fmt.Errorf("failed to json marshalling: %v", err)
	}

	_, err = svc.SubmitTransaction(user, "CreateAsset", string(buf))
	if err != nil {
		return fmt.Errorf("failed to Submit transaction: %v", err)
	}

	return nil
}

func (svc *HLFService) UpdateAsset(user string, item *model.Asset) error {
	buf, err := json.Marshal(*item)
	if err != nil {
		return fmt.Errorf("failed to json marshalling: %v", err)
	}

	_, err = svc.SubmitTransaction(user, "UpdateAsset", string(buf))
	if err != nil {
		return fmt.Errorf("failed to Submit transaction: %v", err)
	}

	return nil
}

func (svc *HLFService) ReadAsset(user string, id int) (*model.Asset, error) {
	buf, err := json.Marshal(id)
	if err != nil {
		return nil, fmt.Errorf("failed to json marshalling: %v", err)
	}

	res, err := svc.SubmitTransaction(user, "ReadAsset", string(buf))
	if err != nil {
		return nil, fmt.Errorf("failed to Submit transaction: %v", err)
	}

	item := &model.Asset{}
	err = json.Unmarshal(res, item)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal result: %v", err)
	}

	return item, nil
}

func (svc *HLFService) SubmitTransaction(user string, name string, args ...string) ([]byte, error) {
	contract, err := svc.getContract(user)
	if err != nil {
		return nil, err
	}

	peers, err := svc.getEndorsingPeers(user)
	if err != nil {
		return nil, err
	}

	txn, err := contract.CreateTransaction(name,
		gateway.WithEndorsingPeers(peers...))

	if err != nil {
		return nil, err
	}

	return txn.Submit(args...)
}
