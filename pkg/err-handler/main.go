package errhandler

import "log"

type ErrorRslt struct {
	RsltCd  string
	RsltMsg string
}

func PanicErr(err error) {
	if err != nil {
		log.Println(err)
		panic(err)
	}
}

func AuthorizedErr(req ...interface{}) ErrorRslt {
	//toDo : logging
	log.Println(req...)
	return ErrorRslt{
		RsltCd:  "01",
		RsltMsg: "Authorized",
	}
}

func NotFoundErr(req ...interface{}) ErrorRslt {
	//toDo : logging
	log.Println(req...)
	return ErrorRslt{
		RsltCd:  "04",
		RsltMsg: "Not Found",
	}
}

func ConflictErr(req ...interface{}) ErrorRslt {
	//toDo : logging
	log.Println(req...)
	return ErrorRslt{
		RsltCd:  "03",
		RsltMsg: "Conflict",
	}
}

func ForbiddenErr(req ...interface{}) ErrorRslt {
	//toDo : logging
	log.Println(req...)
	return ErrorRslt{
		RsltCd:  "09",
		RsltMsg: "Forbidden",
	}
}

func BadRequestErr(req ...interface{}) ErrorRslt {
	log.Println(req...)
	return ErrorRslt{
		RsltCd:  "99",
		RsltMsg: "Bad Request",
	}
}
