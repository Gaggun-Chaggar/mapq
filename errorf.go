package mapq

import "fmt"

type methodErrorfFunc func(methodName string) errorfFunc
type errorfFunc func(format string, args ...any) error

func createErrorf(packageName string) methodErrorfFunc {
	return func(methodName string) errorfFunc {
		return func(format string, args ...any) error {
			return fmt.Errorf("%s.%s: "+format, packageName, methodName)
		}
	}
}

var packageErrorf = createErrorf("jsonq")
