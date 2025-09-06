package private

import "errors"

var ErrDecryptionDEK = errors.New("decrypt DEK from repository failed")
var ErrCreateDEK = errors.New("create DEK encoder from repository data failed")
var ErrGenerateDEK = errors.New("generating new DEK failed ")
var ErrSaveData = errors.New("saving entity failed")
var ErrReplaceKey = errors.New("replace key failed")