package testutil

import (
	"context"
	"example.com/m/util/fileutil"
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var (
	doOnce = &sync.Once{}
)

func initConfig(ctx context.Context) {
	doOnce.Do(func() {
		initConfPath(ctx)
		viper.SetConfigName("config")                            // 文件名（不带扩展名）
		viper.SetConfigType("yaml")                              // 文件类型
		viper.AddConfigPath(strings.Join([]string{"conf"}, "/")) // 搜索路径

		if err := viper.ReadInConfig(); err != nil {
			panic(fmt.Errorf("读取配置失败: %w", err))
		}
	})
}

func initConfPath(ctx context.Context) {
	pwd, _ := os.Getwd()
	homePath := pwd[0:strings.LastIndex(pwd, "go-utils")] + "go-utils"
	if _, fErr := fileutil.CopyFile(
		filepath.Join(homePath, "test/testutil/conf", "config.yaml"), filepath.Join(pwd, "conf", "config.yaml")); fErr != nil {
		fmt.Println("copy chassis.yaml fail")
		panic(fErr)
	}

}

func GetStringConfig(ctx context.Context, module, key string) string {
	initConfig(ctx)
	config := viper.GetString(strings.Join([]string{module, key}, "."))
	return config
}

func GetInt64Config(ctx context.Context, module, key string) int64 {
	initConfig(ctx)
	config := viper.GetInt64(strings.Join([]string{module, key}, "."))
	return config
}
