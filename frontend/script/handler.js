import $ from 'jquery'

$( document ).ready( docReadyHandler )

function signOut() {
  $("#userinfo").addClass('d-none');
  $("#signin").removeClass('d-none');
}

function docReadyHandler() {
  fetch('/auth/self')
    .then((response) => {
  	    return response.json();
    })
  	.then((data) => {
      if ($.isEmptyObject(data)) {
        $("#userinfo").addClass('d-none');
        $("#signin").removeClass('d-none');
      } else {
        $("#userinfo").removeClass('d-none');
        $("#signin").addClass('d-none');
        $("#user_icon").attr("src", data.icon);
        $("#user_name").html(data.name)
      }
  	});
}
