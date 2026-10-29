// Copyright 2014 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package vm

import "github.com/ethereum/go-ethereum/common"

// OpCodeHooks is a set of hooks that can be used to intercept and modify the
// behavior of the EVM when executing certain opcodes.
// The hooks are called before the execution of the respective opcodes.
type OpCodeHooks interface {
	// CallHook is called before executing a CALL, CALLCODE, DELEGATECALL and STATICCALL opcodes.
	CallHook(evm *EVM, caller common.Address, recipient common.Address) error
	// CreateHook is called before executing a CREATE and CREATE2 opcodes.
	CreateHook(evm *EVM, caller common.Address) error
}

// CreatedCodeHooks is an OPTIONAL extension to OpCodeHooks, consulted on the
// code a contract creation produced, before that code is written to state.
//
// WHY A SEPARATE INTERFACE RATHER THAN A THIRD METHOD ON OpCodeHooks.
// Adding a method to OpCodeHooks is a breaking change for every existing
// implementation, including hooks types that assert against that interface in
// downstream repositories. A consumer needing this check would then have to
// fork those too. Detected by type assertion instead, so implementations that
// do not want it are unaffected.
//
// WHY IT IS NEEDED. CreateHook runs BEFORE the opcode and receives only the
// caller, so it can answer "may this account deploy" and cannot answer "is
// this code permitted" - the code does not exist yet. A chain that admits
// contracts by what their code IS, rather than by who deploys them, has no
// point at which to make that decision. This is that point: the runtime code
// exists, and the creation can still be failed.
type CreatedCodeHooks interface {
	// CreatedCodeHook is called with the runtime code a creation produced,
	// after the constructor has run and before the code is written to state.
	// A non-nil error fails the creation; the code is not written.
	CreatedCodeHook(evm *EVM, address common.Address, code []byte) error
}

type NoopOpCodeHooks struct {
}

func (NoopOpCodeHooks) CallHook(evm *EVM, caller common.Address, recipient common.Address) error {
	return nil
}

func (NoopOpCodeHooks) CreateHook(evm *EVM, caller common.Address) error {
	return nil
}

func newNoopOpCodeHooks() OpCodeHooks {
	return NoopOpCodeHooks{}
}

func NewDefaultOpCodeHooks() OpCodeHooks {
	return newNoopOpCodeHooks()
}
