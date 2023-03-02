package funcs

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func RegApp() {
	//_, _, _, port := disco.GetEnv()
	//instanceID, err := sd.ClassicRegister(disco.ConsulRegister{
	//	ServiceName:  "test.testapp",
	//	InstanceHost: "192.168.88.18",
	//	InstancePort: port,
	//	Tags:         []string{disco.TagDataApp},
	//	Meta: map[string]string{
	//		"data_app_id": "testapp",
	//	},
	//})
	//if err != nil {
	//	panic(err)
	//}
	//ProcessExitHook(Dereg, *instanceID)
	//http.HandleFunc(disco.HealthCheckAPI, disco.HealthCheckFunc)
	//err = http.ListenAndServe(":"+port, nil)
	//if err != nil {
	//	panic(err)
	//}

}

func Dereg(arg interface{}) {
	//instanceID := arg.(string)
	//err := sd.ClassicDeregister(instanceID)
	//if err != nil {
	//	log.Println("Deregister ", err)
	//}
}

func ProcessExitHook(fn func(arg interface{}), arg interface{}) {
	exitChan := make(chan os.Signal, 1)
	signal.Notify(exitChan, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
	go exitHandle(exitChan, fn, arg)
}

func exitHandle(exitChan chan os.Signal, fn func(args interface{}), arg interface{}) {
	for sig := range exitChan {
		log.Println("接受到来自系统的信号：", sig)
		fn(arg)
		//如果ctrl+c 关不掉程序，使用os.Exit强行关掉
		os.Exit(1)
	}
}
