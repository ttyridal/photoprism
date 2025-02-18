package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_RawtherapeeBin(t *testing.T) {
	c := NewConfig(CliTestContext())

	assert.Equal(t, "/usr/bin/rawtherapee-cli", c.RawtherapeeBin())
}

func TestConfig_RawtherapeeEnabled(t *testing.T) {
	c := NewConfig(CliTestContext())
	assert.True(t, c.RawtherapeeEnabled())

	c.options.DisableRawtherapee = true
	assert.False(t, c.RawtherapeeEnabled())
}

func TestConfig_DarktableBin(t *testing.T) {
	c := NewConfig(CliTestContext())

	assert.Equal(t, "/usr/bin/darktable-cli", c.DarktableBin())
}

func TestConfig_DarktablePresets(t *testing.T) {
	c := NewConfig(CliTestContext())

	assert.False(t, c.RawPresets())
}

func TestConfig_DarktableEnabled(t *testing.T) {
	c := NewConfig(CliTestContext())
	assert.True(t, c.DarktableEnabled())

	c.options.DisableDarktable = true
	assert.False(t, c.DarktableEnabled())
}

func TestConfig_SipsBin(t *testing.T) {
	c := NewConfig(CliTestContext())

	bin := c.SipsBin()
	assert.Equal(t, "", bin)
}

func TestConfig_SipsEnabled(t *testing.T) {
	c := NewConfig(CliTestContext())
	c.options.DisableSips = true
	assert.False(t, c.SipsEnabled())
}

func TestConfig_HeifConvertBin(t *testing.T) {
	c := NewConfig(CliTestContext())

	bin := c.HeifConvertBin()
	assert.Contains(t, bin, "/bin/heif-convert")
}

func TestConfig_HeifConvertEnabled(t *testing.T) {
	c := NewConfig(CliTestContext())
	assert.True(t, c.HeifConvertEnabled())

	c.options.DisableHeifConvert = true
	assert.False(t, c.HeifConvertEnabled())
}
