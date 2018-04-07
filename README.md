# DICGO

*dicGO* is simplest possible DIC (depenency injection container) written after [minidic](https://github.com/DrBenton/minidic)

## Install

```bash
go get -u github.com/rukavina/dicgo
```

## Usage

```golang
import "github.com/rukavina/dicgo"

//create container
dicgo.NewContainer()

type service1 struct {
	id  string
	num int
}

type service2 struct {
	s1   service1
	name string
}

//add singleton service
dicgo.AddSingleton("s1", func(c Container)interface{}{
    return &s1{
        id: "s1",
        num: 1,
    }
})

//add another service, depending on s1 service
dicgo.AddSingleton("s2", func(c Container)interface{}{
    return &s2{
        //watch out!, unlike Get, Service methods panics if service not found
        s1: c.Service("s1").(*service1)
        name: "s2",
    }
})

//add just value
dicgo.AddValue("config", map[string]interface{}{
    "test": 123,
})

//add factory, non sigleton, always returns new instance
dicgo.AddValue("gen", func(c Container)interface{}{
    return &s1{
        id: "new",
        num: 1,
    }
})
```