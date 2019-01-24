package player

import (
	"errors"
	db "github.com/andreweggleston/DeathByDagger/databaseDagger"
	"github.com/andreweggleston/DeathByDagger/helpers/authority"
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/sirupsen/logrus"
	"time"
)

var ErrPlayerNotFound = errors.New("Player not found")
var ErrPlayerInReportedSlot = errors.New("Player in reported slot")

type Player struct {
	ID          uint   `gorm:"primary_key" json:"id"`
	Sub         string `sql:"not null;unique" json:"sub"`
	CSHUsername string `sql:"not null;unique" json:"cshusername"`
	Name        string `json:"name"`
	Email       string `sql:"not null" json:"email"`
	PhoneNumber string `json:"phonenumber"`
	SlackUserID string `json:"slackuserid"`

	Target           string    `json:"target"`
	Kills            uint      `sql:"not null" json:"kills"`
	CreatedAt        time.Time `json:"createdAt"`
	ProfileUpdatedAt time.Time `json:"-"`

	MarkedForDeath bool `sql:"not null" json:"markedfordeath"`
	Killed         bool `sql:"not null" json:"killed"`

	GlobalData

	Settings postgres.Hstore `json:"-"`

	Role authority.AuthRole `sql:"default:0" json:"-"`
}

type GlobalData struct {
	SafetyItem   string
	KillByDate   string
	Announcement string
}

func NewPlayer(cshusername string) (*Player, error) {
	player := &Player{CSHUsername: cshusername, Kills: 0}

	//check if admin todo

	last := &Player{}

	db.DB.Model(&Player{}).Last(last)

	return player, nil
}


func (p *Player) Alias() string {
	alias := p.GetSetting("siteAlias")
	if alias == "" {
		return p.Name
	}

	return alias
}

func (player *Player) Save() error {
	var err error
	if db.DB.NewRecord(player) {
		err = db.DB.Create(player).Error
	} else {
		err = db.DB.Save(player).Error
	}
	return err
}

func GetPlayerByID(ID uint) (*Player, error) {
	player := &Player{}

	if err := db.DB.First(player, ID).Error; err != nil {
		return nil, err
	}

	return player, nil
}

func GetPlayerByCSHUsername(cshusername string) (*Player, error) {
	var player = Player{}
	err := db.DB.Where("csh_username = ?", cshusername).First(&player).Error
	if err != nil {
		return nil, ErrPlayerNotFound
	}
	return &player, nil
}

func (player *Player) GetSetting(key string) string {
	if player.Settings == nil {
		return ""
	}

	value, ok := player.Settings[key]
	if !ok {
		return ""
	}

	return *value
}
func (player *Player) SetSetting(key string, value string) {
	if player.Settings == nil {
		player.Settings = make(postgres.Hstore)
	}

	if key == "phoneNumber" {
		player.PhoneNumber = value
	}

	player.Settings[key] = &value
	player.Save()
}

func (player *Player) MarkForDeath() {
	player.MarkedForDeath = true
	player.Save()
	//TODO: notify via slack or text
}

func (player *Player) MarkTarget() error {
	target, err := GetPlayerByCSHUsername(player.Target)

	target.MarkForDeath()

	return err
}

func (player *Player) ConfirmOwnMark() {
	if player.MarkedForDeath {
		player.Killed = true
		player.Save()
	}
}

func (player *Player) DenyOwnMark() {
	if player.MarkedForDeath {
		player.MarkedForDeath = false
		player.Save()
	}
}

func (player *Player) TargetIsDead() bool {
	if player.Target != "" {
		target, err := GetPlayerByCSHUsername(player.Target)
		if err != nil {
			logrus.Error(err)
		}
		return target.Killed
	}
	return false

}
func (player *Player) UpdatePlayerData() error {
	defer player.Save()
	if player.TargetIsDead() {
		target, err := GetPlayerByCSHUsername(player.Target)
		if err != nil {
			return err
		}
		player.Target = target.Target
		player.Kills = player.Kills + 1

		//TODO: notify via slack or text

	}
	return nil
}

func (player *Player) SetSlackUsername(username string) {
	defer player.Save()
	player.SlackUsername = username
}