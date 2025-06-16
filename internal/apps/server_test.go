package apps

import (
	"testing"

	"github.com/sbilibin2017/yandex-go-advanced/internal/configs"
	"github.com/stretchr/testify/assert"
)

func TestNewServerApp(t *testing.T) {
	config := &configs.ServerConfig{
		Address:  ":8080",
		LogLevel: "debug",
	}

	server, err := NewServerApp(config)

	assert.NoError(t, err)
	assert.NotNil(t, server)
	assert.Equal(t, config.Address, server.Addr)
	assert.NotNil(t, server.Handler)
}
