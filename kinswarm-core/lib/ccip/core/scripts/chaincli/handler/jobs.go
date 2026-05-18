package handler

import (
	"context"
	"log"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func (k *Keeper) CreateJob(ctx context.Context) {
	k.createJobs(ctx)
}

func (k *Keeper) createJobs(ctx context.Context) {
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

		cl, err := authenticate(ctx, url, email, pwd, lggr)
		if err != nil {
			log.Fatal(err)
		}

		if err = k.createKeeperJob(ctx, cl, k.cfg.RegistryAddress, keeperAddr); err != nil {
			log.Fatal(err)
		}
	}
}
