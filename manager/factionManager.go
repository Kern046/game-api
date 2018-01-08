package manager

import(
  "kalaxia-game-api/database"
  factionModel "kalaxia-game-api/model/faction"
  mapModel "kalaxia-game-api/model/map"
  serverModel "kalaxia-game-api/model/server"
)

func GetFaction(id uint16) *factionModel.Faction {
  faction := &factionModel.Faction{Id: id}
  if err := database.Connection.Select(faction); err != nil {
    return nil
  }
  return faction
}

func GetServerFactions(serverId uint16) []*factionModel.Faction {
  factions := make([]*factionModel.Faction, 0)
  if err := database.Connection.Model(&factions).Where("server_id = ?", serverId).Select(); err != nil {
    return nil
  }
  return factions
}

func GetFactionPlanetChoices(factionId uint16) []*mapModel.Planet {
  planets := make([]*mapModel.Planet, 0)
  faction := &factionModel.Faction{ Id: factionId }
  if err := database.Connection.Select(faction); err != nil {
    return planets
  }
  if _, err := database.
    Connection.
    Query(&planets, `SELECT p.*
      FROM map__maps m
      LEFT JOIN map__systems s ON s.map_id = m.id
      LEFT JOIN map__planets p ON p.system_id = s.id
      LEFT JOIN diplomacy__relations d ON d.planet_id = p.id
      WHERE p.player_id IS NULL AND m.server_id = ? AND d.faction_id = ?
      ORDER BY d.score DESC
      LIMIT 4`, faction.ServerId, faction.Id); err != nil {
    return planets
  }
  for _, planet := range planets {
      planet.Resources = make([]mapModel.PlanetResource, 0)
      if err := database.
        Connection.
        Model(&planet.Resources).
        Where("planet_id = ?", planet.Id).
        Select(); err != nil {
        panic(err)
      }
      relations := GetPlanetRelations(planet.Id)
      r := make([]interface{}, len(relations))
      for i, v := range relations {
          r[i] = v
      }
      planet.Relations = r
  }
  return planets
}

func CreateServerFactions(server *serverModel.Server, factions []interface{}) []*factionModel.Faction {
  results := make([]*factionModel.Faction, 0)
  for _, factionData := range factions {
    data := factionData.(map[string]interface{})
    faction := &factionModel.Faction{
      Name: data["name"].(string),
      Description: data["description"].(string),
      Color: data["color"].(string),
      ServerId: server.Id,
      Server: server,
    }
    if err := database.Connection.Insert(faction); err != nil {
      panic(err)
    }
    results = append(results, faction)
  }
  return results
}