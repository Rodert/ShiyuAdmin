package system

import (
	"github.com/gin-gonic/gin"
	"shiyu-admin-backend/internal/bootstrap"
	"shiyu-admin-backend/internal/release"
	"shiyu-admin-backend/pkg/response"
)

func versionInfo(c *gin.Context) {
	response.Success(c, gin.H{"version": release.Text("VERSION"), "commit": release.Commit,
		"build_time": release.BuildTime, "schema_version": bootstrap.SchemaVersion,
		"release_notes": release.Text("CHANGELOG.md"), "upgrade_guide": release.Text("UPGRADE.md")})
}
