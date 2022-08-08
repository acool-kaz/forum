function func() {
    const password = document.getElementById('Password')
    const confirmPassword = document.getElementById('ConfirmPassword')
    if (confirmPassword.value !== password.value) {
        confirmPassword.style.border = '1px solid red'
        password.style.border = '1px solid red'
    } else {
        confirmPassword.style.border = '1px solid #fff'
        password.style.border = '1px solid #fff'
    }
}