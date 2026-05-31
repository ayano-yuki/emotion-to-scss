package domain

type Style struct {
	Name      string
	ClassName string
	CSS       string
	Line      int
}

type VerificationResult struct {
	InputPath string
	SCSSPath  string
	OK        bool
	Reason    string
}
