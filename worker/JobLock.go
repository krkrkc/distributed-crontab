package worker

import (
	"context"
	"github.com/krkrkc/distributed-crontab/common"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type JobLock struct {
	kv         clientv3.KV
	lease      clientv3.Lease
	jobName    string
	cancleFunc context.CancelFunc
	leaseId    clientv3.LeaseID
	isLock     bool
}

func InitJobLock(jobName string, kv clientv3.KV, lease clientv3.Lease) *JobLock {
	return &JobLock{
		kv:      kv,
		lease:   lease,
		jobName: jobName,
		isLock:  false,
	}
}

func (jobLock *JobLock) TryLock() error {
	cancleCtx, cancleFunc := context.WithCancel(context.TODO())
	jobLock.cancleFunc = cancleFunc
	leaseResp, err := jobLock.lease.Grant(context.TODO(), 5)
	jobLock.leaseId = leaseResp.ID
	if err != nil {
		jobLock.freeLease()
		return err
	}

	keepAliveRespChan, err := jobLock.lease.KeepAlive(cancleCtx, leaseResp.ID)
	go func() {
		for {
			select {
			case keepAliveResp := <-keepAliveRespChan:
				if keepAliveResp == nil {
					return
				}
			}
		}
	}()

	txn := jobLock.kv.Txn(context.TODO())
	lockKey := common.JobLockDir + jobLock.jobName

	txn.If(clientv3.Compare(clientv3.CreateRevision(lockKey), "=", 0)).
		Then(clientv3.OpPut(lockKey, "", clientv3.WithLease(leaseResp.ID))).
		Else(clientv3.OpGet(lockKey))

	txnResp, err := txn.Commit()
	if err != nil {
		jobLock.freeLease()
		return err
	}

	if !txnResp.Succeeded {
		jobLock.freeLease()
		return common.ERR_LOCK_ALREADY_REQUIRED
	}
	jobLock.isLock = true
	return nil
}

func (jobLock *JobLock) UnLock() {
	if jobLock.isLock {
		jobLock.freeLease()
	}
}

func (jobLock *JobLock) freeLease() {
	jobLock.cancleFunc()
	jobLock.lease.Revoke(context.TODO(), jobLock.leaseId)
}
