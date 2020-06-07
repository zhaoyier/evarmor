package model

import (
	"time"

	//3rd party libs
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	//Own libs
	"github.com/ezbuy/ezorm/db"
	. "github.com/ezbuy/ezorm/orm"
)

var _ time.Time

func init() {

	db.SetOnEnsureIndex(initCounterIndex)

	RegisterEzOrmObjByID("gateway", "Counter", newCounterFindByID)
	RegisterEzOrmObjRemove("gateway", "Counter", CounterMgr.RemoveByID)

}

func initCounterIndex() {
	session, collection := CounterMgr.GetCol()
	defer session.Close()

	if err := collection.EnsureIndex(mgo.Index{
		Key:        []string{"Key"},
		Unique:     true,
		Background: true,
		Sparse:     true,
	}); err != nil {
		panic("ensureIndex .Counter Key error:" + err.Error())
	}

}

func newCounterFindByID(id string) (result EzOrmObj, err error) {
	return CounterMgr.FindByID(id)
}

//mongo methods
var (
	insertCB_Counter []func(obj EzOrmObj)
	updateCB_Counter []func(obj EzOrmObj)
)

func CounterAddInsertCallback(cb func(obj EzOrmObj)) {
	insertCB_Counter = append(insertCB_Counter, cb)
}

func CounterAddUpdateCallback(cb func(obj EzOrmObj)) {
	updateCB_Counter = append(updateCB_Counter, cb)
}

func (o *Counter) Id() string {
	return o.ID.Hex()
}

func (o *Counter) Save() (info *mgo.ChangeInfo, err error) {
	session, col := CounterMgr.GetCol()
	defer session.Close()

	isNew := o.isNew

	info, err = col.UpsertId(o.ID, o)
	o.isNew = false

	if isNew {
		CounterInsertCallback(o)
	} else {
		CounterUpdateCallback(o)
	}

	return
}

func (o *Counter) InsertUnique(query interface{}) (saved bool, err error) {
	session, col := CounterMgr.GetCol()
	defer session.Close()

	info, err := col.Upsert(query, db.M{"$setOnInsert": o})
	if err != nil {
		return
	}
	if info.Updated == 0 {
		saved = true
	}
	o.isNew = false
	if saved {
		CounterInsertCallback(o)
	}
	return
}

func CounterInsertCallback(o *Counter) {
	for _, cb := range insertCB_Counter {
		cb(o)
	}
}

func CounterUpdateCallback(o *Counter) {
	for _, cb := range updateCB_Counter {
		cb(o)
	}
}

//foreigh keys

//Collection Manage methods

func (o *_CounterMgr) FindOne(query interface{}, sortFields ...string) (result *Counter, err error) {
	session, col := CounterMgr.GetCol()
	defer session.Close()

	q := col.Find(query)

	_CounterSort(q, sortFields)

	err = q.One(&result)
	return
}

func _CounterSort(q *mgo.Query, sortFields []string) {
	sortFields = XSortFieldsFilter(sortFields)
	if len(sortFields) > 0 {
		q.Sort(sortFields...)
		return
	}

	q.Sort("-_id")
}

func (o *_CounterMgr) Query(query interface{}, limit, offset int, sortFields []string) (*mgo.Session, *mgo.Query) {
	session, col := CounterMgr.GetCol()
	q := col.Find(query)
	if limit > 0 {
		q.Limit(limit)
	}
	if offset > 0 {
		q.Skip(offset)
	}

	_CounterSort(q, sortFields)
	return session, q
}

func (o *_CounterMgr) NQuery(query interface{}, limit, offset int, sortFields []string) (*mgo.Session, *mgo.Query) {
	session, col := CounterMgr.GetCol()
	q := col.Find(query)
	if limit > 0 {
		q.Limit(limit)
	}
	if offset > 0 {
		q.Skip(offset)
	}

	if sortFields = XSortFieldsFilter(sortFields); len(sortFields) > 0 {
		q.Sort(sortFields...)
	}

	return session, q
}
func (o *_CounterMgr) FindOneByKey(Key string) (result *Counter, err error) {
	query := db.M{
		"Key": Key,
	}
	session, q := CounterMgr.NQuery(query, 1, 0, nil)
	defer session.Close()
	err = q.One(&result)
	return
}

func (o *_CounterMgr) MustFindOneByKey(Key string) (result *Counter) {
	result, _ = o.FindOneByKey(Key)
	if result == nil {
		result = CounterMgr.NewCounter()
		result.Key = Key
		result.Save()
	}
	return
}

func (o *_CounterMgr) RemoveByKey(Key string) (err error) {
	session, col := CounterMgr.GetCol()
	defer session.Close()

	query := db.M{
		"Key": Key,
	}
	return col.Remove(query)
}

func (o *_CounterMgr) Find(query interface{}, limit int, offset int, sortFields ...string) (result []*Counter, err error) {
	session, q := CounterMgr.Query(query, limit, offset, sortFields)
	defer session.Close()
	err = q.All(&result)
	return
}

func (o *_CounterMgr) FindAll(query interface{}, sortFields ...string) (result []*Counter, err error) {
	session, q := CounterMgr.Query(query, -1, -1, sortFields)
	defer session.Close()
	err = q.All(&result)
	return
}

func (o *_CounterMgr) Has(query interface{}) bool {
	session, col := CounterMgr.GetCol()
	defer session.Close()

	var ret interface{}
	err := col.Find(query).One(&ret)
	if err != nil || ret == nil {
		return false
	}
	return true
}

func (o *_CounterMgr) Count(query interface{}) (result int) {
	result, _ = o.CountE(query)
	return
}

func (o *_CounterMgr) CountE(query interface{}) (result int, err error) {
	session, col := CounterMgr.GetCol()
	defer session.Close()

	result, err = col.Find(query).Count()
	return
}

func (o *_CounterMgr) FindByIDs(id []string, sortFields ...string) (result []*Counter, err error) {
	ids := make([]bson.ObjectId, 0, len(id))
	for _, i := range id {
		if bson.IsObjectIdHex(i) {
			ids = append(ids, bson.ObjectIdHex(i))
		}
	}
	return CounterMgr.FindAll(db.M{"_id": db.M{"$in": ids}}, sortFields...)
}

func (m *_CounterMgr) FindByID(id string) (result *Counter, err error) {
	session, col := CounterMgr.GetCol()
	defer session.Close()

	if !bson.IsObjectIdHex(id) {
		err = mgo.ErrNotFound
		return
	}
	err = col.FindId(bson.ObjectIdHex(id)).One(&result)
	return
}

func (m *_CounterMgr) RemoveAll(query interface{}) (info *mgo.ChangeInfo, err error) {
	session, col := CounterMgr.GetCol()
	defer session.Close()

	return col.RemoveAll(query)
}

func (m *_CounterMgr) RemoveByID(id string) (err error) {
	session, col := CounterMgr.GetCol()
	defer session.Close()

	if !bson.IsObjectIdHex(id) {
		err = mgo.ErrNotFound
		return
	}
	err = col.RemoveId(bson.ObjectIdHex(id))

	return
}

func (m *_CounterMgr) GetCol() (session *mgo.Session, col *mgo.Collection) {
	if mgoInstances == nil {
		return db.GetCol("", "gateway.Counter")
	}
	return getCol("", "gateway.Counter")
}

//Search

func (o *Counter) IsSearchEnabled() bool {

	return false

}

//end search
