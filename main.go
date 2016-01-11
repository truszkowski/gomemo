package main

import (
	"flag"
	"fmt"
	"github.com/emicklei/go-restful"
	"io/ioutil"
	"net/http"
	"regexp"
	"sync"
)

type Object struct {
	Value []byte
}

type MemoDb struct {
	sync.Mutex
	Objects map[string]Object
}

var objectIdRegexp = regexp.MustCompile("^[a-zA-Z0-9]+$")
var memoDb = MemoDb{Objects: make(map[string]Object, 0)}

func CheckObjectId(objectId string) (string, bool) {
	if len(objectId) < 0 || len(objectId) > 100 {
		return "", false
	}

	if objectIdRegexp.MatchString(objectId) {
		return objectId, true
	}

	return "", false
}

func Help(req *restful.Request, res *restful.Response) {
	res.Write([]byte("Uzyj: `GET /v1/objects/{objects_id}` lub `PUT /v1/objects/{objects_id}`"))
}

func Put(req *restful.Request, res *restful.Response) {
	objectId, ok := CheckObjectId(req.PathParameter("object_id"))
	if !ok {
		res.WriteHeader(400)
		fmt.Println("Niepoprawny object_id")
		return
	}

	value, err := ioutil.ReadAll(req.Request.Body)
	if err != nil {
		res.WriteHeader(500)
		fmt.Println("Nie udany odczyt danych,", err)
		return
	}

	if len(value) > 1024*1024 {
		res.WriteHeader(413)
		fmt.Println("Zbyt duza wartosci,", len(value), "> 1MB")
		return
	}

	memoDb.Lock()
	memoDb.Objects[objectId] = Object{Value: value}
	memoDb.Unlock()
}

func Get(req *restful.Request, res *restful.Response) {
	objectId, ok := CheckObjectId(req.PathParameter("object_id"))
	if !ok {
		res.WriteHeader(400)
		fmt.Println("Niepoprawny object_id")
		return
	}

	memoDb.Lock()
	object, ok := memoDb.Objects[objectId]
	memoDb.Unlock()

	if !ok {
		res.WriteHeader(404)
		fmt.Println("Brak objektu")
		return
	}

	res.Write(object.Value)
}

func main() {
	var listenAddress string
	flag.StringVar(&listenAddress, "listen", "127.0.0.1:1234", "adres nasłuchu dla zapytań HTTP")
	flag.Parse()
	fmt.Println("Nasluch na", listenAddress)

	// obsługa endpointów HTTP przez https://godoc.org/github.com/emicklei/go-restful
	container := restful.NewContainer()
	ws := &restful.WebService{}
	ws.Path("/")
	ws.Route(ws.GET("/").To(Help))
	ws.Route(ws.PUT("/v1/objects/{object_id}").To(Put))
	ws.Route(ws.GET("/v1/objects/{object_id}").To(Get))
	container.Add(ws)

	// nasłuch i obsługa zapytań
	http.ListenAndServe(listenAddress, container)
}
