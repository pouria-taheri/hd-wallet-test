package cmd

import (
	"fmt"
	"git.mazdax.tech/blockchain/hdwallet/cardano"
	cd "git.mazdax.tech/blockchain/hdwallet/cardano/domain"
	"git.mazdax.tech/blockchain/hdwallet/cmd/cli/cliclient"
	"git.mazdax.tech/blockchain/hdwallet/cmd/domain"
	"git.mazdax.tech/blockchain/hdwallet/cmd/grpc/grpcclient"
	grpcServer "git.mazdax.tech/blockchain/hdwallet/cmd/grpc/server"
	"git.mazdax.tech/blockchain/hdwallet/cmd/usecase"
	"git.mazdax.tech/blockchain/hdwallet/coins"
	"git.mazdax.tech/blockchain/hdwallet/config"
	"git.mazdax.tech/blockchain/hdwallet/manager"
	md "git.mazdax.tech/blockchain/hdwallet/manager/domain"
	"git.mazdax.tech/blockchain/hdwallet/monitoring"
	"git.mazdax.tech/blockchain/hdwallet/swagger"
	"git.mazdax.tech/data-layer/loggercore"
	"git.mazdax.tech/data-layer/loggercore/logger"
	"git.mazdax.tech/delivery/ginger"
	"git.mazdax.tech/delivery/handlercore"
	"math/rand"
	"os"
	"strings"
	"time"
)

var app *domain.Application

func Prepare(rootDirectory string) {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "./config.yaml"
	}
	// initialize
	rand.Seed(time.Now().UTC().UnixNano())
	registry, serverConfig := config.ReadConfig(configPath)
	serverConfig.Initialize()
	app = &domain.Application{
		RootDirectory:  rootDirectory,
		ServerModel:    ginger.NewServer(registry),
		ConfigRegistry: registry,
		Config:         serverConfig,
	}
	// Logger
	var err error
	app.Logger, err = loggercore.GetLogger(app.ConfigRegistry.ValueOf("logger"))
	if err != nil {
		panic(err.Error())
	}
}

func initialize() {
	if app.Config.Cmd.Client.Cli.Enabled {
		// cmd client
		app.Client = cliclient.New(app, rootCli)
	}
	if app.Config.Mode == config.AppModeServer {
		app.SecuredConfigFileDetails = config.LoadSecureConfigFiles(app.ConfigRegistry, app.Config)
	}
	if app.Config.Mode == config.AppModeClient {
		grpcLogger := app.Logger.With(logger.Field{
			"module": "grpc client",
		})
		grpcClient := grpcclient.New(app.Config.Cmd, grpcLogger)
		grpcClient.SetInnerClient(app.Client)
		grpcClient.Initialize()

		app.Client = grpcClient
		app.Handler = usecase.NewRemoteHandler(app, app.Client)
	} else {
		app.Handler = usecase.NewHandler(app, app.Client)
	}
}

func run() {
	if app.Config.Mode != config.AppModeServer {
		app.Logger.InfoF("application is not running in server mode")
		return
	}
	app.SecuredConfigFileDetails = config.LoadSecureConfigFiles(app.ConfigRegistry, app.Config)
	if config.IsAnyConfigSecured(app.SecuredConfigFileDetails) {
		app.Logger.InfoF("Config has been secured")

		var srv grpcServer.Model
		if app.Config.Cmd.Server.Grpc.Enabled {
			app.Logger.InfoF("initializing grpc server...")
			grpcServerLogger := app.Logger.With(logger.Field{
				"module": "grpc server",
			})
			srv = grpcServer.NewServer(grpcServerLogger, app.Handler,
				app.Config.Cmd)
			go srv.Start()
			app.Logger.InfoF("waiting for unlock config...")
			app.Handler.WaitForUnlock()
			srv.Stop()
		}
	} else {
		app.Handler.LoadSecureConfigs()
	}
	app.Handler.EnsureSecureConfigLoaded()
	//
	var err error
	// router group
	rg := app.NewRouterGroup("/")
	ig := app.NewRouterGroup("/internal/")
	monitoringGroup := app.NewRouterGroup("/monitoring")

	// account
	accountLogger := app.Logger.With(logger.Field{
		"source": "account",
	})
	var accountManagers = make(map[string]md.AccountManagerModel)
	accountHandler := manager.NewAccountHandler(initHandlerModel())
	for _, coin := range app.Config.AvailableCoins {
		coin = strings.ToLower(coin)
		accountManager := coins.DecideAccountManagerUseCase(coin,
			app.ConfigRegistry, app.SecureConfigs[coin], accountLogger.With(logger.Field{
				"module": strings.ToUpper(coin),
			}))
		if accountManager == nil {
			panic(fmt.Sprintf("account manager for %s not found", coin))
		}
		accountManagers[coin] = accountManager
		accountHandler.RegisterAccountHandler(accountManager)
	}
	handlercore.RegisterRouters(rg, "/assets/:asset_type/account",
		handlercore.ActionCreate, accountHandler)

	// cardano wallet routes
	if accMgr, ok := accountManagers["ada"]; ok {
		cardanoGetAccountsHandler := cardano.NewGetAccountsHandler(initHandlerModel(),
			accMgr.(cd.CardanoWalletModel))
		cardanoWalletHandler := cardano.NewGetWalletHandler(initHandlerModel(),
			accMgr.(cd.CardanoWalletModel))
		cardanoNewAddressHandler := cardano.NewAddressHandler(initHandlerModel(),
			accMgr.(cd.CardanoWalletModel))
		handlercore.RegisterRouters(ig, "/ada/accounts", handlercore.ActionCreate, cardanoGetAccountsHandler)
		handlercore.RegisterRouters(ig, "/ada/wallet", handlercore.ActionCreate, cardanoWalletHandler)
		handlercore.RegisterRouters(ig, "/ada/new_address", handlercore.ActionCreate, cardanoNewAddressHandler)
	}

	// sign
	signHandler := manager.NewSignHandler(initHandlerModel())
	for _, coin := range app.Config.AvailableCoins {
		signerLogger := app.Logger.With(logger.Field{
			"source": "signer",
		})
		signer := coins.DecideSignerUseCase(coin, signerLogger.With(logger.Field{
			"module": strings.ToUpper(coin),
		}), app.ConfigRegistry, accountManagers[coin])
		if signer == nil {
			panic(fmt.Sprintf("signer for %s not found", coin))
		}
		signHandler.RegisterSigner(signer)
	}
	handlercore.RegisterRouters(rg, "/assets/:asset_type/sign",
		handlercore.ActionCreate, signHandler)
	// health
	healthHandler := monitoring.NewHealthHandler(initHandlerModel())
	handlercore.RegisterRouters(monitoringGroup, "/health", handlercore.ActionRead, healthHandler)
	// docs
	if docs, ok := os.LookupEnv("DOCS"); ok && docs == "true" {
		// swagger
		swaggerRg := app.NewRouterGroup("/swagger")
		swagH := swagger.SwagHandler(initNoResponderHandlerModel())
		handlercore.RegisterRouters(swaggerRg, "/*any",
			handlercore.ActionRead, swagH)
	}
	// run
	err = app.Run(app.Config.Http.Address)
	if err != nil {
		panic(err)
	}
}

func initHandlerModel() (handlerModel handlercore.HandlerModel) {
	responder := ginger.NewResponder(nil)
	baseHandler := ginger.NewHandler(app.Logger)
	return handlercore.NewModel(baseHandler, responder)
}
func initNoResponderHandlerModel() (handlerModel handlercore.HandlerModel) {
	baseHandler := ginger.NewHandler(app.Logger)
	responder := ginger.NewEmptyResponder(nil)
	return handlercore.NewModel(baseHandler, responder)
}
