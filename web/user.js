var formulario = document.getElementById('formulario');
var traduccion = document.getElementById('c3d');
var campos = document.getElementById('campos');
var cabecera= document.getElementById('cabecera');
var formulario = document.getElementById('formulario');

var existe = false;
var url = "http://localhost:3000/";


    axios.get(url)
        .then(function(response) {
            campos.innerHTML = "";
            let cont = response.data.length

            cabecera.innerHTML = `<h5 align=center id="r">Procesos en ejecución: ${response.data[cont-1].Pid}</h5>
            <h5 align=center id="s">Procesos suspendidos: ${response.data[cont-1].Name}</h5>
            <h5 align=center id="t">Procesos detenidos: : ${response.data[cont-1].User}</h5>
            <h5 align=center id="z">Procesos zombie: ${response.data[cont-1].State}</h5>
            <h4 align=center id="total">Total de procesos: : ${response.data[cont-1].Ram}</h4>

`

            for(let i=0; i<cont-1; i++){
                campos.innerHTML += `                                <tr>
                <td class="column1">${response.data[i].Pid}</td>
                <td class="column2">${response.data[i].Name}</td>
                <td class="column3">${response.data[i].User}</td>
                <td class="column4">${response.data[i].State}</td>
                <td class="column6">${response.data[i].Ram}</td>
                <td class="column6"> <button class="btn btn-danger" onClick="eliminar(${response.data[i].Pid})">Matar</button></td>
            </tr>
`
            }
        })
        .catch(function(error) {
            console.log(error);
        })
        .then(function() {});


function eliminar(id){
        axios.post(url + 'kill',
                    {
                        Pid: id 
                    })
            .then(function(response) {
                console.log(response);
            })
            .catch(function(error) {
                console.log(error);
            })
            .then(function() {});
    
}

formulario.addEventListener('submit', async function(e) {
    e.preventDefault();

    var datos = new FormData(formulario);
    var pid = datos.get('pid');

    postUser(pid);
})

function postUser(pid) {
    axios.post(url + 'kill/', {
            Pid: pid
        })
        .then(function(response) {
            console.log(response);
        })
        .catch(function(error) {
            console.log(error);
        })
        .then(function() {});
}