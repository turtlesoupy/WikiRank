"use strict";

var tableTemplate = _.template(
    "<table>" + 
      "<thead>" + 
       "<tr><th>Rank</th><th>Name</th><th>Relative</th></tr>" + 
      "</thead>" + 
      "<% _.each(movies, function(u) { %>" +
        "<tr>" + 
          "<td> <%= u.rank %> </td>" + 
          "<td> <a href='<%=wikilink(u.title)%>'> <%= removeBrackets(u.title) %> </td> </a>" + 
          "<td> <%= new Number(u.relativeRank * 100).toPrecision(3) %>% </td>" + 
        "</tr>" + 
      "<% }); %>" +
    "</table>" 
);

function renderTable(movies, selectedUniversity) {
  $("#movie-table-wrapper").html(tableTemplate({
    movies: movies,
    wikilink: wikilink,
    removeBrackets: removeBrackets
  }));
}


function wikilink(t) {
  return "http://en.wikipedia.org/wiki/" + _.escape(t);
}

function removeBrackets(t) {
  return t.replace(/ *\([^)]*\)/g, "");
}

$(document).ready(function() {
  var byYear = _.groupBy(MOVIES, 'releaseYear');
  var years = _.keys(byYear);
  years.sort();
  years.reverse()
  years = _.rest(_.first(years, 101));
  var massByYear = _.map(years, function(year) {
    var bestMovie = _.max(byYear[year], function(m) { return m.pageRank; });
    return {
      year: year,
      bestMovie: _.max(byYear[year], function(m) { return m.pageRank; }),
      mass: bestMovie.pageRank,
      //_.reduce(byYear[year], function(s, movie) { return s + movie.pageRank; }, 0)
    };
  });
  var i;
  for(i = 0; i < MOVIES.length; i++) {
    MOVIES[i].rank = i+1;
    MOVIES[i].relativeRank = MOVIES[i].pageRank / MOVIES[0].pageRank;
  }

  renderTable(_.first(MOVIES, 25));

  var height = 2800;
  var width = 700;
  var labelPadding = 80;
  var lineOffset = 5;
  var chart = d3.select("#movie-histogram-svg")
     .attr("class", "chart")
     .attr("width", width)
     .attr("height", height)
     .style("fill", "steelblue")
   .append("g")
    .attr("transform", "translate(" + labelPadding + ", 5)");

  var x = d3.scale.linear()
         .domain([0, d3.max(massByYear, function(d) { return d.mass; })])
         .range([0, width - labelPadding - 30]);

  var y = d3.scale.ordinal()
    .domain(years)
    .rangeBands([0, height]);


  var displayYears = _.filter(years, function(year) { return year % 4 === (years[0] % 4);});
  var yearOrdinal = d3.scale.ordinal()
    .domain(displayYears)
    .rangeBands([0, height]);

  var bars = chart.append("g").attr("transform", "translate(" + lineOffset + ")");

  var dataReachDelay = function(d, i) {
    return i * 30;
  };

  bars.selectAll("zebra")
      .data(massByYear)
    .enter().append("rect")
      .attr("y", function(d) { return y(d.mass); })
      .attr("height", y.rangeBand())
      .attr("fill", function(d,i) { 
        return i % 2 == 0 ? "#fcfcfc" : "none";
      }).attr("width", width - labelPadding);

  bars.selectAll("bars")
      .data(massByYear)
    .enter().append("rect")
      .attr("y", function(d) { return y(d.mass); })
      .attr("height", y.rangeBand())
      .attr("width", 0)
      .attr("stroke", "white")
      .attr("stroke-width", "1")
      .style("cursor", "pointer")
    .transition()
      .duration(1000)
      .delay(function(d,i) { return i * 30;})
      .attr("width", function(d) { return x(d.mass); })
      ;


  chart.selectAll("movieText")
      .data(massByYear)
    .enter().append("a")
      .attr("xlink:href", function(d) { return wikilink(d.bestMovie.title); })
      .attr("target", "_blank")
    .append("text")
      .attr("x", width - labelPadding)
      .attr("y", function(d) { 
        console.log(d.year);
        return y(d.year) + y.rangeBand() / 1.5;
      })
      .attr("font-size", 12)
      .attr("fill", "#0000A0")
      .attr("text-anchor", "end")
      .attr("opacity", 0)
      .text(function(d) {
        return removeBrackets(d.bestMovie.title);
      })
    .transition()
      .duration(1000)
      .attr("opacity", 1)
      .delay(dataReachDelay);

  var line = d3.svg.line()
    .x(function(d) { return d[0]; })
    .y(function(d) { return d[1]; })
    .interpolate("basis");

  var path = chart.append("path")
    .attr("d", line([[0,0], [0,height]]))
    .attr("stroke", "black")
    .attr("stroke-width", 2)

  var totalLength = path.node().getTotalLength();

  path
    .attr("stroke-dasharray", totalLength + " " + totalLength)
    .attr("stroke-dashoffset", totalLength)
    .transition()
      .duration(massByYear.length * 30)
      .ease("linear")
      .attr("stroke-dashoffset", 0);

  chart.selectAll("yearText")
      .data(displayYears)
    .enter().append("text")
      .attr("x", -50)
      .attr("y", function(d) { return yearOrdinal(d) + y.rangeBand() / 5;})
      .attr("dx", -3)
      .attr("font-size", 14)
      .attr("fill", "black")
      .style("font-weight", "bold")
      .attr("opacity", 0)
      .text(String)
    .transition()
      .duration(100)
      .delay(function(d,i) { return dataReachDelay(d,i) * 4; })
      .attr("opacity", 1);


  chart.selectAll("ticks2")
      .data(displayYears)
    .enter()
      .append("path")
      .attr("d", function(d) {
        return line([[-10,yearOrdinal(d) + 1], [0,yearOrdinal(d) + 1]]);
      }).attr("stroke", "black")
        .attr("stroke-width", 2)
        .attr("opacity", 0)
    .transition()
      .duration(100)
      .delay(function(d,i) { return dataReachDelay(d,i) * 4; })
      .attr("opacity", 1);


});
