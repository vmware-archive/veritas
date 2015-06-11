var LogRenderer = function() {
}

var ansiColorCode = /\[(\d+)m/g

LogRenderer.prototype = {
  prerenderLogs: function(logs) {
    for (var i = 0 ; i < logs.length ; i++) {
      this.prerenderLog(i, logs[i])
    }
  },

  prerenderLog: function(index, log) {
    var cssClass = log.is_lager ? "lager-log source-" + log.log.source.replace(/\./g, "-") + " level-" + log.log.level : "raw-log"
    var dom = this.renderTimestamp(log)
    if (log.is_lager) {
      if (log.log.source == "ginkgo") {
        if (log.log.message.match(/spec\.start/)) {
          cssClass += " ginkgo-start"
        } else {
          cssClass += " ginkgo-end"
        }
        dom += this.renderGinkgo(log)
      } else {
        dom += this.renderLager(log)
      }
      log.searchText = this.lagerSearchText(log)
    } else {
      dom += this.renderRaw(log)
      log.searchText = this.rawSearchText(log)
    }
    log.dom = '<div id="log-' + index + '" class="' + cssClass + '">' + dom + '</div>'
  },

  renderTimestamp: function(log) {
    var absoluteTimestamp, relativeTimestamp
    if (log.is_lager) {
      if (!this.firstTimestamp) {
        this.firstTimestamp = log.log.timestamp
      }
      this.mostRecentTimestamp = log.log.timestamp
      absoluteTimestamp = '<div class="absolute-timestamp">' + formatUnixTimestamp(log.log.timestamp) + '</div>'
      relativeTimestamp = '<div class="relative-timestamp">' + formatRelativeTimestamp(log.log.timestamp - this.firstTimestamp) + '</div>'
      log.timestamp = log.log.timestamp
    } else {
      if (!this.mostRecentTimestamp) {
        absoluteTimestamp = '<div class="absolute-timestamp unknown">???</div>'
        relativeTimestamp = '<div class="relative-timestamp unknown">???</div>'
        log.timestamp = undefined
      } else {
        absoluteTimestamp = '<div class="absolute-timestamp">' + formatUnixTimestamp(this.mostRecentTimestamp) + '</div>'
        relativeTimestamp = '<div class="relative-timestamp">' + formatRelativeTimestamp(this.mostRecentTimestamp - this.firstTimestamp) + '</div>'
        log.timestamp = this.mostRecentTimestamp
      }
    }

    return absoluteTimestamp + relativeTimestamp
  },

  renderLager: function(log) {
    var sourceDom = '<div class="source">' + log.log.source + '</div>'
    var logLevelDom = '<div class="level">[' + log.log.level + ']</div>'
    var sessionDom = '<div class="session">' + log.log.session + '</div>'
    var messageDom = '<div class="message">' + log.log.message + '</div>'
    var errorDom = ""
    var traceDom = ""
    var dataDom = ""
    if (log.log.error) {
      errorDom = '<div class="error">' + log.log.error + '</div>'
    }
    if (log.log.trace) {
      traceDom = '<div class="trace">' + log.log.trace + '</div>'
    }
    if (log.log.data) {
      dataDom = '<div class="data">' + this.trimData(log.log.data) + '</div>'
    }

    return sourceDom + logLevelDom + sessionDom + '<div class="content">' + messageDom + errorDom + traceDom + dataDom + '</div>'
  },

  renderGinkgo: function(log) {
    var sessionDom = '<div class="session">' + log.log.session + '</div>'
    var name = ""
    var nameComponents = log.log.data.summary.name
    for (var i = 1 ; i < nameComponents.length ; i++) {
      var klass = (i == nameComponents.length - 1 ? "it" : (i%2 ? "odd" : "even"))
      name += '<span class="ginkgo-name-' + klass + '">' + nameComponents[i] + ' </span>'
    }
    var messageDom = '<div class="message">' + log.log.message + ": " + name + '</div>'
    var errorDom = ""
    var dataDom = ""
    if (log.log.error) {
      errorDom = '<div class="error"><pre class="ginkgo-error">' + log.log.error + '</pre></div>'
    }

    dataDom = '<div class="data">'
    dataDom += '<div>'+log.log.data.summary.location+'</div>'
    if (log.log.data.summary.run_time) {
      dataDom += '<div>'+formatRelativeTimestamp(log.log.data.summary.run_time)+'</div>'
    }
    dataDom += '</div>'

    return sessionDom + '<div class="content">' + messageDom + errorDom + dataDom + '</div>'
  },

  renderRaw: function(log) {
    return '<div class="message">' + log.raw.replace(ansiColorCode, "") + '</div>'
  },

  trimData: function(data) {
      var keys = _.keys(data).sort()
      var dataDom = '<dl class="big-data">'
      _.each(keys, function(key) {
        dataDom += "<dt>"+key+":</dt>"
        var prettyJSON = JSON.stringify(data[key], null, "  ").replace(/"/g, '')
        dataDom += "<dd><pre>"+prettyJSON+"</pre></dd>"
      })
      dataDom += "</dl>"

      shortData = JSON.stringify(data)
      shortData = shortData.slice(1,-1).replace(/"/g, '').replace(/,/g, ", ")
      return '<div class="small-data">' + shortData + '</div>' + dataDom
  },

  lagerSearchText: function(log) {
    var searchText = "source:"+log.log.source
    searchText += " session:"+log.log.session
    searchText += " message:"+log.log.message
    if (log.log.error) {
      searchText += " error:"+log.log.error
    }
    if (log.log.trace) {
      searchText += " trace:"+log.log.trace
    }
    if (log.log.data) {
      searchText += " data:"+this.trimData(log.log.data)
    }

    return searchText
  },

  rawSearchText: function(log) {
    return "raw:" + log.raw.replace(ansiColorCode, "")
  },
}

function formatUnixTimestamp(timestamp) {
  var date = new Date(timestamp/1e6)
  var month = date.getMonth() + 1
  var day = date.getDate()
  var hours = date.getHours()
  var minutes = date.getMinutes()
  var seconds = date.getSeconds()
  var milliseconds = date.getMilliseconds()

  month = month < 10 ? "0" + month : month
  day = day < 10 ? "0" + day : day
  hours = hours < 10 ? "0" + hours : hours
  minutes = minutes < 10 ? "0" + minutes : minutes
  seconds = seconds < 10 ? "0" + seconds : seconds
  if (milliseconds < 10) {
    milliseconds = "00" + milliseconds
  } else if (milliseconds < 100) {
    milliseconds = "0" + milliseconds
  }

  return month + "/" + day + " " + hours + ":" + minutes + ":" + seconds + "." + milliseconds
}

function formatRelativeTimestamp(nanoseconds) {
  var days = Math.floor(nanoseconds/8.64e13)
  nanoseconds = nanoseconds - days * 8.64e13
  var hours = Math.floor(nanoseconds/3.6e12)
  nanoseconds = nanoseconds - hours * 3.6e12
  var minutes = Math.floor(nanoseconds/6e10)
  nanoseconds = nanoseconds - minutes * 6e10
  var seconds = Math.floor(nanoseconds/1e9)
  nanoseconds = nanoseconds - seconds*1e9
  var milliseconds = Math.floor(nanoseconds/1e6)

  var relativeTimestamp = ""
  if (days > 0) {
    relativeTimestamp += days + "d"
  }

  if (hours > 0) {
    relativeTimestamp += hours + "h"
  }

  if (minutes > 0) {
    relativeTimestamp += minutes + "m"
  }

  relativeTimestamp += seconds + "." + milliseconds + "s"

  return relativeTimestamp
}
