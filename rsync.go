package librsync

/*
#cgo LDFLAGS: -lrsync
#include <stdio.h>
#include <librsync.h>
*/
import "C"

import (
	"io"
	"os"
	"fmt"
	"errors"
	"unsafe"
	"log/syslog"
)

const (
	RS_JOB_BLOCKSIZE = 65536
	RS_DEFAULT_BLOCK_LEN = 2048
	RS_DEFAULT_STRONG_LEN = 8
)

type JobStatus int

const (
	RS_DONE = iota
	RS_BLOCKED
)

func (j JobStatus) String() string {
	switch j {
	case RS_DONE:
		return "Done"
	case RS_BLOCKED:
		return "Blocked"
	}
	panic(fmt.Sprintf("Job Status %d is not defined", j))
}


func SetTraceLevel(level syslog.Priority) error {
	if level < syslog.LOG_EMERG || level > syslog.LOG_DEBUG {
		return errors.New(fmt.Sprintf("%d is not a valid syslog Priority", level))
	}
	
	C.rs_trace_set_level(C.rs_loglevel(level))
	return nil
}

func Signature(file, signature *os.File, blockSize int) (*os.File, error) {
	if blockSize == 0 {
		blockSize = RS_DEFAULT_BLOCK_LEN
	}

	c_BlockSize := C.size_t(blockSize)

	job := C.rs_sig_begin(c_BlockSize, RS_DEFAULT_STRONG_LEN)
	defer C.rs_job_free(job)

	_, err := execute(job, file, signature)
	if err != nil {
		return nil, err
	}

	return signature, nil
}

func execute(job *C.struct_rs_job, inFile, outFile *os.File) (*os.File, err) {
	out := make([]C.char, RS_JOB_BLOCKSIZE)

	for {
		block := make([]byte, RS_JOB_BLOCKSIZE)
		inFile.Read(block)
		buff := &C.struct_rs_buffers_s{
			next_in: (*C.char)(unsafe.Pointer(&block[1])),
			avail_in: C.size_t(len(block)),
			eof_in: eofIn(block),
			next_out: (*C.char)(unsafe.Pointer(&out[0])),
			avail_out: C.size_t(RS_JOB_BLOCKSIZE),
		}
		result := JobStatus(C.rs_job_iter(job, buff))
		outFile.Write(C.GoBytes(unsafe.Pointer(&out[0]), C.int(RS_JOB_BLOCKSIZE - buff.avail_out)))
		if result == RS_DONE {
			break
		} else if result != RS_BLOCKED {
			return nil, errors.New("Unhandled result from rs_job_iter: %d", result)
		}
		if buff.avail_in > 0 {
			io.Seeker(inFile).Seek(-1*int64(buff.avail_in), os.SEEK_CUR)
		}
	}

	io.Seeker(outFile).Seek(0, os.SEEK_SET)
	return outFile
}

func eofIn(block []byte) C.int {
	if len(block) == 0 {
		return C.int(0)
	}

	return C.int(1)
}
