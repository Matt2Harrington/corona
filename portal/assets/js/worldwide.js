// $(document).ready(function () {
//   var dataResults = [];
//   // var timestamps = [];
//   var deaths = [];
//   var countries = [];
//   // var countries = [];
//   $.ajax({
//       url: 'http://localhost:10000/corona',
//       headers: {
//           'Content-Type': 'application/x-www-form-urlencoded'
//       },
//       type: "GET", /* or type:"GET" or type:"PUT" */
//       dataType: "json",
//       data: {
//       },
//       success: function (result) {
//           dataResults.push(result);
//
//       },
//       error: function () {
//           console.log("error");
//       }
//   });
//
//   setTimeout(function() {
//     var i;
//     for (i = 0; i < dataResults[0]['DataList'].length; i++) {
//       countries.push(dataResults[0]['DataList'][i].Country)
//       // timestamps.push(dataResults[0]['DataList'][i].Updated)
//       deaths.push(dataResults[0]['DataList'][i].Deaths)
//       var ctx = document.getElementById('worldwideDeaths').getContext('2d');
//       var chart = new Chart(ctx, {
//           // The type of chart we want to create
//           type: 'pie',
//
//           // The data for our dataset
//           data: {
//               labels: countries,
//               datasets: [{
//                   label: 'Worldwide Death Percentage',
//                   borderColor: 'rgb(255, 99, 132)',
//                   data: deaths
//               }]
//           },
//
//           // Configuration options go here
//           options: {
//             fill: false
//           }
//       });
//     }
//       // console.log(dataResults[0]['DataList'][0].length)
//   },1000);
// })
