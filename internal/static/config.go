// Package static contains static "text" files.
package static

import _ "embed"

// GoExampleConfig is the config used within goreleaser init.
//
//go:embed ease-example.yaml
var EaseExampleConfig []byte
