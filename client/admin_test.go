package client

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"testing"

	"github.com/emicklei/go-restful/v3"

	"github.com/polaris-contrib/polaris-server-remote-plugin-common/log"
)

func TestMain(m *testing.M) {
	cmd := &exec.Cmd{
		Path: "/usr/bin/make",
		Args: append([]string{"/usr/bin/make"}, "build"),
		Dir:  "../",
	}

	err := cmd.Run()
	if err != nil {
		log.Fatal("got error", "error", err.Error())
		os.Exit(-1)
	}

	m.Run()
}

// TestAdminAPI 测试 API server
func TestAdminAPI(t *testing.T) {
	if _, err := Register(
		&Config{
			Name: "remote-rate-limit-server-v1",
			Mode: RumModelLocal,
			Local: LocalConfig{
				MaxProcs: 1,
				Path:     "../remote-rate-limit-server-v1",
			},
		},
	); err != nil {
		log.Fatal("server-v1 register failed", "error", err.Error())
		return
	}

	if _, err := Register(
		&Config{
			Name: "remote-rate-limit-server-v2",
			Mode: RumModelLocal,
			Local: LocalConfig{
				Path: "../remote-rate-limit-server-v2",
			},
		},
	); err != nil {
		log.Fatal("server-v2 register failed", "error", err.Error())
	}

	restful.DefaultContainer.Add(NewResource().WebService())
	adminPort := 9050

	log.Info(fmt.Sprintf("request the admin api using http://localhost:%d", adminPort))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", adminPort), nil); err != nil {
		log.Fatal("plugin admin serve error", "error", err)
	}
}
