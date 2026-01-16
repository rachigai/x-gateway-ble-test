package main

import (
	"log"

	"github.com/paypal/gatt"
)

func main() {
	opts := []gatt.Option{
		gatt.LnxMaxConnections(1),
		gatt.LnxDeviceID(-1, false),
	}
	d, err := gatt.NewDevice(opts...)
	if err != nil {
		log.Fatal(err)
	}

	suuid := gatt.MustParseUUID("82473ACA-A36D-4D49-AE1C-EB17F4C56218")
	cuuid := gatt.MustParseUUID("56A4FA7A-E1D2-4852-A4A0-2866DEC91905")

	pch := make(chan gatt.Peripheral)
	done := make(chan struct{})

	d.Handle(
		gatt.PeripheralDiscovered(
			func(p gatt.Peripheral, a *gatt.Advertisement, rssi int) {
				log.Printf("scanned: id=%s connectable=%t rssi=%d\n", p.ID(), a.Connectable, rssi)
				if p.ID() == "44:1D:64:62:B3:1E" {
					pch <- p
				}
			}))
	d.Handle(
		gatt.PeripheralConnected(
			func(p gatt.Peripheral, err error) {
				log.Printf("connected: %s", p.ID())

				if ss, err := p.DiscoverServices([]gatt.UUID{suuid}); err != nil {
					log.Fatal(err)
				} else {
					for _, s := range ss {
						if s.UUID().Equal(suuid) {
							log.Printf("discovered service: UUID=%s", s.UUID().String())
							if cs, err := p.DiscoverCharacteristics([]gatt.UUID{cuuid}, s); err != nil {
								log.Fatal(err)
							} else {
								for _, c := range cs {
									log.Printf("discovered characteristic: UUID=%s", c.UUID().String())
									if bs, err := p.ReadCharacteristic(c); err != nil {
										log.Print(err)
									} else {
										log.Printf("len=%d, bs=%v", len(bs), bs)
									}
								}
							}
						}
					}
				}
				close(done)
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
