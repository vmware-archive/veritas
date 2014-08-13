var App = Backbone.View.extend({
  events: {
    'click [type="checkbox"]': 'setViewOptions',
    'keyup #filter': 'filter',
  },

  initialize: function() {
    this.logView = new LogView({
        el: this.$("#log-view"),
    })
    this.histogramView = new HistogramView({
      el: this.$("#histogram"),
    })
    this.setViewOptions()
    this.fetchLogs()
  },

  setViewOptions: function() {
    this.$('#log-view').toggleClass('show-absolute-time', this.$('#show-absolute-time').prop('checked'))    
    this.$('#log-view').toggleClass('show-relative-time', this.$('#show-relative-time').prop('checked'))    
    this.$('#log-view').toggleClass('show-data', this.$('#show-data').prop('checked'))    
    this.$('#log-view').toggleClass('show-raw', this.$('#show-raw').prop('checked'))    
  },

  fetchLogs: function() {
    var that = this
    $.get("/data", function(json) {
      that.logs = JSON.parse(json)
      that.prerenderLogs()
      that.renderLogs()
      that.renderHistogram()
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

  filter: _.throttle(function() {
    var filter = this.$("#filter").val()
    this.$("#log-view").toggleClass('filtered', !!filter)
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