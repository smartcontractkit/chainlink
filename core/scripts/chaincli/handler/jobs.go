package handler

import (
	"context"
	"log"

	"github.com/smartcontractkit/chainlink/core/logger"
)

func (k *Keeper) CreateJob(ctx context.Context) {
	k.createJobs()
}

func (k *Keeper) createJobs() {
	lggr, closeLggr := logger.NewLogger()
	logger.Sugared(lggr).ErrorIfFn(closeLggr, "Failed to close logger")

	// Create Keeper Jobs on Nodes for Registry
	for i, keeperAddr := range k.cfg.Keepers {
		url := k.cfg.KeeperURLs[i]
		email := k.cfg.KeeperEmails[i]
		if len(email) == 0 {
			email = defaultChainlinkNodeLogin
		}
		pwd := k.cfg.KeeperPasswords[i]
		if len(pwd) == 0 {
			pwd = defaultChainlinkNodePassword
		}

		cl, err := authenticate(url, email, pwd, lggr)
		if err != nil {
			log.Fatal(err)
		}

		if err = k.createKeeperJob(cl, k.cfg.RegistryAddress, keeperAddr); err != nil {
			log.Fatal(err)
		}
	}
}
