package main

import (
	"testing"

	. "gopkg.in/check.v1"
)

type GoMemoSuite struct {
	memoDb *MemoDb
}

var _ = Suite(&GoMemoSuite{})

func Test(t *testing.T) { TestingT(t) }

func (s *GoMemoSuite) SetUpSuite(c *C) {
	memoDb, err := RunServer("127.0.0.1:1234")
	c.Assert(err, IsNil)
	s.memoDb = memoDb
}

func (s *GoMemoSuite) TearDownSuite(c *C) {
}

func (s *GoMemoSuite) SetUpTest(c *C) {
}

func (s *GoMemoSuite) TearDownTest(c *C) {
}
