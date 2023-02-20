package main

import (
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"os"
	"strings"
)

var (
	file       string
	fileKey    string
	overWrite  bool
	region     string
	bucketName string
	dir        string
	operate    string
)

func init() {

	flag.BoolVar(&overWrite, "ow", false, "Overwrite if s3 has the same file object")
	flag.StringVar(&file, "f", "", "File name or dir name")
	flag.StringVar(&region, "r", "", "Aws region")
	flag.StringVar(&bucketName, "b", "", "Aws bucket name")
	flag.StringVar(&fileKey, "k", "", "Upload full path at bucket")
	flag.StringVar(&dir, "d", "", "Dir Name")
	flag.StringVar(&operate, "op", "", "Operate")
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}
	flag.Parse()
	if operate == "" {
		fmt.Println("Need param -op\n Suppose: list, upload")
	}
	switch operate {
	case "list":
		if dir == "" {
			fmt.Println("Need param -d\n S3 object key prefix, usually as folder name.")
			return
		}
		if region == "" {
			fmt.Println("Need param -r\n Aws region")
			return
		}
		if bucketName == "" {
			fmt.Println("Need param -b\n bucket name")
			return
		}
		var awsConf aws.Config
		awsConf.Credentials = credentials.NewStaticCredentials("", "", "")
		sess, err := session.NewSession()
		if err != nil {
			fmt.Println("New aws session err:", err)
			return
		}
		sess.Config.Region = &region
		s3Svc := s3.New(sess, nil)
		out, err := s3Svc.ListObjects(&s3.ListObjectsInput{
			Bucket: &bucketName,
			Prefix: &dir,
		})
		if err != nil {
			fmt.Println("List Object err:", err)
			return
		}
		if out.Contents == nil || len(out.Contents) == 0 {
			fmt.Println("No found data in this key, maybe this is a empty folder or this folder doesn't exist.")
			return
		}
		fmt.Println("Success")
		res := make([]string, 0)
		for _, p := range out.Contents {
			res = append(res, *p.Key)
		}
		fmt.Println(strings.Join(res, "\n"))

	case "upload":

		if region == "" {
			fmt.Println("Need param -r\n Aws region")
			return
		}
		if bucketName == "" {
			fmt.Println("Need param -b\n bucket name")
			return
		}
		if file == "" {
			fmt.Println("Need param -f\n File full path")
			return
		}
		var awsConf aws.Config
		awsConf.Credentials = credentials.NewStaticCredentials("", "", "")
		sess, err := session.NewSession()
		if err != nil {
			fmt.Println("New aws session err:", err)
			return
		}
		sess.Config.Region = &region
		key := dir + "/" + getFileName(file)
		uploader := s3manager.NewUploader(sess)
		f, err := os.Open(file)
		if err != nil {
			fmt.Println("Open file err:", err)
			return
		}
		defer f.Close()
		_, err = uploader.Upload(&s3manager.UploadInput{
			Body:   f,
			Bucket: &bucketName,
			Key:    &key,
		})
		if err != nil {
			fmt.Println("Upload file err")
			return
		}
		fmt.Println("Success")
	default:
		fmt.Println("UnSuppose param: -op")
		return
	}
}

func printUsage() {

}

func getFileName(path string) string {
	if strings.Contains(path, "/") {
		sli := strings.Split(path, "/")
		return sli[len(sli)-1]
	}
	return path
}
