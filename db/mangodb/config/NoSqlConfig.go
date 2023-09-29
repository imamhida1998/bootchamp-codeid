package config

const (
	DEV     = "DEV"
	STAGING = "STAGING"
	PROD    = "PROD"

	ARIGO         = "ARIGO"
	AHMAD         = "AHMAD"
	HAIKAL        = "HAIKAL"
	IMAM          = "IMAM"
	DATABASEMONGO = "finaltest"
	CURRENT_PHASE = DEV
	ARIGO_PHASE   = ARIGO
	HAIKAL_PHASE  = HAIKAL
	IMAM_MANGODB  = IMAM
)

type MongoConfig struct {
	Host           string
	Port           string
	User           string
	Pwd            string
	Authentication bool
}
var MONGO_CONFIGS map[string]MongoConfig = map[string]MongoConfig{
	DEV: {

		Host: "localhost",
		Port: "27017",
	},
	IMAM: {

		Host: "localhost",
		Port: "27017",
	},
}
