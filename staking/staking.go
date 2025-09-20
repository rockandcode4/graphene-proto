package staking

import (
	"fmt"
	"sync"

	"github.com/rockandcode4/graphene-proto/consensus"
	"github.com/rockandcode4/graphene-proto/state"
)

type Validator struct {
	Address string
	Stake   uint64
	Active  bool
}

type Delegation struct {
	Delegator string
	Validator string
	Amount    uint64
}

type Manager struct {
	st   *state.StateDB
	cons *consensus.Consensus

	mu          sync.Mutex
	validators  map[string]*Validator
	delegations map[string][]*Delegation // validator -> list
}

func NewManager(st *state.StateDB, cons *consensus.Consensus) *Manager {
	return &Manager{st: st, cons: cons, validators: make(map[string]*Validator), delegations: make(map[string][]*Delegation)}
}

func (m *Manager) RegisterValidator(addr string, stake uint64) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	// deduct stake from account
	acct, _ := m.st.GetAccount(addr)
	if acct.Balance < stake {
		return fmt.Errorf("insufficient balance")
	}
	acct.Balance -= stake
	if err := m.st.PutAccount(acct); err != nil {
		return err
	}

	m.validators[addr] = &Validator{Address: addr, Stake: stake, Active: true}
	// notify consensus about validator set (simple replacement)
	var vals []string
	for k, v := range m.validators {
		if v.Active {
			vals = append(vals, k)
		}
	}
	m.cons.SetValidators(vals)
	return nil
}

func (m *Manager) Delegate(delegator, validator string, amount uint64) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	acct, _ := m.st.GetAccount(delegator)
	if acct.Balance < amount {
		return fmt.Errorf("insufficient balance")
	}
	acct.Balance -= amount
	if err := m.st.PutAccount(acct); err != nil {
		return err
	}

	m.delegations[validator] = append(m.delegations[validator], &Delegation{Delegator: delegator, Validator: validator, Amount: amount})
	// increase validator stake in manager (does not change validator's locked stake here for simplicity)
	if v, ok := m.validators[validator]; ok {
		v.Stake += amount
	} else {
		// create placeholder validator entry if not present (will not be active until registered)
		m.validators[validator] = &Validator{Address: validator, Stake: amount, Active: false}
	}
	// update consensus validators list for active ones
	var vals []string
	for k, v := range m.validators {
		if v.Active {
			vals = append(vals, k)
		}
	}
	m.cons.SetValidators(vals)
	return nil
}

func (m *Manager) GetValidators() []*Validator {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := []*Validator{}
	for _, v := range m.validators {
		out = append(out, v)
	}
	return out
}
