package awslambda

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/mholt/archiver"
)

// MemorySizes indicates what sizes are supported by AWS Lambda
var MemorySizes = []int64{128, 256, 512, 1024, 1536}

// Result stores performance information
type Result struct {
	CPUMs int64 `json:"cpu_ms"`
}

// Results maps memory size to result
type Results map[int64]*Result

// Benchmark run banchmarks
func Benchmark() (Results, error) {
	sess := session.Must(session.NewSession())
	iamsvc := iam.New(sess)
	lambdasvc := lambda.New(sess)

	role, err := createRole(iamsvc)
	if err != nil {
		return nil, fmt.Errorf("error creating role: %s", err.Error())
	}
	defer deleteRole(iamsvc)

	err = createFunction(lambdasvc, role)
	if err != nil {
		return nil, fmt.Errorf("error creating function: %s", err.Error())
	}
	defer deleteFunction(lambdasvc)

	res := Results{}
	for _, size := range MemorySizes {
		fmt.Printf("%dMB memory size...", size)

		_, err := lambdasvc.UpdateFunctionConfiguration(&lambda.UpdateFunctionConfigurationInput{
			FunctionName: aws.String(functionName),
			MemorySize:   aws.Int64(size),
		})
		if err != nil {
			return nil, err
		}

		result, err := lambdasvc.Invoke(&lambda.InvokeInput{FunctionName: aws.String(functionName)})
		if err != nil {
			return nil, err
		}

		res[size] = &Result{}
		err = json.Unmarshal(result.Payload, res[size])
		if err != nil {
			return nil, err
		}

		fmt.Println(" done")
	}

	return res, nil
}

const (
	roleName     = "faasperf"
	functionName = "faasperf"
)

func createRole(svc *iam.IAM) (string, error) {
	result, err := svc.CreateRole(&iam.CreateRoleInput{
		AssumeRolePolicyDocument: aws.String(`{
			"Version": "2012-10-17",
			"Statement": [{
				"Effect": "Allow",
				"Principal": { "Service": "lambda.amazonaws.com" },
				"Action": "sts:AssumeRole"
			}]
		}`),
		Path:     aws.String("/"),
		RoleName: aws.String(roleName),
	})

	if err != nil {
		return "", err
	}

	return *result.Role.Arn, nil
}

func deleteRole(svc *iam.IAM) error {
	_, err := svc.DeleteRole(&iam.DeleteRoleInput{
		RoleName: aws.String(roleName),
	})

	return err
}

func createFunction(svc *lambda.Lambda, role string) error {
	buf := &bytes.Buffer{}

	err := archiver.Zip.Write(buf, []string{"./awslambda/handler/handler.js", "./awslambda/handler/node_modules"})
	if err != nil {
		return err
	}

	_, err = svc.CreateFunction(&lambda.CreateFunctionInput{
		Code:         &lambda.FunctionCode{ZipFile: buf.Bytes()},
		FunctionName: aws.String(functionName),
		Handler:      aws.String("handler.perf"),
		MemorySize:   aws.Int64(MemorySizes[0]),
		Role:         aws.String(role),
		Runtime:      aws.String("nodejs6.10"),
		Timeout:      aws.Int64(10),
	})

	return err
}

func deleteFunction(svc *lambda.Lambda) error {
	_, err := svc.DeleteFunction(&lambda.DeleteFunctionInput{
		FunctionName: aws.String(functionName),
	})

	return err
}
