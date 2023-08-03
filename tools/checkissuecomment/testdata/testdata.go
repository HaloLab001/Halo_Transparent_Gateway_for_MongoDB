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

// testing for issue lint
// TODO https://github.com/FerretDB/FerretDB/issues/123
func testCorrect() {
	fmt.Println("checking the lint that should be found which is valid one")
}

// testing for issue lint
func testIncorrect() {
	fmt.Println("checking the lint which has only comment")
}

// TODO https://github.com/FerretDB/FerretDB/issues/123
func testIncorrectOnlycomment(v any) {
	fmt.Println("checking the lint which has only issue linbk as comment")
}
