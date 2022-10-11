package validation

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	InvalidVersionErr               = status.Error(codes.InvalidArgument, "invalid version, must be >= 1")
	InvalidIDErr                    = status.Error(codes.InvalidArgument, "invalid event ID, must be an UUID")
	InconsistentIDErr               = status.Error(codes.InvalidArgument, "the event array have inconsistent IDs, all events of the slice must have the same ID")
	InconsistentVersionErr          = status.Error(codes.InvalidArgument, "the event array have inconsistent version, they must be sequential in the ascending order")
	EmptyEventArrayErr              = status.Error(codes.InvalidArgument, "the input event array is empty")
	BackendAppendErr                = status.Error(codes.Internal, "the backend returned an internal error (EventStore.Append)")
	BackendScanErr                  = status.Error(codes.Internal, "the backend returned an internal error (EventStore.Scan)")
	BackendLatestErr                = status.Error(codes.Internal, "the backend returned an internal error (EventStore.Latest)")
	BackendSetErr                   = status.Error(codes.Internal, "the backend returned an internal error (SnapshotStore.Set)")
	BackendGetErr                   = status.Error(codes.Internal, "the backend returned an internal error (SnapshotStore.Get)")
	ProtoMarshalErr                 = status.Error(codes.Internal, "error while marshaling proto message")
	LatestVersionMismatchErr        = status.Error(codes.FailedPrecondition, "latest version does not match the event version")
	SnapshotNotFoundErr             = status.Error(codes.NotFound, "snapshot not found")
	LatestNotFoundErr               = status.Error(codes.NotFound, "latest version not found")
	InvalidQueryNaturalTimestampErr = status.Error(codes.InvalidArgument, "invalid natural timestamp in the query")
	MissingStartQueryErr            = status.Error(codes.InvalidArgument, "missing start query")
	MissingEndQueryErr              = status.Error(codes.InvalidArgument, "missing start query")
	InvalidVersionRangeErr          = status.Error(codes.InvalidArgument, "end version must be greater than or equal to start version")
	InvalidNaturalTimestampRangeErr = status.Error(codes.InvalidArgument, "end natural timestamp must be after or equal to start natural timestamp")
	InconsistentRangeErr            = status.Error(codes.InvalidArgument, "start and ent parameters must be defined in the same type, as version or as natural timestamp")
)
