"use strict";

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
    $.getJSON("/compare", {
      "things": things
    })
      .done(function(data) {
        endSpin();

        function makeInfluenceResults(page) {
          return "<strong>" + page.Page.Title + "</strong> is most influenced by " + _.map(page.Influencers, function(i) {
            return i.Page.Title;
          }).join(",");
        }

        var pages = data.pages;
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

        $(".influencer.influencerFirst").html("<strong>" + page1.Page.Title + "</strong> is");
        $(".influenceResultsRatio").text(influenceText);
        $(".influencer.influencerSecond").html(copy + " <strong>" + page2.Page.Title + "</strong>");
        $(".influence1Stats").html(makeInfluenceResults(page1));
        $(".influence2Stats").html(makeInfluenceResults(page2));
        $(".influenceResults").show();

      })
      .fail(function(xhr, textStatus, errorThrown) {
        endSpin();
      });
  });
});
