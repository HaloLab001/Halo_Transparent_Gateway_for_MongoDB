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

package backend

type Backend interface {
	Database(params *DatabaseParams) Database
}

func BackendContract(b Backend) Backend {
	return &backendContract{
		b: b,
	}
}

type backendContract struct {
	b Backend
}

type DatabaseParams struct{}

func (bc *backendContract) Database(params *DatabaseParams) Database {
	return bc.b.Database(params)
}

// check interfaces
var (
	_ Backend = (*backendContract)(nil)
)
