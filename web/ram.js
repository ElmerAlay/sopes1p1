var info = document.getElementById('info');
var url = "http://localhost:3000/";

var cont = 1;
var libre=0;
var utilizada=0;
var ret = [];

setInterval(function(){ 
    axios.get(url+"memoria")
        .then(function(response) {
            info.innerHTML = `<h3>Memoria Total: ${response.data.Total_mem} MB</h3>
            <h3>Memoria Libre: ${response.data.Free_mem} MB</h3>
            <h3>Porcentaje de memoria utilizada: ${response.data.Porcent}%</h3>`

            libre = parseInt(response.data.Free_mem)
            utilizada = (parseInt(response.data.Total_mem) - parseInt(response.data.Free_mem))
        })
        .catch(function(error) {
            console.log(error);
        })
        .then(function() {}); 
}, 3000);

function data() {
      ret.push({
        no: cont++,
        value: utilizada,
        val2: libre
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
    ykeys: ['value','val2'],
    // Labels for the ykeys -- will be displayed when you hover over the
    // chart.
    labels: ['Memoria utilizada MB','Memoria libre MB'],
    resize: true,
    parseTime: false,
    hideHover: true,
    lineColors: ['#0101DF','#DF3A01']
  });

  function update() {
    graph.setData(data());
  }

setInterval(update, 3000);