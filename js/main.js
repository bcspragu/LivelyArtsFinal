var fill = d3.scale.category20();
var width, height, conn;

$(function() {
  if (window["WebSocket"]) {
    conn = new WebSocket("ws://localhost:3000/ws");
    conn.onmessage = function(evt) {
      $("body").empty();
      render(JSON.parse(evt.data));
    }
  } 
  render(wordList);
});

function render(list) {
  width = $("body").width();
  height = $("body").height();
  d3.layout.cloud().size([width, height])
      .words(list)
      .padding(5)
      .rotate(function() { return ~~(Math.random() * 2) * 90; })
      .font("Impact")
      .fontSize(function(d) { return d.size; })
      .on("end", draw)
      .start();
}

function draw(words) {
  d3.select("body").append("svg")
      .attr("width", width)
      .attr("height", height)
    .append("g")
      .attr("transform", "translate(" + width/2 + "," + height/2 + ")")
    .selectAll("text")
      .data(words)
    .enter().append("text")
      .style("font-size", function(d) { return d.size + "px"; })
      .style("font-family", "Impact")
      .style("fill", function(d, i) { return fill(i); })
      .attr("text-anchor", "middle")
      .attr("transform", function(d) {
        return "translate(" + [d.x, d.y] + ")rotate(" + d.rotate + ")";
      })
      .text(function(d) { return d.text; });
}
