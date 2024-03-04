package model

import (
	"math"
	"sync"

	mathutils "github.com/heavenston/creeps_server/creeps_lib/math_utils"
)

type ResourceKind string

const (
	Copper    ResourceKind = "copper"
	Food                   = "food"
	Oil                    = "oil"
	Rock                   = "rock"
	Wood                   = "wood"
	WoodPlank              = "woodPlank"
)

type Resources struct {
	Rock      int `json:"rock"`
	Wood      int `json:"wood"`
	Food      int `json:"food"`
	Oil       int `json:"oil"`
	Copper    int `json:"copper"`
	WoodPlank int `json:"woodPlank"`
}

func (res *Resources) OfKind(kind ResourceKind) *int {
	switch kind {
	case Rock:
		return &res.Rock
	case Wood:
		return &res.Wood
	case Food:
		return &res.Food
	case Oil:
		return &res.Oil
	case Copper:
		return &res.Copper
	case WoodPlank:
		return &res.WoodPlank
	}
	return nil
}

// return how many times this resources have the other one
// (4 copper 2 rock for 1 copper 1 rock returns 2)
func (res Resources) EnoughFor(other Resources) float64 {
	div := func(a float64, b float64) float64 {
		if b == 0 {
			return math.Inf(1)
		}
		return a / b
	}

	return mathutils.Min(
		div(float64(res.Rock), float64(other.Rock)),
		div(float64(res.Wood), float64(other.Wood)),
		div(float64(res.Food), float64(other.Food)),
		div(float64(res.Oil), float64(other.Oil)),
		div(float64(res.Copper), float64(other.Copper)),
		div(float64(res.WoodPlank), float64(other.WoodPlank)),
	)
}

func (res *Resources) Remove(other Resources) {
	res.Rock -= other.Rock
	res.Wood -= other.Wood
	res.Food -= other.Food
	res.Oil -= other.Oil
	res.Copper -= other.Copper
	res.WoodPlank -= other.WoodPlank
}

func (res Resources) Sub(other Resources) Resources {
	res.Remove(other)
	return res
}

func (res *Resources) Add(other Resources) {
	res.Rock += other.Rock
	res.Wood += other.Wood
	res.Food += other.Food
	res.Oil += other.Oil
	res.Copper += other.Copper
	res.WoodPlank += other.WoodPlank
}

func (res Resources) Sum(other Resources) Resources {
	res.Add(other)
	return res
}

func (res Resources) Size() int {
	return res.Rock + res.Wood + res.Food + res.Oil + res.Copper + res.WoodPlank
}

type AtomicResources struct {
	lock      sync.RWMutex
	resources Resources
}

func NewAtomicResources(res Resources) AtomicResources {
	return AtomicResources{
		resources: res,
	}
}

func (res *AtomicResources) Load() Resources {
	res.lock.RLock()
	defer res.lock.RUnlock()
	return res.resources
}

// returns the previous value
func (res *AtomicResources) Store(new Resources) Resources {
	res.lock.Lock()
	defer res.lock.Unlock()
	last := res.resources
	res.resources = new
	return last
}

func (res *AtomicResources) Modify(cb func(Resources) Resources) {
	res.lock.Lock()
	defer res.lock.Unlock()
	res.resources = cb(res.resources)
}

func (res *AtomicResources) Sub(other Resources) {
	res.lock.Lock()
	defer res.lock.Unlock()
	res.resources.Remove(other)
}

func (res *AtomicResources) TrySub(other Resources) bool {
	res.lock.Lock()
	defer res.lock.Unlock()
	if res.resources.EnoughFor(other) < 1 {
		return false
	}
	res.resources.Remove(other)
	return true
}

func (res *AtomicResources) Add(other Resources) {
	res.lock.Lock()
	defer res.lock.Unlock()
	res.resources.Add(other)
}
