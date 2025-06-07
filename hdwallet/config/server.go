package config

type Cmd struct {
	Server struct {
		Grpc struct {
			Enabled     bool
			Address     string
			TlsCertFile string
			TlsKeyFile  string
		}
	}
	Client struct {
		Cli struct {
			Enabled bool
		}
		Grpc struct {
			Enabled bool
			Address string
			CAFile  string
		}
	}
}

type AppMode string

const (
	AppModeServer AppMode = "SERVER"
	AppModeClient AppMode = "CLIENT"
)

type Application struct {
	Debug bool
	Mode  AppMode
	Http  struct {
		Address string
	}
	AvailableCoins []string
	Salt           string
	Cmd            Cmd
}

func (c *Application) Initialize() {
	c.Mode = AppModeServer
}
