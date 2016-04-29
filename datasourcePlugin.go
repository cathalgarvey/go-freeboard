package freeboard

import (
	"time"

	"github.com/gopherjs/gopherjs/js"
)

// DsPlugin is a constructed plugin that presents
// the expected interface for freeboard to access
// data.
type DsPlugin interface {
	// Called when new settings are given.
	OnSettingsChanged(*js.Object)
	// A public function we must implement
	// that will be called when the user wants
	// to manually refresh the datasource
	UpdateNow()
	// A public function we must implement that
	// will be called when this instance of this
	// plugin is no longer needed. Do anything
	// you need to cleanup after yourself here.
	OnDispose()
	// An additional function, for sanity's sake,
	// which should return the current internal
	// settings object.
	CurrentSettings() *js.Object
}

// WrapDsPlugin converts the exported function set of a DsPlugin
// into a map consisting of the required JS function set, with
// the other DsPlugin methods presented with lowercase leading
// characters in the JS style. Also present is "plugin", which
// directly references the DsPlugin object.
func WrapDsPlugin(dsp DsPlugin) map[string]interface{} {
	return map[string]interface{}{
		"plugin":            dsp,
		"updateNow":         dsp.UpdateNow,
		"onDispose":         dsp.OnDispose,
		"onSettingsChanged": dsp.OnSettingsChanged,
		"currentSettings":   dsp.CurrentSettings,
	}
}

// MakeUpdateTicker creates a goroutine that polls a DataSource's
// UpdateNow method every few seconds (as provided). It returns a
// channel to close when this goroutine should be stopped.
// This is provided as a helper because of how common these update
// tickers are in freeboard plugins. Just store the ticker channel
// and close it in the OnDispose() method.
func MakeUpdateTicker(dsp DsPlugin, seconds int) chan interface{} {
	closeToKillUpdate := make(chan interface{})
	go func(dsp DsPlugin, seconds int) {
		for {
			select {
			case <-closeToKillUpdate:
				return
			case <-time.After(time.Duration(seconds) * time.Second):
				dsp.UpdateNow()
			}
		}
	}(dsp, seconds)
	return closeToKillUpdate
}

// DsPluginDefinition is a Datasource Plugin
type DsPluginDefinition struct {
	// TypeName should be a unique name for this plugin.
	// Must be a valid JS name. Avoid potential naming conflicts!
	TypeName string

	// The displayed name, need not be unique.
	DisplayName string

	// Front-facing description of this plugin.
	Description string

	// ExternalScripts are outside script URIs required for
	// this plugin. They will be loaded prior to the plugin.
	ExternalScripts []string

	// Settings are the user-facing options for this plugin.
	// They will be converted to a map[string]interface{} for
	// consumption by the plugin.
	Settings []FBSetting

	// NewInstance is called to create a new plugin. This is the
	// Go wrapper. It is passed the calculated settings based on
	// the definition's settings array, as a js.Object; this should
	// be retained. It is also passed one special function, which
	// is a simple go wrapper around the updateCallback function
	// given by the FreeBoard NewInstance function.
	// Notably absent is NewInstanceCallback; this is handled under
	// the Go layer for you. All you have to do is return the
	// prepared DsPlugin-interfacing plugin object.
	NewInstance func(settings *js.Object, updateCallback func(interface{})) DsPlugin
}

// ToFBInterface returns a map for FreeBoard's loadDatasourcePlugin func.
func (dsp DsPluginDefinition) ToFBInterface() map[string]interface{} {
	output := make(map[string]interface{})
	output["type_name"] = dsp.TypeName
	output["display_name"] = dsp.DisplayName
	output["description"] = dsp.Description
	// Exposing an empty array for ExternalScripts breaks freeboard plugins.
	if dsp.ExternalScripts != nil && len(dsp.ExternalScripts) > 0 {
		output["external_scripts"] = dsp.ExternalScripts
	}
	settingSlice := make([]map[string]interface{}, 0, len(dsp.Settings))
	for _, s := range dsp.Settings {
		settingSlice = append(settingSlice, s.ToFBInterface())
	}
	output["settings"] = settingSlice
	output["newInstance"] = func(settings, newInstanceCallback, updateCallback *js.Object) {
		Plugin := dsp.NewInstance(settings, func(i interface{}) { updateCallback.Invoke(i) })
		wrapper := WrapDsPlugin(Plugin)
		newInstanceCallback.Invoke(wrapper)
	}
	return output
}
