package meter

func NewGroup(clients ...Meter) *Group {
	m := &Group{
		clients: make(map[string]Meter),
	}

	for _, c := range clients {
		m.clients[c.Name()] = c
	}

	return m
}

type Group struct {
	clients map[string]Meter
}

func (g Group) Name() string {
	return "group"
}

func (g Group) Get(name string) (Meter, bool) {
	if val, ok := g.clients[name]; ok {
		return val, true
	}

	return nil, false
}
