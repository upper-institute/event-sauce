package validation

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	InvalidVersionErr        = status.Error(codes.InvalidArgument, "invalid version, must be >= 1")
	InvalidIDErr             = status.Error(codes.InvalidArgument, "invalid event ID, must be an UUID")
	InconsistentIDErr        = status.Error(codes.InvalidArgument, "the event array have inconsistent IDs, all events of the slice must have the same ID")
	InconsistentVersionErr   = status.Error(codes.InvalidArgument, "the event array have inconsistent version, they must be sequential in the ascending order")
	EmptyEventArrayErr       = status.Error(codes.InvalidArgument, "the input event array is empty")
	BackendAppendErr         = status.Error(codes.Internal, "the backend returned an internal error (EventStore.Append)")
	BackendLatestErr         = status.Error(codes.Internal, "the backend returned an internal error (EventStore.Latest)")
	BackendSetErr            = status.Error(codes.Internal, "the backend returned an internal error (SnapshotStore.Set)")
	BackendGetErr            = status.Error(codes.Internal, "the backend returned an internal error (SnapshotStore.Get)")
	ProtoMarshalErr          = status.Error(codes.Internal, "error while marshaling proto message")
	LatestVersionMismatchErr = status.Error(codes.FailedPrecondition, "latest version does not match the event version")
	SnapshotNotFoundErr      = status.Error(codes.NotFound, "snapshot not found")
	LatestNotFoundErr        = status.Error(codes.NotFound, "latest version not found")
)
