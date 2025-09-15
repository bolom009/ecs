package ecs_test

import (
	"fmt"
	"testing"

	"github.com/bolom009/ecs"
)

func TestEntityManager_Entities_Should_Have_No_Entity_At_Start(t *testing.T) {
	m := ecs.NewEntityManager()
	if len(m.Entities()) != 0 {
		t.Errorf("EntityManager should have no entity at start, but got %d", len(m.Entities()))
	}
}

func TestEntityManager_Entities_Should_Have_One_Entity_After_Adding_One_Entity(t *testing.T) {
	m := ecs.NewEntityManager()
	m.Add(&ecs.Entity{})
	if len(m.Entities()) != 1 {
		t.Errorf("EntityManager should have one entity, but got %d", len(m.Entities()))
	}
}

func TestEntityManager_Entities_Should_Have_Two_Entities_After_Adding_Two_Entities(t *testing.T) {
	m := ecs.NewEntityManager()
	m.Add(ecs.NewEntity(nil))
	m.Add(ecs.NewEntity(nil))
	if len(m.Entities()) != 2 {
		t.Errorf("EntityManager should have two entities, but got %d", len(m.Entities()))
	}
}

func TestEntityManager_Entities_Should_Have_One_Entity_After_Removing_One_Of_Two_Entities(t *testing.T) {
	m := ecs.NewEntityManager()
	e1 := ecs.NewEntity(nil)
	e2 := ecs.NewEntity(nil)
	m.Add(e1)
	m.Add(e2)
	m.Remove(e2)
	if len(m.Entities()) != 1 {
		t.Errorf("EntityManager should have one entity after removing one, but got %d", len(m.Entities()))
	}
	if m.Entities()[0].Id != 2 {
		t.Errorf("Entity should have correct Id, but got %d", m.Entities()[0].Id)
	}
}

func TestEntityManager_FilterByMask_Should_Return_No_Entity_Out_Of_One(t *testing.T) {
	em := ecs.NewEntityManager()
	e := ecs.NewEntity([]ecs.Component{
		&mockComponent{name: "position", mask: 1},
	})
	em.Add(e)
	filtered := em.FilterByMask(2)
	if len(filtered) != 0 {
		t.Errorf("EntityManager should return no entity, but got %d", len(filtered))
	}
}

func TestEntityManager_FilterByMask_Should_Return_One_Entity_Out_Of_One(t *testing.T) {
	em := ecs.NewEntityManager()
	e := ecs.NewEntity([]ecs.Component{
		&mockComponent{name: "position", mask: 1},
	})
	em.Add(e)
	filtered := em.FilterByMask(1)
	if len(filtered) != 1 {
		t.Errorf("EntityManager should return one entity, but got %d", len(filtered))
	}
}

func TestEntityManager_FilterByMask_Should_Return_One_Entity_Out_Of_Two(t *testing.T) {
	em := ecs.NewEntityManager()
	e1 := ecs.NewEntity([]ecs.Component{
		&mockComponent{name: "position", mask: 1},
	})
	e2 := ecs.NewEntity([]ecs.Component{
		&mockComponent{name: "position", mask: 1},
		&mockComponent{name: "size", mask: 2},
	})
	em.Add(e1, e2)
	filtered := em.FilterByMask(2)
	if len(filtered) != 1 {
		t.Errorf("EntityManager should return one entity, but got %d", len(filtered))
	}
	if filtered[0].Id != 7 {
		t.Errorf("Entity should have correct Id, but got %d", filtered[0].Id)
	}
}

func TestEntityManager_FilterByMask_Should_Return_Two_Entities_Out_Of_Three(t *testing.T) {
	em := ecs.NewEntityManager()
	e1 := ecs.NewEntity([]ecs.Component{
		&mockComponent{name: "position", mask: 1},
	})
	e2 := ecs.NewEntity([]ecs.Component{
		&mockComponent{name: "position", mask: 1},
		&mockComponent{name: "size", mask: 2},
	})
	e3 := ecs.NewEntity([]ecs.Component{
		&mockComponent{name: "position", mask: 1},
		&mockComponent{name: "size", mask: 2},
	})
	em.Add(e1, e2, e3)
	filtered := em.FilterByMask(2)
	if len(filtered) != 2 {
		t.Errorf("EntityManager should return two entities, but got %d", len(filtered))
	}
	if filtered[0].Id != 9 {
		t.Errorf("Entity should have correct Id, but got %d", filtered[0].Id)
	}
	if filtered[1].Id != 10 {
		t.Errorf("Entity should have correct Id, but got %d", filtered[1].Id)
	}
}

func TestEntityManager_FilterByMask_Should_Return_Three_Entities_Out_Of_Three(t *testing.T) {
	em := ecs.NewEntityManager()
	e1 := ecs.NewEntity([]ecs.Component{
		&mockComponent{name: "position", mask: 1},
		&mockComponent{name: "size", mask: 2},
	})
	e2 := ecs.NewEntity([]ecs.Component{
		&mockComponent{name: "position", mask: 1},
		&mockComponent{name: "size", mask: 2},
	})
	e3 := ecs.NewEntity([]ecs.Component{
		&mockComponent{name: "position", mask: 1},
		&mockComponent{name: "size", mask: 2},
		&mockComponent{name: "transform", mask: 4},
	})
	em.Add(e1, e2, e3)
	filtered := em.FilterByMask(1 | 2)
	if len(filtered) != 3 {
		t.Errorf("EntityManager should return three entities, but got %d", len(filtered))
	}
	if filtered[0].Id != 11 {
		t.Errorf("Entity should have correct Id, but got %d", filtered[0].Id)
	}
	if filtered[1].Id != 12 {
		t.Errorf("Entity should have correct Id, but got %d", filtered[1].Id)
	}
	if filtered[2].Id != 13 {
		t.Errorf("Entity should have correct Id, but got %d", filtered[2].Id)
	}
}

func TestEntityManager_FilterByMask(t *testing.T) {
	const (
		MaskPosition = uint64(1 << iota)
		MaskRotation
		MaskSize
		MaskScale
		MaskVelocity
		MaskAgent
	)

	em := ecs.NewEntityManager()
	e1 := ecs.NewEntity([]ecs.Component{
		&mockComponent{name: "position", mask: MaskPosition},
		&mockComponent{name: "size", mask: MaskSize},
		&mockComponent{name: "rotation", mask: MaskRotation},
	})
	e2 := ecs.NewEntity([]ecs.Component{
		&mockComponent{name: "position", mask: MaskPosition},
		&mockComponent{name: "size", mask: MaskSize},
	})
	e3 := ecs.NewEntity([]ecs.Component{
		&mockComponent{name: "position", mask: MaskPosition},
		&mockComponent{name: "size", mask: MaskSize},
	})

	em.Add(e1, e2, e3)
	e1.Remove(MaskRotation)

	filtered := em.FilterByMask(MaskPosition)
	if len(filtered) != 3 {
		t.Errorf("EntityManager should return three entities, but got %d", len(filtered))
	}
}

func TestEntityManager_FilterByMask_Only_One_From_Three(t *testing.T) {
	const (
		MaskPosition = uint64(1 << iota)
		MaskRotation
		MaskSize
	)

	em := ecs.NewEntityManager()
	e1 := ecs.NewEntity([]ecs.Component{
		&mockComponent{name: "position", mask: MaskPosition},
		&mockComponent{name: "size", mask: MaskSize},
		&mockComponent{name: "rotation", mask: MaskRotation},
	})
	e2 := ecs.NewEntity([]ecs.Component{
		&mockComponent{name: "position", mask: MaskPosition},
		&mockComponent{name: "size", mask: MaskSize},
	})
	e3 := ecs.NewEntity([]ecs.Component{
		&mockComponent{name: "position", mask: MaskPosition},
		&mockComponent{name: "size", mask: MaskSize},
	})

	em.Add(e1, e2, e3)

	filtered := em.FilterByMask(MaskRotation)
	if len(filtered) != 1 {
		t.Errorf("EntityManager should return three entities, but got %d", len(filtered))
	}
}

func TestEntityManager_Get_Should_Return_Entity(t *testing.T) {
	em := ecs.NewEntityManager()
	e1 := ecs.NewEntity([]ecs.Component{
		&mockComponent{name: "position", mask: 1},
		&mockComponent{name: "size", mask: 2},
	})
	e2 := ecs.NewEntity([]ecs.Component{
		&mockComponent{name: "position", mask: 1},
		&mockComponent{name: "size", mask: 2},
	})
	em.Add(e1, e2)
	if e := em.Get(20); e == nil {
		t.Error("Entity should not be nil")
	}
	if e := em.Get(21); e == nil {
		t.Error("Entity should not be nil")
	}
}

func BenchmarkEntityManager_FilterByMask(b *testing.B) {
	em := ecs.NewEntityManager()

	entities := make([]*ecs.Entity, 500)
	for i := range entities {
		if i >= 0 && i <= 8 {
			entities[i] = createBenchEntity(8)
		} else if i > 8 && i <= 16 {
			entities[i] = createBenchEntity(16)
		} else if i > 16 && i <= 450 {
			entities[i] = createBenchEntity(24)
		} else if i > 450 {
			entities[i] = createBenchEntity(30)
		}
	}

	em.Add(entities...)

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		em.FilterByMask(536870912)
	}
}

func createBenchEntity(lenComponents int) *ecs.Entity {
	components := make([]ecs.Component, lenComponents)
	for i := 0; i < lenComponents; i++ {
		components[i] = &mockComponent{name: fmt.Sprintf("name_%v", i), mask: uint64(1 << i)}
	}

	return ecs.NewEntity(components)
}
