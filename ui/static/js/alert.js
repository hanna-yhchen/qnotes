const Swal = window.Swal;

var deleteButtons = document.getElementsByName("deleteButton")
deleteButtons.forEach(element => {
  element.addEventListener("click", function () {
    Swal.fire({
      title: 'Are you sure?',
      text: "You won't be able to revert this!",
      icon: 'warning',
      showCancelButton: true,
      confirmButtonColor: "var(--bs-dark)",
      confirmButtonText: 'Yes, delete it!'
    }).then(result => {
      if (result.isConfirmed) {
        const csrf = document.getElementsByName("csrf_token")[0].value
        const formData = new FormData()
        formData.append("csrf_token", csrf)

        fetch(`/note/${element.id}`, {
          method: "DELETE",
          body: formData
        }).then(response => {
          if (!response.ok) {
            Swal.fire({
              title: "Oops...",
              text: "Something went wrong",
              icon: "error",
              confirmButtonColor: "var(--bs-dark)"
            })
          } else {
            Swal.fire({
              title: "Deleted!",
              text: "The note has been deleted.",
              icon: "success",
              confirmButtonColor: "var(--bs-dark)",
              willClose: () => {
                window.location.assign("/")
              }
            })
          }
        })
      }
    })
  })
})
