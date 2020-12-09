package handlers

type CheckResult struct {
	ExitCode int    `json:"exit_code"`
	StdOut   string `json:"stdout"`
	Error    string `json:"error"`
}


