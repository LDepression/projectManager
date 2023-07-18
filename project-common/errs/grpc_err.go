/**
 * @Author: lenovo
 * @Description:
 * @File:  grpc_err
 * @Version: 1.0.0
 * @Date: 2023/07/17 11:45
 */

package errs

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	common "projectManager/project-common"
)

func GrpcError(err *BError) error {
	return status.Error(codes.Code(err.Code), err.Msg)
}

func HandleGrpcError(err error) (common.BusinessCode, string) {
	fromError, _ := status.FromError(err)
	return common.BusinessCode(fromError.Code()), fromError.Message()
}
