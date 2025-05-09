package ecs_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/bolom009/ecs"
)

func BenchmarkEntityManager_Get_With_1_Entity_Id_Found(b *testing.B) {
	m := ecs.NewEntityManager()
	m.Add(ecs.NewEntity("foo", nil))
	for b.Loop() {
		m.Get("foo")
	}
}

func BenchmarkEntityManager_Get_With_1000_Entities_Id_Not_Found(b *testing.B) {
	m := ecs.NewEntityManager()
	for i := 0; i < 1000; i++ {
		m.Add(ecs.NewEntity("foo", nil))
	}
	for b.Loop() {
		m.Get("1000")
	}
}

// BenchmarkEntityManager_Get_With_1000_Entities_Id-16    	168744212	         7.069 ns/op
func BenchmarkEntityManager_Get_With_1000_Entities_Id(b *testing.B) {
	m := ecs.NewEntityManager()
	for i := 0; i < 1000; i++ {
		m.Add(ecs.NewEntity("3d78b074-dae6-419c-be63-6565375e3eba", nil))
	}
	searchID := "a11efca1-e420-4869-a424-95539ce1dad7"
	m.Add(ecs.NewEntity(searchID, nil))

	b.ResetTimer()

	for b.Loop() {
		m.Get(searchID)
	}
}

func BenchmarkEntityManager_FilterByMask_With_1000_Entities(b *testing.B) {
	m := ecs.NewEntityManager()
	for i := 0; i < 1000; i++ {
		m.Add(ecs.NewEntity(fmt.Sprintf("%d", i), []ecs.Component{
			&mockComponent{name: "position", mask: 1},
			&mockComponent{name: "size", mask: 2},
			&mockComponent{name: "velocity", mask: 3},
		}))
	}

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		m.FilterByMask(1 | 2 | 3)
	}
}

func BenchmarkEntity_Get_Should_Return_Component(b *testing.B) {
	entity := ecs.NewEntity("e", generateComponents([]string{
		"position", "rotation", "scale", "material", "security",
		"damage", "agent", "rvo", "move_speed", "aggro", "attack_speed",
		"attack_range", "network_identity", "team", "health", "mana",
		"death_timer", "texture", "melee", "state", "target", "velocity",
		"effects", "pathfinding", "flocking", "follow",
	}))

	for b.Loop() {
		_ = entity.Get(100)
	}
}

func BenchmarkEntity_Get_Should_Remove_Component(b *testing.B) {
	entity := ecs.NewEntity("e", generateComponents([]string{
		"position", "rotation", "scale", "material", "security",
		"damage", "agent", "rvo", "move_speed", "aggro", "attack_speed",
		"attack_range", "network_identity", "team", "health", "mana",
		"death_timer", "texture", "melee", "state", "target", "velocity",
		"effects", "pathfinding", "flocking", "follow",
	}))

	for b.Loop() {
		entity.Remove(100)
	}
}

func BenchmarkEngine_Run(b *testing.B) {
	entityCounts := []int{100, 1000, 10000}
	systemCounts := []int{1, 2, 4}
	for _, systemCount := range systemCounts {
		for _, entityCount := range entityCounts {
			b.Run(fmt.Sprintf("%d system(s) with %d entities", systemCount, entityCount), func(b *testing.B) {
				b.ResetTimer()
				em := ecs.NewEntityManager()
				em.Add(generateEntities(entityCount)...)
				sm := ecs.NewSystemManager()
				sm.Add(generateUseAllEntitiesSystems(systemCount)...)
				engine := ecs.NewDefaultEngine(em, sm)
				engine.Setup()
				defer engine.Teardown()
				for b.Loop() {
					engine.Run()
				}
			})
		}
	}
}

/*
       _   _ _
 _   _| |_(_) |___
| | | | __| | / __|
| |_| | |_| | \__ \
 \__,_|\__|_|_|___/
*/

func generateComponents(entries []string) []ecs.Component {
	components := make([]ecs.Component, len(entries))
	for i, entry := range entries {
		components[i] = &mockComponent{
			name:  entry,
			mask:  uint64(i + 1),
			value: fmt.Sprintf("%s-%d", entry, i+1),
		}
	}

	return components
}

func generateEntities(count int) []*ecs.Entity {
	out := make([]*ecs.Entity, count)
	for i := 0; i < count; i++ {
		out[i] = ecs.NewEntity(
			fmt.Sprintf("e%d", rand.Uint64()),
			[]ecs.Component{
				&mockComponent{mask: 1},
			},
		)
	}
	return out
}

func generateUseAllEntitiesSystems(count int) []ecs.System {
	out := make([]ecs.System, count)
	for i := 0; i < count-1; i++ {
		out[i] = &mockupUseAllEntitiesSystem{}
	}
	out[count-1] = &mockupShouldStopSystem{}
	return out
}

// mockupUseAllEntitiesSystem works on all entities from the defaultEntityManager which represents the worst-case scenario for performance.
type mockupUseAllEntitiesSystem struct{}

func (s *mockupUseAllEntitiesSystem) Process(entityManager ecs.EntityManager) (state int) {
	for range entityManager.FilterByMask(1) {
	}
	return ecs.StateEngineContinue
}
func (s *mockupUseAllEntitiesSystem) Setup() {
}
func (s *mockupUseAllEntitiesSystem) Teardown() {
}

// mockupShouldStopSystem is the last System in the queue and should stop the defaultEngine.
type mockupShouldStopSystem struct{}

func (s *mockupShouldStopSystem) Process(entityManager ecs.EntityManager) (state int) {
	for range entityManager.FilterByMask(1) {
	}
	return ecs.StateEngineStop
}
func (s *mockupShouldStopSystem) Setup() {
}
func (s *mockupShouldStopSystem) Teardown() {
}
