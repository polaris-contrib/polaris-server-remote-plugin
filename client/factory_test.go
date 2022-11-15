package client

import (
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	patches := gomonkey.NewPatches()
	patches = patches.ApplyFunc(newClient, func(config *Config) (*clientImpl, error) {
		return &clientImpl{
			pluginName: config.Name,
		}, nil
	})

	patches = patches.ApplyMethod(reflect.TypeOf(&clientImpl{}), "Check", func(_ *clientImpl) error {
		return nil
	})

	defer patches.Reset()

	var c1, c2, c3 Client
	var err1, err2, err3 error
	t.Run("注册 client-1", func(t *testing.T) {
		c1, err1 = Register(&Config{Name: "client-1"})
		assert.Nil(t, err1)
		assert.NotNil(t, c1)
	})

	t.Run("register client-2 first time", func(t *testing.T) {
		c2, err2 = Register(&Config{Name: "client-2"})
		assert.Nil(t, err2)
		assert.NotNil(t, c2)
	})

	t.Run("register client-2 second time", func(t *testing.T) {
		c3, err3 = Register(&Config{Name: "client-2"})
		assert.Nil(t, err3)
		assert.Equal(t, c2, c3)
		assert.NotEqual(t, c1, c3)
	})
}
