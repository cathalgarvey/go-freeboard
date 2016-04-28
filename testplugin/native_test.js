(function(){
  var TestPlugin = function(settings, updateCallback) {

    var self = this;
		var updateTimer = null;
		var currentSettings = settings;

		function updateRefresh(refreshTime) {
			if (updateTimer) {
				clearInterval(updateTimer);
			}
			updateTimer = setInterval(function () {
				self.updateNow();
			}, refreshTime);
		}

		updateRefresh(5 * 1000);

		this.updateNow = function () {
      var output = {
        "datatext": settings.datatext,
        "animal": "dinosaur"
      }
      updateCallback(output);
		}

		this.onDispose = function () {
			clearInterval(updateTimer);
			updateTimer = null;
		}

		this.onSettingsChanged = function (newSettings) {
			currentSettings = newSettings;
			updateRefresh(currentSettings.refresh * 1000);
			self.updateNow();
		}
  }

  freeboard.loadDatasourcePlugin({
    type_name: "testnative",
    display_name: "Native JS test plugin",
    description: "JS representing what test.go is supposed to do",
    //external_scripts: [], // Decommenting this line breaks everything. WTF.
    settings: [
      {
        name: "datatext",
        display_name: "Data Text",
        description: "Text to provide as data.",
        type: "text",
        default_value: "Foobar"
      }
    ],
    newInstance: function(settings, newInstanceCallback, updateCallback) {
      console.log("Returning new plugin object to newInstanceCallback.");
      newInstanceCallback(new TestPlugin(settings, updateCallback));
    }
  });
})()
//*/
