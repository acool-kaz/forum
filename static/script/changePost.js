const change_tags = document.querySelector('.shown-categories')
const change_category = document.querySelector('.category')
// var usedTags = change_category.value.slice(0, change_category.value.length-1).split(' ')

function func(input) {
    if (input.value == ' ') {
        input.value = ''
        return
    }
    if (input.value[input.value.length - 1] == ' ') {
        var inTag = input.value.slice(0, input.value.length-1)
        console.log(usedTags);
        for (var i = 0; i < usedTags.length; i++) {
            if (usedTags[i] == inTag) {
                input.value = ''
                return
            }
        }
        change_tags.innerHTML += `<div class='change-tag'>` + inTag + `<span class='delete-tag' onclick='del(this)' /></div>`
        usedTags.push(inTag)
        change_category.value += inTag + " "
        input.value = ''
    }
}

function del(e) {
    e.parentElement.style.display = 'none';
    usedTags.splice(usedTags.indexOf(e.parentElement.innerText), 1)
    change_category.value = usedTags.join(' ')
}