package structs

// Data is a structure for holding public data fetched from the Sensu APIs and exposed by the endpoints
type Data struct {
	Aggregates    []interface{}
	Checks        []interface{}
	Clients       []interface{}
	Dc            []*Datacenter
	Events        []interface{}
	Health        Health
	Stashes       []interface{}
	Subscriptions []string
}

// Datacenter is a structure for holding the information about a datacenter
type Datacenter struct {
	Name  string         `json:"name"`
	Info  Info           `json:"info"`
	Stats map[string]int `json:"stats"`
}

// Health is a structure for holding health informaton about Sensu & Uchiwa
type Health struct {
	Sensu  map[string]SensuHealth `json:"sensu"`
	Uchiwa string                 `json:"uchiwa"`
}

// SensuHealth is a structure for holding health information about a specific sensu datacenter
type SensuHealth struct {
	Output string `json:"output"`
}

// Info is a structure for holding the /info API information
type Info struct {
	Redis     redis     `json:"redis"`
	Sensu     sensu     `json:"sensu"`
	Transport transport `json:"transport"`
}

type redis struct {
	Connected bool `json:"connected"`
}

type sensu struct {
	Version string `json:"version"`
}

type transport struct {
	Connected  bool            `json:"connected"`
	Keepalives transportStatus `json:"keepalives"`
	Results    transportStatus `json:"results"`
}

type transportStatus struct {
	Messages  int `json:"messages"`
	Consumers int `json:"consumers"`
}
