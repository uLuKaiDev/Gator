package app

import (
	"database/sql"

	"github.com/uLuKaiDev/Gator/internal/config"
	"github.com/uLuKaiDev/Gator/internal/database"
)

type State struct {
	DB     *database.Queries
	DBConn *sql.DB
	Config *config.Config
}
