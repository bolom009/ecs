package ecs

import "sync/atomic"

var idCounter atomic.Uint64

// Entity is simply a composition of one or more Components with an Id.
type Entity struct {
	Components map[uint64]Component
	Id         uint32 `json:"id"`
	Masked     uint64 `json:"masked"`
}

// Add a component.
func (e *Entity) Add(cn ...Component) {
	for _, c := range cn {
		cMask := c.Mask()
		if e.Masked&cMask == cMask {
			continue
		}

		e.Components[c.Mask()] = c
		e.Masked = e.Masked | cMask
	}
}

// Get a component by its bitmask.
func (e *Entity) Get(mask uint64) Component {
	c, ok := e.Components[mask]
	if ok {
		return c
	}

	return nil
}

// Mask returns a pre-calculated maskSlice to identify the Components.
func (e *Entity) Mask() uint64 {
	return e.Masked
}

// Remove a component by using its maskSlice.
func (e *Entity) Remove(mask uint64) {
	c, ok := e.Components[mask]
	if ok {
		delete(e.Components, mask)
		e.Masked = e.Masked &^ c.Mask()
	}

	//modified := false
	//for i, c := range e.Components {
	//	if c.Mask() == mask {
	//		copy(e.Components[i:], e.Components[i+1:])
	//		e.Components[len(e.Components)-1] = nil
	//		e.Components = e.Components[:len(e.Components)-1]
	//		e.Masked = e.Masked &^ c.Mask()
	//		break
	//	}
	//}
	//if modified {
	//	e.Masked = maskSlice(e.Components)
	//}
}

// NewEntity creates a new entity and pre-calculates the component maskSlice.
func NewEntity(components []Component) *Entity {
	eId := newId()
	e := &Entity{
		Components: make(map[uint64]Component),
		Id:         eId,
		Masked:     maskSlice(components),
	}

	for _, c := range components {
		e.Components[c.Mask()] = c
	}

	return e
}

func newId() uint32 {
	for {
		val := idCounter.Load()
		if idCounter.CompareAndSwap(val, val+1) {
			return uint32(val)
		}
	}
}

func maskSlice(components []Component) uint64 {
	mask := uint64(0)
	for _, c := range components {
		mask = mask | c.Mask()
	}
	return mask
}
