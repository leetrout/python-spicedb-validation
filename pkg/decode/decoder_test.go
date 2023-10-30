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

package decode

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRewriteURL(t *testing.T) {
	tests := []struct {
		name string
		in   url.URL
		out  url.URL
	}{
		{
			name: "gist",
			in: url.URL{
				Scheme: "https",
				Host:   "gist.github.com",
				Path:   "/ecordell/9e2110ac4a1292b899784ed809d44b8f",
			},
			out: url.URL{
				Scheme: "https",
				Host:   "gist.githubusercontent.com",
				Path:   "/ecordell/9e2110ac4a1292b899784ed809d44b8f/raw",
			},
		},
		{
			name: "playground schema",
			in: url.URL{
				Scheme: "https",
				Host:   "play.authzed.com",
				Path:   "/s/KY7TEKLs5_9R/schema",
			},
			out: url.URL{
				Scheme: "https",
				Host:   "play.authzed.com",
				Path:   "/s/KY7TEKLs5_9R/download",
			},
		},
		{
			name: "playground relationships",
			in: url.URL{
				Scheme: "https",
				Host:   "play.authzed.com",
				Path:   "/s/KY7TEKLs5_9R/relationships",
			},
			out: url.URL{
				Scheme: "https",
				Host:   "play.authzed.com",
				Path:   "/s/KY7TEKLs5_9R/download",
			},
		},
		{
			name: "playground assertions",
			in: url.URL{
				Scheme: "https",
				Host:   "play.authzed.com",
				Path:   "/s/KY7TEKLs5_9R/assertions",
			},
			out: url.URL{
				Scheme: "https",
				Host:   "play.authzed.com",
				Path:   "/s/KY7TEKLs5_9R/download",
			},
		},
		{
			name: "playground expected",
			in: url.URL{
				Scheme: "https",
				Host:   "play.authzed.com",
				Path:   "/s/KY7TEKLs5_9R/expected",
			},
			out: url.URL{
				Scheme: "https",
				Host:   "play.authzed.com",
				Path:   "/s/KY7TEKLs5_9R/download",
			},
		},
		{
			name: "pastebin",
			in: url.URL{
				Scheme: "https",
				Host:   "pastebin.com",
				Path:   "/LuCwwBwU",
			},
			out: url.URL{
				Scheme: "https",
				Host:   "pastebin.com",
				Path:   "/raw/LuCwwBwU",
			},
		},
		{
			name: "pastebin raw",
			in: url.URL{
				Scheme: "https",
				Host:   "pastebin.com",
				Path:   "/raw/LuCwwBwU",
			},
			out: url.URL{
				Scheme: "https",
				Host:   "pastebin.com",
				Path:   "/raw/LuCwwBwU",
			},
		},
		{
			name: "direct",
			in: url.URL{
				Scheme: "https",
				Host:   "somethingelse.com",
				Path:   "/any/path",
			},
			out: url.URL{
				Scheme: "https",
				Host:   "somethingelse.com",
				Path:   "/any/path",
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			rewriteURL(&tt.in)
			require.EqualValues(t, tt.out, tt.in)
		})
	}
}
