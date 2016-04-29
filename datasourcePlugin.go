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

type settingType string

var (
	// SettingTextType is used for text input.
	SettingTextType settingType = "text"
	// SettingNumberType is used for number input
	SettingNumberType settingType = "number"
	// SettingCalculatedType is used to permit access to datasources or to calculated JS.
	SettingCalculatedType settingType = "calculated"
	// SettingBooleanType is used to provide a checkbox.
	SettingBooleanType settingType = "boolean"
	// SettingOptionType is used to offer a select-list of options.
	SettingOptionType settingType = "option"
	// SettingArrayType is used to ask for multiple rows of data.
	SettingArrayType settingType = "array"
)

// FBSettingOpt is an option for the "option" type of setting.
type FBSettingOpt struct {
	// Name of the option.
	Name string
	// If not specified, name is used.
	Value string
}

// FBSettingSet is a setting for the "setting" type of setting (O_o)
type FBSettingSet struct {
	Name        string
	DisplayName string
	// Presumably only text or numeric make sense here..
	Type settingType
}

// FBSetting is a settings object.
type FBSetting struct {
	// Name must be a valid JS name and should be unique.
	Name string
	// DisplayName is the name presented to the user.
	DisplayName string
	// Description is what's presented to users.
	Description string
	// Type is the type of the setting.
	Type settingType
	// Options are required for option-type settings.
	Options []FBSettingOpt
	// Settings is required for "array" type settings.
	Settings []FBSettingSet
	// DefaultValues are the default value. Optional. String takes precedence in text.
	DefaultStringValue string
	// DefaltIntValue or DefaultFloatValue can be used as default values for
	// number types; whichever is nonzero is used. If both are nonzero, panic!
	// If both are zero, then the default is left unset.
	DefaultIntValue   int
	DefaultFloatValue float64
}

// ToFBInterface compiles a setting to a map-able representation
// expected by the FreeBoard interface.
func (set FBSetting) ToFBInterface() map[string]interface{} {
	output := make(map[string]interface{})
	output["name"] = set.Name
	output["display_name"] = set.DisplayName
	output["description"] = set.Description
	output["type"] = string(set.Type)
	switch set.Type {
	case SettingTextType, SettingCalculatedType:
		// Assuming that calculated type can have defaults?
		{
			if set.DefaultStringValue != "" {
				output["default_value"] = set.DefaultStringValue
			} else if set.DefaultIntValue != 0 {
				output["default_value"] = set.DefaultIntValue
			}
		}
	case SettingNumberType:
		{
			if set.DefaultIntValue != 0 && set.DefaultFloatValue == 0 {
				output["default_value"] = set.DefaultIntValue
			} else if set.DefaultFloatValue != 0.0 && set.DefaultIntValue == 0 {
				output["default_value"] = set.DefaultFloatValue
			} else if set.DefaultIntValue != 0 && set.DefaultFloatValue != 0.0 {
				panic("Cannot have defaults for both int and float numeric values.")
			}
		}
	case SettingOptionType:
		{
			output["options"] = make([]map[string]string, 0, len(set.Options))
			for _, opt := range set.Options {
				o := make(map[string]string)
				o["name"] = opt.Name
				if opt.Value != "" {
					o["value"] = opt.Value
				} else {
					o["value"] = opt.Name
				}
				output["options"] = append(output["options"].([]map[string]string), o)
			}
		}
	case SettingArrayType:
		{
			output["settings"] = make([]map[string]string, 0, len(set.Settings))
			for _, st := range set.Settings {
				s := make(map[string]string)
				s["name"] = st.Name
				s["display_name"] = st.DisplayName
				s["type"] = string(st.Type)
				output["settings"] = append(output["settings"].([]map[string]string), s)
			}
		}
		//case SettingBooleanType:  // No special handling required?
	default:
		panic("Unknown setting type: " + string(set.Type))
	}
	return output
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
