package dfquery

import (
	"net"
	"strconv"

	"github.com/df-mc/dragonfly/server"
	"github.com/gameparrot/goquery"
	"github.com/gameparrot/goqueryraknet"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

const (
	QueryKeyHostName     = "hostname"
	QueryKeyGameType     = "gametype"
	QueryKeyGameID       = "game_id"
	QueryKeyVersion      = "version"
	QueryKeyServerEngine = "server_engine"
	QueryKeyPlugins      = "plugins"
	QueryKeyMap          = "map"
	QueryKeyNumPlayers   = "numplayers"
	QueryKeyMaxPlayers   = "maxplayers"
	QueryKeyWhitelist    = "whitelist"

	gameType = "MINECRAFTPE"
)

type DfQuery struct {
	OnQueryRequest func(addr net.Addr, info map[string]string, players *[]string)
	q              *goquery.QueryServer
	srv            *server.Server
	statusProvider minecraft.ServerStatusProvider

	cfg       server.Config
	dfVersion string
}

func NewServerWithQuery(c server.Config) (*server.Server, *DfQuery) {
	dfQuery := &DfQuery{q: goquery.New(map[string]string{}, []string{}), dfVersion: getDfVersion(), statusProvider: c.StatusProvider, cfg: c}
	goqueryraknet.CreateGophertunnelNetwork("raknet", dfQuery.q)
	dfQuery.srv = c.New()
	dfQuery.q.SetInfoFunc(dfQuery.handleQuery)
	return dfQuery.srv, dfQuery
}

func (q *DfQuery) handleQuery(addr net.Addr) (map[string]string, []string) {
	playerCount := q.srv.PlayerCount()
	queryInfo := map[string]string{
		QueryKeyHostName:     q.cfg.Name,
		QueryKeyGameType:     gameType,
		QueryKeyVersion:      "v" + protocol.CurrentVersion,
		QueryKeyServerEngine: "Dragonfly " + q.dfVersion,
		QueryKeyPlugins:      "Dragonfly " + q.dfVersion,
		QueryKeyMap:          q.srv.World().Name(),
		QueryKeyNumPlayers:   strconv.Itoa(playerCount),
		QueryKeyMaxPlayers:   strconv.Itoa(q.srv.MaxPlayerCount()),
		QueryKeyWhitelist:    "false",
	}
	playerNames := make([]string, playerCount)
	i := 0
	for pl := range q.srv.Players(nil) {
		playerNames[i] = pl.Name()
		i++
	}

	if q.OnQueryRequest != nil {
		q.OnQueryRequest(addr, queryInfo, &playerNames)
	}

	return queryInfo, playerNames
}
