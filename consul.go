package main

import (
	"github.com/reddec/monexec/monexec"
	"time"
	"fmt"
	"os"
	"log"
	"github.com/hashicorp/consul/api"
)

func updateServiceStatus(event monexec.Event, ttl time.Duration, client *api.Client) error {
	autoDeregistrationTimeout := 2 * ttl
	if autoDeregistrationTimeout < 1*time.Minute {
		autoDeregistrationTimeout = 1 * time.Minute
	}
	switch event.Type {
	case monexec.STARTED:
		err := client.Agent().ServiceRegister(&api.AgentServiceRegistration{
			ID:   event.Executable.GetGUID(),
			Name: event.Executable.ID(),
			Tags: []string{fmt.Sprintf("%v", os.Getpid())},
			Check: &api.AgentServiceCheck{
				DeregisterCriticalServiceAfter: autoDeregistrationTimeout.String(),
			},
		})
		if err != nil {
			log.Println("Can't register service", event.Executable.ID(), "in Consul:", err)
		} else {
			reg := api.AgentCheckRegistration{}
			reg.Name = event.Executable.ID() + ":ttl"
			reg.ID = event.Executable.GetGUID()
			reg.TTL = ttl.String()
			reg.DeregisterCriticalServiceAfter = autoDeregistrationTimeout.String()
			reg.ServiceID = event.Executable.GetGUID()
			err = client.Agent().CheckRegister(&reg)
			if err != nil {
				log.Println("Can't register service TTL check", event.Executable.ID(), "in Consul:", err)
			} else {
				log.Println("Service", event.Executable.ID(), "registered in Consul")
			}
		}
		return err
	case monexec.STOPPED:
		err := client.Agent().ServiceDeregister(event.Executable.GetGUID())
		if err != nil {
			log.Println("Can't deregister service", event.Executable.ID(), "in Consul:", err)
		} else {
			log.Println("Service", event.Executable.ID(), "deregistered in Consul")
		}
		return err
	}
	return nil
}

func runService(filteredEvents <-chan monexec.Event, interval time.Duration, client *api.Client) {
	ttlTimer := time.NewTicker(interval / 2)
	reregisterTimer := time.NewTicker(2 * interval)
	var lastEvent monexec.Event

	defer ttlTimer.Stop()
	defer reregisterTimer.Stop()

	needsToBeUpdate := true
LOOP:
	for {
		select {
		case event, ok := <-filteredEvents:
			if (!ok) {
				break LOOP
			}
			lastEvent = event
			needsToBeUpdate = updateServiceStatus(lastEvent, interval, client) != nil
		case <-reregisterTimer.C:
			if needsToBeUpdate {
				needsToBeUpdate = updateServiceStatus(lastEvent, interval, client) != nil
			}
		case <-ttlTimer.C:
			if lastEvent.Executable != nil {
				err := client.Agent().UpdateTTL(lastEvent.Executable.GetGUID(), "application running", "pass")
				if err != nil {
					log.Println("Can't update TTL for service", lastEvent.Executable.ID(), "in Consul:", err)
				}
				needsToBeUpdate = err != nil
			}
		}
	}
}
