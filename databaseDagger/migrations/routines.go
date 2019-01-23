package migrations

import (
	db "github.com/andreweggleston/DeathByDagger/databaseDagger"
	"github.com/andreweggleston/DeathByDagger/models/player"
	"github.com/sirupsen/logrus"
)

var migrationRoutines = map[uint64]func(){
	5:  updateAllPlayerInfo,
	6:  truncateHTTPSessions,
	8:  setPlayerExternalLinks,
	9:  setPlayerSettings,
	10: dropTableSessions,
	11: dropColumnUpdatedAt,
	13: dropUnusedColumns,
}

func updateAllPlayerInfo() {
	var players []*player.Player
	db.DB.Model(&player.Player{}).Find(&players)

	for _, player := range players {
		player.Save()
	}
}

func truncateHTTPSessions() {
	db.DB.Exec("TRUNCATE TABLE http_sessions")
}

func setPlayerExternalLinks() {
	var players []*player.Player
	db.DB.Model(&player.Player{}).Find(&players)

	for _, player := range players {
		player.Save()
	}
}

// move player_settings values to player.Settings hstore
func setPlayerSettings() {
	rows, err := db.DB.DB().Query("SELECT player_id, key, value FROM player_settings")
	if err != nil {
		logrus.Fatal(err)
	}
	for rows.Next() {
		var playerID uint
		var key, value string

		rows.Scan(&playerID, &key, &value)
		p, _ := player.GetPlayerByID(playerID)
		p.SetSetting(key, value)
	}

	db.DB.Exec("DROP TABLE player_settings")
}

func dropTableSessions() {
	db.DB.Exec("DROP TABLE http_sessions")
}

func dropColumnUpdatedAt() {
	db.DB.Exec("ALTER TABLE players DROP COLUMN updated_at")
}

func dropUnusedColumns() {
	db.DB.Model(&player.Player{}).DropColumn("debug")
}
