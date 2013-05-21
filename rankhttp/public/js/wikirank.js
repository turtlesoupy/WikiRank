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
        console.log(parsedResponse.suggestions);
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

        var ratio = data.pages[0].Rank / data.pages[1].Rank;
        var influenceText, copy;
        console.log(ratio);

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
          influenceText = textRatio + (ratio >= 1.0 ?  "x more" : "x less");
          copy = "influential than";
        }

        $(".influencer.influencerFirst").html("<strong>" + data.pages[0].Title + "</strong> is");
        $(".influenceResultsRatio").text(influenceText);
        $(".influencer.influencerSecond").html(copy + " <strong>" + data.pages[1].Title + "</strong>");
        $(".influenceResults").show();

      })
      .fail(function(xhr, textStatus, errorThrown) {
        endSpin();
      });
  });
});
