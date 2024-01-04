// Copyright 2021 FerretDB Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package conninfo provides access to connection-specific information.
package conninfo

import (
	"context"
	"sync"

	scramutil "github.com/FerretDB/FerretDB/internal/util/scram"
)

// contextKey is a named unexported type for the safe use of context.WithValue.
type contextKey struct{}

// Context key for WithConnInfo/Get.
var connInfoKey = contextKey{}

// ConnInfo represents client connection information.
type ConnInfo struct {
	// the order of fields is weird to make the struct smaller due to alignment

	PeerAddr     string
	username     string // protected by rw
	password     string // protected by rw
	metadataRecv bool   // protected by rw
	BypassAuth   bool
	rw           sync.RWMutex

	sc *scramutil.ScramConversation
}

// New returns a new ConnInfo.
func New() *ConnInfo {
	return new(ConnInfo)
}

// Conv returns stored SCRAM server conversation.
func (connInfo *ConnInfo) Conv() *scramutil.ScramConversation {
	connInfo.rw.RLock()
	defer connInfo.rw.RUnlock()

	return connInfo.sc
}

// SetConv stores the SCRAM server conversation.
func (connInfo *ConnInfo) SetConv(sc *scramutil.ScramConversation) {
	connInfo.rw.RLock()
	defer connInfo.rw.RUnlock()

	connInfo.sc = sc
}

// Username returns stored username.
func (connInfo *ConnInfo) Username() string {
	connInfo.rw.RLock()
	defer connInfo.rw.RUnlock()

	return connInfo.username
}

// Auth returns stored username and password.
func (connInfo *ConnInfo) Auth() (username, password string) {
	connInfo.rw.RLock()
	defer connInfo.rw.RUnlock()

	return connInfo.username, connInfo.password
}

// SetAuth stores username and password.
func (connInfo *ConnInfo) SetAuth(username, password string) {
	connInfo.rw.Lock()
	defer connInfo.rw.Unlock()

	connInfo.username = username
	connInfo.password = password
}

// MetadataRecv returns whatever client metadata was received already.
func (connInfo *ConnInfo) MetadataRecv() bool {
	connInfo.rw.RLock()
	defer connInfo.rw.RUnlock()

	return connInfo.metadataRecv
}

// SetMetadataRecv marks client metadata as received.
func (connInfo *ConnInfo) SetMetadataRecv() {
	connInfo.rw.Lock()
	defer connInfo.rw.Unlock()

	connInfo.metadataRecv = true
}

// Ctx returns a derived context with the given ConnInfo.
func Ctx(ctx context.Context, connInfo *ConnInfo) context.Context {
	return context.WithValue(ctx, connInfoKey, connInfo)
}

// Get returns the ConnInfo value stored in ctx.
func Get(ctx context.Context) *ConnInfo {
	value := ctx.Value(connInfoKey)
	if value == nil {
		panic("connInfo is not set in context")
	}

	connInfo, ok := value.(*ConnInfo)
	if !ok {
		panic("connInfo is set in context with wrong value type")
	}

	if connInfo == nil {
		panic("connInfo is set in context with nil value")
	}

	return connInfo
}
