package keystore

import (
	"context"

	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/plugins"
)

type Manager struct {
	services.Service
	lggr           logger.Logger
	looppRegistrar plugins.RegistrarConfig

	keystores map[string]*loop.KeystoreService
}

func NewLOOPPKeystoreManager(r plugins.RegistrarConfig, lggr logger.Logger) *Manager {
	return &Manager{
		lggr:           lggr.Named("LOOPPKeystoreManager"),
		looppRegistrar: r,
		keystores:      make(map[string]*loop.KeystoreService),
	}
}

func (m *Manager) Register(id string, cmd string) error {
	cmdFn, grpcOpts, err := m.looppRegistrar.RegisterLOOP(plugins.CmdConfig{
		ID:  id,
		Cmd: cmd,
		Env: nil,
	})
	if err != nil {
		return err
		//m.lggr.Errorw("Cannot start Keystore LOOPP", "ID", id, "cmd", cmd, "err", err)
	}
	ks := loop.NewKeystoreService(m.lggr, grpcOpts, cmdFn, nil)
	m.addKeystoreService(ks, id)
	return nil
}

func (m *Manager) addKeystoreService(ks *loop.KeystoreService, id string) error {
	//Check for duplicates
	m.keystores[id] = ks
	return nil
}

func (m *Manager) Start(ctx context.Context) error {
	for _, ks := range m.keystores {
		err := ks.Start(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Manager) Close() error {
	for _, ks := range m.keystores {
		err := ks.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
