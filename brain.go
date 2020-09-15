package ctypes

type BrainExecuteRequest struct {
	Executable *Executable       `json:"executable"`
	Request    *ExecutionRequest `json:"request"`
}

type BrainExecuteResult struct {
	*ExecutionResult
	FinalTree *Context               `json:"final_tree"`
	FinalData map[string]interface{} `json:"final_data"`
}
