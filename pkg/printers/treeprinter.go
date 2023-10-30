// Copyright 2023 Authzed, Inc.
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//        http://www.apache.org/licenses/LICENSE-2.0
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

package printers

import (
	"strings"

	"github.com/leetrout/python-spicedb-validation/pkg/console"
	"github.com/xlab/treeprint"
)

type TreePrinter struct {
	tree treeprint.Tree
}

func NewTreePrinter() *TreePrinter {
	return &TreePrinter{}
}

func (tp *TreePrinter) Child(val string) *TreePrinter {
	if tp.tree == nil {
		tp.tree = treeprint.NewWithRoot(val)
		return tp
	}
	return &TreePrinter{tree: tp.tree.AddBranch(val)}
}

func (tp *TreePrinter) Print() {
	console.Println(tp.String())
}

func (tp *TreePrinter) PrintIndented() {
	lines := strings.Split(tp.String(), "\n")
	indentedLines := make([]string, 0, len(lines))
	for _, line := range lines {
		indentedLines = append(indentedLines, "  "+line)
	}

	console.Println(strings.Join(indentedLines, "\n"))
}

func (tp *TreePrinter) String() string {
	return tp.tree.String()
}
