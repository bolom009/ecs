package ecs

import "github.com/bolom009/ecs/intmap"

type defaultEntityManager struct {
	entities    []*Entity
	mapEntities *intmap.Map[uint32, *Entity]
}

// NewEntityManager creates a new defaultEntityManager and returns its address.
func NewEntityManager() *defaultEntityManager {
	return &defaultEntityManager{
		entities:    make([]*Entity, 0),
		mapEntities: intmap.New[uint32, *Entity](100),
	}
}

// Add entries to the manager.
func (m *defaultEntityManager) Add(entities ...*Entity) {
	m.entities = append(m.entities, entities...)
	for _, entity := range entities {
		m.mapEntities.Put(entity.Id, entity)
	}
}

// Entities returns all the entities.
func (m *defaultEntityManager) Entities() []*Entity {
	return m.entities
}

// FilterByMask returns the mapped entities, which Components mask matched.
func (m *defaultEntityManager) FilterByMask(mask uint64) (entities []*Entity) {
	// Allocate the worst-case amount of memory (all entities needed).
	entities = make([]*Entity, len(m.entities))
	index := 0
	for _, e := range m.entities {
		// Use the pre-calculated Components maskSlice.
		observed := e.Mask()
		// Add the entity to the filter list, if all Components are found.
		if observed&mask == mask {
			// Direct access
			entities[index] = e
			index++
		}
	}
	// Return only the needed slice.
	return entities[:index]
}

// Get a specific entity by Id.
func (m *defaultEntityManager) Get(id uint32) *Entity {
	if v, ok := m.mapEntities.Get(id); ok {
		return v
	}

	return nil
}

// Remove a specific entity.
func (m *defaultEntityManager) Remove(entity *Entity) {
	for i, e := range m.entities {
		if e.Id == entity.Id {
			copy(m.entities[i:], m.entities[i+1:])
			m.entities[len(m.entities)-1] = nil
			m.entities = m.entities[:len(m.entities)-1]
			m.mapEntities.Del(e.Id)
			break
		}
	}
}
