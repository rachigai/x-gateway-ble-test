package main

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/go-ble/ble"
	"github.com/go-ble/ble/linux"
)

func main() {
	d, err := linux.NewDeviceWithName("hci0")
	if err != nil {
		log.Fatal(err)
	}
	ble.SetDefaultDevice(d)

	// ctx := ble.WithSigHandler(context.WithTimeout(context.Background(), 10*time.Second))
	// d.Scan(ctx, true, func(a ble.Advertisement) {
	// 	log.Printf("advertising: addr=%s connectable=%t rssi=%d", a.Addr().String(), a.Connectable(), a.RSSI())
	// })

	filter := func(a ble.Advertisement) bool {
		return strings.ToUpper(a.Addr().String()) == "44:1D:64:62:B3:1E"
	}

	if cli, err := ble.Connect(
		ble.WithSigHandler(context.WithTimeout(context.Background(), 20*time.Second)),
		filter); err != nil {
		log.Fatal(err)
	} else {
		log.Printf("connected: addr=%s", cli.Addr().String())

		if ss, err := cli.DiscoverServices([]ble.UUID{ble.MustParse("82473ACA-A36D-4D49-AE1C-EB17F4C56218")}); err != nil {
			log.Print(err)
		} else {
			for _, s := range ss {
				log.Printf("service=%s", s.UUID.String())
				if cs, err := cli.DiscoverCharacteristics([]ble.UUID{ble.MustParse("56A4FA7A-E1D2-4852-A4A0-2866DEC91905")}, s); err != nil {
					log.Print(err)
				} else {
					for _, c := range cs {
						if bs, err := cli.ReadLongCharacteristic(c); err != nil {
							log.Print(err)
						} else {
							log.Printf("    characteristic=%s len=%d value=%v", c.UUID.String(), len(bs), bs)
						}
					}
				}

			}
		}

	}

}
