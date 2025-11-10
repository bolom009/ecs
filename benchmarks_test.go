package ecs_test

import (
	"fmt"
	"math/rand/v2"
	"testing"

	"github.com/bolom009/ecs"
	"github.com/mlange-42/ark-tools/app"
	arkecs "github.com/mlange-42/ark/ecs"
)

func BenchmarkEntityManager_Get_With_1_Entity_Id_Found(b *testing.B) {
	m := ecs.NewEntityManager()
	m.Add(ecs.NewEntity(nil))
	for b.Loop() {
		m.Get(1)
	}
}

func BenchmarkEntityManager_Get_With_1000_Entities_Id_Not_Found(b *testing.B) {
	m := ecs.NewEntityManager()
	for i := 0; i < 1000; i++ {
		m.Add(ecs.NewEntity(nil))
	}
	for b.Loop() {
		m.Get(1000)
	}
}

// BenchmarkEntityManager_Get_With_1000_Entities_Id-16    	168744212	         7.069 ns/op
func BenchmarkEntityManager_Get_With_1000_Entities_Id(b *testing.B) {
	m := ecs.NewEntityManager()
	for i := 0; i < 1000; i++ {
		m.Add(ecs.NewEntity(nil))
	}

	m.Add(ecs.NewEntity(nil))

	b.ResetTimer()

	for b.Loop() {
		m.Get(1001)
	}
}

func BenchmarkEntityManager_FilterByMask_With_1000_Entities(b *testing.B) {
	m := ecs.NewEntityManager()
	for i := 0; i < 1000; i++ {
		m.Add(ecs.NewEntity([]ecs.Component{
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
	entity := ecs.NewEntity(generateComponents([]string{
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
	entity := ecs.NewEntity(generateComponents([]string{
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

// BenchmarkEngine_Run/2_system(s)_with_1000_entities-16         	   48218	     24332 ns/op
func BenchmarkEngine_Run(b *testing.B) {
	entityCounts := []int{1000}
	systemCounts := []int{1}
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

type mockupMovementSystem struct{}

func (s *mockupMovementSystem) Process(entityManager ecs.EntityManager) (state int) {
	for _, e := range entityManager.FilterByMask(1 | 2) {
		pos := e.Get(1).(*position)
		vel := e.Get(2).(*velocity)
		pos.x += vel.x * 0.33
		pos.y += vel.y * 0.33
	}

	return ecs.StateEngineContinue
}
func (s *mockupMovementSystem) Setup()    {}
func (s *mockupMovementSystem) Teardown() {}

// BenchmarkTest-16    	     273	   4622808 ns/op	      48 B/op	       3 allocs/op
// BenchmarkTest-16    	        448	        2632051 ns/op	  802816 B/op	       1 allocs/op
// BenchmarkArk-16    	    	5017	    227828  ns/op	       0 B/op	       0 allocs/op
// BenchmarkArk-16    	    	7854	    152535 ns/op	       0 B/op	       0 allocs/op
func BenchmarkTest(b *testing.B) {
	em := ecs.NewEntityManager(100000)
	em.Add(generateEntities(100000)...)

	sm := ecs.NewSystemManager()
	sm.Add(&mockupMovementSystem{})

	engine := ecs.NewDefaultEngine(em, sm)
	engine.Setup()

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		engine.Tick()
	}
}

type arkMoveSystem struct {
	filter *arkecs.Filter2[position, velocity]
}

func (s *arkMoveSystem) Initialize(w *arkecs.World) {
	s.filter = s.filter.New(w)
}

// Update the system.
func (s *arkMoveSystem) Update(w *arkecs.World) {
	query := s.filter.Query()
	for query.Next() {
		pos, vel := query.Get()
		pos.x += vel.x * 0.33
		pos.y += vel.y * 0.33
	}
}

// Finalize the system.
func (s *arkMoveSystem) Finalize(w *arkecs.World) {}

func BenchmarkArk(b *testing.B) {
	arkManager := app.New(100000).Seed(123)
	arkManager.AddSystem(&arkMoveSystem{})

	// Create a component mapper.
	mapper := arkecs.NewMap2[position, velocity](&arkManager.World)
	mapper.NewBatch(100000, &position{x: 1, y: 1}, &velocity{x: 1, y: 1})

	arkManager.Initialize()

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		arkManager.Update()
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
		out[i] = ecs.NewEntity([]ecs.Component{
			&position{x: 1, y: 1},
			&velocity{x: 1, y: 1},
		})
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

type position struct {
	x, y float64
}

func (p *position) Mask() uint64 {
	return 1
}

type velocity struct {
	x, y float64
}

func (p *velocity) Mask() uint64 {
	return 2
}

type data struct {
	thingy int
	dingy  float64
	mingy  bool
	numgy  int
}

func (p *data) Mask() uint64 {
	return 3
}

// mockupUseAllEntitiesSystem works on all entities from the defaultEntityManager which represents the worst-case scenario for performance.
type mockupUseAllEntitiesSystem struct{}

func (s *mockupUseAllEntitiesSystem) Process(entityManager ecs.EntityManager) (state int) {
	dt := 0.25
	for _, e := range entityManager.FilterByMask(1 | 2) {
		pos := e.Get(1).(*position)
		dir := e.Get(2).(*velocity)
		pos.x += dir.x * dt
		pos.y += dir.y * dt
	}

	for _, e := range entityManager.FilterByMask(3) {
		data := e.Get(3).(*data)
		data.thingy = (data.thingy + 1) % 1000000
		data.dingy += 0.0001 * dt
		data.mingy = !data.mingy
		data.numgy = rand.Int()
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
