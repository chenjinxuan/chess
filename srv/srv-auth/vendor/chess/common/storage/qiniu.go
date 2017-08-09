package storage

import (
	"bytes"
	"golang.org/x/net/context"
	"qiniupkg.com/api.v7/kodo"
)

var Qiniu = new(QiniuSdk)

type QiniuSdk struct{}

func (q *QiniuSdk) SetMac(ak, sk string) {
	kodo.SetMac(ak, sk)
}

func (q *QiniuSdk) GenUptoken(bucket string) string {
	//创建一个Client
	c := kodo.New(0, nil)

	//设置上传的策略
	policy := &kodo.PutPolicy{
		Scope: bucket,
		//设置Token过期时间
		Expires: 3600,
	}
	//生成一个上传token
	return c.MakeUptoken(policy)
}

func (q *QiniuSdk) Put(bucketName, key string, data []byte) error {
	//创建一个Client
	c := kodo.New(0, nil)

	bucket := c.Bucket(bucketName)
	ctx := context.Background()

	reader := bytes.NewReader(data)

	err := bucket.Put(ctx, nil, key, reader, int64(reader.Len()), nil)
	return err
}
