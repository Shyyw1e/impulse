package domain

type TrialStatus string

const(
	StatusSuccess TrialStatus = "SUCCESS"
	StatusFail TrialStatus = "FAIL"
	StatusDisqual TrialStatus = "DISQUAL"
)