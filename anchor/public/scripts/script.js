function Search() {
    var input, filter, table, tr, td, i;
    input = document.getElementById("search");
    filter = input.value.toUpperCase();
    table = document.getElementById("table");
    tr = table.getElementsByTagName("tr");
    for (i = 0; i < tr.length; i++) {
        td = tr[i].getElementsByTagName("td")[0];
        if (td) {
            if (td.innerHTML.toUpperCase().indexOf(filter) > -1) {
                tr[i].style.display = "";
            } else {
                tr[i].style.display = "none";
            }
        }
    }
}

function Live() {
    var filter, table, tr, td, i;
    filter = "fa-check";
    table = document.getElementById("table");
    tr = table.getElementsByTagName("tr");
    for (i = 0; i < tr.length; i++) {
        td = tr[i].getElementsByTagName("td")[6];
        if (td) {
            if (td.innerHTML.indexOf(filter) > -1) {
                tr[i].style.display = "";
            } else {
                tr[i].style.display = "none";
            }
        }
    }
}

function Stopped() {
    var filter, table, tr, td, i;
    filter = "fa-times";
    table = document.getElementById("table");
    tr = table.getElementsByTagName("tr");
    for (i = 0; i < tr.length; i++) {
        td = tr[i].getElementsByTagName("td")[6];
        if (td) {
            if (td.innerHTML.indexOf(filter) > -1) {
                tr[i].style.display = "";
            } else {
                tr[i].style.display = "none";
            }
        }
    }
}

function All() {
    var  table, tr, i;
    table = document.getElementById("table");
    tr = table.getElementsByTagName("tr");
    for (i = 0; i < tr.length; i++) {
        tr[i].style.display = "";
    }
}

function Copy() {
    var copyText = document.getElementById("key");
    copyText.select();
    document.execCommand("copy");
    document.getSelection().removeAllRanges();
}

function syntaxHighlight(json) {
    json = json.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;');
    return json.replace(/("(\\u[a-zA-Z0-9]{4}|\\[^u]|[^\\"])*"(\s*:)?|\b(true|false|null)\b|-?\d+(?:\.\d*)?(?:[eE][+\-]?\d+)?)/g, function (match) {
        var cls = 'number';
        if (/^"/.test(match)) {
            if (/:$/.test(match)) {
                cls = 'key';
            } else {
                cls = 'string';
            }
        } else if (/true|false/.test(match)) {
            cls = 'boolean';
        } else if (/null/.test(match)) {
            cls = 'null';
        }
        console.log(json)
        return '<span class="' + cls + '">' + match + '</span>';
    });
}

function deleteService(btn) {
  row = btn.parentNode.parentElement.parentElement.children
btn.parents
  name = row[1].innerText
  namespace = row[2].innerText
  $.ajax({
      type: "DELETE",
      url: "/services",
      contentType:"application/json",
      data: JSON.stringify({namespace: namespace, name: name}),//参数列表
      dataType:"json",
      success: function(result){
				$("#result").text("成功删除服务")
        $('#success-result').modal()
      },
      error: function(result){
				$("#result").text("删除服务失败")
        $('#danger-result').modal()
      }
  });
  $('#success-result').on('hidden.bs.modal', function () {
     location.reload();
  })
  $('#danger-result').on('hidden.bs.modal', function () {
     location.reload();
  })
}

function deleteDeployment(btn) {
  row = btn.parentNode.parentElement.parentElement.children
btn.parents
  name = row[1].innerText
  namespace = row[2].innerText
  $.ajax({
      type: "DELETE",
      url: "/deployments",
      contentType:"application/json",
      data: JSON.stringify({namespace: namespace, name: name}),//参数列表
      dataType:"json",
      success: function(result){
				$("#result").text("成功删除部署")
        $('#success-result').modal()
      },
      error: function(result){
				$("#result").text("删除部署失败")
        $('#danger-result').modal()
      }
  });
  $('#success-result').on('hidden.bs.modal', function () {
     location.reload();
  })
  $('#danger-result').on('hidden.bs.modal', function () {
     location.reload();
  })
}
