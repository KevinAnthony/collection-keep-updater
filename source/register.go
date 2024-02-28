package source

import (
	"fmt"
	"sync"

	"github.com/kevinanthony/collection-keep-updater/types"

	"github.com/spf13/cobra"
)

type ConfigCallback struct {
	SetFlagsFunc                func(cmd *cobra.Command)
	SourceSettingFromConfigFunc func(data map[string]interface{}) types.ISourceSettings
	SourceSettingFromFlagsFunc  func(cmd *cobra.Command, original types.ISourceSettings) (types.ISourceSettings, error)
	GetIDFromURL                func(url string) (string, error)
}

var (
	_callbacks     map[types.SourceType]*ConfigCallback
	_callbackMutex sync.RWMutex
)

func RegisterConfigCallbacks(key types.SourceType, callback *ConfigCallback) {
	_callbackMutex.Lock()
	defer _callbackMutex.Unlock()

	if _callbacks == nil {
		_callbacks = map[types.SourceType]*ConfigCallback{}
	}
	fmt.Println(callback)
	_callbacks[key] = callback
}

func GetCallback(key types.SourceType) *ConfigCallback {
	_callbackMutex.RLock()
	defer _callbackMutex.RUnlock()

	callback, found := _callbacks[key]
	if !found {
		return nil
	}

	return callback
}

func SetSourceFlags(cmd *cobra.Command) {
	_callbackMutex.RLock()
	defer _callbackMutex.RUnlock()
	c := _callbacks
	fmt.Println(c)
	for _, cb := range _callbacks {
		cb.SetFlagsFunc(cmd)
	}
}
