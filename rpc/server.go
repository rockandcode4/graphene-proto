package rpc

import (
	"fmt"
	"log"
	"net/http"

	gorpc "github.com/gorilla/rpc"
	jsonrpc "github.com/gorilla/rpc/json"
	"github.com/rockandcode4/graphene-proto/consensus"
	"github.com/rockandcode4/graphene-proto/staking"
)

type Server struct {
	cons    *consensus.Consensus
	stake   *staking.Manager
	httpSrv *http.Server
	port    int
}

func NewServer(cons *consensus.Consensus, stake *staking.Manager, port int) (*Server, error) {
	s := &Server{cons: cons, stake: stake, port: port}
	rpcS := gorpc.NewServer()
	rpcS.RegisterCodec(jsonrpc.NewCodec(), "application/json")
	api := &API{cons: cons, stake: stake}
	if err := rpcS.RegisterService(api, "Graphene"); err != nil {
		return nil, err
	}
	mux := http.NewServeMux()
	mux.Handle("/rpc", rpcS)
	s.httpSrv = &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: mux}
	return s, nil
}

func (s *Server) Start() {
	go func() {
		log.Printf("RPC server listening on %s", s.httpSrv.Addr)
		if err := s.httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Println("rpc listen:", err)
		}
	}()
}

func (s *Server) Stop() {
	_ = s.httpSrv.Close()
}

type API struct {
	cons  *consensus.Consensus
	stake *staking.Manager
}

type SendArgs struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Amount uint64 `json:"amount"`
}
type SendReply struct {
	Ok    bool   `json:"ok"`
	Error string `json:"error,omitempty"`
}

func (a *API) SendTx(r *http.Request, args *SendArgs, reply *SendReply) error {
	if err := a.cons.SubmitTx(args.From, args.To, args.Amount); err != nil {
		reply.Ok = false
		reply.Error = err.Error()
		return nil
	}
	reply.Ok = true
	return nil
}

type BalanceArgs struct {
	Address string `json:"address"`
}
type BalanceReply struct {
	Balance uint64 `json:"balance"`
}

func (a *API) GetBalance(r *http.Request, args *BalanceArgs, reply *BalanceReply) error {
	b, err := a.cons.GetBalance(args.Address)
	if err != nil {
		return err
	}
	reply.Balance = b
	return nil
}

type RegisterValidatorArgs struct {
	Address string `json:"address"`
	Stake   uint64 `json:"stake"`
}
type GenericReply struct {
	Ok    bool   `json:"ok"`
	Error string `json:"error,omitempty"`
}

func (a *API) RegisterValidator(r *http.Request, args *RegisterValidatorArgs, reply *GenericReply) error {
	if err := a.stake.RegisterValidator(args.Address, args.Stake); err != nil {
		reply.Ok = false
		reply.Error = err.Error()
		return nil
	}
	reply.Ok = true
	return nil
}

type DelegateArgs struct {
	Delegator string `json:"delegator"`
	Validator string `json:"validator"`
	Amount    uint64 `json:"amount"`
}

func (a *API) Delegate(r *http.Request, args *DelegateArgs, reply *GenericReply) error {
	if err := a.stake.Delegate(args.Delegator, args.Validator, args.Amount); err != nil {
		reply.Ok = false
		reply.Error = err.Error()
		return nil
	}
	reply.Ok = true
	return nil
}
