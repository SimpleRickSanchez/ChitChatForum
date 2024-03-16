const links = document.querySelectorAll('#myList .show-input');
const inputs = document.querySelectorAll('.input-container');
let activeinput = null;
links.forEach(link => {
    link.addEventListener('click', function (event) {
        event.preventDefault();
        if (activeinput === this.getAttribute('ruuid')) {
            activeinput = null;
            inputs.forEach(function (input) {
                input.style.display = 'none';
            });

        } else {
            if (activeinput !== null) {
                inputs.forEach(function (inputContainer) {
                    inputContainer.style.display = 'none';
                });
            }
            const group = this.closest('.group');
            const inputContainer = group.querySelector('.input-container');
            // inputContainer.style.display = inputContainer.style.display === 'none' ? 'block' : 'none';
            inputContainer.style.display = 'block';

            const action = inputContainer.querySelector(".action");
            action.value = this.getAttribute('action');
            const tuuid = inputContainer.querySelector(".tuuid");
            tuuid.value = this.getAttribute('tuuid');
            const puuid = inputContainer.querySelector(".puuid");
            puuid.value = this.getAttribute('puuid');
            const ruuid = inputContainer.querySelector(".ruuid");
            ruuid.value = this.getAttribute('ruuid');
            activeinput = this.getAttribute('ruuid');

            const inputelem = inputContainer.querySelector("textarea");
            inputelem.focus();
            if (!inputelem.value) {
                inputelem.placeholder = this.getAttribute('data-text');
            }
        }
    });
});
const forms = document.querySelectorAll('.input-container form');
forms.forEach(form => {
    const summitButton = form.querySelector('.btn-primary');
    summitButton.addEventListener('click', async function (event) {
        event.preventDefault();
        form.closest('.input-container').display = 'none';
        const contentInput = form.querySelector(".inputext");
        content = contentInput.value;

        const tuuid = form.querySelector('.tuuid').value;
        const puuid = form.querySelector('.puuid').value;
        const ruuid = form.querySelector('.ruuid').value;
        const action = form.querySelector('.action').value;
        let obj = {
            content, tuuid, puuid, ruuid, action
        }
        try {
            const response = await fetch('/create', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(obj)
            });

            if (response.ok) {
                const data = await response.json();
                const scrollPosition = {
                    x: window.scrollX || window.pageXOffset,
                    y: window.scrollY || window.pageYOffset
                };
                window.location.reload();
                window.addEventListener('load', () => {
                    window.scrollTo(scrollPosition.x, scrollPosition.y);
                });
            } else {
                console.error('Failed to fetch result', response.msg);
            }
        } catch (error) {
            console.error('Error fetching result', error);
        }
    });


});