package ctypes

// API Error Codes
const (
	ErrGenericError             = 0
	ErrResourceNotFound         = 404
	ErrInputValidation          = 501
	ErrDatabaseIssue            = 405
	ErrRedisFailure             = 406
	ErrGrooveFailure            = 407
	ErrNotAuthenticated         = 408
	ErrInvalidToken             = 409
	ErrMongoFailure             = 410
	ErrRateLimitExceeded        = 429
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
	ErrInsufficientPermissions  = 855
	ErrMissingOrgHeader         = 901
	ErrMissingBotHeader         = 902
	ErrMissingEnvHeader         = 903
	ErrInvalidHeaderFormat      = 904
)
