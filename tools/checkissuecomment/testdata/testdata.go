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

// Package testdata provides vet tool test data.
package testdata

import (
	"fmt"
)

// TODO https://github.com/FerretDB/FerretDB/issues/2275
func testcorrect() {
	// TODO https://github.com/FerretDB/FerretDB/issues/2612
	fmt.Println("checking the lint which has proper TODO")
}

// issue to be resolved below
// TODO https://github.com/FerretDB/FerretDB/issues/123
func testcorrectTODOcomment(v any) {
	fmt.Println("checking the lint which has proper TODO with comment")
}
