package hmongo

import (
	"context"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

// Mongodb 数据库事务操作
func Transaction(txname string, client *mongo.Client, resove func(sessionContext mongo.SessionContext) error) error {
	tlog := log.WithField("txname", txname)
	return client.UseSession(context.Background(), func(sessionContext mongo.SessionContext) (err error) {
		defer func() {
			if err != nil {
				e := sessionContext.AbortTransaction(sessionContext)
				if e != nil {
					tlog.Error("AbortTransaction Error:" + e.Error())
				}
				return
			}
			e := sessionContext.CommitTransaction(sessionContext)
			if e != nil {
				tlog.Error("CommitTransaction Error:" + e.Error())
				return
			}
			tlog.Info("CommitTransaction success")
		}()
		e := sessionContext.StartTransaction()
		if e != nil {
			tlog.Error("StartTransaction Error:" + e.Error())
			return
		}
		err = resove(sessionContext)
		return
	})
}
