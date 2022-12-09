function myFunction() {
    document.getElementById("dropDown").classList.toggle('show')
}
  
window.onclick = function(e) {
    if (e.target.classList.value != 'bi bi-person') {
        document.getElementById("dropDown").classList.remove('show')
    }
}