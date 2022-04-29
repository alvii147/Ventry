var deliveredAtElem = document.getElementById("delivered_at");
deliveredAtElem.min = new Date().toISOString().split("T")[0];

var itemsSelected = {};
const selectItem = (itemId, elemId) => {
  let itemElem = document.getElementById(elemId);
  if (itemElem.checked) {
    itemsSelected[itemId] = true;
  }
  else {
    delete itemsSelected[itemId];
  }
}
