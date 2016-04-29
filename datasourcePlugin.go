package freeboard

import "github.com/gopherjs/gopherjs/js"

// DsPlugin is a constructed plugin that presents
// the expected interface for freeboard to access
// data.
type DsPlugin interface {
	// Called when new settings are given.
	onSettingsChanged(map[string]interface{})
	// A public function we must implement
	// that will be called when the user wants
	// to manually refresh the datasource
	updateNow()
	// A public function we must implement that
	// will be called when this instance of this
	// plugin is no longer needed. Do anything
	// you need to cleanup after yourself here.
	onDispose()
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

	// NewInstance is called to create a new plugin. It is
	// passed the calculated settings based on the definition's
	// settings array. This should be kept by the plugin.
	// It is also passed two special functions:
	// * newInstanceCallback should be called at the end of NewInstance
	//   with the new plugin. Remember may have to use js.New to construct
	//   your new plugin from a constructor function?
	// * updateCallback should be called with new data whenever it's
	//   ready for freeboard. This should be kept by the new instance.
	//NewInstance func(settings map[string]interface{}, newInstanceCallback func(DsPlugin), updateCallback func(interface{}))
	//NewInstance func(settings map[string]interface{}, newInstanceCallback func(map[string]interface{}), updateCallback func(interface{}))
	//NewInstance func(settings map[string]interface{}, newInstanceCallback func(*js.Object), updateCallback func(interface{}))
	NewInstance func(settings, newInstanceCallback, updateCallback *js.Object)
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
	output["newInstance"] = dsp.NewInstance
	return output
}
