package worker

import (
	"context"
	"github.com/krkrkc/distributed-crontab/common"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type LogSink struct {
	client         *mongo.Client
	logCollection  *mongo.Collection
	logChan        chan *common.JobLog
	autoCommitChan chan *common.LogBatch
}

var (
	G_logSink *LogSink
)

func InitLogSink() error {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(G_config.MongodbUri), options.Client().SetConnectTimeout(time.Duration(G_config.MongodbConnectionTimeout)*time.Millisecond))
	if err != nil {
		return err
	}

	G_logSink = &LogSink{
		client:         client,
		logCollection:  client.Database("cron").Collection("log"),
		logChan:        make(chan *common.JobLog, 1000),
		autoCommitChan: make(chan *common.LogBatch, 1000),
	}
	go G_logSink.writeLoop()
	return nil
}

func (logSink *LogSink) writeLoop() {
	var logBatch *common.LogBatch
	var commitTimer *time.Timer
	for {
		select {
		case log := <-logSink.logChan:
			if logBatch == nil {
				logBatch = &common.LogBatch{}
				commitTimer = time.AfterFunc(1*time.Second, /*time.Duration(G_config.JobLogCommitTimeout)*time.Millisecond*/
					func(batch *common.LogBatch) func() {
						return func() {
							logSink.autoCommitChan <- batch
						}
					}(logBatch))
			}
			logBatch.Logs = append(logBatch.Logs, log)
			if len(logBatch.Logs) >= G_config.JobLogBatchSize {
				logSink.saveLogs(logBatch)
				logBatch = nil
				commitTimer.Stop()
			}
		case timeoutBatch := <-logSink.autoCommitChan:
			if timeoutBatch != logBatch {
				continue
			}
			logSink.saveLogs(timeoutBatch)
			logBatch = nil
		}
	}
}

func (logSink *LogSink) saveLogs(batch *common.LogBatch) {
	logSink.logCollection.InsertMany(context.TODO(), batch.Logs)
}

func (logSink *LogSink) Append(jobLog *common.JobLog) {
	select {
	case logSink.logChan <- jobLog:
	default:
	}
}
