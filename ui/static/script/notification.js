function showModal(id) {
    // document.querySelector('.modal').style.display = 'flex'
    document.getElementById(id).style.display = 'flex'
}

function closeModal(id) {
    // document.querySelector('.modal').style.display = 'none'
    document.getElementById(id).style.display = 'none'

}
 
const modals = document.querySelectorAll(".modal-content")

modals.forEach(element => {
    element.addEventListener("click", (e) => { e.stopPropagation() })
});