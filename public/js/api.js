const baseURL = "http://localhost:3000"

const deletePost = id => {
  fetch(`${baseURL}/deletePost/${id}`, {
    method: "DELETE"
  }).then(() => {
    location.reload()
  }
  ).catch(console.error)

}

