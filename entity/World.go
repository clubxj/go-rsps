package entity

import (
	"log"
	"sync"
	"time"
)

var world *World
func WorldProvider() *World {
	if world == nil {
		world = CreateWorld()
	}
	return world
}

type World struct {
	Regions map[uint16]*Region
	lock sync.Mutex
}

func CreateWorld() *World {
	return &World{
		Regions: make(map[uint16]*Region),
	}
}

func (w *World) Tick() {
	for {
		<- time.After(600 * time.Millisecond)
		w.lock.Lock()
		for k, v := range w.Regions {
			v.Tick()
			if v.MarkedForDeletion {
				delete(w.Regions, k)
			}
		}
		w.lock.Unlock()
	}
}
func (w *World) GetRegion(id uint16) *Region {
	w.lock.Lock()
	if w.Regions[id] == nil {
		w.Regions[id] = CreateRegion(id)
	}
	w.lock.Unlock()
	return w.Regions[id]
}

func (w *World) AddPlayerToRegion(player *Player) {
	previousRegion := player.Region
	regionId := GetRegionIdByPosition(player.Position)
	region := w.GetRegion(regionId)
	player.Region = region
	if regionId == 12850 {
		log.Printf("adding player to 12850")
	}

	newAdj := region.GetAdjacentIds()
	var prevAdj []uint16
	if previousRegion != nil {
		prevAdj = previousRegion.GetAdjacentIds()
	}

	for _, x := range newAdj {
		var exists bool
		for _, y := range prevAdj {
			if x == y {
				exists = true
			}
		}
		if !exists {
			r := world.GetRegion(x)
			r.OnEnter(player)
		}
	}

	for _, x := range prevAdj {
		var exists bool
		for _, y := range newAdj {
			if x == y {
				exists = true
			}
		}
		if !exists {
			r := world.GetRegion(x)
			r.OnLeave(player)
		}
	}
}

func (w *World) GetRegionForPlayer(player *Player) *Region {
	regionId := GetRegionIdByPosition(player.Position)
	return w.Regions[regionId]
}

