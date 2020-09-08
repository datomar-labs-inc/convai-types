package ctypes

// API Error Codes
const (
	ErrResourceNotFound         = 404
	ErrDatabaseIssue            = 405
	ErrDispatchUnmarshal        = 466
	ErrDispatchNotFound         = 467
	ErrHandlerFailure           = 515
	ErrTemplateExecutionFailure = 516
	ErrConfigUnmarshal          = 468
	ErrInvalidConfig            = 469
	ErrInvalidPath              = 470
	ErrFailedToParseArgument    = 480
	ErrFailedToCallPackage      = 481
	ErrPackageMissingLink       = 482
)
