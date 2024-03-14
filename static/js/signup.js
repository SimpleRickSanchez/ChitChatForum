document.addEventListener('DOMContentLoaded', function () {
    const emailInput = document.getElementById('email');
    const nameInput = document.getElementById('name');
    const passwordInput = document.getElementById('password');
    const loginButton = document.getElementById('login-btn');
    const saltInput = document.getElementById('salt');

    nameInput.addEventListener('input', debounce(async function (event) {
        document.getElementById("msg-display1").textContent = "";
    }, 500));

    emailInput.addEventListener('input', debounce(async function (event) {
        document.getElementById("msg-display2").textContent = "";
    }, 500));

    passwordInput.addEventListener('input', debounce(async function (event) {
        document.getElementById("msg-display3").textContent = "";
    }, 500));

    loginButton.addEventListener('click', async function (event) {
        event.preventDefault();

        const email = emailInput.value;
        const name = nameInput.value;
        const password = passwordInput.value;
        const salt = saltInput.value;
        const pwdmd5 = md5(password + salt);
        var emailValid = false;
        var nameValid = false;
        var pwdValid = false;
        if (name.length > 0) {
            nameValid = true;
        } else {
            nameValid = false;
            document.getElementById("msg-display1").textContent = "Name empty!";
        }
        if (email.length > 0) {
            try {
                const response = await fetch('/user/checkexists', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({ email })
                });
                if (response.ok) {
                    const data = await response.json();
                    console.log(data)
                    document.getElementById("msg-display2").textContent = data.msg;
                    if (data && data.emailValid) {
                        emailValid = true;
                    } else {
                        emailValid = false;
                    }
                } else {
                    console.error('Failed to fetch info', response.msg);
                }
            } catch (error) {
                console.error('Error fetching info', error);
            }
        } else {
            emailValid = false;
            document.getElementById("msg-display2").textContent = "Email empty!";
        }

        if (password.length >= 6) {
            pwdValid = true;
        } else {
            pwdValid = false;
            document.getElementById("msg-display3").textContent = "Password must be >= 6 characters!";
        }

        if (emailValid && nameValid && pwdValid) {
            try {
                const regResponse = await fetch('/user/signup_account', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({ email, name, pwdmd5 })
                });
                if (regResponse.ok) {
                    const regData = await regResponse.json();
                    document.getElementById("msg-display4").textContent = regData.msg;
                    if (regData && regData.success) {
                        setTimeout(function () {
                            window.location.replace('/');
                        }, 2000);
                    }
                } else {
                    console.error('Failed to fetching signup_account', regResponse.msg);
                }
            } catch (error) {
                console.error('Error signup_account', error);
            }
        }
    });
    function md5(str) {
        return CryptoJS.MD5(str).toString();
    }
    function debounce(func, wait) {
        let timeout;
        return function () {
            const context = this;
            const args = arguments;
            clearTimeout(timeout);
            timeout = setTimeout(function () {
                func.apply(context, args);
            }, wait);
        };
    }
});