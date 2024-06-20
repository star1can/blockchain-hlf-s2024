package hlf

import (
	"encoding/json"
	"fmt"
	"github.com/hlf-mipt/saac-v2-core/pkg/model"
)

func (svc *HLFService) CreateAsset(user string, item *model.Asset) error {
	contract, err := svc.getContract(user)
	if err != nil {
		return err
	}

	buf, err := json.Marshal(*item)
	if err != nil {
		return fmt.Errorf("failed to json marshalling: %v", err)
	}

	_, err = contract.SubmitTransaction("CreateAsset", string(buf))
	if err != nil {
		return fmt.Errorf("failed to Submit transaction: %v", err)
	}

	return nil
}

func (svc *HLFService) UpdateAsset(user string, item *model.Asset) error {
	contract, err := svc.getContract(user)
	if err != nil {
		return err
	}

	buf, err := json.Marshal(*item)
	if err != nil {
		return fmt.Errorf("failed to json marshalling: %v", err)
	}

	_, err = contract.SubmitTransaction("UpdateAsset", string(buf))
	if err != nil {
		return fmt.Errorf("failed to Submit transaction: %v", err)
	}

	return nil
}

func (svc *HLFService) ReadAsset(user string, id int) (*model.Asset, error) {
	contract, err := svc.getContract(user)
	if err != nil {
		return nil, err
	}

	buf, err := json.Marshal(id)
	if err != nil {
		return nil, fmt.Errorf("failed to json marshalling: %v", err)
	}

	res, err := contract.SubmitTransaction("ReadAsset", string(buf))
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
