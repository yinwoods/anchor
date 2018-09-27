// javascript

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
    return '<span class="' + cls + '">' + match + '</span>';
  });
}

function deletePod(btn) {
  row = btn.parentNode.parentElement.parentElement.children
  name = row[1].innerText
  namespace = row[2].innerText
  $.ajax({
    type: "DELETE",
    url: "/pods",
    contentType:"application/json",
    data: JSON.stringify({namespace: namespace, name: name}),//参数列表
    dataType:"json",
    success: function(result){
      $("#success-result").text("成功删除Pod")
      $("#modal-success").modal()
    },
    error: function(result){
      $("#danger-result").text("删除Pod失败")
      $("#modal-danger").modal()
    }
  });
  $("#modal-success").on('hidden.bs.modal', function () {
    location.reload();
  })
  $("#modal-danger").on('hidden.bs.modal', function () {
    location.reload();
  })
}

function updatePod(btn) {
  content = btn.parentElement.parentElement.children[1].firstElementChild.lastElementChild.value
  $.ajax({
    type: "PUT",
    url: "/pods",
    contentType:"application/json",
    data: JSON.stringify({body: content}),//参数列表
    dataType:"json",
    success: function(result){
      $("#update").modal("hide")
      $("#success-result").text("更新Pod成功")
      $("#modal-success").modal()
    },
    error: function(result){
      $("#update").modal("hide")
      $("#danger-result").text("更新Pod失败")
      $("#modal-danger").modal()
    }
  });
  $("#modal-success").on('hidden.bs.modal', function () {
    location.replace("/pods/" + namespace + "/" + name)
  })
  $("#modal-danger").on('hidden.bs.modal', function () {
    location.reload();
  })
}

function showPodConfigModal(btn) {
  row = btn.parentNode.parentElement.parentElement.children
  btn.parents
  namespace = row[2].innerText
  name = row[1].innerText
  $.ajax({
    type: "GET",
    url: "/api/pods/" + namespace + "/" + name,
    contentType:"application/json",
    dataType:"json",
    success: function(result){
      $("#update").modal("show")
      $("#updateTextArea").val(result["result"])
    },
    error: function(result){
      $("#danger-result").text("获取pod信息失败")
      $("#modal-danger").modal()
    }
  });
}

function deleteContainer(btn) {
  row = btn.parentNode.parentElement.parentElement.children
  cid = row[1].innerText
  $.ajax({
    type: "DELETE",
    url: "/containers",
    contentType:"application/json",
    data: JSON.stringify({cid: cid}),//参数列表
    dataType:"json",
    success: function(result){
      $("#success-result").text("成功删除容器")
      $("#modal-success").modal()
    },
    error: function(result){
      $("#danger-result").text("删除容器失败")
      $("#modal-danger").modal()
    }
  });
  $("#modal-success").on('hidden.bs.modal', function () {
    location.reload();
  })
  $("#modal-danger").on('hidden.bs.modal', function () {
    location.reload();
  })
}

function deleteService(btn) {
  row = btn.parentNode.parentElement.parentElement.children
  name = row[1].innerText
  namespace = row[2].innerText
  console.log(namespace)
  $.ajax({
    type: "DELETE",
    url: "/services",
    contentType:"application/json",
    data: JSON.stringify({namespace: namespace, name: name}),//参数列表
    dataType:"json",
    success: function(result){
      $("#success-result").text("成功删除服务")
      $("#modal-success").modal()
    },
    error: function(result){
      $("#danger-result").text("删除服务失败")
      $("#modal-danger").modal()
    }
  });
  $("#modal-success").on('hidden.bs.modal', function () {
    location.reload();
  })
  $("#modal-danger").on('hidden.bs.modal', function () {
    location.reload();
  })
}

function updateService(btn) {
  content = btn.parentElement.parentElement.children[1].firstElementChild.lastElementChild.value
  $.ajax({
    type: "PUT",
    url: "/services",
    contentType:"application/json",
    data: JSON.stringify({body: content}),//参数列表
    dataType:"json",
    success: function(result){
      $("#update").modal("hide")
      $("#success-result").text("更新服务成功")
      $("#modal-success").modal()
    },
    error: function(result){
      $("#update").modal("hide")
      $("#danger-result").text("更新服务失败")
      $("#modal-danger").modal()
    }
  });
  $("#modal-success").on('hidden.bs.modal', function () {
    location.replace("/services/" + namespace + "/" + name)
  })
  $("#modal-danger").on('hidden.bs.modal', function () {
    location.reload();
  })
}

function showServiceConfigModal(btn) {
  row = btn.parentNode.parentElement.parentElement.children
  btn.parents
  namespace = row[2].innerText
  name = row[1].innerText
  $.ajax({
    type: "GET",
    url: "/api/services/" + namespace + "/" + name,
    contentType:"application/json",
    dataType:"json",
    success: function(result){
      $("#update").modal("show")
      $("#updateTextArea").val(result["result"])
    },
    error: function(result){
      $("#danger-result").text("获取服务信息失败")
      $("#modal-danger").modal()
    }
  });
}

function deleteDeployment(btn) {
  row = btn.parentNode.parentElement.parentElement.children
  name = row[1].innerText
  namespace = row[2].innerText
  $.ajax({
    type: "DELETE",
    url: "/deployments",
    contentType:"application/json",
    data: JSON.stringify({namespace: namespace, name: name}),//参数列表
    dataType:"json",
    success: function(result){
      $("#success-result").text("成功删除部署")
      $("#modal-success").modal()
    },
    error: function(result){
      $("#danger-result").text("删除部署失败")
      $("#modal-danger").modal()
    }
  });
  $("#modal-success").on('hidden.bs.modal', function () {
    location.reload();
  })
  $("#modal-danger").on('hidden.bs.modal', function () {
    location.reload();
  })
}

function updateDeployment(btn) {
  content = btn.parentElement.parentElement.children[1].firstElementChild.lastElementChild.value
  $.ajax({
    type: "PUT",
    url: "/deployments",
    contentType:"application/json",
    data: JSON.stringify({body: content}),//参数列表
    dataType:"json",
    success: function(result){
      $("#update").modal("hide")
      $("#success-result").text("更新部署成功")
      $("#modal-success").modal()
    },
    error: function(result){
      $("#update").modal("hide")
      $("#danger-result").text("更新部署失败")
      $("#modal-danger").modal()
    }
  });
  $("#modal-success").on('hidden.bs.modal', function () {
    location.replace("/deployments/" + namespace + "/" + name)
  })
  $("#modal-danger").on('hidden.bs.modal', function () {
    location.reload();
  })
}

function showDeploymentConfigModal(btn) {
  row = btn.parentNode.parentElement.parentElement.children
  btn.parents
  namespace = row[2].innerText
  name = row[1].innerText
  $.ajax({
    type: "GET",
    url: "/api/deployments/" + namespace + "/" + name,
    contentType:"application/json",
    dataType:"json",
    success: function(result){
      $("#update").modal("show")
      $("#updateTextArea").val(result["result"])
    },
    error: function(result){
      $("#danger-result").text("获取部署信息失败")
      $("#modal-danger").modal()
    }
  });
}
