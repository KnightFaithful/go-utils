package testutil

type Module = string

const (
	ModuleEs    Module = "es"
	ModuleMysql Module = "mysql"
	ModuleRedis Module = "redis"
)

const (
	ConfigESHost = "host"

	ConfigMysqlHost     = "host"
	ConfigMysqlPort     = "port"
	ConfigMysqlUser     = "user"
	ConfigMysqlPassword = "password"
	ConfigMysqlDatabase = "db"

	ConfigRedisHost = "host"
)

type CID = string

const (
	CID_ID CID = "ID"
	CID_MY CID = "MY"
	CID_PH CID = "PH"
	CID_SG CID = "SG"
	CID_TH CID = "TH"
	CID_TW CID = "TW"
	CID_VN CID = "VN"
)

type ENV = string

const (
	ENV_TEST ENV = "TEST"
)
