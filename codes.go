package ctypes

// API Error Codes
const (
	ErrResourceNotFound         = 404
	ErrInputValidation          = 501
	ErrDatabaseIssue            = 405
	ErrRedisFailure             = 406
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
