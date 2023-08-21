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

// Package main contains linter for todo issue comments.
package main

import (
	"bufio"
	"os"
	"regexp"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/singlechecker"
)

// todoIssueRegex matches potential TODO representations of an issue.
var todoIssueRegex = regexp.MustCompile(`(?i)(?:issue[s]?|TODO|#):?\/?\s*\d+`)

// validTodoIssueRegex matches correct TODO comment with the issue number.
var validTodoIssueRegex = regexp.MustCompile(`^(.*?)\/\/\s*TODO https://github\.com/FerretDB/FerretDB/issues/\d+$`)

var analyzer = &analysis.Analyzer{
	Name: "checkissuecomment",
	Doc:  "check for TODO comments with issue links",
	Run:  run,
}

func main() {
	singlechecker.Main(analyzer)
}

// run analyses the presence of TODO issue in code.
func run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		fileName := pass.Fset.File(file.Pos()).Name()
		f, err := os.Open(fileName)
		if err != nil {
			return nil, err
		}

		defer f.Close()

		scanner := bufio.NewScanner(f)
		lineNumber := 1

		for scanner.Scan() {
			line := scanner.Text()

			if strings.Contains(line, "// TODO") {
				if todoIssueRegex.MatchString(line) &&
					!validTodoIssueRegex.MatchString(line) {
					pass.Reportf(file.Pos(), "TODO comments mentioning issues not satisfied the pattern (// TODO <issue_url>): %s on line number: %d", strings.TrimSpace(line), lineNumber)
				}
			}

			lineNumber++
		}

		if err := scanner.Err(); err != nil {
			return nil, err
		}
	}

	return nil, nil
}
