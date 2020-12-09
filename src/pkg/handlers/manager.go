package handlers

import (
	"gocheck/pkg/config"
	"log"
	"time"
)

type managerChannels struct {
	control chan string
	notify chan string
}

type Manager interface {
	Start()
	UpdateConfig()
}

type HandlerManager struct {
	workerChannels map[string]*managerChannels
	resultChannels *managerChannels
	daemonChannels *managerChannels
	workingConfig config.Config
}

func (m *HandlerManager) Start() {
	//
}

func (m *HandlerManager) UpdateConfig() {

}

func NewHandlerManager(c config.Config) *HandlerManager {
	var manager HandlerManager
	manager.workingConfig = c
	manager.daemonChannels = &managerChannels{control: make(chan string, 1), notify: make(chan string, 1)}
	manager.resultChannels = &managerChannels{control: make(chan string, 1), notify: make(chan string, 1)}

	// TODO - UpdateConfig()
	return &manager
}

func killChecks(chanMap map[string]*managerChannels) {
	// Add wait group
	for name, channels := range chanMap {
		 channels.control <- KILL
		 timer := time.NewTimer(10 * time.Second)
		 select {
		 case status := <- channels.notify:
		 	if status != STOPNORM {
		 		log.Printf("error when stopping check %s: %s", name, status)
			}
			timer.Stop()
		 case timeout := <- timer.C:
		 	log.Printf("attempt to stop check %s timed out at %s", name, timeout)
		 }
	}
	// End wait group
}

func createChecks(chanMap map[string]*workerChannels, r resultChannels)


// Prepare/ensure result handlers are available
// For each check
	// Setup the channels
	// Start the check handler