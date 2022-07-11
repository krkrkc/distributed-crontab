package worker

import (
	"context"
	"fmt"
	"github.com/krkrkc/distributed-crontab/common"
	clientv3 "go.etcd.io/etcd/client/v3"
	"net"
	"time"
)

type Register struct {
	client  *clientv3.Client
	kv      clientv3.KV
	lease   clientv3.Lease
	localIP string
}

var (
	G_register *Register
)

func InitRegister() error {
	config := clientv3.Config{
		Endpoints:   G_config.EtcdEndpoints,
		DialTimeout: time.Duration(G_config.EtcdDialTimeout) * time.Millisecond,
	}

	localIP, err := getLocalIp()
	if err != nil {
		return err
	}

	client, err := clientv3.New(config)
	if err != nil {
		fmt.Println(err)
		return err
	}

	kv := clientv3.NewKV(client)
	lease := clientv3.NewLease(client)

	G_register = &Register{
		client:  client,
		kv:      kv,
		lease:   lease,
		localIP: localIP,
	}
	go G_register.keepOnline()
	return nil
}

func getLocalIp() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, addr := range addrs {
		ipNet, isIpNet := addr.(*net.IPNet)
		if isIpNet && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String(), nil
			}
		}
	}
	return "", common.ERR_NO_LOCAL_IP_FOUND
}

func (register *Register) keepOnline() {
	regKey := common.JobWorkerDir + register.localIP
	for {
		leaseGrantResp, err := register.lease.Grant(context.TODO(), 10)
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}

		cancleCtx, cancleFunc := context.WithCancel(context.TODO())
		keepAliveChan, err := register.lease.KeepAlive(cancleCtx, leaseGrantResp.ID)
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}

		_, err = register.kv.Put(context.TODO(), regKey, "", clientv3.WithLease(leaseGrantResp.ID))
		if err != nil {
			cancleFunc()
			time.Sleep(1 * time.Second)
			continue
		}

		for {
			select {
			case keepAliveResp := <-keepAliveChan:
				if keepAliveResp != nil {
					break
				}
			}
		}
	}
}
