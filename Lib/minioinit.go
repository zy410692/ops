package Lib

import (
	"context"
	"fmt"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const (
	useSSL bool = false
)

var (
	client *minio.Client
	err    error
)

func InitMinio(endpoint string, accessKeyID string, secretAccessKey string) {
	minioclient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL})
	if err != nil {
		log.Println("minio连接错误: ", err)
	}
	log.Printf("%#v\n", client)
	client = minioclient

}

func CreateBucket(bucketName string) {

	exists, _ := client.BucketExists(context.TODO(), bucketName)
	if exists {
		log.Printf("bucket: %s已经存在", bucketName)

	} else {

		err = client.MakeBucket(context.TODO(), bucketName, minio.MakeBucketOptions{Region: "cn-south-1", ObjectLocking: false})
		if err != nil {
			log.Println(err)
		}
		log.Printf("Successfully created %s\n", bucketName)
	}
}

func listBucket() {
	buckets, _ := client.ListBuckets(context.Background())
	for _, bucket := range buckets {
		fmt.Println(bucket)
	}
}

func GetMinioClient() *minio.Client {
	return client
}
