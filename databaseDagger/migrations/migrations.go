package migrations

import (
	"github.com/andreweggleston/DeathByDagger/databaseDagger"
	"github.com/andreweggleston/DeathByDagger/models/player"
	"sync"
)

var once = new(sync.Once)

func Do() {
	databaseDagger.DB.Exec("CREATE EXTENSION IF NOT EXISTS hstore")
	databaseDagger.DB.AutoMigrate(&player.Player{})
	databaseDagger.DB.AutoMigrate(&Constant{})

	once.Do(checkSchema)
}
