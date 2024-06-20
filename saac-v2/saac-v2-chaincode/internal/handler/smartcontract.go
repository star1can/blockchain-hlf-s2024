package handler

import (
	"errors"
	"fmt"
	"github.com/hlf-mipt/saac-v2-chaincode-go/internal/repository"
	"github.com/hlf-mipt/saac-v2-core/pkg/model"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"reflect"
	"strconv"
)

func toOwner(clientId, mspId string) string {
	return fmt.Sprintf("%v@%v", clientId, mspId)
}

// CreateAsset issues a new asset to the world state with given details.
func (s *SmartContract) CreateAsset(ctx contractapi.TransactionContextInterface, item *model.Asset) error {
	exists, err := s.AssetExists(ctx, item.ID)
	if err != nil {
		s.log.Error(err)
		return err
	}

	if exists {
		err = fmt.Errorf("the asset %s already exists", item.ID)
		s.log.Error(err)
		return err
	}

	mspId, clientID, err := GetClient(ctx)
	if err != nil {
		s.log.Error(err)
		return err
	}

	item.Owner = toOwner(clientID, mspId)
	item.Type = fmt.Sprint(reflect.TypeOf(*item))

	s.log.Debugf("trying to create asset: %v", item)
	err = s.repository.Assets().SetPrivate(ctx, mspId, item)
	if err != nil {
		s.log.Error(err)
		return err
	}

	return nil
}

// ReadAsset returns the asset stored in the world state with given id.
func (s *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, id int) (*model.Asset, error) {
	mspId, _, err := GetClient(ctx)

	asset, err := s.repository.Assets().GetPrivate(ctx, mspId, strconv.Itoa(id))
	if err != nil {
		if errors.Is(err, repository.NotExistsErr) {
			err = errors.New(fmt.Sprintf("the asset %v does not exist in collection: %v", id, mspId))
		}

		s.log.Error(err)
		return nil, err
	}

	return asset, nil
}

// UpdateAsset updates an existing asset in the world state with provided parameters.
func (s *SmartContract) UpdateAsset(ctx contractapi.TransactionContextInterface, item *model.Asset) error {
	asset, err := s.ReadAsset(ctx, item.ID)
	if err != nil {
		s.log.Error(err)
		return err
	}

	isGranted, err := s.isAccessGranted(ctx, asset)
	if err != nil {
		s.log.Error(err)
		return err
	}

	if !isGranted {
		err = errors.New("forbidden")
		s.log.Error(err)
		return err
	}

	mspId, _, err := GetClient(ctx)
	if err != nil {
		s.log.Error(err)
		return err
	}

	err = s.repository.Assets().SetPrivate(ctx, mspId, asset)
	if err != nil {
		s.log.Error(err)
		return err
	}

	return nil
}

// AssetExists returns true when asset with given ID exists in world state
func (s *SmartContract) AssetExists(ctx contractapi.TransactionContextInterface, id int) (bool, error) {
	mspId, _, err := GetClient(ctx)
	if err != nil {
		s.log.Error(err)
		return false, err
	}

	_, err = s.repository.Assets().GetPrivate(ctx, mspId, strconv.Itoa(id))
	if err != nil {
		if errors.Is(err, repository.NotExistsErr) {
			return false, nil
		}
		s.log.Error(err)
		return false, err
	}

	return true, nil
}

func (s *SmartContract) isAccessGranted(ctx contractapi.TransactionContextInterface, item *model.Asset) (bool, error) {
	mspId, clientId, err := GetClient(ctx)
	s.log.Debugf("clientId: %v; mspId: %v", clientId, mspId)
	if err != nil {
		s.log.Error(err)
		return false, err
	}
	owner := toOwner(clientId, mspId)
	s.log.Debugf("currOwner: %v; actualOwner: %v", owner, item.Owner)
	if item.Owner != owner {
		return false, nil
	}

	return true, nil
}
