package global

import (
	"database/sql"
	"backend-go/pkg/setting"
)

var (
	DB     *sql.DB
	Config *setting.Config
) 