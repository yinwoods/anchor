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

function getResourceNameByType(type) {
  resourceName = "unknown"
  switch (type) {
    case "containers":
      resourceName = "容器"
      break
    case "images":
      resourceName = "镜像"
      break
    case "networks":
      resourceName = "网络"
      break
    case "pods":
      resourceName = "容器组";
      break;
    case "services":
      resourceName = "服务";
      break;
    case "deployments":
      resourceName = "部署"
      break
    case "ups":
      resourceName = "供电设备"
      break
    case "refs":
      resourceName = "制冷设备"
      break
  }
  return resourceName
}

function create(btn, type) {
  body = btn.parentNode.previousElementSibling.childNodes[1].children[1].value;
  resourceName = getResourceNameByType(type)
  $.ajax({
    type: "POST",
    url: "/" + type,
    contentType: "application/json",
    data: JSON.stringify({"body": body}),
    success: function(result) {
      $("#create").modal("hide")
      $("#success-result").text("成功创建" + resourceName);
      $("#modal-success").modal();
    },
    error: function(result) {
      $("#create").modal("hide")
      $("#danger-result").text("创建" + resourceName + "失败");
      $("#modal-danger").modal();
    }
  });
  $("#modal-success").on('hidden.bs.modal', function () {
    location.reload();
  })
  $("#modal-danger").on('hidden.bs.modal', function () {
    location.reload();
  })
}

function remove(btn, type) {
  row = btn.parentNode.parentElement.parentElement.children
  name = row[1].innerText
  namespace = row[2].innerText
  resourceName = getResourceNameByType(type)

  $.ajax({
    type: "DELETE",
    url: "/" + type,
    contentType:"application/json",
    data: JSON.stringify({"namespace": namespace, "name": name}),//参数列表
    dataType:"json",
    success: function(result){
      $("#success-result").text("成功删除" + resourceName)
      $("#modal-success").modal()
    },
    error: function(result){
      $("#danger-result").text("删除" + resourceName + "失败")
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

function update(btn, type) {
  content = btn.parentElement.parentElement.children[1].firstElementChild.lastElementChild.value
  resourceName = getResourceNameByType(type)

  $.ajax({
    type: "PUT",
    url: "/" + type,
    contentType:"application/json",
    data: JSON.stringify({"body": content}),//参数列表
    dataType:"json",
    success: function(result){
      $("#update").modal("hide")
      $("#success-result").text("更新" + resourceName + "成功")
      $("#modal-success").modal()
    },
    error: function(result){
      $("#update").modal("hide")
      $("#danger-result").text("更新" + resourceName + "失败")
      $("#modal-danger").modal()
    }
  });
  $("#modal-success").on('hidden.bs.modal', function () {
    location.replace("/" + type + "/" + namespace + "/" + name)
  })
  $("#modal-danger").on('hidden.bs.modal', function () {
    location.reload();
  })
}

function showConfigModal(btn, type) {
  row = btn.parentNode.parentElement.parentElement.children
  namespace = row[2].innerText
  name = row[1].innerText
  resourceName = getResourceNameByType(type)

  $.ajax({
    type: "GET",
    url: "/api/" + type + "/" + namespace + "/" + name,
    contentType:"application/json",
    dataType:"json",
    success: function(result){
      $("#update").modal("show")
      $("#updateTextArea").val(result["result"])
    },
    error: function(result){
      $("#danger-result").text("获取" + resourceName + "信息失败")
      $("#modal-danger").modal()
    }
  });
}

function showConfigModalByID(btn, type) {
  row = btn.parentNode.parentElement.parentElement.children
  link = row[1].firstElementChild.href.split("/")
  id = link[link.length - 1]
  resourceName = getResourceNameByType(type)

  $.ajax({
    type: "GET",
    url: "/api/" + type + "/" + id,
    contentType:"application/json",
    dataType:"json",
    success: function(result){
      label = document.getElementById("updateTextArea").previousElementSibling
      label.innerText += result["id"]
      $("#update").modal("show")
      $("#updateTextArea").val(result["result"])
    },
    error: function(result){
      $("#danger-result").text("获取" + resourceName + "信息失败")
      $("#modal-danger").modal()
    }
  });
}

function updateByID(btn, type) {
  content = btn.parentElement.parentElement.children[1].firstElementChild.lastElementChild.value
  resourceName = getResourceNameByType(type)
  label = document.getElementById("updateTextArea").previousElementSibling
  id = label.innerText

  $.ajax({
    type: "PUT",
    url: "/" + type,
    contentType:"application/json",
    data: JSON.stringify({"id": id, "body": content}),//参数列表
    dataType:"json",
    success: function(result){
      $("#update").modal("hide")
      $("#success-result").text("更新" + resourceName + "成功")
      $("#modal-success").modal()
    },
    error: function(result){
      $("#update").modal("hide")
      $("#danger-result").text("更新" + resourceName + "失败")
      $("#modal-danger").modal()
    }
  });
  $("#modal-success").on('hidden.bs.modal', function () {
    location.replace("/" + type + "/" + id)
  })
  $("#modal-danger").on('hidden.bs.modal', function () {
    location.reload();
  })
}

function removeByID(btn, type) {
  row = btn.parentNode.parentElement.parentElement.children
  link = row[1].firstElementChild.href.split("/")
  id = link[link.length - 1]
  resourceName = getResourceNameByType(type)
  $.ajax({
    type: "DELETE",
    url: "/" + type,
    contentType:"application/json",
    data: JSON.stringify({"id": id}),//参数列表
    dataType:"json",
    success: function(result){
      $("#success-result").text("成功删除" + resourceName)
      $("#modal-success").modal()
    },
    error: function(result){
      $("#danger-result").text("删除" + resourceName + "失败")
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

function showImageConfigModal(btn, type) {
  row = btn.parentNode.parentElement.parentElement.children
  mid = row[1].innerText
  resourceName = getResourceNameByType(type)

  $.ajax({
    type: "GET",
    url: "/api/" + type + "/" + mid,
    contentType:"application/json",
    dataType:"json",
    success: function(result){
      $("#update").modal("show")
      $("#updateTextArea").val(result["result"])
    },
    error: function(result){
      $("#danger-result").text("获取" + resourceName + "信息失败")
      $("#modal-danger").modal()
    }
  });
}
