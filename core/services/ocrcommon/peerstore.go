package ocrcommon

import (
	"context"
	"fmt"
	"strings"
	"time"

	p2ppeer "github.com/libp2p/go-libp2p-core/peer"
	p2ppeerstore "github.com/libp2p/go-libp2p-core/peerstore"
	"github.com/libp2p/go-libp2p-peerstore/pstoremem"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/recovery"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type (
	P2PPeer struct {
		ID        string
		Addr      string
		PeerID    string
		CreatedAt time.Time
		UpdatedAt time.Time
	}

	Pstorewrapper struct {
		utils.StartStopOnce
		Peerstore     p2ppeerstore.Peerstore
		peerID        string
		q             pg.Q
		writeInterval time.Duration
		ctx           context.Context
		ctxCancel     context.CancelFunc
		chDone        chan struct{}
		lggr          logger.SugaredLogger
	}
)

// NewPeerstoreWrapper creates a new database-backed peerstore wrapper scoped to the given jobID
// Multiple peerstore wrappers should not be instantiated with the same jobID
func NewPeerstoreWrapper(db *sqlx.DB, writeInterval time.Duration, peerID p2pkey.PeerID, lggr logger.Logger, cfg pg.QConfig) (*Pstorewrapper, error) {
	ctx, cancel := context.WithCancel(context.Background())
	namedLogger := lggr.Named("PeerStore")
	q := pg.NewQ(db, namedLogger, cfg)

	return &Pstorewrapper{
		utils.StartStopOnce{},
		pstoremem.NewPeerstore(),
		peerID.Raw(),
		q,
		writeInterval,
		ctx,
		cancel,
		make(chan struct{}),
		logger.Sugared(namedLogger),
	}, nil
}

func (p *Pstorewrapper) Start() error {
	return p.StartOnce("PeerStore", func() error {
		err := p.readFromDB()
		if err != nil {
			return errors.Wrap(err, "could not start peerstore wrapper")
		}
		go p.dbLoop()
		return nil
	})
}

func (p *Pstorewrapper) dbLoop() {
	defer close(p.chDone)
	ticker := time.NewTicker(utils.WithJitter(p.writeInterval))
	defer ticker.Stop()
	for {
		select {
		case <-p.ctx.Done():
			return
		case <-ticker.C:
			recovery.WrapRecover(p.lggr, func() {
				if err := p.WriteToDB(); err != nil {
					p.lggr.Errorw("Error writing peerstore to DB", "err", err)
				}
			})
		}
	}
}

func (p *Pstorewrapper) Close() error {
	return p.StopOnce("PeerStore", func() error {
		p.ctxCancel()
		<-p.chDone
		return p.Peerstore.Close()
	})
}

func (p *Pstorewrapper) readFromDB() error {
	peers, err := p.getPeers()
	if err != nil {
		return err
	}
	for _, peer := range peers {
		peerID, err := p2ppeer.Decode(peer.ID)
		if err != nil {
			return errors.Wrapf(err, "unexpectedly failed to decode peer ID '%s'", peer.ID)
		}
		peerAddr, err := ma.NewMultiaddr(peer.Addr)
		if err != nil {
			return errors.Wrapf(err, "unexpectedly failed to decode peer multiaddr '%s'", peer.Addr)
		}
		p.Peerstore.AddAddr(peerID, peerAddr, p2ppeerstore.PermanentAddrTTL)
	}
	return nil
}

func (p *Pstorewrapper) getPeers() (peers []P2PPeer, err error) {
	rows, err := p.q.WithOpts(pg.WithParentCtx(p.ctx)).Query(`SELECT id, addr FROM p2p_peers WHERE peer_id = $1`, p.peerID)
	if err != nil {
		return nil, errors.Wrap(err, "error querying peers")
	}
	defer p.lggr.ErrorIfFn(rows.Close, "Error closing p2p_peers rows")

	peers = make([]P2PPeer, 0)

	for rows.Next() {
		peer := P2PPeer{}
		if err = rows.Scan(&peer.ID, &peer.Addr); err != nil {
			return nil, errors.Wrap(err, "unexpected error scanning row")
		}
		peers = append(peers, peer)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return peers, nil
}

func (p *Pstorewrapper) WriteToDB() error {
	q := p.q.WithOpts(pg.WithParentCtx(p.ctx))
	err := q.Transaction(func(tx pg.Queryer) error {
		_, err := tx.Exec(`DELETE FROM p2p_peers WHERE peer_id = $1`, p.peerID)
		if err != nil {
			return errors.Wrap(err, "delete from p2p_peers failed")
		}
		peers := make([]P2PPeer, 0)
		for _, pid := range p.Peerstore.PeersWithAddrs() {
			addrs := p.Peerstore.Addrs(pid)
			for _, addr := range addrs {
				p := P2PPeer{
					ID:     pid.String(),
					Addr:   addr.String(),
					PeerID: p.peerID,
				}
				peers = append(peers, p)
			}
		}
		valueStrings := []string{}
		valueArgs := []interface{}{}
		for _, p := range peers {
			valueStrings = append(valueStrings, "(?, ?, ?, NOW(), NOW())")
			valueArgs = append(valueArgs, p.ID)
			valueArgs = append(valueArgs, p.Addr)
			valueArgs = append(valueArgs, p.PeerID)
		}

		/* #nosec G201 */
		stmt := fmt.Sprintf("INSERT INTO p2p_peers (id, addr, peer_id, created_at, updated_at) VALUES %s", strings.Join(valueStrings, ","))
		stmt = sqlx.Rebind(sqlx.DOLLAR, stmt)
		_, err = tx.Exec(stmt, valueArgs...)
		return errors.Wrap(err, "insert into p2p_peers failed")
	})
	return errors.Wrap(err, "could not write peers to DB")
}
