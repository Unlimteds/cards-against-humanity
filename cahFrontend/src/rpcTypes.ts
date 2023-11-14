// Code generated by tygo. DO NOT EDIT.

//////////
// source: connection.go

export interface GameMessage {
  Message: string;
  GameId: string /* uuid */;
  PlayerId: string /* uuid */;
}
export interface WsConnection {
  Conn?: any /* websocket.Conn */;
  PlayerId: string /* uuid */;
  GameID: string /* uuid */;
  JoinTime: string /* RFC3339 */;
  LastPingTime: string /* RFC3339 */;
  WsRecieve: any;
  WsBroadcast: any;
}

//////////
// source: connection_manager.go

export interface GameConnection {
  PlayerConnectionMap: { [key: string /* uuid */]: WsConnection | undefined};
}
export interface GlobalConnectionManager {
  /**
   * Maps a game ID to the game connection pool
   */
  GameConnectionMap: { [key: string /* uuid */]: GameConnection | undefined};
}

//////////
// source: game_metrics.go

export interface Metrics {
}

//////////
// source: games_endpoints.go


//////////
// source: helpers.go

export interface ApiError {
  error: string;
}

//////////
// source: resources_endpoints.go
