package Lib

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"strings"
	"unsafe"

	"github.com/google/uuid"
	"github.com/minio/madmin-go/v3"
	iampolicy "github.com/minio/pkg/iam/policy"
)

var (
	mdmclient *madmin.AdminClient
)

func Initialize(url string, miniouser string, miniopasswd string) {
	// Use a secure connection.
	ssl := false

	// Initialize minio client object.
	mdmClnt, err := madmin.New(url, miniouser, miniopasswd, ssl)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Fetch service status.
	st, err := mdmClnt.ServerInfo(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	log.Println(st.Servers)

	mdmclient = mdmClnt

}
func ExistUser(madmClnt *madmin.AdminClient, user string) bool {
	users, err := madmClnt.ListUsers(context.Background())
	if err != nil {
		log.Fatalln(err)
	}
	_, ok := users[user]
	if ok {
		return true
	} else {
		return false
	}

}

func CreateUser(madmClnt *madmin.AdminClient, newuser string, password string) {

	// if err = madmClnt.AddUser(context.Background(), newuser, "6sbT.05~"); err != nil {
	// 	log.Println(err)
	// }
	//更改为随机密码
	if err = madmClnt.AddUser(context.Background(), newuser, password); err != nil {
		log.Println(err)
	}
}

func MyStructToBytes(s *iampolicy.Policy) []byte {
	var sizeofPolicy = int(unsafe.Sizeof(iampolicy.Policy{}))
	var x reflect.SliceHeader
	x.Len = sizeofPolicy
	x.Cap = sizeofPolicy
	x.Data = uintptr(unsafe.Pointer(s))
	return *(*[]byte)(unsafe.Pointer(&x))
}

func SetPolicySingle(madmClnt *madmin.AdminClient, newbucket string, newuser string) {
	json_base := `
	{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Effect": "Allow",
				"Action": [
					"s3:DeleteObject",
					"s3:GetBucketLocation",
					"s3:GetObject",
					"s3:ListAllMyBuckets",
					"s3:ListBucket",
					"s3:PutObject"
				],
				"Resource": [
					"arn:aws:s3:::my-bucketname/*"
				]
			}
		]
	}
	`

	policy_name := newbucket + uuid.New().String()
	body_json := strings.Replace(json_base, "my-bucketname", newbucket, 1)

	// policy, err := iampolicy.ParseConfig(strings.NewReader(body_json))
	// log.Println(body_json)
	// if err != nil {
	// 	log.Println("解析json")
	// 	log.Println(err)
	// }
	log.Printf("添加policy%s\n", policy_name)
	ctx := context.Background()
	if err = madmClnt.AddCannedPolicy(ctx, policy_name, []byte(body_json)); err != nil {

		log.Println(err)
	}
	log.Println("policy绑定用户和桶 %s-%s", policy_name, newuser)
	if err = madmClnt.SetPolicy(ctx, policy_name, newuser, false); err != nil {

		log.Println(err)
	}

}

func GetMinioAdmClient() *madmin.AdminClient {
	return mdmclient
}
