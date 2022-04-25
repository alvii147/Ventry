const exportCSV = () => {
  link = document.createElement("a");
  link.href = "/export";
  link.download = "inventory.csv";
  link.click()
}

const deleteItem = (itemId) => {
  console.log("HELLO")
  fetch(
    `/items/${itemId}`,
    {
      method: "DELETE",
    },
  )
  .then(response => {
    window.location.href = "/";
  })
  .catch(error => {
    console.log(error);
  });
}
