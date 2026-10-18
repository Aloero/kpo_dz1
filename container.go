package main

type serviceContainer struct {
	providers map[string]func(*serviceContainer) interface{}
	instances map[string]interface{}
}

func newServiceContainer() *serviceContainer {
	return &serviceContainer{
		providers: make(map[string]func(*serviceContainer) interface{}),
		instances: make(map[string]interface{}),
	}
}

func (c *serviceContainer) register(name string, provider func(*serviceContainer) interface{}) {
	c.providers[name] = provider
}

func (c *serviceContainer) resolve(name string) interface{} {
	if instance, ok := c.instances[name]; ok {
		return instance
	}
	provider, ok := c.providers[name]
	if !ok {
		panic("сервис не зарегистрирован: " + name)
	}
	c.instances[name] = provider(c)
	return c.instances[name]
}
