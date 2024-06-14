package lib

import (
	"errors"
	"fmt"
	"runtime"
)

func ErrBuilder(errs ...error) error {
	pc, _, line, _ := runtime.Caller(1)
	function := runtime.FuncForPC(pc)
	thisErrMsg := fmt.Errorf("%s at line %d", function.Name(), line)

	allErrors := append([]error{thisErrMsg}, errs...)

	return errors.Join(allErrors...)
}
