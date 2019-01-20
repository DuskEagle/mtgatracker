package main

type Transaction struct {
	TransactionID    string `json:"transactionId"`
	Timestamp        string `json:"timestamp"`
	GreToClientEvent struct {
		GreToClientMessages []struct {
			Type             string `json:"type"`
			SystemSeatIds    []int  `json:"systemSeatIds"`
			MsgID            int    `json:"msgId"`
			GameStateID      int    `json:"gameStateId"`
			GameStateMessage struct {
				Type        string `json:"type"`
				GameStateID int    `json:"gameStateId"`
				GameObjects []struct {
					InstanceID       int      `json:"instanceId"`
					GrpID            int      `json:"grpId"`
					Type             string   `json:"type"`
					ZoneID           int      `json:"zoneId"`
					Visibility       string   `json:"visibility"`
					OwnerSeatID      int      `json:"ownerSeatId"`
					ControllerSeatID int      `json:"controllerSeatId"`
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
						TargetID      int  `json:"targetId"`
						DamageOrdered bool `json:"damageOrdered"`
					} `json:"attackInfo"`
					Name         int   `json:"name"`
					Abilities    []int `json:"abilities"`
					OverlayGrpID int   `json:"overlayGrpId"`
				} `json:"gameObjects"`
				TurnInfo    struct {
					Phase          string `json:"phase"`
					Step           string `json:"step"`
					TurnNumber     int    `json:"turnNumber"`
					ActivePlayer   int    `json:"activePlayer"`
					PriorityPlayer int    `json:"priorityPlayer"`
					DecisionPlayer int    `json:"decisionPlayer"`
					NextPhase      string `json:"nextPhase"`
					NextStep       string `json:"nextStep"`
				} `json:"turnInfo"`
				Annotations []struct {
					ID          int      `json:"id"`
					AffectorID  int      `json:"affectorId"`
					AffectedIds []int    `json:"affectedIds"`
					Type        []string `json:"type"`
					Details     []struct {
						Key        string `json:"key"`
						Type       string `json:"type"`
						ValueInt32 []int  `json:"valueInt32"`
					} `json:"details,omitempty"`
					AllowRedaction bool `json:"allowRedaction,omitempty"`
				} `json:"annotations"`
				PrevGameStateID int `json:"prevGameStateId"`
				Update  string `json:"update"`
				Actions []struct {
					SeatID int `json:"seatId"`
					Action struct {
						ActionType string `json:"actionType"`
						InstanceID int    `json:"instanceId"`
					} `json:"action"`
				} `json:"actions"`
				Zones []struct {
					ZoneId int `json:"zoneId"`
					Type string `json:"type"`
					Visibility string `json:"visibility"`
					ObjectInstanceIds []int `json:"objectInstanceIds"`
					OwnerSeatId int `json:"ownerSeatId"`
				} `json:"zones"`
			} `json:"gameStateMessage"`
		} `json:"greToClientMessages"`
	} `json:"greToClientEvent"`
}