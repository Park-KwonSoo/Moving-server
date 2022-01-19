package errhandler

type ErrorRslt struct {
	RsltCd  string
	RsltMsg string
}

func AuthorizedErr() ErrorRslt {
	return ErrorRslt{
		RsltCd:  "01",
		RsltMsg: "Authorized",
	}
}

func NotFoundErr() ErrorRslt {
	return ErrorRslt{
		RsltCd:  "04",
		RsltMsg: "Not Found",
	}
}

func ConflictErr() ErrorRslt {
	return ErrorRslt{
		RsltCd:  "09",
		RsltMsg: "Conflict",
	}
}
