package red_packet

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

type Packet struct {
	Id     int     `json:"id"` //指向一个红包
	Total  float64 `json:"total"`
	Number int     `json:"number"`
}

var RedPackets = make(map[int]*Packets)

type Packets struct {
	OnePacket chan float64
}

func NewPackets(size int) *Packets {
	return &Packets{OnePacket: make(chan float64, size)}

}

func (p *Packet) CreatRed(min float64) []float64 {

	var reds []float64
	total := p.Total
	for i := 1; i < p.Number; i++ {
		//保证即使一个红包是最大的了,后面剩下的红包,每个红包也不会小于最小值
		max := total - min*float64(p.Number-i) //算最多能得多少

		k := (p.Number - i) / 2
		//保证最后两个人拿的红包不超出剩余红包
		if p.Number-i <= 2 {
			k = p.Number - i
		}
		//最大的红包限定的平均线上下 这里的平均线=（min+max）/2
		max = max / float64(k)

		//设置随机种子
		rand.Seed(time.Now().UnixNano())
		//保证每个红包大于最小值,又不会大于最大值
		money := (int)(min*100 + rand.Float64()*(max*100-min*100+1)) //min*100 - max*100
		value := float64(money) / 100
		//保留两位小数

		value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)

		t := int64(total*100 - value*100)
		total = float64(t) / 100

		reds = append(reds, value)
	}
	reds = append(reds, total)
	return reds
}
