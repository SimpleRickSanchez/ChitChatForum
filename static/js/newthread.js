const dropbtn = document.getElementById('topics');
dropbtn.addEventListener('click', function(event){
  event.preventDefault();
} );
const topics = document.querySelectorAll('.topic');
topics.forEach(topic => {
  topic.addEventListener('click', function (event) {
    event.preventDefault();
    const choseone = document.getElementById('topics')
    choseone.textContent = this.textContent
    choseone.setAttribute("topicid", this.getAttribute("topicid"))
    document.getElementById("msg").textContent = ""
  })
});
const summitButton = document.getElementById('summit');
const title = document.getElementById('title');
const content = document.getElementById('content');
summitButton.addEventListener('click', async function (event) {
  event.preventDefault();
  const choseone = document.getElementById('topics')
  if (choseone.getAttribute("topicid") === '0') {
    document.getElementById("msg").textContent = "Please choose a topic."
    return
  }

  if (title.value.length === 0) {
    document.getElementById("msg").textContent = "Please enter the title."
    return
  }
  if (title.value.length > 300) {
    document.getElementById("msg").textContent = "Title too long. Max 300 characters."
    return
  }

  if (content.value.length === 0) {
    document.getElementById("msg").textContent = "Please enter some content."
    return
  }
  if (content.value.length > 1000) {
    document.getElementById("msg").textContent = "Please enter some content. Max 1000 characters."
    return
  }

  try {
    let obj = {};
    obj["title"] = title.value;
    obj["content"] = content.value;
    obj["topic_id"] = choseone.getAttribute("topicid");
    obj["action"] = "thread"
    const response = await fetch('/create', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },

      body: JSON.stringify(obj)
    });
    if (response.ok) {
      const data = await response.json();
      document.getElementById("msg").textContent = data.msg;
      if (data && data.success) {
        setTimeout(function () {
          window.location.replace('/thread?uuid=' + data.uuid);
        }, 2000);
      }
    } else {
      console.error('Failed to create a thread', response.msg);
    }
  } catch (error) {
    console.error('Error create thread', error);
  }

});
title.addEventListener('input', debounce(async function (event) {
  document.getElementById("msg").textContent = "";
}, 500))
content.addEventListener('input', debounce(async function (event) {
  document.getElementById("msg").textContent = "";
}, 500))
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