var info = document.getElementById('info');
var url = "http://localhost:3000/";

var cont = 1;
var cpu=0;
var ret = [];

setInterval(function(){ 
    axios.get(url+"cpu")
        .then(function(response) {
            info.innerHTML = `<h3>%CPU utilizado: ${response.data.CPU}%</h3>`

            cpu = response.data.CPU
        })
        .catch(function(error) {
            console.log(error);
        })
        .then(function() {}); 
}, 3000);

function data() {
      ret.push({
        no: cont++,
        value: cpu
      });

    return ret;
  }

var graph = new Morris.Line({
    // ID of the element in which to draw the chart.
    element: 'myfirstchart',
    data: data(),
    // The name of the data record attribute that contains x-values.
    xkey: 'no',
    // A list of names of data record attributes that contain y-values.
    ykeys: ['value'],
    // Labels for the ykeys -- will be displayed when you hover over the
    // chart.'
    labels: ['%CPU utilizado'],
    resize: true,
    parseTime: false,
    hideHover: true,
    lineColors: ['#DF3A01']
  });

  function update() {
    graph.setData(data());
  }

setInterval(update, 3000);