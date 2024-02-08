package dependencyinjection

import (
	"github.com/samber/do"
)

var injector = do.DefaultInjector

func Register[T any](value T) {
	do.ProvideValue(injector, value)
}

func Get[T any]() T {
	return do.MustInvoke[T](injector)
}
