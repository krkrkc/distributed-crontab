package master

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/krkrkc/distributed-crontab/common"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

type JobMgr struct {
	client *clientv3.Client
	kv     clientv3.KV
	lease  clientv3.Lease
}

var (
	G_jobMgr *JobMgr
)

func InitJobMgr() (err error) {
	var (
		config clientv3.Config
		client *clientv3.Client
		kv     clientv3.KV
		lease  clientv3.Lease
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

	G_jobMgr = &JobMgr{
		client: client,
		kv:     kv,
		lease:  lease,
	}

	return
}

func (j *JobMgr) SaveJob(job *common.Job) (oldJob *common.Job, err error) {
	var (
		jobKey   string
		jobValue []byte
		putResp  *clientv3.PutResponse
	)

	jobKey = "/cron/jobs/" + job.Name
	if jobValue, err = json.Marshal(job); err != nil {
		return
	}
	if putResp, err = j.kv.Put(context.TODO(), jobKey, string(jobValue), clientv3.WithPrevKV()); err != nil {
		return
	}

	if putResp.PrevKv != nil {
		json.Unmarshal(putResp.PrevKv.Value, &oldJob)
	}
	return
}

func (j *JobMgr) DeleteJob(jobName string) (*common.Job, error) {
	jobKey := common.JobSaveDir + jobName
	deleteResp, err := j.kv.Delete(context.TODO(), jobKey, clientv3.WithPrevKV())
	if err != nil {
		return nil, err
	}

	if len(deleteResp.PrevKvs) > 0 {
		var oldJob common.Job
		err = json.Unmarshal(deleteResp.PrevKvs[0].Value, &oldJob)
		if err != nil {
			return nil, err
		}

		return &oldJob, nil
	}

	return nil, nil
}
