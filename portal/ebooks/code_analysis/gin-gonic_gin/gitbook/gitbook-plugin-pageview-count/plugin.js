var gitbook = window.gitbook;

/*
<!-- Start of CuterCounter Code -->
<a href="http://www.cutercounter.com/" target="_blank"><img src="http://www.cutercounter.com/hit.php?id=gmvufxqck&nd=1&style=116" border="0" alt="visitor counter"></a>
<!-- End of CuterCounter Code -->
*/


function requestCount(targetUrl) {
  return $.ajax('https://hitcounter.pythonanywhere.com/count',{
    data:{url:encodeURIComponent(targetUrl)},
  })
}

require(["gitbook", "jQuery"], function (gitbook, $) {

  function resetViewCount() {
    var bookHeader = $('.book-header')
    var lastChild = bookHeader.children().last()

    var renderWrapper = $('<div class="page-view-wrapper dropdown pull-left">\
        <span class="page-view-counter btn" title="阅读数">-</span>\
      </div>')

    if(lastChild.length){
      renderWrapper.insertBefore(lastChild)
    }else{
      bookHeader.append(renderWrapper)
    }
    

    requestCount(location.href).then(function(data){
      renderWrapper.find('.page-view-counter').html(data)
    })
  }

  gitbook.events.bind("page.change", resetViewCount)
});
