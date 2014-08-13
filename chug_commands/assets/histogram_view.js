var HistogramView = Backbone.View.extend({
    bins: 200,

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

  clearFilter: function() {
    this.$(".filter").remove()
  },

  filterLogs: function(logs, inputVisibleIndices) {
    var visibleIndices = []
    for (i = 0 ; i < inputVisibleIndices.length ; i++) {
        if (logs[inputVisibleIndices[i]].is_lager) {
            visibleIndices.push(inputVisibleIndices[i])
        }
    }

    if (visibleIndices.length < 2) {
        return
    }

    var counts = this.binUp(logs, visibleIndices)
    this.renderBins(counts, "filter")
  },
})
