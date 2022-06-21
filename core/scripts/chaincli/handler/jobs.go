package handler

import (
	"context"
	"log"
)

func (k *Keeper) CreateJob(ctx context.Context) {
	k.createJobs()
}

func (k *Keeper) createJobs() {
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
		err := k.createKeeperJobOnExistingNode(url, email, pwd, k.cfg.RegistryAddress, keeperAddr)
		if err != nil {
			log.Printf("Keeper Job not created for keeper %d: %s %s\n", i, url, keeperAddr)
			log.Println("Please create it manually")
		}
	}
}
