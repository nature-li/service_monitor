package job

type CheckCmd struct {
	id         int32
	serviceId  int32
	localCheck bool
	checkShell string
	operator   string
	checkValue string
	goodMatch  bool
}

func NewCheckCmd(id, serviceId, localCheck int32, checkShell, operator, checkValue string, goodMatch int) *CheckCmd {
	return &CheckCmd{
		id:         id,
		serviceId:  serviceId,
		localCheck: localCheck == 1,
		checkShell: checkShell,
		operator:   operator,
		checkValue: checkValue,
		goodMatch:  goodMatch == 1,
	}
}
