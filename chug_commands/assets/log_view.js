var LogView = Backbone.View.extend({
  events: {
    'scroll': 'scroll',
  },

  initialize: function(options, delegate) {
    this.showRaw = undefined
    this.delegate = delegate
    this.scrolledTimestamp = 0
  },

  renderLogs: function(logs) {
    var dom = ""
    _.each(logs, function(log) {
      dom += log.dom
    })
    this.$el[0].innerHTML = dom
    this.count = logs.length
    this.logs = logs
    this.computeTimestampsAtEdgesAndNotifyDelegate()
  },

  setShow: function(showProperty, enabled) {
    this.$el.toggleClass(showProperty, !!enabled)
    if (showProperty == "show-raw" && this.showRaw !== !!enabled) {
      this.visibleSubset = undefined
    }
    this.showRaw = !!enabled
    if (this.logs) {
      this.scrollToTimestamp(this.scrolledTimestamp)
    }
  },

  clearFilter: function() {
    this.$el.removeClass("filtered")
    this.$(".show").removeClass("show")
    this.visibleSubset = undefined
    this.filteredIndices = undefined
    this.computeTimestampsAtEdgesAndNotifyDelegate()
  },

  filterLogs: function(visibleIndices) {
    this.visibleSubset = undefined
    this.filteredIndices = visibleIndices
    this.$el.addClass("filtered")
    for (i = 0 ; i < visibleIndices.length ; i++) {
        this.$("#log-"+visibleIndices[i]).addClass("show")
    }
    this.computeTimestampsAtEdgesAndNotifyDelegate()
  },

  scroll: function() {
    this.computeTimestampsAtEdgesAndNotifyDelegate()
  },

  computeTimestampsAtEdgesAndNotifyDelegate: function() {
    var top = 0
    var bottom = top + this.$el.height()
    if (!this.visibleSubset) {
      this.computeVisibleSubset()
    }
    if (this.visibleSubset.length == 0) {
      return
    }
    var idxTop = this.findEntryNear(top, 0, this.visibleSubset.length)
    var idxBottom = this.findEntryNear(bottom, 0, this.visibleSubset.length)
    this.scrolledTimestamp = this.logs[this.visibleSubset[Math.floor((idxTop + idxBottom)/2)]].timestamp
    this.delegate.updateVisibleTimestampRange(this.logs[this.visibleSubset[idxTop]].timestamp, this.logs[this.visibleSubset[idxBottom]].timestamp)
  },

  findEntryNear: function(offset, a, b) {
    if (b - a < 2) {
      return a
    }
    var midpoint = Math.floor((a + b) / 2)
    var midPointTop = this.$("#log-" + this.visibleSubset[midpoint]).position().top
    var midPointBottom = this.$("#log-" + this.visibleSubset[midpoint]).height() + midPointTop
    if (midPointTop <= offset && offset <= midPointBottom) {
      return midpoint
    }
    if (offset < midPointTop) {
      return this.findEntryNear(offset, a, midpoint)
    } else {
      return this.findEntryNear(offset, midpoint, b)
    }
  },

  computeVisibleSubset: function() {
    this.visibleSubset = []
    var length = this.logs.length
    var index = function(i) { return i }
    if (this.filteredIndices) {
      length = this.filteredIndices.length
      var fi = this.filteredIndices
      index = function(i) { return fi[i]}
    }
    for (i = 0 ; i < length ; i++) {
      idx = index(i)
      if (this.logs[idx].is_lager || this.showRaw) {
        this.visibleSubset.push(idx)
      }
    }
  },

  scrollToTimestamp: function(timestamp) {
    this.scrolledTimestamp = timestamp
    if (!this.visibleSubset) {
      this.computeVisibleSubset()
    }
    if (this.visibleSubset.length == 0) {
      return
    }

    for (i = 0 ; i < this.visibleSubset.length ; i++) {
      idx = this.visibleSubset[i]
      if (this.logs[idx].timestamp >= timestamp) {
        var top = $("#log-"+idx).position().top
        var height = $("#log-"+idx).height()
        var midpoint = top + height/2 + this.$el.scrollTop()
        this.$el.scrollTop(midpoint - this.$el.height()/2)
        return
      }
    }
  },
})
