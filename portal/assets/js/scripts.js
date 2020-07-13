$(document).ready(function () {
  var dataResults = [];
  var timestamps = [];
  var cases = [];
  var deaths = [];
  // var countries = [];
  $.ajax({
      url: 'http://localhost:10000/corona',
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
      if (dataResults[0]['DataList'][i].Country === 'USA') {
        timestamps.push(dataResults[0]['DataList'][i].Updated)
        cases.push(dataResults[0]['DataList'][i].Cases)
        var ctx = document.getElementById('myChart').getContext('2d');
        var chart = new Chart(ctx, {
            // The type of chart we want to create
            type: 'line',

            // The data for our dataset
            data: {
                labels: timestamps,
                datasets: [{
                    label: 'United States Cases of COVID-19',
                    borderColor: 'rgb(19, 235, 162)',
                    data: cases,
                    fill: false
                }]
            },

            // Configuration options go here
            options: {
              fill: false
            }
        });
      } else {
        // console.log(dataResults[0]['DataList'][i].ID)
      }
    }
      // console.log(dataResults[0]['DataList'][0].length)
  },200);
})
