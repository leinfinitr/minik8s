package etcd

import (
	"context"
	"fmt"
	"time"

	etcd "go.etcd.io/etcd/client/v3"
)
type EtcdClientWrapper struct {
	etcdClient *etcd.Client
}
func NewEtcdClient(endpoints []string,timeout time.Duration) (*EtcdClientWrapper, error) {
	cli,err := etcd.New(etcd.Config{
		Endpoints: endpoints,
		DialTimeout: timeout,
	})
	if err != nil {
		return nil,fmt.Errorf("etcd.New err:%v",err)
	}
	timeoutCtx ,cancel := context.WithTimeout(context.Background(),timeout)
	defer cancel()
	_,err = cli.Status(timeoutCtx,endpoints[0])
	if err != nil {
		return nil,fmt.Errorf("cli.Status err:%v",err)
	}
	return &EtcdClientWrapper{etcdClient:cli},nil
}

func (c *EtcdClientWrapper) Put(key,value string) error {
	ctx := context.Background()
	_,err := c.etcdClient.Put(ctx,key,value)
	if err != nil {
		return fmt.Errorf("cli.Put err:%v",err)
	}
	return nil
}

func (c *EtcdClientWrapper) Get(key string) (string,error) {
	ctx := context.Background()
	resp,err := c.etcdClient.Get(ctx,key)
	if err != nil {
		return "",fmt.Errorf("cli.Get err:%v",err)
	}
	if len(resp.Kvs) == 0 {
		return "",nil
	}
	return string(resp.Kvs[0].Value),nil
}

func (c *EtcdClientWrapper) Delete(key string) error {
	ctx := context.Background()
	_,err := c.etcdClient.Delete(ctx,key)
	if err != nil {
		return fmt.Errorf("cli.Delete err:%v",err)
	}
	return nil
}

func (c *EtcdClientWrapper) Watch(key string) etcd.WatchChan {
	return c.etcdClient.Watch(context.Background(),key)
}

func (c *EtcdClientWrapper) PrefixGet(key string) ([]string,error) {
	ctx := context.Background()
	resp,err := c.etcdClient.Get(ctx,key,etcd.WithPrefix())
	if err != nil {
		return nil,fmt.Errorf("cli.Get err:%v",err)
	}
	var values []string
	for _,kv := range resp.Kvs {
		values = append(values,string(kv.Value))
	}
	return values,nil
}