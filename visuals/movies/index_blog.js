"use strict";

var MASS_BY_YEAR = [{"year":"2012","bestMovie":{"pageRank":0.0005255874454256719,"position":67,"releaseYear":2012,"title":"The Dark Knight Rises","rank":67,"relativeRank":0.28607952190887476},"mass":0.013904217931578217},{"year":"2011","bestMovie":{"pageRank":0.0003507907580171381,"position":222,"releaseYear":2011,"title":"The Artist (film)","rank":222,"relativeRank":0.19093692822574587},"mass":0.018182358464402356},{"year":"2010","bestMovie":{"pageRank":0.0006952509926623833,"position":28,"releaseYear":2010,"title":"The Social Network","rank":28,"relativeRank":0.3784281251741829},"mass":0.019649135137841212},{"year":"2009","bestMovie":{"pageRank":0.0009080078686210272,"position":9,"releaseYear":2009,"title":"Avatar (2009 film)","rank":9,"relativeRank":0.49423261382169975},"mass":0.023668022223914696},{"year":"2008","bestMovie":{"pageRank":0.0007753077035025826,"position":21,"releaseYear":2008,"title":"The Dark Knight (film)","rank":21,"relativeRank":0.4220033394645709},"mass":0.022163334399995925},{"year":"2007","bestMovie":{"pageRank":0.0005100625364260203,"position":71,"releaseYear":2007,"title":"Juno (film)","rank":71,"relativeRank":0.27762924673020856},"mass":0.02286003246513391},{"year":"2006","bestMovie":{"pageRank":0.000552892243606405,"position":57,"releaseYear":2006,"title":"Casino Royale (2006 film)","rank":57,"relativeRank":0.30094164176608723},"mass":0.026007114653911197},{"year":"2005","bestMovie":{"pageRank":0.0006563947090344212,"position":34,"releaseYear":2005,"title":"Brokeback Mountain","rank":34,"relativeRank":0.3572784817795614},"mass":0.024442253613779254},{"year":"2004","bestMovie":{"pageRank":0.0008534922160955119,"position":13,"releaseYear":2004,"title":"The Passion of the Christ","rank":13,"relativeRank":0.4645595081438829},"mass":0.026060081042908186},{"year":"2003","bestMovie":{"pageRank":0.0007814081757842581,"position":20,"releaseYear":2003,"title":"The Lord of the Rings: The Return of the King","rank":20,"relativeRank":0.4253238529375413},"mass":0.01947601494371774},{"year":"2002","bestMovie":{"pageRank":0.0006008783543018752,"position":46,"releaseYear":2002,"title":"Minority Report (film)","rank":46,"relativeRank":0.32706068955823586},"mass":0.021880354165087913},{"year":"2001","bestMovie":{"pageRank":0.0005055990466165774,"position":73,"releaseYear":2001,"title":"Shrek","rank":73,"relativeRank":0.2751997499036693},"mass":0.01975538602744521},{"year":"2000","bestMovie":{"pageRank":0.0006148803046036405,"position":42,"releaseYear":2000,"title":"Gladiator (2000 film)","rank":42,"relativeRank":0.33468201172447726},"mass":0.015563069118688719},{"year":"1999","bestMovie":{"pageRank":0.000883778551204347,"position":10,"releaseYear":1999,"title":"The Matrix","rank":10,"relativeRank":0.48104449145867717},"mass":0.015402604330087404},{"year":"1998","bestMovie":{"pageRank":0.0006057609849983526,"position":45,"releaseYear":1998,"title":"Saving Private Ryan","rank":45,"relativeRank":0.3297183265841252},"mass":0.014032893620401775},{"year":"1997","bestMovie":{"pageRank":0.0014246323597266529,"position":4,"releaseYear":1997,"title":"Titanic (1997 film)","rank":4,"relativeRank":0.7754335608919133},"mass":0.014898646333549762},{"year":"1996","bestMovie":{"pageRank":0.00040118229840597656,"position":147,"releaseYear":1996,"title":"Independence Day (1996 film)","rank":147,"relativeRank":0.2183652618135377},"mass":0.010991465614386771},{"year":"1995","bestMovie":{"pageRank":0.0005756117366529739,"position":50,"releaseYear":1995,"title":"Braveheart","rank":50,"relativeRank":0.3133079601881539},"mass":0.012445604756259136},{"year":"1994","bestMovie":{"pageRank":0.000797798570199574,"position":18,"releaseYear":1994,"title":"The Lion King","rank":18,"relativeRank":0.43424521557479745},"mass":0.011853718180143974},{"year":"1993","bestMovie":{"pageRank":0.0007156766532785275,"position":24,"releaseYear":1993,"title":"Schindler's List","rank":24,"relativeRank":0.38954590067395145},"mass":0.010889305733265454},{"year":"1992","bestMovie":{"pageRank":0.000548595938206439,"position":60,"releaseYear":1992,"title":"Aladdin (1992 Disney film)","rank":60,"relativeRank":0.2986031441373979},"mass":0.009292220913195445},{"year":"1991","bestMovie":{"pageRank":0.0007221620281487321,"position":23,"releaseYear":1991,"title":"The Silence of the Lambs (film)","rank":23,"relativeRank":0.39307591829217164},"mass":0.008161458429621194},{"year":"1990","bestMovie":{"pageRank":0.0005028139012810725,"position":74,"releaseYear":1990,"title":"Total Recall (1990 film)","rank":74,"relativeRank":0.2736837832401531},"mass":0.007997331544167539},{"year":"1989","bestMovie":{"pageRank":0.0004553411424474579,"position":107,"releaseYear":1989,"title":"The Little Mermaid (1989 film)","rank":107,"relativeRank":0.247844155088806},"mass":0.009191204424699824},{"year":"1988","bestMovie":{"pageRank":0.0006845553311083534,"position":31,"releaseYear":1988,"title":"Who Framed Roger Rabbit","rank":31,"relativeRank":0.3726064302868595},"mass":0.009290006868327882},{"year":"1987","bestMovie":{"pageRank":0.00036819055891433236,"position":198,"releaseYear":1987,"title":"RoboCop","rank":198,"relativeRank":0.20040771518099268},"mass":0.008033587415945337},{"year":"1986","bestMovie":{"pageRank":0.0004362172495448376,"position":121,"releaseYear":1986,"title":"Top Gun","rank":121,"relativeRank":0.23743493739109783},"mass":0.007844036716226946},{"year":"1985","bestMovie":{"pageRank":0.0003619966548128734,"position":205,"releaseYear":1985,"title":"Brazil (1985 film)","rank":205,"relativeRank":0.19703634636403072},"mass":0.006642498541828006},{"year":"1984","bestMovie":{"pageRank":0.0006243224615333268,"position":40,"releaseYear":1984,"title":"Ghostbusters","rank":40,"relativeRank":0.3398214186181208},"mass":0.008997364790944646},{"year":"1983","bestMovie":{"pageRank":0.0008453476862268645,"position":15,"releaseYear":1983,"title":"Star Wars Episode VI: Return of the Jedi","rank":15,"relativeRank":0.4601264052772265},"mass":0.006512902530378996},{"year":"1982","bestMovie":{"pageRank":0.0010437494988207443,"position":7,"releaseYear":1982,"title":"Blade Runner","rank":7,"relativeRank":0.5681173707896211},"mass":0.008740797677365974},{"year":"1981","bestMovie":{"pageRank":0.0007878317970523471,"position":19,"releaseYear":1981,"title":"Raiders of the Lost Ark","rank":19,"relativeRank":0.42882025780278743},"mass":0.007335964219682221},{"year":"1980","bestMovie":{"pageRank":0.0009547553008994389,"position":8,"releaseYear":1980,"title":"Star Wars Episode V: The Empire Strikes Back","rank":8,"relativeRank":0.5196774435889793},"mass":0.007275994604332796},{"year":"1979","bestMovie":{"pageRank":0.0007153813312577393,"position":25,"releaseYear":1979,"title":"Apocalypse Now","rank":25,"relativeRank":0.3893851556195337},"mass":0.006536091793389534},{"year":"1978","bestMovie":{"pageRank":0.0005634515752904434,"position":54,"releaseYear":1978,"title":"Superman (film)","rank":54,"relativeRank":0.3066891317149774},"mass":0.006131089802198406},{"year":"1977","bestMovie":{"pageRank":0.0014119647235574068,"position":5,"releaseYear":1977,"title":"Star Wars Episode IV: A New Hope","rank":5,"relativeRank":0.7685385116844908},"mass":0.006189458909270794},{"year":"1976","bestMovie":{"pageRank":0.0006065498460814711,"position":43,"releaseYear":1976,"title":"Rocky","rank":43,"relativeRank":0.33014770708678975},"mass":0.0052879172496973995},{"year":"1975","bestMovie":{"pageRank":0.0006443531573301719,"position":35,"releaseYear":1975,"title":"Jaws (film)","rank":35,"relativeRank":0.35072421305687024},"mass":0.0062938336401546485},{"year":"1974","bestMovie":{"pageRank":0.0005547365703649048,"position":56,"releaseYear":1974,"title":"The Godfather Part II","rank":56,"relativeRank":0.30194551680516474},"mass":0.0048612497764804535},{"year":"1973","bestMovie":{"pageRank":0.0008603318616304854,"position":12,"releaseYear":1973,"title":"The Exorcist (film)","rank":12,"relativeRank":0.468282356818639},"mass":0.006401932554720887},{"year":"1972","bestMovie":{"pageRank":0.0012199286521431154,"position":6,"releaseYear":1972,"title":"The Godfather","rank":6,"relativeRank":0.664012446654598},"mass":0.0052852674129822445},{"year":"1971","bestMovie":{"pageRank":0.0004894004579099235,"position":84,"releaseYear":1971,"title":"A Clockwork Orange (film)","rank":84,"relativeRank":0.2663827879439997},"mass":0.004746664549275269},{"year":"1970","bestMovie":{"pageRank":0.00037208350688216666,"position":188,"releaseYear":1970,"title":"Patton (film)","rank":188,"relativeRank":0.20252666361316496},"mass":0.004830752725152791},{"year":"1969","bestMovie":{"pageRank":0.00039314097318744024,"position":164,"releaseYear":1969,"title":"Midnight Cowboy","rank":164,"relativeRank":0.2139883336847284},"mass":0.0049756247410463516},{"year":"1968","bestMovie":{"pageRank":0.0008021518652251407,"position":17,"releaseYear":1968,"title":"2001: A Space Odyssey (film)","rank":17,"relativeRank":0.4366147329034197},"mass":0.006470615760690326},{"year":"1967","bestMovie":{"pageRank":0.00043152740406621226,"position":124,"releaseYear":1967,"title":"The Graduate","rank":124,"relativeRank":0.23488223419388765},"mass":0.004618446876305415},{"year":"1966","bestMovie":{"pageRank":0.0005968002123124114,"position":47,"releaseYear":1966,"title":"The Good, the Bad and the Ugly","rank":47,"relativeRank":0.32484093921835205},"mass":0.005264437091936303},{"year":"1965","bestMovie":{"pageRank":0.0005586351505105817,"position":55,"releaseYear":1965,"title":"The Sound of Music (film)","rank":55,"relativeRank":0.3040675308561195},"mass":0.005296192010460844},{"year":"1964","bestMovie":{"pageRank":0.0006642781551807858,"position":33,"releaseYear":1964,"title":"Mary Poppins (film)","rank":33,"relativeRank":0.3615694756459004},"mass":0.005295640055640922},{"year":"1963","bestMovie":{"pageRank":0.00048175619009850146,"position":87,"releaseYear":1963,"title":"Cleopatra (1963 film)","rank":87,"relativeRank":0.26222197988081647},"mass":0.0044444785489663585},{"year":"1962","bestMovie":{"pageRank":0.0006681451238935123,"position":32,"releaseYear":1962,"title":"Lawrence of Arabia (film)","rank":32,"relativeRank":0.3636742834570486},"mass":0.003418263835479012},{"year":"1961","bestMovie":{"pageRank":0.00037761753345157356,"position":179,"releaseYear":1961,"title":"West Side Story (film)","rank":179,"relativeRank":0.20553885823270107},"mass":0.0033682651826264886},{"year":"1960","bestMovie":{"pageRank":0.0006216193278309128,"position":41,"releaseYear":1960,"title":"Psycho (1960 film)","rank":41,"relativeRank":0.33835009124153925},"mass":0.004707510562333933},{"year":"1959","bestMovie":{"pageRank":0.0007009758231696897,"position":27,"releaseYear":1959,"title":"Ben-Hur (1959 film)","rank":27,"relativeRank":0.38154417520314265},"mass":0.0043293547927564955},{"year":"1958","bestMovie":{"pageRank":0.0003999947325593788,"position":148,"releaseYear":1958,"title":"South Pacific (1958 film)","rank":148,"relativeRank":0.2177188645820459},"mass":0.0036694870168385595},{"year":"1957","bestMovie":{"pageRank":0.0004395911689498288,"position":117,"releaseYear":1957,"title":"The Bridge on the River Kwai","rank":117,"relativeRank":0.23927137632954554},"mass":0.003314132967599763},{"year":"1956","bestMovie":{"pageRank":0.0003685922467285267,"position":197,"releaseYear":1956,"title":"The Ten Commandments (1956 film)","rank":197,"relativeRank":0.2006263555972383},"mass":0.004224431113756356},{"year":"1955","bestMovie":{"pageRank":0.0003024256217882789,"position":302,"releaseYear":1955,"title":"The Ladykillers","rank":302,"relativeRank":0.16461157519490308},"mass":0.0041240026584372595},{"year":"1954","bestMovie":{"pageRank":0.0004853626093666852,"position":85,"releaseYear":1954,"title":"On the Waterfront","rank":85,"relativeRank":0.2641849695013341},"mass":0.0037356320202265154},{"year":"1953","bestMovie":{"pageRank":0.00046888668598982333,"position":100,"releaseYear":1953,"title":"Peter Pan (1953 film)","rank":100,"relativeRank":0.2552170530800381},"mass":0.003952523504653269},{"year":"1952","bestMovie":{"pageRank":0.000428333269120733,"position":128,"releaseYear":1952,"title":"Singin' in the Rain","rank":128,"relativeRank":0.23314365271507195},"mass":0.002751240621098053},{"year":"1951","bestMovie":{"pageRank":0.00037067717656260907,"position":192,"releaseYear":1951,"title":"Alice in Wonderland (1951 film)","rank":192,"relativeRank":0.20176119193196992},"mass":0.0029243475647161804},{"year":"1950","bestMovie":{"pageRank":0.0004300311147321812,"position":126,"releaseYear":1950,"title":"All About Eve","rank":126,"relativeRank":0.23406779743166578},"mass":0.002806554569439382},{"year":"1949","bestMovie":{"pageRank":0.0003983198163676058,"position":150,"releaseYear":1949,"title":"The Third Man","rank":150,"relativeRank":0.2168072004478469},"mass":0.002778923696962474},{"year":"1948","bestMovie":{"pageRank":0.0003872853782032674,"position":171,"releaseYear":1948,"title":"Bicycle Thieves","rank":171,"relativeRank":0.210801107984907},"mass":0.0030481155785741765},{"year":"1947","bestMovie":{"pageRank":0.00030751635366448857,"position":290,"releaseYear":1947,"title":"Miracle on 34th Street","rank":290,"relativeRank":0.16738248259382854},"mass":0.0017603673185133473},{"year":"1946","bestMovie":{"pageRank":0.0005502386641643766,"position":59,"releaseYear":1946,"title":"It's a Wonderful Life","rank":59,"relativeRank":0.2994972869879629},"mass":0.0029219053811666146},{"year":"1945","bestMovie":{"pageRank":0.00016625577911073345,"position":976,"releaseYear":1945,"title":"The Lost Weekend (film)","rank":976,"relativeRank":0.09049374032149012},"mass":0.0023655028674150025},{"year":"1944","bestMovie":{"pageRank":0.00024113721975315352,"position":494,"releaseYear":1944,"title":"Meet Me in St. Louis","rank":494,"relativeRank":0.13125203263854052},"mass":0.002780075365327741},{"year":"1943","bestMovie":{"pageRank":0.00021351632950448952,"position":606,"releaseYear":1943,"title":"The Song of Bernadette (film)","rank":606,"relativeRank":0.11621786250033322},"mass":0.0031726581963694533},{"year":"1942","bestMovie":{"pageRank":0.0008717381571018678,"position":11,"releaseYear":1942,"title":"Casablanca (film)","rank":11,"relativeRank":0.4744908528236409},"mass":0.0034923343190362585},{"year":"1941","bestMovie":{"pageRank":0.0018372075076142217,"position":1,"releaseYear":1941,"title":"Citizen Kane","rank":1,"relativeRank":1},"mass":0.004981955519034979},{"year":"1940","bestMovie":{"pageRank":0.00077505970958493,"position":22,"releaseYear":1940,"title":"Fantasia (film)","rank":22,"relativeRank":0.4218683553015818},"mass":0.004737369425162225},{"year":"1939","bestMovie":{"pageRank":0.0017739509822159226,"position":2,"releaseYear":1939,"title":"The Wizard of Oz (1939 film)","rank":2,"relativeRank":0.9655691993766978},"mass":0.005835166448335492},{"year":"1938","bestMovie":{"pageRank":0.00034402221594473045,"position":230,"releaseYear":1938,"title":"Olympia (1938 film)","rank":230,"relativeRank":0.18725278147348423},"mass":0.0024429209210234896},{"year":"1937","bestMovie":{"pageRank":0.0008501286705833058,"position":14,"releaseYear":1937,"title":"Snow White and the Seven Dwarfs (1937 film)","rank":14,"relativeRank":0.46272871576018865},"mass":0.0026848586870983023},{"year":"1936","bestMovie":{"pageRank":0.0002074510082457919,"position":640,"releaseYear":1936,"title":"The Great Ziegfeld","rank":640,"relativeRank":0.1129164818813448},"mass":0.0021484330239650495},{"year":"1935","bestMovie":{"pageRank":0.00047107735829328185,"position":95,"releaseYear":1935,"title":"Triumph of the Will","rank":95,"relativeRank":0.2564094454986296},"mass":0.003106509318076534},{"year":"1934","bestMovie":{"pageRank":0.0002922502577067247,"position":332,"releaseYear":1934,"title":"It Happened One Night","rank":332,"relativeRank":0.1590730804743106},"mass":0.0019096874734707632},{"year":"1933","bestMovie":{"pageRank":0.000396286889566884,"position":155,"releaseYear":1933,"title":"King Kong (1933 film)","rank":155,"relativeRank":0.21570066958930403},"mass":0.002662468519928841},{"year":"1932","bestMovie":{"pageRank":0.0001862101511062614,"position":785,"releaseYear":1932,"title":"Grand Hotel (film)","rank":785,"relativeRank":0.10135499138476303},"mass":0.0013952354982691996},{"year":"1931","bestMovie":{"pageRank":0.0003065031627152852,"position":291,"releaseYear":1931,"title":"Frankenstein (1931 film)","rank":291,"relativeRank":0.16683099837388918},"mass":0.0025039852192866256},{"year":"1930","bestMovie":{"pageRank":0.0004554884250111195,"position":105,"releaseYear":1930,"title":"The Blue Angel","rank":105,"relativeRank":0.24792432162582004},"mass":0.002057696087884327},{"year":"1929","bestMovie":{"pageRank":0.00016174930220479256,"position":1025,"releaseYear":1929,"title":"The Show of Shows","rank":1025,"relativeRank":0.08804084543222801},"mass":0.001252842564481625},{"year":"1928","bestMovie":{"pageRank":0.00010269644067249991,"position":2155,"releaseYear":1928,"title":"The Singing Fool","rank":2155,"relativeRank":0.05589811724961892},"mass":0.00089654716517078},{"year":"1927","bestMovie":{"pageRank":0.0006382772173013916,"position":36,"releaseYear":1927,"title":"Metropolis (film)","rank":36,"relativeRank":0.34741705259535527},"mass":0.0020307219401442663},{"year":"1926","bestMovie":{"pageRank":0.0003458442143460125,"position":228,"releaseYear":1926,"title":"The General (1926 film)","rank":228,"relativeRank":0.18824450309106464},"mass":0.0007593347713125997},{"year":"1925","bestMovie":{"pageRank":0.00019906930358372694,"position":692,"releaseYear":1925,"title":"The Phantom of the Opera (1925 film)","rank":692,"relativeRank":0.10835428374785831},"mass":0.0009690232422088751},{"year":"1924","bestMovie":{"pageRank":0.00020245740361040081,"position":665,"releaseYear":1924,"title":"Greed (film)","rank":665,"relativeRank":0.11019844126007838},"mass":0.0007967860540542872},{"year":"1923","bestMovie":{"pageRank":0.00017859506821638714,"position":847,"releaseYear":1923,"title":"The Ten Commandments (1923 film)","rank":847,"relativeRank":0.09721006880072508},"mass":0.0004656107536788885},{"year":"1922","bestMovie":{"pageRank":0.00046912417244525585,"position":99,"releaseYear":1922,"title":"Nosferatu","rank":99,"relativeRank":0.2553463179858521},"mass":0.0009102756438808401},{"year":"1921","bestMovie":{"pageRank":0.00012497831572685908,"position":1573,"releaseYear":1921,"title":"The Sheik (film)","rank":1573,"relativeRank":0.0680262383039979},"mass":0.00040764727193633296},{"year":"1920","bestMovie":{"pageRank":0.00033686033414781605,"position":237,"releaseYear":1920,"title":"The Cabinet of Dr. Caligari","rank":237,"relativeRank":0.18335453820633432},"mass":0.0007331650385380932},{"year":"1919","bestMovie":{"pageRank":0.00020056839787589302,"position":678,"releaseYear":1919,"title":"Broken Blossoms","rank":678,"relativeRank":0.10917024726093626},"mass":0.0003813666229638665},{"year":"1918","bestMovie":{"pageRank":0.00004048370475480583,"position":7640,"releaseYear":1918,"title":"Hearts of the World","rank":7640,"relativeRank":0.022035455759364678},"mass":0.00017030549655048513},{"year":"1917","bestMovie":{"pageRank":0.0000616275459613803,"position":4565,"releaseYear":1917,"title":"Cleopatra (1917 film)","rank":4565,"relativeRank":0.03354414006363885},"mass":0.000229113667921619},{"year":"1916","bestMovie":{"pageRank":0.0004908499534185591,"position":82,"releaseYear":1916,"title":"Intolerance (film)","rank":82,"relativeRank":0.2671717546244799},"mass":0.0006674634161069272},{"year":"1915","bestMovie":{"pageRank":0.0008282898898679914,"position":16,"releaseYear":1915,"title":"The Birth of a Nation","rank":16,"relativeRank":0.45084177287278776},"mass":0.0009159432259717657},{"year":"1914","bestMovie":{"pageRank":0.0001535461353541795,"position":1140,"releaseYear":1914,"title":"Cabiria","rank":1140,"relativeRank":0.08357582620243748},"mass":0.0005671253601640353},{"year":"1913","bestMovie":{"pageRank":0.0000446688247309195,"position":6802,"releaseYear":1913,"title":"Atlantis (1913 film)","rank":6802,"relativeRank":0.02431343468050921},"mass":0.00008326575686520897}];

function wikilink(t) {
  return "http://en.wikipedia.org/wiki/" + _.escape(t);
}

function removeBrackets(t) {
  return t.replace(/ *\([^)]*\)/g, "");
}

$(document).ready(function() {
  var massByYear = MASS_BY_YEAR;
  var years = _.pluck(massByYear, 'year');

  var height = 3000;
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
         .range([0, 420]);

  var y = d3.scale.ordinal()
    .domain(years)
    .rangeBands([0, height]);


  var displayYears = _.filter(years, function(year) { return year % 4 === (years[0] % 4);});
  console.log(displayYears);
  var yearOrdinal = d3.scale.ordinal()
    .domain(displayYears)
    .rangeBands([0, height]);

  var bars = chart.append("g").attr("transform", "translate(" + lineOffset + ")");

  var dataReachDelay = function(d, i) {
    return i * 30;
  };

  bars.selectAll("rect")
      .data(massByYear)
    .enter().append("rect")
      .attr("y", function(d) { return y(d.mass); })
      .attr("height", y.rangeBand())
      .attr("width", 0)
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
      .attr("x", width - 100)
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

  chart.selectAll("ticks")
      .data(displayYears)
    .enter()
      .append("path")
      .attr("d", function(d) {
        return line([[0,yearOrdinal(d) + 1], [width,yearOrdinal(d) + 1]]);
      }).attr("stroke", "gray")
        .attr("stroke-width", 0.2);

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
