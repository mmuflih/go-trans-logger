package trslog

import (
	"encoding/json"
	"fmt"
	"time"

	mgopaginator "github.com/mmuflih/mgo-paginator"
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
	ReadLog(query map[string]interface{}) ([]*TrsLog, error)
	Logs(query map[string]interface{}, page, size int) *mgopaginator.PaginatorResponse
}

type transLog struct {
	col *mgo.Collection
}

func NewTransactionLog(db *mgo.Database) TransactionLog {
	fmt.Println("======+++> Initial Transaction LOG")
	return &transLog{db.C("transaction_logs")}
}

func (tl transLog) getLastID() int64 {
	sl := new(TrsLog)
	err := tl.col.Find(nil).Sort("-id").One(&sl)
	if err != nil {
		return 1
	}
	return sl.ID + 1
}

func (t transLog) WriteLog(data *TrsLogData) {
	tl := &TrsLog{
		ID:       t.getLastID(),
		UserID:   data.UserID,
		RefType:  data.RefType,
		RefID:    data.RefID,
		Action:   data.Action,
		NewValue: data.NewValue,
		Details:  data.Details,
		ActionAt: time.Now().Unix(),
	}

	err := t.col.Insert(tl)
	if err != nil {
		jsonData, _ := json.Marshal(tl)
		fmt.Println("Error inserting log data => ", string(jsonData))
	}
}

func (t transLog) ReadLog(query map[string]interface{}) ([]*TrsLog, error) {
	var items []*TrsLog
	err := t.col.Find(query).All(&items)
	if err != nil {
		return nil, err
	}

	var newItems []*TrsLog
	for _, d := range items {
		d.ActionDateAt = time.Unix(d.ActionAt, 0)
		newItems = append(newItems, d)
	}
	return newItems, nil
}

func (t transLog) Logs(query map[string]interface{}, page, size int) *mgopaginator.PaginatorResponse {
	var items []*TrsLog
	qu := t.col.Find(query)

	paginate := mgopaginator.Paginator{
		Query: qu,
		Page:  page,
		Size:  size,
		Sort:  "-action_at",
	}

	resp := paginate.Paginate(&items)

	var newItems []*TrsLog
	for _, d := range items {
		d.ActionDateAt = time.Unix(d.ActionAt, 0)
		newItems = append(newItems, d)
	}

	return resp
}

type TrsLogData struct {
	UserID   int64
	RefType  string
	RefID    interface{}
	Action   string
	NewValue interface{}
	Details  interface{}
}

type TrsLog struct {
	MID          bson.ObjectId `bson:"_id,omitempty" json:"_id,omitempty"`
	ID           int64         `bson:"id" json:"id"`
	UserID       int64         `bson:"user_id" json:"user_id"`
	RefType      string        `bson:"ref_type" json:"ref_type"`
	RefID        interface{}   `bson:"ref_id" json:"ref_id"`
	Action       string        `bson:"action" json:"action"`
	NewValue     interface{}   `bson:"new_value" json:"new_value"`
	Details      interface{}   `bson:"details" json:"details"`
	ActionAt     int64         `bson:"action_at" json:"action_at"`
	ActionDateAt time.Time     `bson:"-" json:"action_date_at"`
}
