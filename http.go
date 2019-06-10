package randgenloader

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/sessions"
	"net/http"
)

var store = sessions.NewCookieStore([]byte(""))

func init() {
	// session 过期时间  24h
	store.MaxAge(24 * 60 * 60)
}

const sessionName = "randgen"
const loaderKey = "loader"

type resultStruct struct {
	Resp      interface{}
	ErrorInfo string
}

func resultSuccess(resp interface{}) *resultStruct {
	return &resultStruct{resp, ""}
}

func resultError(errorInfo string) *resultStruct {
	return &resultStruct{nil, errorInfo}
}

func MustJosnMarshal(v interface{}) []byte {
	bytes, _ := json.Marshal(v)
	return bytes
}

// paramKeys: parameters must have and can not empty in request
// mustAccessFirst: the url must access first
// business:   business logic func must return a resultStruct
func base(w http.ResponseWriter, r *http.Request, paramKeys []string, mustAccessFirst string, business func(params []string, session *sessions.Session) *resultStruct) {
	session, err := store.Get(r, sessionName)
	if err != nil {
		panic("session decode error")
	}

	if mustAccessFirst != "" && session.IsNew {
		w.Write(MustJosnMarshal(resultError(
			fmt.Sprintf("please access %q first", mustAccessFirst))))
		return
	}

	r.ParseForm()

	paramValues := make([]string, len(paramKeys))

	// 参数非空校验
	for i, paramKey := range paramKeys {
		if len(r.Form[paramKey]) == 0 {
			w.Write(MustJosnMarshal(resultError(
				fmt.Sprintf("param %q required", paramKey))))
			return
		}
		value := r.Form[paramKey][0]
		if value == "" {
			w.Write(MustJosnMarshal(resultError(
				fmt.Sprintf("param %q can not be empty", paramKey))))
			return
		}
		paramValues[i] = value
	}

	resp := business(paramValues, session)

	session.Save(r, w)

	w.Write(MustJosnMarshal(resp))
}

func Init(w http.ResponseWriter, r *http.Request) {
	base(w, r, []string{"testname"}, "",
		func(params []string, session *sessions.Session) *resultStruct {
			testname := params[0]
			loader := &RandgenLoader{}
			loader.Init(testname)
			session.Values[loaderKey] = loader
			return resultSuccess("OK")
		})
}

func LoadData(w http.ResponseWriter, r *http.Request) {
	base(w, r, []string{"zz", "yy"}, "/init",
		func(params []string, session *sessions.Session) *resultStruct {
			zzContent := params[0]
			yyContent := params[1]
			loader := session.Values[loaderKey].(*RandgenLoader)
			sqls, err := loader.LoadData(zzContent, yyContent)
			if err != nil {
				return resultError(err.Error())
			}
			return resultSuccess(sqls)
		})
}

func Query(w http.ResponseWriter, r *http.Request) {
	base(w, r, nil, "/init",
		func(params []string, session *sessions.Session) *resultStruct {
			loader := session.Values[loaderKey].(RandgenLoader)
			if loader.cachedQueries == nil {
				return resultError("please access '/loaddata' first")
			}
			return resultSuccess(loader.Query())
		})
}

func Compare(w http.ResponseWriter, r *http.Request) {
	base(w, r, []string{"mysql", "tidb"}, "/init",
		func(params []string, session *sessions.Session) *resultStruct {
			r1 := params[0]
			r2 := params[1]
			loader := session.Values[loaderKey].(RandgenLoader)
			comment, consistent := loader.Compare(r1, r2)
			return resultSuccess(struct {
				string
				bool
			}{comment, consistent})
		})
}
