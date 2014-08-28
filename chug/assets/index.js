  (function() {
  var js = [];
  js.push('/assets/jquery.js');
  js.push('/assets/underscore.js');
  js.push('/assets/backbone.js');

  js.push('/assets/renderer.js');
  js.push('/assets/log_view.js');
  js.push('/assets/histogram_view.js');
  js.push('/assets/app.js');

  head.js.apply(head, js);
})()

head.ready(function() {
  window.app = new App({
    el: $("#app")
  })
});