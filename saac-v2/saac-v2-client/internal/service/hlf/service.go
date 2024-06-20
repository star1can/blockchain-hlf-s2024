package hlf

import (
	"context"
	"fmt"
	"github.com/hlf-mipt/saac-v2-client/internal/config"
	fabricsdkconfig "github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type HLFService struct {
	log        *logrus.Entry
	cfg        config.HLF
	wallet     sync.Map
	contractGW sync.Map
}

func NewHLFService(lc fx.Lifecycle,
	cfg *config.Config,
	log *logrus.Entry,
) *HLFService {
	srv := &HLFService{
		log:        log,
		cfg:        cfg.HLF,
		wallet:     sync.Map{},
		contractGW: sync.Map{},
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if err := srv.Start(ctx); err != nil {
				return err
			}
			return nil
		},
	})

	return srv
}

func (svc *HLFService) Start(ctx context.Context) error {
	err := os.Setenv("DISCOVERY_AS_LOCALHOST", "true")
	if err != nil {
		return fmt.Errorf("Error setting DISCOVERY_AS_LOCALHOST environment variable: %v", err)
	}

	go svc.initWallet("User1@org1.example.com", "Org1MSP", filepath.Join(
		svc.cfg.OrganizationsFolderPath,
		"peerOrganizations",
		"org1.example.com",
		"connection-org1.yaml",
	))

	go svc.initWallet("User1@org2.example.com", "Org2MSP", filepath.Join(
		svc.cfg.OrganizationsFolderPath,
		"peerOrganizations",
		"org2.example.com",
		"connection-org2.yaml",
	))
	return nil
}

func (svc *HLFService) initWallet(userName, mspId, ccpPath string) {
	wallet, err := gateway.NewFileSystemWallet(svc.cfg.WalletPath)
	if err != nil {
		svc.log.Errorf("failed to create wallet: %v", err)
		return
	}

	if !wallet.Exists(userName) {
		err = svc.populateWallet(wallet, userName, mspId)
		if err != nil {
			svc.log.Errorf("failed to populate wallet contents: %v", err)
			return
		}
	}

	gw, err := gateway.Connect(
		gateway.WithConfig(fabricsdkconfig.FromFile(filepath.Clean(ccpPath))),
		gateway.WithIdentity(wallet, userName),
	)
	if err != nil {
		svc.log.Errorf("failed to connect to gateway: %v", err)
		return
	}

	network, err := gw.GetNetwork(svc.cfg.ChannelName)
	if err != nil {
		svc.log.Errorf("failed to get network: %v", err)
		return
	}

	contract := network.GetContract(svc.cfg.ChaincodeName)

	svc.contractGW.Store(userName, contract)
}

func (svc *HLFService) getContract(user string) (contract *gateway.Contract, err error) {
	item, ok := svc.contractGW.Load(user)
	if !ok || item == nil {
		err = fmt.Errorf("user =[%v] not found", user)
		return
	}

	contract = item.(*gateway.Contract)
	return
}

func (svc *HLFService) populateWallet(wallet *gateway.Wallet, userName, mspid string) error {
	parts := strings.Split(userName, "@")
	if len(parts) != 2 {
		return fmt.Errorf("wrond username=[%v]", userName)
	}

	credPath := filepath.Join(
		svc.cfg.OrganizationsFolderPath,
		"peerOrganizations",
		parts[1],
		"users",
		userName,
		"msp",
	)

	certPath := filepath.Join(credPath, "signcerts", "cert.pem")
	// read the certificate pem
	cert, err := os.ReadFile(filepath.Clean(certPath))
	if err != nil {
		return err
	}

	keyDir := filepath.Join(credPath, "keystore")
	// there's a single file in this dir containing the private key
	files, err := os.ReadDir(keyDir)
	if err != nil {
		return err
	}
	if len(files) != 1 {
		return fmt.Errorf("keystore folder should have contain one file")
	}
	keyPath := filepath.Join(keyDir, files[0].Name())
	key, err := os.ReadFile(filepath.Clean(keyPath))
	if err != nil {
		return err
	}

	identity := gateway.NewX509Identity(mspid, string(cert), string(key))

	return wallet.Put(userName, identity)
}
