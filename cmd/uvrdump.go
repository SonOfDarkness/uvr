package main

import (
	"github.com/brutella/can"
	"github.com/brutella/uvr"
	"log"
)

func readOutlet(outlet uvr.Outlet, client *uvr.Client) (descr string, mode string, val string) {
	if value, err := client.Read(outlet.Description); err == nil {
		descr = value.(string)
	}

	if value, err := client.Read(outlet.Mode); err == nil {
		mode = value.(string)
	}

	if value, err := client.Read(outlet.State); err == nil {
		val = value.(string)
	}

	return
}

func readInlet(inlet uvr.Inlet, client *uvr.Client) (descr string, state string, val float32) {
	if value, err := client.Read(inlet.Description); err == nil {
		descr = value.(string)
	}

	if value, err := client.Read(inlet.State); err == nil {
		state = value.(string)
	}

	if value, err := client.Read(inlet.Value); err == nil {
		if float, ok := value.(float32); ok == true {
			val = float
		}
	}

	return
}

func readOutlets(client *uvr.Client) {
	outlets := []uvr.Outlet{
		uvr.NewOutlet(0x1),
		uvr.NewOutlet(0x2),
		uvr.NewOutlet(0x3),
		uvr.NewOutlet(0x4),
		uvr.NewOutlet(0x5),
		uvr.NewOutlet(0x6),
		uvr.NewOutlet(0x7),
		uvr.NewOutlet(0x8),
		uvr.NewOutlet(0x9),
		uvr.NewOutlet(0xa),
		uvr.NewOutlet(0xb),
	}

	log.Printf("+---------+-----------------+--------+------+")
	log.Printf("| Ausgang | Bezeichnung     | Mode   | Wert |")
	log.Printf("+---------+-----------------+--------+------+")
	for index, outlet := range outlets {
		descr, mode, val := readOutlet(outlet, client)
		log.Printf("| %-7d | %-15s | %-6s | %-4s |", index+1, descr, mode, val)
	}
	log.Printf("+---------+-----------------+--------+------+")
}

func readInlets(client *uvr.Client) {
	inlets := []uvr.Inlet{
		uvr.NewInlet(0x1),
		uvr.NewInlet(0x2),
		uvr.NewInlet(0x3),
		uvr.NewInlet(0x4),
		uvr.NewInlet(0x5),
		uvr.NewInlet(0x6),
		uvr.NewInlet(0x7),
		uvr.NewInlet(0x8),
		uvr.NewInlet(0x9),
		uvr.NewInlet(0xa),
		uvr.NewInlet(0xb),
		uvr.NewInlet(0xc),
		uvr.NewInlet(0xd),
		uvr.NewInlet(0xe),
	}

	log.Printf("+---------+-----------------+--------+------")
	log.Printf("| Eingang | Bezeichnung     | Status | Wert")
	log.Printf("+---------+-----------------+--------+------")
	for index, inlet := range inlets {
		descr, state, val := readInlet(inlet, client)
		log.Printf("| %-7d | %-15s | %-6s | %.1f", index+1, descr, state, val)
	}
	log.Printf("+---------+-----------------+--------+------")
}

func HandleCANopen(frame can.Frame) {
	log.Printf("%X % X\n", frame.ID, frame.Data)
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	bus, err := can.NewBusForInterfaceWithName("can0")

	if err != nil {
		log.Fatal(err)
	}
	// bus.SubscribeFunc(HandleCANopen)
	go bus.ConnectAndPublish()

	nodeID := uint8(0x10)
	uvrID := uint8(0x1)

	c := uvr.NewClient(nodeID, bus)
	c.Connect(uvrID)

	readInlets(c)
	readOutlets(c)

	c.Disconnect(uvrID)
	bus.Disconnect()
}