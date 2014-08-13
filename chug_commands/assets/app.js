var App = Backbone.View.extend({
  events: {
    'click [type="checkbox"]': 'setViewOptions',
    'keyup #filter': 'filter',
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
    var filter = this.$("#filter").val()
    this.$("#histogram").toggleClass('filtered', !!filter)
    this.logView.clearFilter()
    this.histogramView.clearFilter()

    if (!filter) {
      return
    }

    var filters = filter.split(" ")
    var regularExpressions = []
    for (var i = 0 ; i < filters.length ; i++) {
        regularExpressions[i] = new RegExp(filters[i])
    }
    var visibleIndices = []
    _.each(this.logs, function(log, index) {
        for (var i = 0 ; i < regularExpressions.length ; i++) {
            if (!regularExpressions[i].test(log.searchText)) {
                return
            }            
        }
        visibleIndices.push(index)
    })

    this.logView.filterLogs(visibleIndices)
    this.histogramView.filterLogs(this.logs, visibleIndices)
  }, 300),
})