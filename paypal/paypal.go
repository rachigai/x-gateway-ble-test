package main

import (
	"log"

	"github.com/paypal/gatt"
)

func MustBeUUID(s string) gatt.UUID {
	if uuid, err := gatt.ParseUUID(s); err != nil {
		panic(err)
	} else {
		return uuid
	}
}

func main() {
	opts := []gatt.Option{
		gatt.LnxMaxConnections(1),
		gatt.LnxDeviceID(-1, false),
	}
	d, err := gatt.NewDevice(opts...)
	if err != nil {
		log.Fatal(err)
	}

	service := gatt.NewService(MustBeUUID("82473ACA-A36D-4D49-AE1C-EB17F4C56218"))
	cuuid := MustBeUUID("56A4FA7A-E1D2-4852-A4A0-2866DEC91905")

	pch := make(chan gatt.Peripheral)
	done := make(chan struct{})

	d.Handle(
		gatt.PeripheralDiscovered(
			func(p gatt.Peripheral, a *gatt.Advertisement, rssi int) {
				//ch <- &scanned{p: p, a: a, rssi: rssi}
				log.Printf("scanned: id=%s connectable=%t rssi=%d\n", p.ID(), a.Connectable, rssi)

				if p.ID() == "44:1D:64:62:B3:1E" {
					pch <- p
				}
			}))
	d.Handle(
		gatt.PeripheralConnected(
			func(p gatt.Peripheral, err error) {
				log.Printf("connected: %s", p.ID())
				c := gatt.NewCharacteristic(cuuid, service, gatt.Property(0), 0, 0)
				if bs, err := p.ReadCharacteristic(c); err != nil {
					log.Print(err)
				} else {
					log.Print(bs)
				}
				done <- struct{}{}
			}))
	d.Init(func(d gatt.Device, s gatt.State) {
		switch s {
		case gatt.StatePoweredOn:
			log.Printf("gatt-state: powered on")
		default:
			log.Printf("gatt-state: stop scanning")
			d.StopScanning()
		}
	})

	go func() {
		p := <-pch
		d.StopScanning()
		d.Connect(p)
	}()

	d.Scan([]gatt.UUID{}, true)

	<-done
}
