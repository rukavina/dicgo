package dicgo

import "testing"

func TestValue(t *testing.T) {
	c := NewContainer()

	c.SetValue("val", "test")

	injectedValue, ok := c.Get("val")
	if !ok || injectedValue != "test" {
		t.Errorf("Expected \"test\", got %+v", injectedValue)
	}
}

func TestSingleton(t *testing.T) {
	c := NewContainer()

	c.SetSingleton("service", func(cont Container) interface{} {
		return "test"
	})

	injectedValue, ok := c.Get("service")
	if !ok || injectedValue != "test" {
		t.Errorf("Expected \"test\", got %+v", injectedValue)
	}
}

type service1 struct {
	id  string
	num int
}

type service2 struct {
	s1   service1
	name string
}

func TestSingletonDep(t *testing.T) {
	c := NewContainer()

	c.SetSingleton("s1", func(cont Container) interface{} {
		return "hello"
	})

	c.SetSingleton("s2", func(cont Container) interface{} {
		s1, ok := cont.Get("s1")
		if !ok {
			return ""
		}
		return s1.(string) + " world"
	})

	injectedValue, ok := c.Get("s2")
	if !ok || injectedValue != "hello world" {
		t.Errorf("Expected \"hello world\", got %+v", injectedValue)
	}
}

func TestSingletonPointer(t *testing.T) {
	c := NewContainer()

	c.SetSingleton("s1", func(cont Container) interface{} {
		return &service1{
			id:  "s1",
			num: 1,
		}
	})

	injectedValue, ok := c.Get("s1")
	s1 := injectedValue.(*service1)
	if !ok || s1.id != "s1" {
		t.Errorf("Expected \"s1\", got %+v", s1)
	}

	s1.num = 2

	injectedValue, ok = c.Get("s1")
	s := injectedValue.(*service1)
	if s != s1 {
		t.Errorf("Not the same pointer: %v, %v", s1, 2)
	}
	if !ok || s.num != 2 {
		t.Errorf("Expected num \"2\", got %d", s.num)
	}
}

func TestFactgory(t *testing.T) {
	c := NewContainer()

	c.SetFactory("s1", func(cont Container) interface{} {
		return &service1{
			id:  "s1",
			num: 1,
		}
	})

	injectedValue, ok := c.Get("s1")
	s1 := injectedValue.(*service1)
	if !ok || s1.id != "s1" {
		t.Errorf("Expected \"s1\", got %+v", s1)
	}

	s1.num = 2

	injectedValue, ok = c.Get("s1")
	s := injectedValue.(*service1)
	if s == s1 {
		t.Errorf("The same pointer: %v, %v", s1, 2)
	}
	if !ok || s.num == 2 {
		t.Errorf("Expected num \"1\", got %d", s.num)
	}
}

func TestHasDel(t *testing.T) {
	c := NewContainer()

	c.SetSingleton("s1", func(cont Container) interface{} {
		return &service1{
			id:  "s1",
			num: 1,
		}
	})

	ok := c.Has("s1")
	if !ok {
		t.Error("Has no s1 service")
	}
	ok = c.Has("s2")
	if ok {
		t.Error("Has non existing s2 service")
	}
	s2, ok := c.Get("s2")
	if ok {
		t.Errorf("Got non existing s2 service, %+v", s2)
	}
	c.Del("s2")
	c.Del("s1")
	ok = c.Has("s1")
	if ok {
		t.Error("Still has deleted s1 service")
	}
	s1, ok := c.Get("s1")
	if ok {
		t.Errorf("Got deleted s1 service, %+v", s1)
	}
}
