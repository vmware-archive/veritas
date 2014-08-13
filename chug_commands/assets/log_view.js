var LogView = Backbone.View.extend({
  events: {
    'scroll': 'scroll',
  },

  renderLogs: function(logs) {
    var dom = ""
    _.each(logs, function(log) {
      dom += log.dom
    })
    this.$el[0].innerHTML = dom
    this.count = logs.length
  },

  clearFilter: function() {
    this.$(".show").removeClass("show")        
  },

  filterLogs: function(visibleIndices) {
    for (i = 0 ; i < visibleIndices.length ; i++) {
        this.$("#log-"+visibleIndices[i]).addClass("show")
    }
  },

  scroll: function() {
    var top = this.$el.scrollTop()
    var bottom = top + this.$el.height()
  },
})
