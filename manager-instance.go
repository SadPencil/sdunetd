package main

import (
	"context"
	"github.com/SadPencil/sdunetd/sdunet"
	"github.com/SadPencil/sdunetd/setting"
	"time"
)

var _manager *sdunet.Manager

func getManager(ctx context.Context, settings *setting.Settings) (*sdunet.Manager, error) {
	if _manager == nil {
		networkInterface := ""
		if settings.Network.StrictMode {
			networkInterface = settings.Network.Interface
		}

		manager, err := sdunet.GetManager(ctx,
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

		_manager = &manager
	}
	return _manager, nil
}
