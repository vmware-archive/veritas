var HistogramView = Backbone.View.extend({
   bins: 200,
   events: {
    "mouseenter #histogram-hover": "addTooltip",
    "mousemove #histogram-hover": "updateTooltip",
    "mouseleave #histogram-hover": "removeTooltip",
    "mousedown #histogram-hover": "enableScrolling",
    "mouseup #histogram-hover": "disableScrolling",
   },

   initialize: function(options, delegate) {
    this.delegate = delegate
   },

  renderHistogram: function(logs) {
    var visibleIndices = []
    for (i = 0 ; i < logs.length ; i++) {
        if (logs[i].is_lager) {
            visibleIndices.push(i)
        }
    }

    if (visibleIndices.length < 2) {
        return
    }

    var firstIndex = visibleIndices[0]
    var lastIndex = visibleIndices[visibleIndices.length - 1]
    this.minTime = logs[firstIndex].log.timestamp
    this.maxTime = logs[lastIndex].log.timestamp
    this.dt = (this.maxTime - this.minTime) / this.bins

    var counts = this.binUp(logs, visibleIndices)

    this.largest = _.max(counts)
    this.renderBins(counts, "base")
    this.$el.append('<div id="visible-range-top">')
    this.$el.append('<div id="visible-range-bottom">')
    this.$el.append('<div id="histogram-hover">')
  },

  binUp: function(logs, visibleIndices) {
    var counts = []
    for (i = 0 ; i < this.bins ; i++) {
        counts[i] = 0
    }

    var bin = 0
    for (i = 0 ; i < visibleIndices.length ; i++) {
        while (logs[visibleIndices[i]].log.timestamp > this.minTime + this.dt*(bin+1)) {
            bin += 1
        }
        counts[bin] += 1
    }

    return counts
  },

  renderBins: function(counts, klass) {
    var spacing = 0
    var height = (1 - (this.bins + 1) * spacing)/this.bins
    for (i = 0 ; i < this.bins ; i++) {
        var bin = $('<div class="' + klass + '">')
        bin.css({
            "top": ((spacing * (i + 1) + height * i)*100) + "%",
            "left": 0,
            "width": ((counts[i] / this.largest)*95) + "%",
            "height": (height*100) + "%",
        })
        this.$el.append(bin)
    }
  },

  addFilteredBins: function(logs, indices, klass) {
    var visibleIndices = []
    for (i = 0 ; i < indices.length ; i++) {
        if (logs[indices[i]].is_lager) {
            visibleIndices.push(indices[i])
        }
    }

    if (visibleIndices.length < 2) {
        return
    }

    var counts = this.binUp(logs, visibleIndices)
    this.renderBins(counts, klass)
  },

  clearFilter: function() {
    this.$(".filter").remove()
    this.$el.removeClass("filtered")
  },

  filterLogs: function(logs, visibleIndices) {
    this.addFilteredBins(logs, visibleIndices, "filter")
    this.$el.addClass("filtered")
  },

  clearHighlight: function() {
    this.$(".highlight").remove()
    this.$el.removeClass("highlighted")
  },

  highlightLogs: function(logs, highlightedIndices) {
    this.addFilteredBins(logs, highlightedIndices, "highlight")
    this.$el.addClass("highlighted")
  },


  yPercentageForTimestamp: function(timestamp) {
    if (timestamp == undefined) {
        return 0
    }
    return ((timestamp - this.minTime) / (this.maxTime - this.minTime)) * 100.0
  },

  timestampFromYCoordinate: function(y) {
    return (y / this.$el.height()) * (this.maxTime - this.minTime) + this.minTime
  },

  updateVisibleTimestampRange: function(top, bottom) {
    var yTop = this.yPercentageForTimestamp(top)
    var yBottom = this.yPercentageForTimestamp(bottom)
    this.$("#visible-range-top").css({
        "height": yTop + "%",
    })
    this.$("#visible-range-bottom").css({
        "top": yBottom + "%",
    })
  },

  addTooltip: function(e) {
    this.$el.append('<div id="histogram-hover-time">')
    this.updateTooltip(e)
  },

  updateTooltip: function(e) {
    var height = 30
    var padding = 10
    var totalHeight = height
    var top = Math.min(Math.max(e.offsetY-totalHeight/2, 0), this.$el.height()-height)
    this.$("#histogram-hover-time").css({
        top:top,
        height:height,
        padding:padding,
        lineHeight:(height - 2 * padding)+"px",
    })
    var timestamp = this.timestampFromYCoordinate(e.offsetY)
    this.$("#histogram-hover-time").text(formatRelativeTimestamp(timestamp - this.minTime) + " | " + formatUnixTimestamp(timestamp))
    if (this.scrollingEnabled == true) {
      this.tellDelegateToScrollTo(e.offsetY)
    }
  },

  removeTooltip: function(e) {
    this.$("#histogram-hover-time").remove()
    this.scrollingEnabled = false
  },

  enableScrolling: function(e) {
    this.scrollingEnabled = true
    this.tellDelegateToScrollTo(e.offsetY)
  },

  disableScrolling: function(e) {
    this.scrollingEnabled = false
  },

  tellDelegateToScrollTo: function(offsetY) {
    this.delegate.scrollToTimestamp(this.timestampFromYCoordinate(offsetY))
  }
})
