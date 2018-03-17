package controller

import(
    "encoding/json"
    "io"
    "io/ioutil"
    "net/http"
    "kalaxia-game-api/manager"
    "kalaxia-game-api/security"
)

func CreateServer(w http.ResponseWriter, r *http.Request) {
    var body []byte
    var err error
    if body, err = ioutil.ReadAll(io.LimitReader(r.Body, 1048576)); err != nil {
        panic(err)
    }
    if err = r.Body.Close(); err != nil {
        panic(err)
    }
    jsonData := security.Decrypt(r.Header.Get("Application-Key"), r.Header.Get("Application-Iv"), body)
    var data map[string]interface{}
    if err = json.Unmarshal(jsonData, &data); err != nil {
        panic(err)
    }
    server := manager.CreateServer(
        data["name"].(string),
        data["type"].(string),
        data["signature"].(string),
    )
    factions := manager.CreateServerFactions(server, data["factions"].([]interface{}))
    manager.GenerateMap(server, factions, uint16(data["map_size"].(float64)))
    w.WriteHeader(http.StatusCreated)
    w.Write([]byte(""))
}