package config

import (
	"os"
	"time"

	"github.com/ethereum/go-ethereum/log"
	"github.com/spf13/viper"
)

// Default values of configuration options
const (
	// this defines default application name
	defApplicationName = "Artion"

	// EmptyAddress defines an empty address
	EmptyAddress = "0x0000000000000000000000000000000000000000"

	// defServerBind holds default API server binding address
	defServerBind = "localhost:16761"

	// default set of timeouts for the server
	defReadTimeout     = 2
	defWriteTimeout    = 15
	defIdleTimeout     = 1
	defHeaderTimeout   = 1
	defResolverTimeout = 30
	defMaxParserMemory = int64(10 * 1024 * 1024) // 10MB

	// defLoggingLevel holds default Logging level
	// See `godoc.org/github.com/op/go-logging` for the full format specification
	// See `golang.org/pkg/time/` for time format specification
	defLoggingLevel = "INFO"

	// defLoggingFormat holds default format of the Logger output
	defLoggingFormat = "%{color}%{level:-8s} %{module}/%{shortfile}::%{shortfunc}()%{color:reset}: %{message}"

	// defLachesisUrl holds default Lachesis connection string
	defLachesisUrl = "~/.lachesis/data/lachesis.ipc"

	// defIpfsUrl holds default IPFS connection string
	defIpfsUrl = "localhost:5001"

	defIpfsFileCacheDir = "/tmp/artion/"

	// defSkipHttpGateways tells whether to skip known HTTP-to-IPFS gateways
	defSkipHttpGateways = true

	// defMongoUrl holds default MongoDB connection string for local database
	defMongoUrl = "mongodb://localhost:27017"

	// defMongoDatabase holds the default name of the MongoDB local database
	defMongoDatabase = "artion"

	// defSharedMongoUrl holds default MongoDB connection string for shared/replicated database
	defSharedMongoUrl = "mongodb://localhost:27017"

	// defSharedMongoDatabase holds the default name of the shared/replicated MongoDB database
	defSharedMongoDatabase = "artionshared"

	// defCacheEvictionTime holds default time for in-memory eviction periods
	defCacheEvictionTime = 15 * time.Minute

	// defCacheMax size represents the default max size of the cache in MB
	defCacheMaxSize = 2048

	// defApiStateOrigin represents the default origin used for API state syncing
	defApiStateOrigin = "https://localhost"

	// defAuthBearerSecret holds the default bearer secret
	defAuthBearerSecret = "0x0123456789"

	// defAuthNonceSecret holds the default nonce secret
	defAuthNonceSecret = "0xABCDEF"

	// defWrappedFTMContract is the default address of the wFTM contract
	defWrappedFTMContract = "0x21be370d5312f44cb42ce377bc9b8a0cef1a4c83"
)

// defCorsAllowOrigins holds CORS default allowed origins.
var defCorsAllowOrigins = []string{"*"}

// applyDefaults sets default values for configuration options.
func applyDefaults(cfg *viper.Viper) {
	// set simple details
	cfg.SetDefault(keyAppName, defApplicationName)

	port := os.Getenv("PORT")
	cfg.SetDefault(keyBindAddress, "0.0.0.0:"+port)
	cfg.SetDefault(keyLoggingLevel, defLoggingLevel)
	cfg.SetDefault(keyLoggingFormat, defLoggingFormat)
	cfg.SetDefault(keyLachesisUrl, getEnv(envLachesisUrl, defLachesisUrl))
	cfg.SetDefault(keyIpfsUrl, getEnv(envIpfsUrl, defIpfsUrl))
	cfg.SetDefault(keySkipHttpGateways, defSkipHttpGateways)

	cfg.SetDefault(keyIpfsGateway, getEnv(envIpfsGateway, ""))
	cfg.SetDefault(keyIpfsGatewayBearer, getEnv(envIpfsGatewayBearer, ""))
	cfg.SetDefault(keyPinataJwt, getEnv(envPinataJwt, ""))
	cfg.SetDefault(keyIpfsFileCacheDir, getEnv(envIpfsFileCacheDir, defIpfsFileCacheDir))

	cfg.SetDefault(keyMongoUrl, getEnv(envMongoUrl, defMongoUrl))
	cfg.SetDefault(keyMongoDatabase, defMongoDatabase)
	cfg.SetDefault(keySharedMongoUrl, getEnv(envMongoUrl, defMongoUrl))
	cfg.SetDefault(keySharedMongoDatabase, defSharedMongoDatabase)
	cfg.SetDefault(keyApiStateOrigin, defApiStateOrigin)

	// in-memory cache
	cfg.SetDefault(keyCacheEvictionTime, defCacheEvictionTime)
	cfg.SetDefault(keyCacheMaxSize, defCacheMaxSize)

	// server related
	cfg.SetDefault(keyTimeoutRead, defReadTimeout)
	cfg.SetDefault(keyTimeoutWrite, defWriteTimeout)
	cfg.SetDefault(keyTimeoutHeader, defHeaderTimeout)
	cfg.SetDefault(keyTimeoutIdle, defIdleTimeout)
	cfg.SetDefault(keyTimeoutResolver, defResolverTimeout)
	cfg.SetDefault(keyMaxParserMemory, defMaxParserMemory)

	// cors
	cfg.SetDefault(keyCorsAllowOrigins, defCorsAllowOrigins)

	// auth
	cfg.SetDefault(keyAuthBearerSecret, getEnv(envAuthBearerSecret, defAuthBearerSecret))
	cfg.SetDefault(keyAuthNonceSecret, getEnv(envAuthNonceSecret, defAuthNonceSecret))

	// contracts
	cfg.SetDefault(keyWrappedFTM, defWrappedFTMContract)
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	log.Debug("getEnv", "key", key, "value", value)
	if len(value) == 0 {
		return fallback
	}
	return value
}
