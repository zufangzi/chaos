package mongo

import (
	"encoding/json"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"opensource/chaos/background/server/dto/model"
	"opensource/chaos/background/utils"
	"reflect"
	"runtime/debug"
	"strings"
)

var session *mgo.Session
var db string

func MongoClose() {
	defer session.Close()
}

func MongoInit() {
	var err error
	db = utils.Param.MongoDB
	session, err = mgo.Dial(utils.Path.MongoUrl)
	utils.AssertPanic(err)
	session.SetMode(mgo.Monotonic, true)
}

func Reflect(freshCondition model.Mongo) map[string]interface{} {
	return FullReflect(freshCondition, "")
}

// 考虑到一种情况，就是条件真的0，那么需要额外一个参数来进行过滤
// 这只是提供的一种简单条件的处理办法，如果比较复杂的，还是用mongo的语法直接来写吧
func FullReflect(freshCondition model.Mongo, args ...string) map[string]interface{} {
	result := bson.M{}

	skipMap := make(map[string]interface{})
	for _, v := range args {
		skipMap[v] = 1
	}
	s := reflect.ValueOf(&freshCondition).Elem()
	typeOfT := s.Type()
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		name := strings.ToLower(typeOfT.Field(i).Name)

		if skipMap[name] == nil {
			realType := f.Type().String()
			if realType == "string" && f.String() == "" {
				continue
			}
			if (realType == "float64" || realType == "float") && f.Float() == 0.0 {
				continue
			}
			if realType == "int" && f.Int() == 0 {
				continue
			}

			if f.String() == "<[]string Value>" {
				continue
			}
		}
		result[strings.ToLower(typeOfT.Field(i).Name)] = f.String()

	}
	fmt.Println("conditon result: ", result)
	return result
}

func Safe(collection string, f func(*mgo.Collection)) {
	s := session.Clone()
	defer func() {
		s.Close()
		if e, ok := recover().(error); ok {
			log.Println("[CHAOS]catchable mongo error occur. " + e.Error())
			debug.PrintStack()
		}
	}()
	c := s.DB(db).C(collection)
	f(c)
}

func query(collection string, result interface{}, condition model.Mongo) {
	Safe(collection, func(c *mgo.Collection) {
		bsonMap := Reflect(condition)
		err := c.Find(&bsonMap).One(result)
		utils.AssertPrint(err)
		b, _ := json.Marshal(result)
		fmt.Println("[CHAOS]query result is: ", string(b))
	})
}

func queryById(collection string, result interface{}, id string) {
	condition := model.Mongo{}
	condition.BizId = id
	query(collection, result, condition)
}

func update(collection string, condition model.Mongo, docs interface{}) {
	Safe(collection, func(c *mgo.Collection) {
		b, _ := json.Marshal(docs)
		fmt.Println("[CHAOS]update data is: ", string(b))
		err := c.Update(Reflect(condition), bson.M{"$set": docs})
		utils.AssertPrint(err)
	})
}

func updateById(collection string, data interface{}, id string) {
	condition := model.Mongo{}
	condition.BizId = id
	update(collection, condition, data)
}

func insert(collection string, docs interface{}) {
	Safe(collection, func(c *mgo.Collection) {
		fmt.Println(&docs)
		err := c.Insert(&docs)
		utils.AssertPanic(err)
	})
}

func delById(collection string, id string) {
	condition := model.Mongo{}
	condition.BizId = id
	del(collection, condition)
}

func del(collection string, condition model.Mongo) {
	Safe(collection, func(c *mgo.Collection) {
		err := c.Remove(Reflect(condition))
		utils.AssertPanic(err)
	})
}
