package config

var Public = new(PublicConfig)

type PublicConfig struct {
	QiniuAccessKey          string
	QiniuSecretKey          string
	QiniuBucketJiafenImages string
	QiniuBucketJiafenCard   string
	QiniuDomainJiafenImages string
	QiniuDomainJiafenCard   string
}

func (c *PublicConfig) Import() error {
	var err error
	c.QiniuAccessKey, err = ConsulClient.Key("public/qiniu/access_key", "")
	if err != nil {
		return err
	}

	c.QiniuAccessKey, err = ConsulClient.Key("public/qiniu/secret_key", "")
	if err != nil {
		return err
	}

	c.QiniuBucketJiafenCard, err = ConsulClient.Key("public/qiniu/bucket_jiafen-card", "")
	if err != nil {
		return err
	}

	c.QiniuBucketJiafenImages, err = ConsulClient.Key("public/qiniu/bucket_jiafen-images", "")
	if err != nil {
		return err
	}

	c.QiniuDomainJiafenCard, err = ConsulClient.Key("public/qiniu/domain_jiafen-card", "")
	if err != nil {
		return err
	}

	c.QiniuDomainJiafenImages, err = ConsulClient.Key("public/qiniu/domain_jiafen-images", "")
	if err != nil {
		return err
	}

	return nil
}
