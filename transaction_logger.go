package trslog

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/mmuflih/go-trans-logger/paginator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/**
 * Created by Muhammad Muflih Kholidin
 * at 2021-03-10 18:24:38
 * https://github.com/mmuflih
 * muflic.24@gmail.com
 **/

type TransactionLog interface {
	WriteLog(data *TrsLogData)
	ReadLog(query map[string]interface{}, limit int64) ([]*TrsLog, error)
	Logs(query primitive.D, page, size int) (*paginator.PaginatorResponse, error)
	GetType() []bson.M
}

type transLog struct {
	col *mongo.Collection
}

func NewTransactionLog(db *mongo.Database) TransactionLog {
	fmt.Println("======+++> Initial Transaction LOG")
	return &transLog{db.Collection("transaction_logs")}
}

func (tl transLog) getLastID() int64 {
	sl := new(TrsLog)
	sort := bson.D{{Key: "id", Value: -1}}
	fo := options.FindOne().SetSort(sort)
	err := tl.col.FindOne(context.TODO(), bson.M{}, fo).Decode(&sl)
	if err != nil {
		return 1
	}
	return sl.ID + 1
}

func (t transLog) WriteLog(data *TrsLogData) {
	tl := &TrsLog{
		ID:       t.getLastID(),
		User:     data.User,
		RefType:  data.RefType,
		RefID:    data.RefID,
		Action:   data.Action,
		NewValue: data.NewValue,
		Details:  data.Details,
		ActionAt: time.Now().Unix(),
	}

	_, err := t.col.InsertOne(context.TODO(), tl)
	if err != nil {
		jsonData, _ := json.Marshal(tl)
		fmt.Println("Error inserting log data => ", string(jsonData))
	}
}

func (t transLog) ReadLog(query map[string]interface{}, limit int64) ([]*TrsLog, error) {
	var items []*TrsLog
	sort := bson.D{{Key: "action_at", Value: -1}}
	fo := options.Find().SetSort(sort).SetLimit(limit)
	cursor, err := t.col.Find(context.TODO(), query, fo)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(context.TODO())
	if err := cursor.All(context.TODO(), &items); err != nil {
		fmt.Println("Error decoding documents:", err)
		return nil, err
	}

	var newItems []*TrsLog
	for _, d := range items {
		d.ActionDateAt = time.Unix(d.ActionAt, 0)
		newItems = append(newItems, d)
	}
	return newItems, nil
}

func (t transLog) Logs(query primitive.D, page, size int) (*paginator.PaginatorResponse, error) {
	var items []*TrsLog

	paginate := paginator.MPaginator{
		Collection: &mongo.Collection{},
		Filter:     query,
		Page:       0,
		Size:       0,
		Sort:       map[string]int{"action_at": -1},
	}

	resp, _ := paginate.GetPaginator(context.TODO(), &items)

	var newItems []*TrsLog
	for _, d := range items {
		d.ActionDateAt = time.Unix(d.ActionAt, 0)
		newItems = append(newItems, d)
	}
	resp.Data = newItems

	return resp, nil
}

func (t transLog) GetType() []bson.M {
	items := []bson.M{}
	pipeline := []bson.M{
		{
			"$group": bson.M{"_id": "$ref_type"},
		},
	}
	cursor, err := t.col.Aggregate(context.TODO(), pipeline)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer cursor.Close(context.TODO())

	err = cursor.All(context.TODO(), &items)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return items
}

type TrsLogData struct {
	User     interface{}
	RefType  string
	RefID    interface{}
	Action   string
	NewValue interface{}
	Details  interface{}
}

type TrsLog struct {
	MID          primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	ID           int64              `bson:"id" json:"id"`
	User         interface{}        `bson:"user_id" json:"user_id"`
	RefType      string             `bson:"ref_type" json:"ref_type"`
	RefID        interface{}        `bson:"ref_id" json:"ref_id"`
	Action       string             `bson:"action" json:"action"`
	NewValue     interface{}        `bson:"new_value" json:"new_value"`
	Details      interface{}        `bson:"details" json:"details"`
	ActionAt     int64              `bson:"action_at" json:"action_at"`
	ActionDateAt time.Time          `bson:"-" json:"action_date_at"`
}
