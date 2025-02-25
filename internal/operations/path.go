package operations

import (
	stdpath "path"
	"strings"

	"github.com/alist-org/alist/v3/internal/driver"
	"github.com/alist-org/alist/v3/internal/errs"
	"github.com/alist-org/alist/v3/pkg/utils"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// ActualPath Get the actual path
// !!! maybe and \ in the path when use windows local
func ActualPath(storage driver.Additional, rawPath string) string {
	if i, ok := storage.(driver.IRootFolderPath); ok {
		rawPath = stdpath.Join(i.GetRootFolderPath(), rawPath)
	}
	return utils.StandardizePath(rawPath)
}

// GetStorageAndActualPath Get the corresponding storage and actual path
// for path: remove the virtual path prefix and join the actual root folder if exists
func GetStorageAndActualPath(rawPath string) (driver.Driver, string, error) {
	rawPath = utils.StandardizePath(rawPath)
	if strings.Contains(rawPath, "..") {
		return nil, "", errors.WithStack(errs.RelativePath)
	}
	storage := GetBalancedStorage(rawPath)
	if storage == nil {
		return nil, "", errors.Errorf("can't find storage with rawPath: %s", rawPath)
	}
	log.Debugln("use storage: ", storage.GetStorage().MountPath)
	virtualPath := utils.GetActualVirtualPath(storage.GetStorage().MountPath)
	actualPath := strings.TrimPrefix(rawPath, virtualPath)
	actualPath = ActualPath(storage.GetAddition(), actualPath)
	return storage, actualPath, nil
}
