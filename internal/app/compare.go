package app

import "errors"

// ErrCompareNotImplemented marks the intentionally deferred compare workflow.
var ErrCompareNotImplemented = errors.New("compare is not implemented in milestone one")
