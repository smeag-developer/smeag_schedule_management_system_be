package grpc

import (
	"log"

	"github.com/jinzhu/copier"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ValidateStatusMessage(err error) {
	st, ok := status.FromError(err)

	if ok {
		switch st.Code() {
		case codes.AlreadyExists:
			log.Printf("[grpc] already exists: %v", st.Message())
		case codes.InvalidArgument:
			log.Printf("[grpc] Invalid Argument: %v", st.Message())
		case codes.Unavailable:
			log.Printf("grpc] Service Unavailable: %v", st.Message())
		case codes.DeadlineExceeded:
			log.Printf("grpc] Deadline Exceeded: %v", st.Message())
		case codes.Aborted:
			log.Printf("[grpc] Aborted: %v", st.Message())
		default:
			log.Printf("Unhandled gRPC error: %v", st.Message())
		}
	} else {
		log.Printf("cannot able to create notification: %v", err)
	}
}

func ConvertToProto[TSource any, TDest any](
	sources *[]TSource,
	idSetter func(dest *TDest, src TSource),
) ([]*TDest, error) {
	result := make([]*TDest, 0, len(*sources))

	for _, src := range *sources {
		dest := new(TDest)
		if err := copier.CopyWithOption(dest, src, copier.Option{DeepCopy: true}); err != nil {
			return nil, status.Errorf(codes.Internal, "failed to copy struct: %v", err)
		}
		idSetter(dest, src)
		result = append(result, dest)
	}

	return result, nil
}
