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

package registry

import (
	"github.com/FerretDB/FerretDB/internal/backends/sqlite"
	"github.com/FerretDB/FerretDB/internal/handlers"
	handler "github.com/FerretDB/FerretDB/internal/handlers/sqlite"
)

// init registers "sqlite" handler.
func init() {
	registry["sqlite"] = func(opts *NewHandlerOpts) (handlers.Interface, error) {
		b, err := sqlite.NewBackend(&sqlite.NewBackendParams{
			URI: opts.SQLiteURL,
			L:   opts.Logger.Named("sqlite"),
			P:   opts.StateProvider,
		})
		if err != nil {
			return nil, err
		}

		handlerOpts := &handler.NewOpts{
			Backend: b,
			URI:     opts.SQLiteURL,

			L:             opts.Logger.Named("sqlite"),
			ConnMetrics:   opts.ConnMetrics,
			StateProvider: opts.StateProvider,

			DisableFilterPushdown:    opts.DisableFilterPushdown,
			EnableUnsafeSortPushdown: opts.EnableUnsafeSortPushdown,
			EnableOplog:              opts.EnableOplog,
		}

		return handler.New(handlerOpts)
	}
}
