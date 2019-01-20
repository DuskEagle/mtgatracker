package main

type Transaction struct {
	TransactionId    string `json:"transactionId"`
	Timestamp        string `json:"timestamp"`
	GreToClientEvent struct {
		GreToClientMessages []struct {
			Type             string `json:"type"`
			SystemSeatIds    []int  `json:"systemSeatIds"`
			MsgId            int    `json:"msgId"`
			GameStateId      int    `json:"gameStateId"`
			GameStateMessage struct {
				Type        string `json:"type"`
				GameStateId int    `json:"gameStateId"`
				GameObjects *[]struct {
					InstanceId       int      `json:"instanceId"`
					GrpId            int      `json:"grpId"`
					Type             string   `json:"type"`
					ZoneId           int      `json:"zoneId"`
					Visibility       string   `json:"visibility"`
					OwnerSeatId      int      `json:"ownerSeatId"`
					ControllerSeatId int      `json:"controllerSeatId"`
					CardTypes        []string `json:"cardTypes"`
					Subtypes         []string `json:"subtypes"`
					Color            []string `json:"color"`
					Power            struct {
						Value int `json:"value"`
					} `json:"power"`
					Toughness struct {
						Value int `json:"value"`
					} `json:"toughness"`
					IsTapped    bool   `json:"isTapped"`
					AttackState string `json:"attackState"`
					BlockState  string `json:"blockState"`
					AttackInfo  struct {
						TargetId      int  `json:"targetId"`
						DamageOrdered bool `json:"damageOrdered"`
					} `json:"attackInfo"`
					Name         int   `json:"name"`
					Abilities    []int `json:"abilities"`
					OverlayGrpId int   `json:"overlayGrpId"`
				} `json:"gameObjects,omitempty"`
				TurnInfo    *struct {
					Phase          string `json:"phase"`
					Step           string `json:"step"`
					TurnNumber     int    `json:"turnNumber"`
					ActivePlayer   int    `json:"activePlayer"`
					PriorityPlayer int    `json:"priorityPlayer"`
					DecisionPlayer int    `json:"decisionPlayer"`
					NextPhase      string `json:"nextPhase"`
					NextStep       string `json:"nextStep"`
				} `json:"turnInfo,omitempty"`
				GameInfo    *struct {
					MatchId            string `json:"matchID"`
					GameNumber         int    `json:"gameNumber"`
					Stage              string `json:"stage"`
					Type               string `json:"type"`
					Variant            string `json:"variant"`
					MatchState         string `json:"matchState"`
					MatchWinCondition  string `json:"matchWinCondition"`
					MaxTimeoutCount    int    `json:"maxTimeoutCount"`
					MaxPipCount        int    `json:"maxPipCount"`
					TimeoutDurationSec int    `json:"timeoutDurationSec"`
					Results            []struct {
						Scope         string `json:"scope"`
						Result        string `json:"result"`
						WinningTeamId int    `json:"winningTeamId"`
						Reason        string `json:"reason"`
					} `json:"results"`
					SuperFormat  string `json:"superFormat"`
					MulliganType string `json:"mulliganType"`
				} `json:"gameInfo"`
				Annotations *[]struct {
					Id          int      `json:"id"`
					AffectorId  int      `json:"affectorId"`
					AffectedIds []int    `json:"affectedIds"`
					Type        []string `json:"type"`
					Details     []struct {
						Key        string `json:"key"`
						Type       string `json:"type"`
						ValueInt32 []int  `json:"valueInt32"`
					} `json:"details,omitempty"`
					AllowRedaction bool `json:"allowRedaction,omitempty"`
				} `json:"annotations,omitempty"`
				PrevGameStateId int `json:"prevGameStateId"`
				Update  string `json:"update"`
				Actions *[]struct {
					SeatId int `json:"seatId"`
					Action struct {
						ActionType string `json:"actionType"`
						InstanceId int    `json:"instanceId"`
					} `json:"action"`
				} `json:"actions,omitempty"`
				Zones *[]struct {
					ZoneId int `json:"zoneId"`
					Type string `json:"type"`
					Visibility string `json:"visibility"`
					ObjectInstanceIds []int `json:"objectInstanceIds"`
					OwnerSeatId int `json:"ownerSeatId"`
				} `json:"zones,omitempty"`
			} `json:"gameStateMessage"`
		} `json:"greToClientMessages"`
	} `json:"greToClientEvent"`
	MatchGameRoomStateChangedEvent *struct {
	    GameRoomInfo struct {
	      	GameRoomConfig struct {
		        EventId         string `json:"eventId"`
		        ReservedPlayers []struct {
		      		UserId         string `json:"userId"`
		      		PlayerName     string `json:"playerName"`
		      		SystemSeatId   int    `json:"systemSeatId"`
		      		TeamId         int    `json:"teamId"`
		      		ConnectionInfo struct {
		        		ConnectionState string `json:"connectionState"`
		    		} `json:"connectionInfo"`
		        	CourseID string `json:"courseId"`
		        } `json:"reservedPlayers"`
		        MatchId     string `json:"matchId"`
		      	Players   []struct {
		    		UserId       string `json:"userId"`
		    		SystemSeatId int    `json:"systemSeatId"`
		  		} `json:"players"`
	  		} `json:"gameRoomConfig"`
    	} `json:"gameRoomInfo"`
  	} `json:"matchGameRoomStateChangedEvent"`
}