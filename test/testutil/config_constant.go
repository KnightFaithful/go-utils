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
