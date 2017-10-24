(function() {
  Polymer({
    is: "px-sample-app",

    properties: {
      /**
       * Used by the px-app-nav to automatically select the first item.
       * @property selected
       */
      selected: {
        type: Array,
        value: function() {
          return ["dashboard"];
        }
      },
      /**
       * Used by the px-context-browser and px-breadcrumbs as the available items array.
       * @property items
       */
      items: {
        type: Array,
        value: function() {
          return [
            {
              label: "North America",
              id: "North_America",
              children: [
                {
                  label: "United States",
                  id: "United_States",
                  children: [
                    {
                      label: "California",
                      id: "California",
                      children: [
                        {
                          label: "San Diego",
                          id: "San Diego"
                        },
                        {
                          label: "San Francisco",
                          id: "San Francisco"
                        },
                        {
                          label: "San Mateo",
                          id: "San_Mateo"
                        }
                      ]
                    }
                  ]
                }
              ]
            }];
        }
      },
      /**
       * Used by the px-context-browser to select an asset in the hierarchy.
       * @property selectedRoute
       */
      selectedRoute: {
        type: Array,
        value: function() {
          return ["North_America", "United_States", "California", "San Diego"];
        },
        observer: "getSelection"
      },
      /**
       * Used as the title of the dashboard page.
       * @property selectedAsset
       */
      selectedAsset: {
        type: String,
        value: ""
      }
    },
    /**
     * Used by the dom-if to test equality.
     * @param {Array} route
     * @param {String} string
     */
    isEqual(route, string) {
      return route[0] === string;
    },
    /**
     * Gets the selected asset from the context browser
     * to use as the title of the dashboard page.
     */
    getSelection(newValue) {
      this.selectedAsset = this.$.cb.selected.label;
    }
  });
})();
