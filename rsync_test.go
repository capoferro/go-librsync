package librsync_test

import (
	"testing"
	"log/syslog"

	"gh.riotgames.com/jkiehl/go-librsync"

	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }
type RsyncSuite struct{}
var _ = Suite(&RsyncSuite{})

func (s *RsyncSuite) Test_TraceLevel_Set(c *C) {
	librsync.SetTraceLevel(syslog.LOG_DEBUG)
}

func (s *RsyncSuite) Test_TraceLevel_SetInvalid(c *C) {
	librsync.SetTraceLevel(syslog.Priority(999))
}