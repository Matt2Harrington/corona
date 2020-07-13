$(document).ready(function () {
  var dataResults = [];
  // var timestamps = [];
  var deaths = [];
  var countries = [];
  // var countries = [];
  $.ajax({
      url: 'http://localhost:10000/new',
      headers: {
          'Content-Type': 'application/x-www-form-urlencoded'
      },
      type: "GET", /* or type:"GET" or type:"PUT" */
      dataType: "json",
      data: {
      },
      success: function (result) {
          dataResults.push(result);

      },
      error: function () {
          console.log("error");
      }
  });

  setTimeout(function() {
    var i;
    for (i = 0; i < dataResults[0]['DataList'].length; i++) {
        if (dataResults[0]['DataList'][i].Country != 'World' || dataResults[0]['DataList'][i].Deaths < 1000){
          countries.push(dataResults[0]['DataList'][i].Country)
          // timestamps.push(dataResults[0]['DataList'][i].Updated)
          deaths.push(dataResults[0]['DataList'][i].Deaths)
          var ctx = document.getElementById('worldwideDeaths').getContext('2d');
          var chart = new Chart(ctx, {
              // The type of chart we want to create
              type: 'bar',

              // The data for our dataset
              data: {
                  labels: countries,
                  datasets: [{
                      label: 'Worldwide Death Percentage',
                      borderColor: 'rgb(19, 235, 162)',
                      data: deaths,
                      backgroundColor: 'rgb(19, 235, 162)',
                      hoverBackgroundColor: 'rgb(19, 235, 200)',
                      maxBarThickness: 8,
                      minBarLength: 2
                  }]
              },

              // Configuration options go here
              options: {
                tooltips: {
                  mode: 'point'
                },
                fill: false,
                events: ['click']
              }
          });
      }
    }
      // console.log(dataResults[0]['DataList'][0].length)
  },200);
})
