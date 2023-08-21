package config

import (
	"github.com/jesse0michael/pulse/internal/server"
	"github.com/jesse0michael/pulse/internal/service"
)

type Config struct {
	Server  server.Config
	Service service.Config
}
