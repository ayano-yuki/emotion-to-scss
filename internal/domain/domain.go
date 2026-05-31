package domain

type Style struct {
	Name        string
	ClassName   string
	CSS         string
	Line        int
	Expressions []Expression
}

type Expression struct {
	Source      string
	Placeholder string
}

type Verification struct {
	Equivalent bool
	Reason     string
}
