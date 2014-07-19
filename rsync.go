package librsync

/*
#cgo LDFLAGS: -lrsync
#include <stdio.h>
#include <librsync.h>
*/
import "C"

import (
	"fmt"
	"errors"
	"log/syslog"
)

func SetTraceLevel(level syslog.Priority) error {
	if level < syslog.LOG_EMERG || level > syslog.LOG_DEBUG {
		return errors.New(fmt.Sprintf("%d is not a valid syslog Priority", level))
	}
	
	C.rs_trace_set_level(C.rs_loglevel(level))
	return nil
}
