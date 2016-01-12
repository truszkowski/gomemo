package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	. "gopkg.in/check.v1"
)

type GoMemoSuite struct {
	memoDb  *MemoDb
	address string
}

var _ = Suite(&GoMemoSuite{})

func Test(t *testing.T) { TestingT(t) }

func (s *GoMemoSuite) SetUpSuite(c *C) {
	s.address = "127.0.0.1:1234"
	memoDb, err := RunServer(s.address)
	c.Assert(err, IsNil)
	s.memoDb = memoDb
}

func (s *GoMemoSuite) TearDownSuite(c *C) {
}

func (s *GoMemoSuite) SetUpTest(c *C) {
}

func (s *GoMemoSuite) TearDownTest(c *C) {
}

func (s *GoMemoSuite) url(objectId string) string {
	return "http://" + s.address + "/v1/objects/" + objectId
}

func (s *GoMemoSuite) put(c *C, objectId string, reader io.Reader) *http.Response {
	req, err := http.NewRequest("PUT", s.url(objectId), reader)
	c.Assert(err, IsNil)

	res, err := http.DefaultClient.Do(req)
	c.Assert(err, IsNil)

	return res
}

func (s *GoMemoSuite) TestPutGet(c *C) {
	// na poczatku nie ma
	r, err := http.Get(s.url("a1"))
	c.Assert(err, IsNil)
	c.Assert(r.StatusCode, Equals, 404)

	// dodajemy
	value := "a1 test"
	r = s.put(c, "a1", strings.NewReader(value))
	c.Assert(r.StatusCode, Equals, 201)

	// juz powinno byc
	r, err = http.Get(s.url("a1"))
	c.Assert(err, IsNil)
	c.Assert(r.StatusCode, Equals, 200)

	b, err := ioutil.ReadAll(r.Body)
	c.Assert(err, IsNil)
	c.Assert(string(b), Equals, value)

	// nadpisanie danych
	value = "a1 test, dwa"
	r = s.put(c, "a1", strings.NewReader(value))
	c.Assert(r.StatusCode, Equals, 201)

	// powinna byc nowa wartosc
	r, err = http.Get(s.url("a1"))
	c.Assert(err, IsNil)
	c.Assert(r.StatusCode, Equals, 200)

	b, err = ioutil.ReadAll(r.Body)
	c.Assert(err, IsNil)
	c.Assert(string(b), Equals, value)
}
