package trslog

import (
	"encoding/json"
	"fmt"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

/**
 * Created by Muhammad Muflih Kholidin
 * at 2021-03-10 18:24:38
 * https://github.com/mmuflih
 * muflic.24@gmail.com
 **/

type TransactionLog interface {
	WriteLog(data *TrsLogData)
}

type transLog struct {
	col *mgo.Collection
}

func NewTransactionLog(db *mgo.Database) TransactionLog {
	return &transLog{db.C("transaction_logs")}
}

func (tl transLog) getLastID() int64 {
	sl := new(trsLog)
	err := tl.col.Find(nil).Sort("-id").One(&sl)
	if err != nil {
		return 1
	}
	return sl.ID + 1
}

func (t transLog) WriteLog(data *TrsLogData) {
	tl := &trsLog{
		ID:       t.getLastID(),
		UserID:   data.UserID,
		RefType:  data.RefType,
		RefID:    data.RefID,
		Action:   data.Action,
		OldValue: data.OldValue,
		NewValue: data.NewValue,
		ActionAt: time.Now().Unix(),
	}

	err := t.col.Insert(tl)
	if err != nil {
		jsonData, _ := json.Marshal(tl)
		fmt.Println("Error inserting log data => ", string(jsonData))
	}
}

type TrsLogData struct {
	UserID   int64
	RefType  string
	RefID    interface{}
	Action   string
	OldValue interface{}
	NewValue interface{}
}

type trsLog struct {
	MID      bson.ObjectId `bson:"_id,omitempty" json:"_id,omitempty"`
	ID       int64         `bson:"id" json:"id"`
	UserID   int64         `bson:"user_id" json:"user_id"`
	RefType  string        `bson:"ref_type" json:"ref_type"`
	RefID    interface{}   `bson:"ref_id" json:"ref_id"`
	Action   string        `bson:"action" json:"action"`
	OldValue interface{}   `bson:"old_value" json:"old_value"`
	NewValue interface{}   `bson:"new_value" json:"new_value"`
	ActionAt int64         `bson:"action_at" json:"action_at"`
}
