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
	UpdateFunc        func(interface{})
	settings          *js.Object
	closeToKillUpdate chan interface{}
}

func (tp *TestPlugin) onSettingsChanged(settings *js.Object) {
	tp.settings = settings
	tp.updateNow()
	//	tp.UpdateFunc(tp.settings)
}

func (tp *TestPlugin) updateNow() {
	data := map[string]string{
		"animal":   "dinosaur",
		"datatext": tp.settings.Get("datatext").String(),
	}
	//tp.UpdateFunc(js.MakeWrapper(data))
	tp.UpdateFunc(data)
}

func (tp *TestPlugin) onDispose() {
	close(tp.closeToKillUpdate)
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

	NewInstance: func(settings, newInstanceCallback, updateCallback *js.Object) {
		pr("In NewInstance")
		pl := new(TestPlugin)
		pl.closeToKillUpdate = make(chan interface{})
		pr("Made new TestPlugin, assigning settings.")
		pl.settings = settings
		pr("Creating updatefunc closure.")
		//pl.UpdateFunc = func(i *js.Object) { updateCallback.Call("call", js.Undefined, i) }
		pl.UpdateFunc = func(i interface{}) { updateCallback.Call("call", js.Undefined, i) }
		// The update func.
		pr("Creating heartbeat update goroutine")
		go func(pl *TestPlugin) {
			for {
				select {
				case <-pl.closeToKillUpdate:
					return
				case <-time.After(5 * time.Second):
					//pl.UpdateFunc(pl.settings)
					pl.updateNow()
				}
			}
		}(pl)
		pr("Making return wrapper")
		wrapper := map[string]interface{}{
			"Plugin":            pl,
			"updateNow":         pl.updateNow,
			"onDispose":         pl.onDispose,
			"onSettingsChanged": pl.onSettingsChanged,
		}
		pr("Returning wrapper through newInstanceCallback")
		js.Global.Set("testPluginInstance", wrapper)
		newInstanceCallback.Call("call", js.Undefined, wrapper)
	},
}

func main() {
	pr("Registering plugin")
	freeboard.LoadDatasourcePlugin(TestDefinition)
}
