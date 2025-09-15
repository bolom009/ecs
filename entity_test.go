package ecs_test

import (
	"fmt"
	"testing"

	"github.com/bolom009/ecs"
)

func TestEntity_NewEntity_Should_Create_A_Correct_Mask(t *testing.T) {
	entity := ecs.NewEntity([]ecs.Component{
		&mockComponent{name: "position", mask: 1},
	})
	if entity.Mask() != 1 {
		t.Errorf("Entity mask should be 1, but got %d", entity.Mask())
	}
}

func TestEntity_Add_Should_Work_With_Multiple_Components(t *testing.T) {
	entity := ecs.NewEntity([]ecs.Component{
		&mockComponent{name: "position", mask: 1},
	})
	entity.Add(&mockComponent{name: "velocity", mask: 2})
	if entity.Mask() != 3 {
		t.Errorf("Entity mask should be 3, but got %d", entity.Mask())
	}
}

func TestEntity_Add_Should_Not_Add_Existing_Component(t *testing.T) {
	entity := ecs.NewEntity([]ecs.Component{
		&mockComponent{name: "position", mask: 1},
		&mockComponent{name: "velocity", mask: 2},
	})
	entity.Add(&mockComponent{name: "velocity", mask: 2})
	if entity.Components.Len() != 2 {
		t.Errorf("Component count should be 2, but got %d", entity.Components.Len())
	}
}

func TestEntity_Get_Should_Return_Component(t *testing.T) {
	entity := ecs.NewEntity([]ecs.Component{
		&mockComponent{name: "position", mask: 1},
	})
	component := entity.Get(1)
	if component == nil {
		t.Error("Component should not be nil")
	}
}

func BenchmarkEntity_GetMany_Should_Return_Component(b *testing.B) {
	list := make([]ecs.Component, 50)
	for i := range list {
		inc := 1 << (i & 63)
		list[i] = &mockComponent{
			name: fmt.Sprintf("position%v", i),
			mask: uint64(inc),
		}
	}

	entity := ecs.NewEntity(list)
	mask := uint64(1 << (len(list) - 1&63))

	b.ResetTimer()
	b.ReportAllocs()

	for b.Loop() {
		_ = entity.Get(mask)
	}
}

func TestEntity_Remove_Should_Work_With_Multiple_Components(t *testing.T) {
	entity := ecs.NewEntity([]ecs.Component{
		&mockComponent{name: "position", mask: 1},
		&mockComponent{name: "size", mask: 2},
		&mockComponent{name: "velocity", mask: 4},
	})
	entity.Remove(4)
	if entity.Mask() != 3 {
		t.Errorf("Entity mask should be 1, but got %d", entity.Mask())
	}
}

/*
       _   _ _
 _   _| |_(_) |___
| | | | __| | / __|
| |_| | |_| | \__ \
 \__,_|\__|_|_|___/
*/

type mockComponent struct {
	mask  uint64
	name  string
	value interface{}
}

func (c *mockComponent) Mask() uint64 { return c.mask }

func (c *mockComponent) Name() string { return c.name }
