var App = Backbone.View.extend({
  events: {
    'click [type="checkbox"]': 'setViewOptions',
    'keyup #filter': 'filter',
    'keyup #exclude': 'filter',
    'keyup #highlight': 'highlight',
  },

  initialize: function() {
    this.logView = new LogView({
        el: this.$("#log-view"),
    }, this)
    this.histogramView = new HistogramView({
      el: this.$("#histogram"),
    }, this)
    this.setViewOptions()
    this.fetchLogs()
  },

  setViewOptions: function() {
    this.logView.setShow("show-absolute-time", this.$('#show-absolute-time').prop('checked'))
    this.logView.setShow("show-relative-time", this.$('#show-relative-time').prop('checked'))
    this.logView.setShow("show-data", this.$('#show-data').prop('checked'))
    if (this.$('#show-data').prop('checked')) {
      this.$('#show-big-data').removeAttr("disabled");
      this.logView.setShow("show-big-data", this.$('#show-big-data').prop('checked'))
    } else {
      this.logView.setShow("show-big-data", false)
      this.$('#show-big-data').attr("disabled", true);
    }
    this.logView.setShow("show-raw", this.$('#show-raw').prop('checked'))
  },

  fetchLogs: function() {
    var that = this
    $.get("/data", function(json) {
      that.logs = JSON.parse(json)
      that.prerenderLogs()
      that.renderHistogram()
      that.renderLogs()
    })
  },

  prerenderLogs: function() {
    var renderer = new LogRenderer()
    renderer.prerenderLogs(this.logs)
  },

  renderLogs: function() {
    this.logView.renderLogs(this.logs)
  },

  renderHistogram: function() {
    this.histogramView.renderHistogram(this.logs)
  },

  updateVisibleTimestampRange: function(top, bottom) {
    this.histogramView.updateVisibleTimestampRange(top, bottom)
  },

  scrollToTimestamp: function(timestamp) {
    this.logView.scrollToTimestamp(timestamp)
  },

  filter: _.throttle(function() {
    this.logView.clearFilter()
    this.histogramView.clearFilter()

    var filter = this.$("#filter").val()
    var exclude = this.$("#exclude").val()
    if (!filter && !exclude) {
      this.visibleIndices = undefined
      this.unthrottledHighlight()
      return
    }

    this.visibleIndices = this.selectIndices(filter, exclude)

    this.logView.filterLogs(this.visibleIndices)
    this.histogramView.filterLogs(this.logs, this.visibleIndices)
    this.unthrottledHighlight()
  }, 300),

  highlight: _.throttle(function() {
    this.unthrottledHighlight()
  }, 300),

  unthrottledHighlight: function() {
    this.logView.clearHighlight()
    this.histogramView.clearHighlight()

    var filter = this.$("#highlight").val()
    if (!filter) {
      return
    }

    var highlightedIndices = this.selectIndices(filter, "", this.visibleIndices)

    this.logView.highlightLogs(highlightedIndices)
    this.histogramView.highlightLogs(this.logs, highlightedIndices)
  },

  selectIndices: function(filter, exclude, subset) {
    var filters = filter.split(" ")
    var filterRegularExpressions = []
    for (var i = 0 ; i < filters.length ; i++) {
      if (filters[i]) {
        filterRegularExpressions[i] = new RegExp(filters[i])
      }
    }

    var exclusions = exclude.split(" ")
    var exclusionRegularExpressions = []
    for (var i = 0 ; i < exclusions.length ; i++) {
      if (exclusions[i]) {
        exclusionRegularExpressions[i] = new RegExp(exclusions[i])
      }
    }

    console.log(filters, exclusions)
    var length = this.logs.length
    var index = function(i) { return i }
    if (subset) {
      length = subset.length
      index = function(i) { return subset[i] }
    }

    var selectedIndices = []
    for (i = 0 ; i < length ; i++) {
      var idx = index(i)
      var found = true
      for (var j = 0 ; j < filterRegularExpressions.length ; j++) {
          if (!filterRegularExpressions[j].test(this.logs[idx].searchText)) {
              found = false
              break
          }
      }
      for (var j = 0 ; j < exclusionRegularExpressions.length ; j++) {
          if (exclusionRegularExpressions[j].test(this.logs[idx].searchText)) {
              found = false
              break
          }
      }
      if (found) {
        selectedIndices.push(idx)
      }
    }

    return selectedIndices
  },

})
