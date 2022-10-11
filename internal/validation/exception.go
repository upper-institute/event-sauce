package validation

import "google.golang.org/grpc/status"

func FallbackGRPCError(err error, grpcErr error) error {
	if _, ok := status.FromError(err); ok {
		return err
	}

	return grpcErr
}
