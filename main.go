package main

import "github.com/markgunnels/gizmo-test/service"

func main() {

	for {
		service.Init()
		service.Log.Info("starting subscriber process")
		
		if err := service.Run(); err != nil {
			service.Log.Error("subscriber encountered a fatal error: ", err)

		}
	}

	service.Log.Info("subscriber process shutting down")
}
