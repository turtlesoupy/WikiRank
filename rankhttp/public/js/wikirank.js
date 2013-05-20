$(document).ready(function() {
  $("input.first-object, input.second-object").typeahead({
    name: "objects",
    remote: {
      url: "/named_entity_suggestions/q?%QUERY",
      filter: function(parsedResponse) {
        return parsedResponse.suggestions;
      }
    }
  });
});
