package utils

import (
	"buf.build/go/protovalidate"
	"errors"
	"github.com/aldngrha/ecommerce-be/pb/common"
	"google.golang.org/protobuf/proto"
)

func CheckValidations(req proto.Message) ([]*common.ValidationError, error) {
	if err := protovalidate.Validate(req); err != nil {
		var validationError *protovalidate.ValidationError
		if errors.As(err, &validationError) {
			var validationErrorResponse []*common.ValidationError = make([]*common.ValidationError, 0)
			for _, violation := range validationError.Violations {
				validationErrorResponse = append(validationErrorResponse, &common.ValidationError{
					Field:   *violation.Proto.Field.Elements[0].FieldName,
					Message: *violation.Proto.Message,
				})
			}
			return validationErrorResponse, nil
		}
		return nil, err
	}
	return nil, nil
}
