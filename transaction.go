package hmongo

import (
	"context"
	"errors"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"strings"
)

// Mongodb 数据库事务操作
func Transaction(requestId string, client *mongo.Client, resove func(sessionContext mongo.SessionContext) error) error {
	if requestId == "" {
		requestId = Rand32()
	}
	tlog := log.WithField("requestId", requestId)
	e := client.UseSession(context.Background(), func(sessionContext mongo.SessionContext) (err error) {
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
	if e != nil {
		if strings.HasPrefix(e.Error(), "(LockTimeout)") {
			e = errors.New("系统繁忙请稍后重试.[LockTimeout]")
		}
		if strings.HasPrefix(e.Error(), "(WriteConflict)") {
			e = errors.New("系统繁忙请稍后重试.[WriteConflict]")
		}
	}
	return e
}
