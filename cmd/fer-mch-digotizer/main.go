package main

import (
	"log"

	"github.com/sbinet-alice/fer"
	"github.com/sbinet-alice/fer/config"
)

type digitizer struct {
	cfg    config.Device
	idatac chan fer.Msg
	odatac chan fer.Msg
}

func (dev *digitizer) Configure(cfg config.Device) error {
	dev.cfg = cfg
	return nil
}

func (dev *digitizer) Init(ctl fer.Controler) error {
	idatac, err := ctl.Chan("mch-hits", 0)
	if err != nil {
		return err
	}

	dev.idatac = idatac
	odatac, err := ctl.Chan("mch-digotizer", 0)
	if err != nil {
		return err
	}

	dev.odatac = odatac
	return nil
}

func (dev *digitizer) Run(ctl fer.Controler) error {
	for {
		select {
		case data := <-dev.idatac:
			out := append([]byte(nil), data.Data...)
			dev.odatac <- fer.Msg{Data: out}
			ctl.Printf("received %q\n", len(data.Data))
		case <-ctl.Done():
			return nil
		}
	}
}

func main() {
	err := fer.Main(&digitizer{})
	if err != nil {
		log.Fatal(err)
	}
}
