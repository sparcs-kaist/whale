package main // import "github.com/sparcs-kaist/whale"

import (
	"github.com/sparcs-kaist/whale"
	"github.com/sparcs-kaist/whale/bolt"
	"github.com/sparcs-kaist/whale/cli"
	"github.com/sparcs-kaist/whale/cron"
	"github.com/sparcs-kaist/whale/crypto"
	"github.com/sparcs-kaist/whale/file"
	"github.com/sparcs-kaist/whale/http"
	"github.com/sparcs-kaist/whale/jwt"

	"log"
)

func initCLI() *whale.CLIFlags {
	var cli whale.CLIService = &cli.Service{}
	flags, err := cli.ParseFlags(whale.APIVersion)
	if err != nil {
		log.Fatal(err)
	}

	err = cli.ValidateFlags(flags)
	if err != nil {
		log.Fatal(err)
	}
	return flags
}

func initFileService(dataStorePath string) whale.FileService {
	fileService, err := file.NewService(dataStorePath, "")
	if err != nil {
		log.Fatal(err)
	}
	return fileService
}

func initStore(dataStorePath string) *bolt.Store {
	store, err := bolt.NewStore(dataStorePath)
	if err != nil {
		log.Fatal(err)
	}

	err = store.Open()
	if err != nil {
		log.Fatal(err)
	}

	err = store.MigrateData()
	if err != nil {
		log.Fatal(err)
	}
	return store
}

func initJWTService(authenticationEnabled bool) whale.JWTService {
	if authenticationEnabled {
		jwtService, err := jwt.NewService()
		if err != nil {
			log.Fatal(err)
		}
		return jwtService
	}
	return nil
}

func initCryptoService() whale.CryptoService {
	return &crypto.Service{}
}

func initEndpointWatcher(endpointService whale.EndpointService, externalEnpointFile string, syncInterval string) bool {
	authorizeEndpointMgmt := true
	if externalEnpointFile != "" {
		authorizeEndpointMgmt = false
		log.Println("Using external endpoint definition. Endpoint management via the API will be disabled.")
		endpointWatcher := cron.NewWatcher(endpointService, syncInterval)
		err := endpointWatcher.WatchEndpointFile(externalEnpointFile)
		if err != nil {
			log.Fatal(err)
		}
	}
	return authorizeEndpointMgmt
}

func initSettings(authorizeEndpointMgmt bool, flags *whale.CLIFlags) *whale.Settings {
	return &whale.Settings{
		HiddenLabels:       *flags.Labels,
		Logo:               *flags.Logo,
		Analytics:          !*flags.NoAnalytics,
		Authentication:     !*flags.NoAuth,
		EndpointManagement: authorizeEndpointMgmt,
	}
}

func retrieveFirstEndpointFromDatabase(endpointService whale.EndpointService) *whale.Endpoint {
	endpoints, err := endpointService.Endpoints()
	if err != nil {
		log.Fatal(err)
	}
	return &endpoints[0]
}

func main() {
	flags := initCLI()

	fileService := initFileService(*flags.Data)

	store := initStore(*flags.Data)
	defer store.Close()

	jwtService := initJWTService(!*flags.NoAuth)

	cryptoService := initCryptoService()

	authorizeEndpointMgmt := initEndpointWatcher(store.EndpointService, *flags.ExternalEndpoints, *flags.SyncInterval)

	settings := initSettings(authorizeEndpointMgmt, flags)

	if *flags.Endpoint != "" {
		var endpoints []whale.Endpoint
		endpoints, err := store.EndpointService.Endpoints()
		if err != nil {
			log.Fatal(err)
		}
		if len(endpoints) == 0 {
			endpoint := &whale.Endpoint{
				Name:          "primary",
				URL:           *flags.Endpoint,
				TLS:           *flags.TLSVerify,
				TLSCACertPath: *flags.TLSCacert,
				TLSCertPath:   *flags.TLSCert,
				TLSKeyPath:    *flags.TLSKey,
			}
			err = store.EndpointService.CreateEndpoint(endpoint)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			log.Println("Instance already has defined endpoints. Skipping the endpoint defined via CLI.")
		}
	}

	var server whale.Server = &http.Server{
		BindAddress:            *flags.Addr,
		AssetsPath:             *flags.Assets,
		Settings:               settings,
		TemplatesURL:           *flags.Templates,
		AuthDisabled:           *flags.NoAuth,
		SSOID:                  *flags.SSOID,
		SSOKey:                 *flags.SSOKey,
		EndpointManagement:     authorizeEndpointMgmt,
		UserService:            store.UserService,
		EndpointService:        store.EndpointService,
		ResourceControlService: store.ResourceControlService,
		CryptoService:          cryptoService,
		JWTService:             jwtService,
		FileService:            fileService,
	}

	log.Printf("Starting Whale on %s", *flags.Addr)
	err := server.Start()
	if err != nil {
		log.Fatal(err)
	}
}
