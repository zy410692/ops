package Lib

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strings"
	"unsafe"

	"github.com/google/uuid"
	"github.com/minio/madmin-go/v3"
	iampolicy "github.com/minio/pkg/iam/policy"
	"github.com/zy410692/ops/models"
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
	log.Printf("policy绑定用户和桶 %s-%s", policy_name, newuser)
	if err = madmClnt.SetPolicy(ctx, policy_name, newuser, false); err != nil {

		log.Println(err)
	}

}

func SetPolicyMult(madmClnt *madmin.AdminClient, newbucket []string, newuser string) error {
	//定义全局变量
	//用户添加的桶bucket
	addedBuckets := make(map[string]bool)
	//用户存在的bucket
	existingBuckets := make(map[string]bool)
	//选择不覆盖模式
	override := false
	//添加的桶-切片
	bucket_slice := make([]string, 0)

	//1 首先检查用户现有策略,先默认不覆盖已有策略

	ctx := context.Background()
	info, err := madmClnt.GetUserInfo(ctx, newuser)
	if err != nil {
		return fmt.Errorf("获取用户信息失败: %v", err)
	}

	existingPolicies := info.PolicyName
	if len(existingPolicies) > 0 {
		if !override {
			log.Printf("用户 %s 已有策略 %s，且未选择覆盖", newuser, existingPolicies)
		}
		//2 获取用户已有的桶权限

		policyInfo, err := madmClnt.InfoCannedPolicy(ctx, existingPolicies)
		if err == nil {
			var policy models.MinioPolicy
			if err = json.Unmarshal(policyInfo, &policy); err == nil {
				for _, statement := range policy.Statement {
					for _, resource := range statement.Resource {
						// 提取桶名
						parts := strings.Split(resource, ":")
						if len(parts) > 0 {
							bucketName := strings.TrimSuffix(parts[len(parts)-1], "/*")
							existingBuckets[bucketName] = true
						}
					}
				}
			}
		}

		// 如果选择覆盖，先删除现有策略绑定
		// if err := madmClnt.SetPolicy(ctx, "", newuser, false); err != nil {
		// 	return fmt.Errorf("删除现有策略绑定失败: %v", err)
		// }

		// 尝试删除旧策略（如果是专门为该用户创建的）
		// if strings.HasPrefix(existingPolicies, "policy-") {
		// 	_ = madmClnt.RemoveCannedPolicy(ctx, existingPolicies)
		// }

		// 如果选择覆盖，先删除现有策略绑定
		// if err := madmClnt.SetPolicy(ctx, "", newuser, false); err != nil {
		// 	return fmt.Errorf("删除现有策略绑定失败: %v", err)
		// }

		// 尝试删除旧策略（如果是专门为该用户创建的）
		// if strings.HasPrefix(existingPolicies, "policy-") {
		// 	_ = madmClnt.RemoveCannedPolicy(ctx, existingPolicies)
		// }

		if !override {
			for bucket := range existingBuckets {
				if !addedBuckets[bucket] {
					bucket_slice = append(bucket_slice, fmt.Sprintf("arn:aws:s3:::%s", bucket))
					bucket_slice = append(bucket_slice, fmt.Sprintf("arn:aws:s3:::%s/*", bucket))
					addedBuckets[bucket] = true
				}
			}
		}

	}
	// 3 生成新的策略

	policy_name := newuser + "-mult-" + uuid.New().String()
	log.Printf("一个用户多个桶添加policy%s\n", policy_name)
	// 构建资源列表（合并现有和新的桶）

	//如果不是覆盖模式

	for _, bucket1 := range newbucket {
		log.Println(bucket1)
		if !addedBuckets[bucket1] {
			bucket_slice = append(bucket_slice, fmt.Sprintf("arn:aws:s3:::%s", bucket1))
			bucket_slice = append(bucket_slice, fmt.Sprintf("arn:aws:s3:::%s/*", bucket1))
			addedBuckets[bucket1] = true
		}
	}

	policy := models.MinioPolicy{
		Version: "2012-10-17",
		Statement: []struct {
			Effect   string   `json:"Effect"`
			Action   []string `json:"Action"`
			Resource []string `json:"Resource"`
		}{
			{
				Effect: "Allow",
				Action: []string{"s3:DeleteObject", "s3:GetBucketLocation", "s3:GetObject",
					"s3:ListAllMyBuckets", "s3:ListBucket", "s3:PutObject"},
				Resource: bucket_slice,
			},
		},
	}

	jsonData, err := json.Marshal(policy)
	if err != nil {
		return fmt.Errorf("策略转换为 JSON 失败: %v", err)
	}

	if err = madmClnt.AddCannedPolicy(ctx, policy_name, jsonData); err != nil {
		return fmt.Errorf("添加策略失败: %v", err)
	}

	if err = madmClnt.SetPolicy(ctx, policy_name, newuser, false); err != nil {
		log.Printf("policy绑定用户和桶 %s-%s 失败", policy_name, newuser)
		return fmt.Errorf("绑定策略失败: %v", err)
	}
	return nil

}

func GetMinioAdmClient() *madmin.AdminClient {
	return mdmclient
}
