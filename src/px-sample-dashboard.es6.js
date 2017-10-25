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
      },
      chartData: {
        type: Array,
        value: function() {
          return [];
        }
      },
      _mapData: {
        type: Object,
        value: function() {
          return this.mapData;
        }
      },
      _costOfInvestment: {
        type: Number,
        value: 0
      },
      _powerConsumption: {
        type: Number,
        value: 0
      },
      _renewableOutput: {
        type: Number,
        value: 0
      },
      _carbonOffset: {
        type: Number,
        value: 0
      },
      _savingsResponse: {
        type: Object,
        value: function() {
          return {};
        }
      },
      _carbonNuetralDate: {
        type: Number,
        value: 2035,
        observer: 'sliderChanged'
      },
      _residentialSliderPercentage: {
        type: Number,
        value: 0,
        observer: 'sliderChanged'
      },
      _commercialSliderPercentage: {
        type: Number,
        value: 0,
        observer: 'sliderChanged'
      },
      _industrialSliderPercentage: {
        type: Number,
        value: 0,
        observer: 'sliderChanged'
      },
      pieChartData: {
        type: Array,
        value: function() {
          return [{
            "x": 43.4,
            "y": "2018 Energy & Utility Budget",
            "colorIndex": 2,
            "backgroundColor": "rgb(123,188,0)"
          }, {
            "x": 1,
            "y": "Additional Recommended Subsidies",
            "colorIndex": 1,
            "backgroundColor": "rgb(226,141,23)"
          }, {
            "x": 15,
            "y": "Other Sources",
            "colorIndex": 0,
            "backgroundColor": "rgb(90,191,248)"
          }]
        }
      },
      _inc1: {
        type: Number,
        value: 0
      },
      _inc2: {
        type: Number,
        value: 0
      },
      _inc3: {
        type: Number,
        value: 0
      },
      _inc4: {
        type: Number,
        value: 0
      },
      _inc5: {
        type: Number,
        value: 0
      },
      _inc6: {
        type: Number,
        value: 0
      },
      _inc7: {
        type: Number,
        value: 0
      },
      _inc8: {
        type: Number,
        value: 0
      },
      _inc9: {
        type: Number,
        value: 0
      },
      _inc1Perc: {
        type: Number,
        value: 0
      },
      _inc2Perc: {
        type: Number,
        value: 0
      },
      _inc3Perc: {
        type: Number,
        value: 0
      },
      _inc4Perc: {
        type: Number,
        value: 0
      },
      _inc5Perc: {
        type: Number,
        value: 0
      },
      _inc6Perc: {
        type: Number,
        value: 0
      },
      _inc7Perc: {
        type: Number,
        value: 0
      },
      _inc8Perc: {
        type: Number,
        value: 0
      },
      _inc9Perc: {
        type: Number,
        value: 0
      }
    },
    attached() {
      this.set('_mapData', Object.assign({}, this.mapData));
      var markers = this.$.markers;
      markers.set('data', this._mapData);
      markers.update();
    },
    updateMap(evt) {
      this.debounce('mapUpdate', function() {
        if(evt.detail.value === 0) return;
        var percent;
        if(evt.target.id === "res") {
          percent = 100000 * parseInt(evt.detail.value) / 100;
        }
        else {
          percent = 1000 * parseInt(evt.detail.value) / 100;
        }
        var arr = this.mapData.features.slice(0, percent);
        this.set('_mapData.type', 'FeatureCollection');
        this.set('_mapData.features', arr);
        var markers = this.$.markers;
        markers.set('data', this._mapData);
        markers.update();
      },100);
    },
    handleSavingsResponse(e) {
      this._powerConsumption = (this._savingsResponse.PredictedConsumptionKWH / 1000000).toLocaleString();
      this._renewableOutput = (this._savingsResponse.SolarEnergyGenerationKWH / 1000000).toLocaleString();
      this._costOfInvestment = (this._savingsResponse.InitialCost / 1000000).toLocaleString();
      this._carbonOffset = (this._savingsResponse.SolarOffsetKWH / 1000000).toLocaleString();
      this.toggleClass('green', parseInt(this._carbonOffset) < 0, this.$.kpi4);
    },
    sliderChanged(e) {
      var base = "https://savings-app.run.aws-usw02-pr.ice.predix.io";

      var res = (this._residentialSliderPercentage / 100) * 500000;
      var com = (this._commercialSliderPercentage / 100) * 10000;
      var ind = (this._industrialSliderPercentage / 100) * 2000;
      var cn = this._carbonNuetralDate;

      this.$.savingsCalc.url = base + "/v1/savings?com=" + com + "&res=" + res + "&ind=" + ind + "&targetYear=" + cn;
      this.$.savingsCalc.generateRequest();

      this._adjustEconomicNumbers(res, this._residentialSliderPercentage);
    },
    _adjustEconomicNumbers(res, perc) {
      if(perc < 20) {
        if((res * .65) >= 79409) {
          this._inc9 = 79409;
          this._inc9Perc = (100).toFixed(2);
        } else {
          this._inc9 = (res * .75).toFixed(0);
          this._inc9Perc = ( ( ( res * .75) /  79409 ) * 100).toFixed(2);
        }

        this._inc8 = (res * .20).toFixed(0);
        this._inc8Perc = ( ( ( res * .20) /  169548 ) * 100).toFixed(2);

        this._inc7 = (res * .13).toFixed(0);
        this._inc7Perc = ( ( ( res * .13) / 139502 ) * 100).toFixed(2);

        this._inc6 = (res * .02);
        this._inc6Perc = ( ( ( res * .02) / 181352 ) * 100).toFixed(2);

      } else if(perc >= 20 && perc < 30) {
        if((res * .50) >= 79409) {
          this._inc9 = 79409;
          this._inc9Perc = (100).toFixed(2);
        } else {
          this._inc9 = (res * .50).toFixed(0);
          this._inc9Perc = ( ( ( res * .50) /  79409 ) * 100).toFixed(2);
        }

        this._inc8 = (res * .25).toFixed(0);
        this._inc8Perc = ( ( ( res * .25) /  169548 ) * 100).toFixed(2);

        this._inc7 = (res * .17).toFixed(0);
        this._inc7Perc = ( ( ( res * .17) / 139502 ) * 100).toFixed(2);

        this._inc6 = (res * .08).toFixed(0);
        this._inc6Perc = ( ( ( res * .08) / 181352 ) * 100).toFixed(2);

      } else if(perc >= 30 && perc < 50) {
        if((res * .28) >= 79409) {
          this._inc9 = 79409;
          this._inc9Perc = (100).toFixed(2);
        } else {
          this._inc9 = (res * .28).toFixed(0);
          this._inc9Perc = ( ( ( res * .28) /  79409 ) * 100).toFixed(2);
        }

        this._inc8 = (res * .32).toFixed(0);
        this._inc8Perc = ( ( ( res * .32) /  169548 ) * 100).toFixed(2);

        this._inc7 = (res * .25).toFixed(0);
        this._inc7Perc = ( ( ( res * .25) / 139502 ) * 100).toFixed(2);

        this._inc6 = (res * .15).toFixed(0);
        this._inc6Perc = ( ( ( res * .15) / 181352 ) * 100).toFixed(2);
      }

    }
  });
})();
