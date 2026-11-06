package release

import (
	"embed"
	"strings"
)

//go:embed VERSION CHANGELOG.md UPGRADE.md
var files embed.FS

var Commit = "development"
var BuildTime = "unknown"

func Text(name string) string { b, _ := files.ReadFile(name); return strings.TrimSpace(string(b)) }
