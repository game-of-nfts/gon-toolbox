package verifier

type (
	Request struct {
		TaskNo string
		User   UserInfo
		Params any
	}

	Respone struct {
		TaskNo   string
		TeamName string
		Point    int32
		Reason   string
	}

	Verifier interface {
		Do(req Request, res chan<- *Respone)
		BuildParams(params [][]string) (any, error)
	}

	UserInfo struct {
		TeamName string
		Github   string
		Address  map[string]string
	}
)
