package red_packet

import (
	"fmt"
	"testing"
)

func TestPacket_CreatRed(t *testing.T) {
	p := Packet{
		Id:     0,
		Total:  100,
		Number: 10,
	}
	fmt.Println(p.CreatRed(0.1))
}
