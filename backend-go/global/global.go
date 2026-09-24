package global

import (
	"database/sql"
	"backend-go/pkg/setting"
	"github.com/redis/go-redis/v9"
)

var (
	DB     *sql.DB
	Config *setting.Config
	Redis *redis.Client
) 