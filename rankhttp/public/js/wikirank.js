"use strict";

$(document).ready(function() {
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
  });
});
