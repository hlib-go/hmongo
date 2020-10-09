package hmongo

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"os"
	"testing"
	"time"
)

var (
	HM_MONGO_CONNECTION_STRING = os.Getenv("HM_MONGO_CONNECTION_STRING")
)

func TestConnect(t *testing.T) {
	t.Log(HM_MONGO_CONNECTION_STRING)
	c, _ := DB(HM_MONGO_CONNECTION_STRING)
	t.Log(c)
}

func TestBsonEncode(t *testing.T) {
	// [{"$match":{"actId":"1234567890","createdAt":{"$gt":"2020-01-02T00:00:00Z"}}}]
	//db.hm_rights_order.aggregate([{$match:{actId:"1234567890",createdAt:{$gt:ISODate("2020-01-02T00:00:00")}}}])
	tp, _ := time.Parse(time.RFC3339, "2020-01-02T00:00:00Z")
	var query = bson.A{bson.M{"$match": bson.M{"actId": "1234567890", "createdAt": bson.M{"$gt": tp}}}}

	bt, btBytes, err := bson.MarshalValue(query)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(bt.String())
	t.Log("bson.MarshalValue btBytes = ", string(btBytes))

	var q1 bson.A
	bson.Unmarshal(btBytes, &q1)

	t.Log("---------------------------------------------------------")
	sjson, _ := json.Marshal(query)
	t.Log("sjson=", string(sjson))

	var qa bson.A
	json.Unmarshal(sjson, &qa)
	qaBytes, _ := json.Marshal(qa)
	t.Log("qa=", string(qaBytes))

	var a = qa

	c, _ := DB(HM_MONGO_CONNECTION_STRING)

	cursor, err := c.Collection("hm_rights_order").Aggregate(nil, a)
	if err != nil {
		t.Error(err)
		return
	}
	var rmap []map[string]interface{}
	err = cursor.All(nil, &rmap)
	if err != nil {
		t.Error(err)
		return
	}
	bytes, err := json.Marshal(rmap)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(string(bytes))
}

func TestJsonQuery(t *testing.T) {
	//query := `[{$match: {"actId":"1234567890","prodId":"10"}}},{"$project": {"_id":0,"mobile":1}}]`
	//"createdAt":{"$gt":"ISODate(\"2020-07-02T00:00:00\")","$lt":"ISODate(\"2020-07-02T23:59:59\")"}
	query := `[
			  {
				"$match": {
				  "actId": "1234567890",
				  "createdAt": {
					"$gt": {
					  "$date": "2020-01-02T00:00:00Z"
					}
				  }
				}
			  }
			]`
	query = `
[ {"$match":{"actId":"1234567890","createdAt":{"$gt": {"$date":"2020-01-02T00:00:00Z"},"$lt": {"$date":"2022-07-02T23:59:59Z"}}}},   {"$project":{"_id":0,"orderId":1,"mobile":1,"prodName":1,"createdAt": {"$dateToString":{"format":"%Y-%m-%d %H:%M:%S","date":"$createdAt","timezone":"+08:00"}}} }]
`
	//ss:="[ {\"$match\":{\"actId\":\"1234567890\",\"createdAt\":{\"$gt\": {\"$date\":\"2020-01-02T00:00:00Z\"},\"$lt\": {\"$date\":\"2022-07-02T23:59:59Z\"}}}},   {\"$project\":{\"_id\":0,\"orderId\":1,\"mobile\":1,\"prodName\":1,\"createdAt\": {\"$dateToString\":{\"format\":\"%Y-%m-%d %H:%M:%S\",\"date\":\"$createdAt\",\"timezone\":\"+08:00\"}}} }]"
	var a bson.A
	err := bson.UnmarshalExtJSON([]byte(query), false, &a)
	if err != nil {
		t.Error(err)
		return
	}

	c, _ := DB(HM_MONGO_CONNECTION_STRING)

	cursor, err := c.Collection("hm_rights_order").Aggregate(nil, a)
	if err != nil {
		t.Error(err)
		return
	}
	var rmap []map[string]interface{}
	err = cursor.All(nil, &rmap)
	if err != nil {
		t.Error(err)
		return
	}
	bytes, err := json.Marshal(rmap)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(string(bytes))
}
