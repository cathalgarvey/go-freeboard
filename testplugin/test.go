package main

import (
	"time"

	"github.com/cathalgarvey/go-freeboard"
	"github.com/gopherjs/gopherjs/js"
)

var (
	cons *js.Object
)

func init() {
	cons = js.Global.Get("console")
}

func pr(s ...string) {
	args := make([]interface{}, 0, len(s))
	for _, a := range s {
		args = append(args, a)
	}
	cons.Call("log", args...)
}

// TestPlugin is me noodling around with the freeboard interface.
type TestPlugin struct {
	UpdateFunc func(*js.Object)
	settings   *js.Object
}

// Called when new settings are given.
//func (tp *TestPlugin) onSettingsChanged(settings map[string]interface{}) {
func (tp *TestPlugin) onSettingsChanged(settings *js.Object) {
	tp.settings = settings
	tp.UpdateFunc(tp.settings)
}

// A public function we must implement
// that will be called when the user wants
// to manually refresh the datasource
func (tp *TestPlugin) updateNow() {
	pr("Update called.")
	tp.UpdateFunc(tp.settings)
	pr("Update delivered.")
}

// A public function we must implement that
// will be called when this instance of this
// plugin is no longer needed. Do anything
// you need to cleanup after yourself here.
func (tp *TestPlugin) onDispose() {
}

// TestDefinition defines a plugin that provides some user-set text.
var TestDefinition = freeboard.DsPluginDefinition{
	// TypeName should be a unique name for this plugin.
	// Must be a valid JS name. Avoid potential naming conflicts!
	TypeName: "testplugin1",

	// The displayed name, need not be unique.
	DisplayName: "test plugin",

	// Front-facing description of this plugin.
	Description: "This is a test plugin",

	// ExternalScripts are outside script URIs required for
	// this plugin. They will be loaded prior to the plugin.
	ExternalScripts: []string{},

	// Settings are the user-facing options for this plugin.
	// They will be converted to a map[string]interface{} for
	// consumption by the plugin.
	Settings: []freeboard.FBSetting{
		freeboard.FBSetting{
			// Name must be a valid JS name and should be unique.
			Name: "datatext",
			// DisplayName is the name presented to the user.
			DisplayName: "Data Text",
			// Description is what's presented to users.
			Description: "Text to provide as data.",
			// Type is the type of the setting.
			Type: freeboard.SettingTextType,
			// DefaultValues are the default value. Optional. String takes precedence.
			DefaultStringValue: "Foobar",
		},
	},

	// NewInstance is called to create a new plugin. It is
	// passed the calculated settings based on the definition's
	// settings array. This should be kept by the plugin.
	// It is also passed two special functions:
	// * newInstanceCallback should be called at the end of NewInstance
	//   with the new plugin. Remember may have to use js.New to construct
	//   your new plugin from a constructor function?
	// * updateCallback should be called with new data whenever it's
	//   ready for freeboard. This should be kept by the new instance.
	NewInstance: func(settings, newInstanceCallback, updateCallback *js.Object) {
		pr("In NewInstance")
		pl := new(TestPlugin)
		pr("Made new TestPlugin, assigning settings.")
		pl.onSettingsChanged(settings)
		pr("Assigning updatefunc.")
		pl.UpdateFunc = func(i *js.Object) { updateCallback.Call("apply", i) }
		go func(pl *TestPlugin) {
			for {
				pl.UpdateFunc(pl.settings)
				time.Sleep(5 * time.Second)
			}
		}(pl)
		pr("Making wrapper")
		wrapper := js.MakeWrapper(pl)
		pr("Returning wrapper through newInstanceCallback")
		newInstanceCallback.Call("apply", wrapper)
	},
}

func main() {
	freeboard.LoadDatasourcePlugin(TestDefinition)
}
