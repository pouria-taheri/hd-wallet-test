package domain

type Request struct {
	DerivationPath
	KeyScope
	// Private returns either a public or private derived extended key
	// based on the flag state
	Private bool `json:"private"`
}
