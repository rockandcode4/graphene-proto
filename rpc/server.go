package rpc

import (
    "fmt"
    "log"
    "net"
    "net/http"

    "github.com/gorilla/rpc"
    jsonrpc "github.com/gorilla/rpc/json"
    "github.com/rockandcode4/graphene-proto/consensus"
)

type Server struct {
    cons *consensus.Consensus
    httpSrv *http.Server
    port int
}

func NewServer(cons *consensus.Consensus, port int) (*Server, error) {
    s := &Server{cons: cons, port: port}
    rpcS := rpc.NewServer()
    rpcS.RegisterCodec(jsonrpc.NewCodec(), "application/json")
    api := &API{cons: cons}
    if err := rpcS.RegisterService(api, "Graphene"); err != nil {
        return nil, err
    }
    mux := http.NewServeMux()
    mux.Handle("/rpc", rpcS)
    s.httpSrv = &http.Server{ Addr: fmt.Sprintf(":%d", port), Handler: mux }
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
    cons *consensus.Consensus
}

type SendArgs struct {
    From string `json:"from"`
    To   string `json:"to"`
    Amount uint64 `json:"amount"`
}

type SendReply struct {
    Ok bool `json:"ok"`
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

type BalanceArgs struct { Address string `json:"address"` }
type BalanceReply struct { Balance uint64 `json:"balance"` }

func (a *API) GetBalance(r *http.Request, args *BalanceArgs, reply *BalanceReply) error {
    b, err := a.cons.GetBalance(args.Address)
    if err != nil { return err }
    reply.Balance = b
    return nil
}
