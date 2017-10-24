(function () {
  Polymer({
    is: 'px-sample-dashboard',
    properties: {
      title: {
        type: String,
        value: ""
      },
      mapData: {
        type: Object,
        value: function() {
          return {};
        }
      }
    }
  });
})();
