package main

import (
	"github.com/cathalgarvey/go-freeboard"
	"github.com/gopherjs/gopherjs/js"
)

// CatsPlugin is me noodling around with the freeboard interface.
type CatsPlugin struct {
	UpdateFunc        func(interface{})
	settings          *js.Object
	closeToKillUpdate chan interface{}
}

// CurrentSettings satisfies the freeboard.DsPlugin interface.
func (tp *CatsPlugin) CurrentSettings() *js.Object {
	return tp.settings
}

// OnSettingsChanged satisfies the freeboard.DsPlugin interface.
func (tp *CatsPlugin) OnSettingsChanged(settings *js.Object) {
	tp.settings = settings
	tp.UpdateNow()
}

// UpdateNow satisfies the freeboard.DsPlugin interface.
func (tp *CatsPlugin) UpdateNow() {
	data := map[string]interface{}{
		"animal":   tp.settings.Get("animal"),
		"datatext": tp.settings.Get("datatext"),
		"refine":   tp.settings.Get("refine"),
	}
	tp.UpdateFunc(data)
}

// OnDispose satisfies the freeboard.DsPlugin interface.
func (tp *CatsPlugin) OnDispose() {
	close(tp.closeToKillUpdate)
}

// TestDefinition defines a plugin that provides some user-set text.
var TestDefinition = freeboard.DsPluginDefinition{
	TypeName:    "catsplugin",
	DisplayName: "Cats",
	Description: "This is a demo Golang plugin about cats",
	Settings: []freeboard.FBSetting{
		freeboard.FBSetting{
			Name:               "catname",
			DisplayName:        "Favourite Cat Name",
			Description:        "What would you call your favourite cat?",
			Type:               freeboard.SettingTextType,
			DefaultStringValue: "Meow",
		},
		freeboard.FBSetting{
			Name:        "animal",
			DisplayName: "Animal",
			Description: "Favourite animal.",
			Type:        freeboard.SettingOptionType,
			Options: []freeboard.FBSettingOpt{
				freeboard.FBSettingOpt{Name: "Tiger"},
				freeboard.FBSettingOpt{Name: "Lion"},
				freeboard.FBSettingOpt{Name: "Tigon"},
				freeboard.FBSettingOpt{Name: "Liger"},
			},
		},
		freeboard.FBSetting{
			Name:        "refine",
			DisplayName: "Refined Animal Preference",
			Description: "More details on what kinda cat you like",
			Type:        freeboard.SettingArrayType,
			Settings: []freeboard.FBSettingSet{
				freeboard.FBSettingSet{
					Name:        "preferred_number",
					DisplayName: "Preferred number per cage",
					Type:        freeboard.SettingNumberType,
				},
				freeboard.FBSettingSet{
					Name:        "preferred_colour",
					DisplayName: "Preferred cat colour",
					Type:        freeboard.SettingTextType,
				},
			},
		},
	},
	NewInstance: func(settings *js.Object, updateCallback func(interface{})) freeboard.DsPlugin {
		pl := new(CatsPlugin)
		pl.settings = settings
		pl.UpdateFunc = updateCallback
		pl.closeToKillUpdate = freeboard.MakeUpdateTicker(pl, 5)
		return pl
	},
}

func main() {
	println("Registering plugin")
	freeboard.FB.LoadGoDatasourcePlugin(TestDefinition)
}
