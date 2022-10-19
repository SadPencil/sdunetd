package main

import (
	"context"
	"github.com/SadPencil/sdunetd/sdunet"
	"github.com/SadPencil/sdunetd/setting"
	"time"
)

var _manager *sdunet.Manager
var _manager_cancel context.CancelFunc

func terminateManager() {
	if _manager_cancel != nil {
		_manager_cancel()
	}
	_manager = nil
	_manager_cancel = nil
}

func getManager(settings *setting.Settings) (*sdunet.Manager, error) {
	if _manager == nil {
		networkInterface := ""
		if settings.Network.StrictMode {
			networkInterface = settings.Network.Interface
		}

		manager, err := sdunet.GetManager(
			settings.Account.Scheme,
			settings.Account.AuthServer,
			settings.Account.Username,
			networkInterface,
		)
		if err != nil {
			return nil, err
		}
		manager.Timeout = time.Duration(settings.Network.Timeout) * time.Second
		manager.MaxRetryCount = int(settings.Network.MaxRetryCount)
		manager.RetryWait = time.Duration(settings.Network.RetryIntervalSec) * time.Second
		manager.Logger = verboseLogger

		ctx, cancel := context.WithCancel(context.Background())
		manager.Context = ctx

		_manager = &manager
		_manager_cancel = cancel
	}
	return _manager, nil
}
