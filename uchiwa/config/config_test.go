package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	_, err := Load("../foo.bar")
	assert.NotNil(t, err, "should return an error when file does not exist")

	_, err = Load("../uchiwa.go")
	assert.NotNil(t, err, "should return an error when it cannot parse a file")

	conf, err := Load("../../fixtures/config_test.json")
	assert.Nil(t, err, "got unexpected error: %s", err)

	// private config
	assert.NotEqual(t, "*****", conf.Uchiwa.User, "Uchiwa user in private config shouldn't be masked")
	assert.NotEqual(t, "*****", conf.Uchiwa.Pass, "Uchiwa pass in private config shouldn't be masked")
	for i := range conf.Sensu {
		assert.NotEqual(t, "*****",conf.Sensu[i].User, "Sensu APIs user in private config shouldn't be masked")
		assert.NotEqual(t, "*****", conf.Sensu[i].Pass, "Sensu APIs pass in private config shouldn't be masked")
	}

	// public config
	public := conf.GetPublic()
	assert.Equal(t, "*****", public.Uchiwa.User, "Uchiwa user in public config should be masked")
	assert.Equal(t, "*****", public.Uchiwa.Pass, "Uchiwa pass in public config should be masked")
	for i := range public.Sensu {
		assert.Equal(t, "*****", public.Sensu[i].User, "Sensu APIs user in public config should be masked")
		assert.Equal(t, "*****", public.Sensu[i].Pass, "Sensu APIs pass in public config should be masked")
	}

}
