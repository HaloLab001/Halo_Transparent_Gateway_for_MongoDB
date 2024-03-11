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

// Package session provides access to session registry.
package session

import (
	"time"

	"github.com/google/uuid"

	"github.com/FerretDB/FerretDB/internal/types"
)

// Session represents a session.
type Session struct {
	id       types.Binary
	lastUsed time.Time
	expired  bool
}

// newSession returns a new session.
func newSession(id uuid.UUID) *Session {
	sessionID := types.Binary{Subtype: types.BinaryUUID, B: id[:]}
	return &Session{
		id:       sessionID,
		lastUsed: time.Now(),
	}
}