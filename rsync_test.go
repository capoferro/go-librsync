package librsync_test

import (
	"io"
	"io/ioutil"
	"os"
	"crypto/rand"
	"testing"
	"log/syslog"

	"github.com/capoferro/go-librsync"

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
	err := librsync.SetTraceLevel(syslog.Priority(999))
	c.Assert(err, NotNil)
	c.Assert(err.Error(), Equals, "999 is not a valid syslog Priority")
}

func (s *RsyncSuite) Test_Signature(c *C) {
	randomBytes := make([]byte, librsync.RS_JOB_BLOCKSIZE*3) //int(math.Pow(1024, 2)))
	rand.Read(randomBytes)
	inFile, _ := ioutil.TempFile("/tmp", "go-librsync_in")
	ioutil.WriteFile(inFile.Name(), randomBytes, 0644)
	io.Seeker(inFile).Seek(0, os.SEEK_SET)
	outFile, _ := ioutil.TempFile("/tmp", "go-librsync_out")
	librsync.Signature(inFile, outFile, 0)
	signature, _ := ioutil.ReadAll(outFile)
	c.Assert(len(signature) > 0, Equals, true)
}

