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
	Cid    string
	Env    string
}

type ContextKey string

const (
	ContextKeyHost           ContextKey = "key_host"
	ContextKeyCookie         ContextKey = "key_cookie"
	ContextKeyWorkforceDbId  ContextKey = "key_workforce_db_id"
	ContextKeyNetworkDbId    ContextKey = "key_network_db_id"
	ContextKeyBasicDbId      ContextKey = "key_basic_db_id"
	ContextKeyCid            ContextKey = "key_cid"
	ContextKeyClockInBuffer  ContextKey = "Key_clock_in_buffer"
	ContextKeyClockOutBuffer ContextKey = "Key_clock_out_buffer"
	ContextKeyNeedLog        ContextKey = "Key_need_log"
	ContextKeyEnv            ContextKey = "Key_env"
	ContextSQLCommon         ContextKey = "Key_sql_common"
	ContextRedisCommon       ContextKey = "Key_redis_common"
)

var cid2WorkforceDbIdMap = map[string]int64{
	CID_ID: 2043,
	CID_MY: 2044,
	CID_PH: 2045,
	CID_SG: 2046,
	CID_TH: 2047,
	CID_TW: 2048,
	CID_VN: 2049,
}

var cid2NetworkDbIdMap = map[string]int64{
	CID_VN: 4895,
}

var cid2BasicDbIdMap = map[string]int64{
	CID_MY: 3301,
}

func NewContext(req NewContextRequest) context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, ContextKeyHost, req.Host)
	ctx = context.WithValue(ctx, ContextKeyCookie, req.Cookie)
	ctx = context.WithValue(ctx, ContextKeyWorkforceDbId, cid2WorkforceDbIdMap[req.Cid])
	ctx = context.WithValue(ctx, ContextKeyNetworkDbId, cid2NetworkDbIdMap[req.Cid])
	ctx = context.WithValue(ctx, ContextKeyBasicDbId, cid2BasicDbIdMap[req.Cid])
	ctx = context.WithValue(ctx, ContextKeyCid, req.Cid)
	ctx = context.WithValue(ctx, ContextKeyNeedLog, true)
	ctx = context.WithValue(ctx, ContextKeyEnv, req.Env)
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

func GetCid(ctx context.Context) string {
	return GetValueByCtxString(ctx, ContextKeyCid)
}

func GetEnv(ctx context.Context) string {
	return GetValueByCtxString(ctx, ContextKeyEnv)
}

func DisableLog(ctx context.Context) context.Context {
	ctx = SetBoolValueToCtx(ctx, ContextKeyNeedLog, false)
	return ctx
}

func IsEnableLog(ctx context.Context) bool {
	return GetValueByCtxBool(ctx, ContextKeyNeedLog)
}
