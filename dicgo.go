package dicgo

import (
	"fmt"
)

// Container is main DIC interface
type Container interface {
	// SetService new service definition
	SetService(id string, fact Factory, singleton bool)
	// SetSingleton adds new singleton service
	SetSingleton(id string, fact Factory)
	// SetFactory adds new non sigleton service factory
	SetFactory(id string, fact Factory)
	// SetValue adds non-service value
	SetValue(id string, val interface{})
	// Get return service
	Get(id string) (interface{}, bool)
	// Service returns service, or panics if not defined
	Service(id string) interface{}
	// Raw returns raw service factory
	Raw(id string) (func(c Container) interface{}, bool)
	// Has returns if container has service definition
	Has(id string) bool
	// Del deletes service from the container
	Del(id string)
}

// Factory service "constructor"
type Factory func(c Container) interface{}

type def struct {
	id        string
	fact      Factory
	singleton bool
}

type container struct {
	defs     map[string]*def
	services map[string]interface{}
}

// NewContainer creates and returns new DIC
func NewContainer() Container {
	return &container{
		defs:     map[string]*def{},
		services: map[string]interface{}{},
	}
}

func (c *container) SetService(id string, fact Factory, singleton bool) {
	c.defs[id] = &def{
		id:        id,
		fact:      fact,
		singleton: singleton,
	}
}

func (c *container) SetSingleton(id string, fact Factory) {
	c.SetService(id, fact, true)
}

func (c *container) SetFactory(id string, fact Factory) {
	c.SetService(id, fact, false)
}

func (c *container) SetValue(id string, val interface{}) {
	c.SetService(id, nil, true)
	c.services[id] = val
}

func (c *container) Get(id string) (interface{}, bool) {
	def, ok := c.defs[id]
	if !ok {
		return nil, ok
	}
	//non sigleton
	if !def.singleton {
		return def.fact(c), true
	}
	//cache singleton service
	service, ok := c.services[id]
	if ok {
		return service, true
	}
	c.services[id] = def.fact(c)
	return c.services[id], true
}

func (c *container) Service(id string) interface{} {
	s, ok := c.Get(id)
	if !ok {
		panic(fmt.Errorf("Service [%s] not defined", id))
	}
	return s
}

func (c *container) Raw(id string) (func(c Container) interface{}, bool) {
	def, ok := c.defs[id]
	if !ok {
		return nil, ok
	}
	if def == nil {
		return nil, false
	}
	return def.fact, true
}

func (c *container) Del(id string) {
	_, ok := c.services[id]
	if ok {
		delete(c.services, id)
	}
	_, ok = c.defs[id]
	if ok {
		delete(c.defs, id)
	}
}

func (c *container) Has(id string) bool {
	_, ok := c.defs[id]
	return ok
}
