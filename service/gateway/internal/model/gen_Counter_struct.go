package model

import "gopkg.in/mgo.v2/bson"

import "time"

var _ time.Time

type Counter struct {
	ID    bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Key   string        `bson:"Key" json:"Key"`
	Val   int32         `bson:"Val" json:"Val"`
	isNew bool
}

const (
	CounterMgoFieldID  = "_id"
	CounterMgoFieldKey = "Key"
	CounterMgoFieldVal = "Val"
)
const (
	CounterMgoSortFieldIDAsc  = "_id"
	CounterMgoSortFieldIDDesc = "-_id"
)

func (p *Counter) GetNameSpace() string {
	return "gateway"
}

func (p *Counter) GetClassName() string {
	return "Counter"
}

type _CounterMgr struct {
}

var CounterMgr *_CounterMgr

// Get_CounterMgr returns the orm manager in case of its name starts with lower letter
func Get_CounterMgr() *_CounterMgr { return CounterMgr }

func (m *_CounterMgr) NewCounter() *Counter {
	rval := new(Counter)
	rval.isNew = true
	rval.ID = bson.NewObjectId()

	return rval
}
