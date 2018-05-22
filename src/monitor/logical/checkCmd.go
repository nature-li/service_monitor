package logical

type CheckCmd struct {
	id         int
	localCheck bool
	checkShell string
	operator   string
	checkValue string
	goodMatch  bool
}
