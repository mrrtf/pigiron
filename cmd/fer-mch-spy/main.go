package main

import (
	"log"

	"github.com/sbinet-alice/fer"
	"github.com/sbinet-alice/fer/config"
)

type spy struct {
	cfg    config.Device
	idatac chan fer.Msg
}

func (dev *spy) Configure(cfg config.Device) error {
	dev.cfg = cfg
	return nil
}

func (dev *spy) Init(ctl fer.Controler) error {
	idatac, err := ctl.Chan("data1", 0)
	if err != nil {
		return err
	}
	dev.idatac = idatac
	return nil
}

func (dev *spy) Run(ctl fer.Controler) error {
	for {
		select {
		case data := <-dev.idatac:
			ctl.Printf("received %d bytes", len(data.Data))
		case <-ctl.Done():
			return nil
		}
	}
}

func main() {
	err := fer.Main(&spy{})
	if err != nil {
		log.Fatal(err)
	}
}
