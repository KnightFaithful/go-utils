package testutil

import (
	"context"
	"example.com/m/util/mysqlclient"
	"example.com/m/util/redisclient"
	"gorm.io/gorm"
)

type NewContextRequest struct {
	Host   string
	Cookie string
}

type ContextKey string

const (
	ContextKeyNeedLog  ContextKey = "Key_need_log"
	ContextSQLCommon   ContextKey = "Key_sql_common"
	ContextRedisCommon ContextKey = "Key_redis_common"
)

func NewContext(req NewContextRequest) context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, ContextKeyNeedLog, true)
	ctx = bindMySQLContext(ctx, mysqlclient.InitDB(&mysqlclient.MySQLConfig{
		Host:         GetStringConfig(ctx, ModuleMysql, ConfigMysqlHost),
		Port:         GetInt64Config(ctx, ModuleMysql, ConfigMysqlPort),
		User:         GetStringConfig(ctx, ModuleMysql, ConfigMysqlUser),
		Password:     GetStringConfig(ctx, ModuleMysql, ConfigMysqlPassword),
		DBName:       GetStringConfig(ctx, ModuleMysql, ConfigMysqlDatabase),
		Charset:      "utf8mb4",
		ParseTime:    true,
		Loc:          "Local",
		MaxIdleConns: 10,
		MaxOpenConns: 100,
	}))
	ctx = bindRedisContext(ctx, redisclient.NewRedisClient(GetStringConfig(ctx, ModuleRedis, ConfigRedisHost)))
	return ctx
}

func bindMySQLContext(ctx context.Context, db *gorm.DB) context.Context {
	return context.WithValue(ctx, ContextSQLCommon, db)
}

func bindRedisContext(ctx context.Context, db *redisclient.RedisClient) context.Context {
	return context.WithValue(ctx, ContextRedisCommon, db)
}

func GetDBCommon(ctx context.Context) *gorm.DB {
	return ctx.Value(ContextSQLCommon).(*gorm.DB)
}

func GetRedisCommon(ctx context.Context) *redisclient.RedisClient {
	return ctx.Value(ContextRedisCommon).(*redisclient.RedisClient)
}

func GetValueByCtxString(ctx context.Context, key ContextKey) string {
	return ctx.Value(key).(string)
}

func GetValueByCtxInt64(ctx context.Context, key ContextKey) int64 {
	return ctx.Value(key).(int64)
}

func GetValueByCtxBool(ctx context.Context, key ContextKey) bool {
	return ctx.Value(key).(bool)
}

func SetInt64ValueToCtx(ctx context.Context, key ContextKey, value int64) context.Context {
	return context.WithValue(ctx, key, value)
}

func SetBoolValueToCtx(ctx context.Context, key ContextKey, value bool) context.Context {
	return context.WithValue(ctx, key, value)
}

func DisableLog(ctx context.Context) context.Context {
	ctx = SetBoolValueToCtx(ctx, ContextKeyNeedLog, false)
	return ctx
}

func IsEnableLog(ctx context.Context) bool {
	return GetValueByCtxBool(ctx, ContextKeyNeedLog)
}
