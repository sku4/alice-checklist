package version

// Version is the application build version. It is injected at build time via
// ldflags (e.g. -ldflags "-X github.com/sku4/alice-checklist/pkg/version.Version=<sha>")
// and defaults to "dev" for local builds. It is logged on startup so a freshly
// deployed version can be told apart from a running older one.
var Version = "dev"
