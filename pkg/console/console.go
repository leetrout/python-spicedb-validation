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

package console

import (
	"fmt"
	"os"
)

// Printf defines an (overridable) function for printing to the console via stdout.
var Printf = func(format string, a ...any) {
	fmt.Printf(format, a...)
}

// Errorf defines an (overridable) function for printing to the console via stderr.
var Errorf = func(format string, a ...any) {
	_, err := fmt.Fprintf(os.Stderr, format, a...)
	if err != nil {
		panic(err)
	}
}

// Println prints a line with optional values to the console.
func Println(values ...any) {
	for _, value := range values {
		Printf("%v\n", value)
	}
}
