package sign_service

import (
	"context"
	"fmt"
	"time"

	"fil-kms/app/config"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-jsonrpc"
	"github.com/filecoin-project/go-state-types/builtin/v10/market"
	"github.com/filecoin-project/go-state-types/crypto"
	"github.com/filecoin-project/lotus/api"
	"github.com/filecoin-project/lotus/api/client"
	cliutil "github.com/filecoin-project/lotus/cli/util"
	logging "github.com/ipfs/go-log/v2"
	"golang.org/x/xerrors"
)

var log = logging.Logger("sign_service")

type IWallet interface {
	WalletSign(ctx context.Context, signer address.Address, toSign []byte, meta api.MsgMeta) (*crypto.Signature, error)
	WalletList(context.Context) ([]address.Address, error)
	Close()
}

type RemoteWallet struct {
	api.Wallet
	closer jsonrpc.ClientCloser
}

func NewRemoteWallet(ctx context.Context, url, token string) (IWallet, error) {
	ai := cliutil.APIInfo{
		Addr:  url,
		Token: []byte(token),
	}

	url, err := ai.DialArgs("v0")
	if err != nil {
		return nil, err
	}

	wapi, closer, err := client.NewWalletRPCV0(ctx, url, ai.AuthHeader())
	if err != nil {
		return nil, xerrors.Errorf("creating jsonrpc client: %w", err)
	}

	return &RemoteWallet{
		Wallet: wapi,
		closer: closer,
	}, nil
}

func (w *RemoteWallet) Close() {
	w.closer()
}

type LimitedWallet struct {
	wapi IWallet
	cfg  *config.Config

	ch chan signReq
}

type signReq struct {
	toSign   []byte
	meta     api.MsgMeta
	singer   address.Address
	proposal *market.DealProposal
	res      chan result
}

type result struct {
	err error
	sig *crypto.Signature
}

func NewLimitedWallet(wapi IWallet, cfg *config.Config) *LimitedWallet {
	return &LimitedWallet{
		wapi: wapi,
		cfg:  cfg,

		ch: make(chan signReq, 10),
	}
}

func (w *LimitedWallet) Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case req := <-w.ch:
			proposal := req.proposal
			filter := w.cfg.GetFilter(proposal.Client, proposal.Provider)
			if filter == nil {
				req.res <- result{err: fmt.Errorf("not allowed to sign deals")}
				continue
			}
			now := time.Now()
			if !filter.Start.Time().IsZero() && filter.Start.Time().After(now) {
				req.res <- result{err: fmt.Errorf("start time: %v > %s", filter.Start, now)}
				continue
			}
			if !filter.End.Time().IsZero() && filter.End.Time().Before(now) {
				req.res <- result{err: fmt.Errorf("end time: %v < %s", filter.End, now)}
				continue
			}

			if filter.GetLimit() < filter.Used+int64(proposal.PieceSize) {
				req.res <- result{err: fmt.Errorf("exceed limit: %d < %d", filter.GetLimit(), filter.Used+int64(proposal.PieceSize))}
				continue
			}
			sig, err := w.wapi.WalletSign(ctx, req.singer, req.toSign, req.meta)
			if err != nil {
				req.res <- result{err: err}
				continue
			}

			if err := w.cfg.SaveFilter(proposal.Client, proposal.Provider, filter); err != nil {
				req.res <- result{err: err}
				continue
			}

			req.res <- result{sig: sig}
		}
	}
}

func (w *LimitedWallet) Stop() {
	w.wapi.Close()
}

func (w *LimitedWallet) WalletSign(ctx context.Context, signer address.Address, toSign []byte, meta api.MsgMeta, proposal *market.DealProposal) (*crypto.Signature, error) {
	req := signReq{
		toSign:   toSign,
		meta:     meta,
		singer:   signer,
		proposal: proposal,
		res:      make(chan result, 1),
	}
	w.ch <- req
	res := <-req.res
	close(req.res)
	if res.err != nil {
		log.Errorf("sign failed, signer: %s, provider: %s, err: %v", signer, proposal.Provider, res.err)
		return nil, fmt.Errorf("sign failed, signer: %s, provider: %s, err: %v", signer, proposal.Provider, res.err)
	}
	return res.sig, nil
}

func (w *LimitedWallet) WalletList(ctx context.Context) ([]address.Address, error) {
	return w.wapi.WalletList(ctx)
}
