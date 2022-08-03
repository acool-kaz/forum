const tags = document.querySelector('.shown-categories')
const category = document.querySelector('.category')
var usedTags = []

function func(input) {
    if (input.value == ' ') {
        input.value = ''
        return
    }
    if (input.value[input.value.length - 1] == ' ') {
        for (var i = 0; i < usedTags.length; i++) {
            if (usedTags[i] == input.value) {
                input.value = ''
                return
            }
        }
        tags.innerHTML += `<div class='tag'>` + input.value + `<span class='close' onclick='del(this)' /></div>`
        usedTags.push(input.value)
        category.value += input.value
        input.value = ''
    }
}

function del(e) {
    e.parentElement.style.display = 'none';
    usedTags.splice(usedTags.indexOf(e.parentElement.innerText), 1)
    category.value = usedTags.join('')
}