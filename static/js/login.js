document.addEventListener('DOMContentLoaded', function () {
    const emailInput = document.getElementById('email');
    const passwordInput = document.getElementById('password');
    const loginButton = document.getElementById('login-btn');

    emailInput.addEventListener('input', debounce(async function () {
        const email = this.value;
        try {
            const response = await fetch('/user/check', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ email })
            });

            if (response.ok) {
                const data = await response.json();
                document.getElementById("msg-display").textContent = data.msg;
            } else {
                console.error('Failed to fetch info', response.msg);
            }
        } catch (error) {
            console.error('Error fetching info', error);
        }
    }, 1000));

    loginButton.addEventListener('click', async function (event) {
        event.preventDefault();

        const email = emailInput.value;
        const password = passwordInput.value;

        try {
            const saltResponse = await fetch('/user/salt', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ email })
            });

            if (saltResponse.ok) {
                const saltData = await saltResponse.json();
                const salt = saltData.salt;
                const pwdmd5 = md5(password + salt);
                try {
                    const authResponse = await fetch('/user/authenticate', {
                        method: 'POST',
                        headers: {
                            'Content-Type': 'application/json'
                        },
                        body: JSON.stringify({ email, pwdmd5 })
                    });
                    if (authResponse.ok) {
                        const authData = await authResponse.json();
                        document.getElementById("msg-display").textContent = authData.msg;
                        if (authData && authData.success) {
                            setTimeout(function () {
                                window.location.replace('/');
                            }, 2000);
                        }
                    } else {
                        console.error('Failed to fetching authenticate', authResponse.msg);
                    }
                } catch (error) {
                    console.error('Error authenticate', error);
                }
            } else {
                console.error('Failed to fetch salt', saltResponse.msg);
            }
        } catch (error) {
            console.error('Error fetching salt or hashing password', error);
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