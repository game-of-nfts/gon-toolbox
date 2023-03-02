package verifier

const (
	ChainIdAbbreviationIris     = "i"
	ChainIdAbbreviationStars    = "s"
	ChainIdAbbreviationJuno     = "j"
	ChainIdAbbreviationUptick   = "u"
	ChainIdAbbreviationOmniflix = "o"
)

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
	}

	UserInfo struct {
		TeamName string
		Github   string
		Address  map[string]string
	}
)
