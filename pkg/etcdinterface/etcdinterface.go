package etcdinterface

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/etcdserver/api/v3rpc/rpctypes"
	"golang.org/x/net/context"
)

const (
	etcdEndpoints = "kless-etcd.kless:2379"
)

// EtcdInterface is just a thin wrapper around etcd for now...
type EtcdInterface struct {
}

// SetValue adds content to etcd
func (e *EtcdInterface) SetValue(key string, value string) error {

	fmt.Printf("Entering SetValue\n")

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{etcdEndpoints},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	fmt.Printf("Got etcd client, performing Put\n")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	_, err = cli.Put(ctx, key, value)
	cancel()
	if err != nil {
		switch err {
		case context.Canceled:
			fmt.Printf("ctx is canceled by another routine: %v\n", err)
		case context.DeadlineExceeded:
			fmt.Printf("ctx is attached with a deadline is exceeded: %v\n", err)
		case rpctypes.ErrEmptyKey:
			fmt.Printf("client-side error: %v\n", err)
		default:
			fmt.Printf("bad cluster endpoints, which are not etcd servers: %v\n", err)
		}
	}

	fmt.Printf("Leaving SetValue\n")

	return err
}

// GetValue retrieves content from etcd
func (e *EtcdInterface) GetValue(key string) (value string, err error) {

	fmt.Printf("Entering GetValue\n")

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"kless-etcd.kless:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	fmt.Printf("Got etcd client, performing Get\n")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	resp, err := cli.Get(ctx, key)
	cancel()
	if err != nil {
		switch err {
		case context.Canceled:
			fmt.Printf("ctx is canceled by another routine: %v\n", err)
		case context.DeadlineExceeded:
			fmt.Printf("ctx is attached with a deadline is exceeded: %v\n", err)
		case rpctypes.ErrEmptyKey:
			fmt.Printf("client-side error: %v\n", err)
		default:
			fmt.Printf("bad cluster endpoints, which are not etcd servers: %v\n", err)
		}
	}

	for _, ev := range resp.Kvs {
		value = string(ev.Value)
		fmt.Printf("Retrieved from etcd key = %s value = %s\n", key, value)
	}

	fmt.Printf("Leaving GetValue\n")

	return value, err
}

// Delete removes content from etcd
func (e *EtcdInterface) Delete(key string) error {

	fmt.Printf("Entering Delete\n")

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"kless-etcd.kless:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	fmt.Printf("Got etcd client, performing Delete\n")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	_, err = cli.Delete(ctx, key)
	cancel()
	if err != nil {
		switch err {
		case context.Canceled:
			fmt.Printf("ctx is canceled by another routine: %v\n", err)
		case context.DeadlineExceeded:
			fmt.Printf("ctx is attached with a deadline is exceeded: %v\n", err)
		case rpctypes.ErrEmptyKey:
			fmt.Printf("client-side error: %v\n", err)
		default:
			fmt.Printf("bad cluster endpoints, which are not etcd servers: %v\n", err)
		}
	}

	fmt.Printf("Leaving Delete\n")

	return err
}

// GetKeysFromPrefix retrieves keys from etcd
func (e *EtcdInterface) GetKeysFromPrefix(prefix string) (keys []string, err error) {

	fmt.Printf("Entering GetKeysFromPrefix\n")

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"kless-etcd.kless:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	fmt.Printf("Got etcd client, performing Get\n")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	resp, err := cli.Get(ctx, prefix, clientv3.WithPrefix())
	cancel()
	if err != nil {
		switch err {
		case context.Canceled:
			fmt.Printf("ctx is canceled by another routine: %v\n", err)
		case context.DeadlineExceeded:
			fmt.Printf("ctx is attached with a deadline is exceeded: %v\n", err)
		case rpctypes.ErrEmptyKey:
			fmt.Printf("client-side error: %v\n", err)
		default:
			fmt.Printf("bad cluster endpoints, which are not etcd servers: %v\n", err)
		}
	}

	var key string
	var keyNoPrefix string
	var value string

	for _, ev := range resp.Kvs {
		key = string(ev.Key)
		keyNoPrefix = strings.TrimPrefix(key, prefix)
		value = string(ev.Value)
		fmt.Printf("Retrieved from etcd prefix = %s key = %s keyNoPrefix = %s value = %s\n", prefix, key, keyNoPrefix, value)

		keys = append(keys, keyNoPrefix)
	}

	fmt.Printf("Leaving GetKeysFromPrefix\n")

	return keys, nil
}

// GetValuesFromPrefix retrieves values from etcd
func (e *EtcdInterface) GetValuesFromPrefix(prefix string) (values []string, err error) {

	fmt.Printf("Entering GetValuesFromPrefix\n")

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"kless-etcd.kless:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	fmt.Printf("Got etcd client, performing Get\n")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	resp, err := cli.Get(ctx, prefix, clientv3.WithPrefix())
	cancel()
	if err != nil {
		switch err {
		case context.Canceled:
			fmt.Printf("ctx is canceled by another routine: %v\n", err)
		case context.DeadlineExceeded:
			fmt.Printf("ctx is attached with a deadline is exceeded: %v\n", err)
		case rpctypes.ErrEmptyKey:
			fmt.Printf("client-side error: %v\n", err)
		default:
			fmt.Printf("bad cluster endpoints, which are not etcd servers: %v\n", err)
		}
	}

	var key string
	var value string

	for _, ev := range resp.Kvs {
		key = string(ev.Key)
		value = string(ev.Value)
		fmt.Printf("Retrieved from etcd prefix = %s key = %s value = %s\n", prefix, key, value)
		values = append(values, value)
	}

	fmt.Printf("Leaving GetValuesFromPrefix\n")

	return values, nil
}

// GetKeysValuesFromPrefix retrieves keys and values from etcd
func (e *EtcdInterface) GetKeysValuesFromPrefix(prefix string) (values map[string]string, err error) {

	fmt.Printf("Entering GetValuesFromPrefix\n")

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"kless-etcd.kless:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	fmt.Printf("Got etcd client, performing Get\n")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	resp, err := cli.Get(ctx, prefix, clientv3.WithPrefix())
	cancel()
	if err != nil {
		switch err {
		case context.Canceled:
			fmt.Printf("ctx is canceled by another routine: %v\n", err)
		case context.DeadlineExceeded:
			fmt.Printf("ctx is attached with a deadline is exceeded: %v\n", err)
		case rpctypes.ErrEmptyKey:
			fmt.Printf("client-side error: %v\n", err)
		default:
			fmt.Printf("bad cluster endpoints, which are not etcd servers: %v\n", err)
		}
	}

	values = make(map[string]string)

	var key string
	var value string

	for _, ev := range resp.Kvs {
		key = string(ev.Key)
		value = string(ev.Value)
		fmt.Printf("Retrieved from etcd prefix = %s key = %s value = %s\n", prefix, key, value)
		values[strings.TrimPrefix(key, prefix)] = value
	}

	fmt.Printf("Leaving GetKeysValuesFromPrefix\n")

	return values, nil
}
