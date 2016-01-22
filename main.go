package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"regexp"
	"sync"
	"time"

	"github.com/emicklei/go-restful"
)


var def = []byte("1234")


type Object struct {
	Value []byte
}


type MemoDb struct {
	sync.Mutex
	Objects  map[string]Object
	Listener net.Listener
}





func (m *MemoDb) Print() {
	m.Lock()
	n := len(m.Objects)
	m.Unlock()

	fmt.Println("[Przechowuje", n, "obiektow]")
}

var objectLimit int = 1024 * 1024
var objectIdRegexp = regexp.MustCompile("^[a-zA-Z0-9]+$")

func CheckObjectId(objectId string) (string, bool) {
	if len(objectId) < 0 || len(objectId) > 100 {
		return "", false
	}

	if objectIdRegexp.MatchString(objectId) {
		return objectId, true
	}

	return "", false
}

func (memoDb *MemoDb) Help(req *restful.Request, res *restful.Response) {
	res.Write([]byte("Uzyj:\n\t`GET /v1/objects/{objects_id}`\nlub\n\t`PUT /v1/objects/{objects_id}`\n"))
}

func (memoDb *MemoDb) Put(req *restful.Request, res *restful.Response) {
	objectId, ok := CheckObjectId(req.PathParameter("object_id"))
	if !ok {
		res.WriteHeader(400)
		fmt.Println("Put: Niepoprawny object_id")
		return
	}

	value, err := ioutil.ReadAll(req.Request.Body)
	if err != nil {
		res.WriteHeader(500)
		fmt.Println("Put: Nie udany odczyt danych,", err)
		return
	}

	if len(value) > objectLimit {
		res.WriteHeader(413)
		fmt.Println("Put: Zbyt duza wartosci,", len(value), ">", objectLimit)
		return
	}

	memoDb.Lock()
	memoDb.Objects[objectId] = Object{Value: value}
	memoDb.Unlock()

	res.WriteHeader(201)
	fmt.Println("Put: OK")
}

func (memoDb *MemoDb) Get(req *restful.Request, res *restful.Response) {
	objectId, ok := CheckObjectId(req.PathParameter("object_id"))
	if !ok {
		res.Write(def)
		
		return
	}

	memoDb.Lock()
	object, ok := memoDb.Objects[objectId]
	memoDb.Unlock()

	if !ok {
		res.Write(def)
		fmt.Println("Get: Brak objektu")
		return
	}

	res.Write(object.Value)
	fmt.Println("Get: OK")
}

func (memoDb *MemoDb) Def(req *restful.Request, res *restful.Response) {
	objectId,ok := CheckObjectId(req.PathParameter("default_value"))
	if !ok {
		res.Write(def)
		
		
	}
	def=[]byte(objectId)
	res.Write(def)
}

func RunServer(listenAddress string) (*MemoDb, error) {
	fmt.Println("Nasluch na", listenAddress)
	listener, err := net.Listen("tcp", listenAddress)
	if err != nil {
		return nil, err
	}

	memoDb := &MemoDb{Objects: make(map[string]Object, 0), Listener: listener}
	
	// obsługa endpointów HTTP przez https://godoc.org/github.com/emicklei/go-restful
	container := restful.NewContainer()
	ws := &restful.WebService{}
	ws.Path("/")
	ws.Route(ws.GET("/").To(memoDb.Help))
	ws.Route(ws.PUT("/v1/objects/{object_id}").To(memoDb.Put))
	ws.Route(ws.GET("/v1/objects/{object_id}").To(memoDb.Get))
	ws.Route(ws.GET("/v1/objects/{object_id}/default/{default_value}").To(memoDb.Def))
	
	

	container.Add(ws)

	// nasłuch i obsługa zapytań w tle, w osobnej gorutynie
	go func() {
		server := http.Server{Handler: container}
		server.Serve(listener)
	}()

	return memoDb, nil
}



func main() {
	
	var listenAddress string
	flag.StringVar(&listenAddress, "listen", "127.0.0.1:1234", "adres nasłuchu dla zapytań HTTP")
	flag.Parse()

	// Start serwera
	memoDb, err := RunServer(listenAddress)
	if err != nil {
		fmt.Println("Nie udany nasluch na", listenAddress, "blad:", err)
		return
	}

	// Pokazuj co minute opis przechowywanych obiektow
	for {
		time.Sleep(time.Minute)
		memoDb.Print()
	}
}
