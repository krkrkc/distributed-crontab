package worker

import (
	"context"
	"fmt"
	"github.com/krkrkc/distributed-crontab/common"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

type JobMgr struct {
	client  *clientv3.Client
	kv      clientv3.KV
	lease   clientv3.Lease
	watcher clientv3.Watcher
}

var (
	G_jobMgr *JobMgr
)

func InitJobMgr() (err error) {
	var (
		config  clientv3.Config
		client  *clientv3.Client
		kv      clientv3.KV
		lease   clientv3.Lease
		watcher clientv3.Watcher
	)

	config = clientv3.Config{
		Endpoints:   G_config.EtcdEndpoints,
		DialTimeout: time.Duration(G_config.EtcdDialTimeout) * time.Millisecond,
	}

	if client, err = clientv3.New(config); err != nil {
		fmt.Println(err)
		return
	}

	kv = clientv3.NewKV(client)
	lease = clientv3.NewLease(client)
	watcher = clientv3.NewWatcher(client)

	G_jobMgr = &JobMgr{
		client:  client,
		kv:      kv,
		lease:   lease,
		watcher: watcher,
	}

	err = G_jobMgr.watchJobs()
	if err != nil {
		fmt.Println(err)
	}
	return
}

func (j *JobMgr) watchJobs() error {
	getResp, err := j.kv.Get(context.TODO(), common.JobSaveDir, clientv3.WithPrefix())
	if err != nil {
		return err
	}

	for _, kvPair := range getResp.Kvs {
		job, err := common.UnpackJob(kvPair.Value)
		if err == nil {
			//TODO
			jobEvent := common.BuildJobEvent(common.JobEventSave, job)
			jobEvent = jobEvent
		}
	}

	go func() {
		watchStartRevision := getResp.Header.Revision + 1
		watchChan := j.watcher.Watch(context.TODO(), common.JobSaveDir, clientv3.WithRev(watchStartRevision), clientv3.WithPrefix())
		for watchResp := range watchChan {
			for _, watchEvent := range watchResp.Events {
				switch watchEvent.Type {
				case mvccpb.PUT:
					job, err := common.UnpackJob(watchEvent.Kv.Value)
					if err != nil {
						continue
					}
					jobEvent := common.BuildJobEvent(common.JobEventSave, job)
					//TODO
					fmt.Println(jobEvent)
				case mvccpb.DELETE:
					jobName := common.ExtractJobName(string(watchEvent.Kv.Value))
					job := &common.Job{Name: jobName}
					jobEvent := common.BuildJobEvent(common.JobEventDelete, job)
					//TODO
					fmt.Println(jobEvent)
				}
			}
		}
	}()
	return nil
}
