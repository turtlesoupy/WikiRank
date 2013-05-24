"use strict";

var viewHelpers = {
  linkWiki: function(title) {
    return "<a class='wikilink' href='http://en.wikipedia.org/wiki/" + _.escape(title) + "'>" + _.escape(title) +  "</a>";
  },
  percentInfluence: function(influence, totalRank) {
    return (influence / totalRank).toPrecision(2) + "%";
  }
};

var resultStringTemplate = _.template(
    "<div><strong><%= linkWiki(firstTitle) %></strong> is</div>" +
    "<div class='ratio'><%= ratioString %></div>" +
    "<div><%= influenceCopy %> <strong><%= linkWiki(secondTitle) %></strong></div>"
);

var resultTableTemplate = _.template(
    "<h4><%=page.Title%></h4>" +
    "<p class='muted'>Ranked <%=page.Order%> overall</p>" +
    "<table class='table table-striped'>" +
    "<thead>" +
      "<tr>" +
      "<th>Influencer</th><th>Influence</th>" +
      "</tr>" +
    "</thead>" +
    "<tbody>" +
      "<% _.each(influencers, function(influencer) { %>" + 
        "<tr>" +
          "<td><%=console.log(influencer)%><%=linkWiki(influencer.Page.Title)%></td><td><%=percentInfluence(influencer.Influence, page.Rank)%></td>" +
        "</tr>" +
      "<% }); %>" +
    "</tbody>" +
    "</table>"
);

$(document).ready(function() {
  var spinner = new Spinner({
    lines: 15,
    length: 0,
    width: 12,
    radius: 30
  });

  function startSpin() {  
    spinner.spin($(".spinner").get(0)); 
    $(".influenceLoader").show();
  }
  function endSpin() {  
    spinner.stop(); 
    $(".influenceLoader").hide();
  }

  $("input.first-object, input.second-object").typeahead({
    name: "objects",
    remote: {
      url: "/named_entity_suggestions?q=%QUERY",
      filter: function(parsedResponse) {
        var suggestions = _.pluck(parsedResponse.suggestions, "Title")
        return suggestions;
      }
    }
  }).bind("typeahead:selected", function(e) {
    var otherSelector = "input.second-object";
    if($(this).hasClass("second-object")) {
      otherSelector = "input.first-object";
    }
    
    if($(otherSelector).val() === "") {
      return; 
    }

    var things = [
      $(otherSelector).val(),
      $(this).val()
    ];

    startSpin();
    $(".influenceResults").hide();
    $.getJSON("/things", {
      "things": things
    })
      .done(function(data) {
        endSpin();

        function makeInfluenceResults(page) {
          return "<strong>" + page.Page.Title + "</strong> is most influenced by " + _.map(page.Influencers, function(i) {
            return i.Page.Title;
          }).join(",");
        }

        var pages = data;
        console.log(pages);
        var page1 = pages[0];
        var page2 = pages[1];
        var ratio = page1.Page.Rank / page2.Page.Rank;
        if(ratio < 1) {
          ratio = 1/ratio;
          var tmp = page1;
          page1 = page2;
          page2 = tmp;
        }

        var influenceText, copy;

        if(ratio >= 0.95 && ratio <= 1.05) {
          influenceText = "about as";
          copy = "influential as"
        } else {
          var textRatio;
          if(ratio < 10) {
            textRatio = Number(ratio.toPrecision(2));
          } else {
            textRatio = Number(ratio.toPrecision(1));
          }
          influenceText = textRatio + "x more";
          copy = "influential than";
        }

        $(".machine-results .result-string").html(resultStringTemplate(_.extend({
          firstTitle: page1.Page.Title,
          secondTitle: page2.Page.Title,
          ratioString: influenceText, 
          influenceCopy: copy
        }, viewHelpers)));

        _.each(pages, function(page, i) {
          var n = ".machine-results .result-table" + (i == 0 ? "-first" : "-second");
          console.log(n);
          $(".machine-results .result-table" + (i == 0 ? "-first" : "-second")).html(resultTableTemplate(_.extend({
            page: page.Page,
            influencers: page.Influencers
          }, viewHelpers)));
        });


        $(".machine-results").show();

      })
      .fail(function(xhr, textStatus, errorThrown) {
        endSpin();
      });
  });
});
