const images = document.querySelectorAll('.img-container img')

var counter = 0
images[counter].style.display = 'block'

const nextBtn = document.getElementById('nextBtn')
const prevBtn = document.getElementById('prevBtn')

nextBtn.addEventListener('click', ()=>{
    images[counter].style.display = 'none'
    counter = (counter+1)%images.length
    images[counter].style.display = 'block'
});

prevBtn.addEventListener('click', ()=>{
    images[counter].style.display = 'none'
    counter--
    if (counter < 0) {
        counter = images.length-1
    }
    images[counter].style.display = 'block'
});