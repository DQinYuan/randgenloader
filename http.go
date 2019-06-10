package randgenloader

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"github.com/alexedwards/scs/v2"
)

var session *scs.Session

func init() {
	// session 过期时间  24h
	session = scs.NewSession()
	session.Lifetime = 24 * time.Hour

	gob.Register(RandgenLoader{})
}

func StartServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/init", Init)
	mux.HandleFunc("/loaddata", LoadData)
	mux.HandleFunc("/query", Query)
	mux.HandleFunc("/compare", Compare)
	err := http.ListenAndServe(":9080", session.LoadAndSave(mux))
	if err != nil {
		fmt.Println(err)
	}
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

// if bool return true, string will be value
// else  string will be info
func checkReqFormNotEmpty(r *http.Request, key string) (string, bool) {
	if len(r.Form[key]) == 0 {
		return fmt.Sprintf("param %q required", key), false
	}
	value := r.Form[key][0]
	if value == "" {
		return fmt.Sprintf("param %q can not be empty", key), false
	}

	return value, true
}

// paramKeys: parameters must have and can not empty in request
// optionalKeys:   parameters => default value
// mustAccessFirst: the url must access first
// business:   business logic func must return a resultStruct
func base(w http.ResponseWriter, r *http.Request, paramKeys []string,
	optionalKeys map[string]string, mustAccessFirst string,
	business func(params []string, optionals map[string]string, r *http.Request) *resultStruct) {

	sessionExists := session.Exists(r.Context(), loaderKey)

	if mustAccessFirst != "" && !sessionExists {
		w.Write(MustJosnMarshal(resultError(
			fmt.Sprintf("please access %q first", mustAccessFirst))))
		return
	}

	r.ParseForm()

	paramValues := make([]string, len(paramKeys))

	// 参数非空校验
	for i, paramKey := range paramKeys {
		v, ok := checkReqFormNotEmpty(r, paramKey)
		if !ok {
			w.Write(MustJosnMarshal(resultError(v)))
			return
		}
		paramValues[i] = v
	}

	for k, _ := range optionalKeys {
		userValue, ok := checkReqFormNotEmpty(r, k)
		if ok {
			optionalKeys[k] = userValue
		}
	}

	resp := business(paramValues, optionalKeys, r)

	w.Write(MustJosnMarshal(resp))
}

func Init(w http.ResponseWriter, r *http.Request) {
	base(w, r, []string{"testname"}, nil, "",
		func(params []string, optionals map[string]string, request *http.Request) *resultStruct {
			testname := params[0]
			loader := &RandgenLoader{}
			loader.Init(testname)
			session.Put(r.Context(), loaderKey, loader)
			return resultSuccess("OK")
		})
}

func LoadData(w http.ResponseWriter, r *http.Request) {
	const db = "db"
	const queries = "queries"
	base(w, r, []string{"zz", "yy"},
		map[string]string{db: "test", queries: "1000",},"/init",
		func(params []string, optionals map[string]string, request *http.Request) *resultStruct {
			zzContent := params[0]
			yyContent := params[1]
			loader := session.Get(r.Context(), loaderKey).(RandgenLoader)
			sqls, err := (&loader).LoadData(zzContent, yyContent,
				optionals[db], optionals[queries])
			if err != nil {
				return resultError(err.Error())
			}

			session.Put(r.Context(), loaderKey, loader)

			return resultSuccess(sqls)
		})
}

func Query(w http.ResponseWriter, r *http.Request) {
	base(w, r, nil, nil, "/init",
		func(params []string, optionals map[string]string, request *http.Request) *resultStruct {
			loader := session.Get(request.Context(), loaderKey).(RandgenLoader)
			if loader.CachedQueries == nil {
				return resultError("please access '/loaddata' first")
			}
			return resultSuccess(loader.Query())
		})
}

func Compare(w http.ResponseWriter, r *http.Request) {
	base(w, r, []string{"mysql", "tidb"}, nil, "/init",
		func(params []string, optionals map[string]string, request *http.Request) *resultStruct {
			r1 := params[0]
			r2 := params[1]
			loader := session.Get(r.Context(), loaderKey).(RandgenLoader)
			comment, consistent := loader.Compare(r1, r2)
			return resultSuccess(struct {
				Comment string
				Consistent bool
			}{comment, consistent})
		})
}
