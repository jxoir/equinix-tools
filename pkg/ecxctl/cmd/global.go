// Copyright © 2018 Juan Manuel Irigaray <jirigaray@gmail.com>
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

package cmd

// GlobalFlags are flags defined globally and are inherited to all sub-commands
type GlobalFlags struct {
	UserName               string
	UserPassword           string
	APICredentialGrantType string
	EcxAPIHost             string
	PlaygroundToken        string
	PlaygroundAPIEndpoint  string

	Debug bool
	NoSSL bool

	EquinixAPISecret string
	EquinixAPIId     string
}
