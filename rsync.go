package librsync

/*
#cgo LDFLAGS: -lrsync
#include <stdio.h>
#include <librsync.h>
*/
import "C"

import (
//	"fmt"
	"log/syslog"
)

func SetTraceLevel(level syslog.Priority) {
	C.rs_trace_set_level(C.rs_loglevel(level))
}
