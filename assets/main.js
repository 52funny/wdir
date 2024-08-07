let button = document.getElementById('navbar-button');
let menu = document.getElementById('navbar-menu');
let show = false;
button.addEventListener('click', function (e) {
  show = !show;
  if (show) {
    e.target.classList.add('is-active');
    menu.classList.add('is-active');
  } else {
    e.target.classList.remove('is-active');
    menu.classList.remove('is-active');
  }
});

// 获取文件icon 属性
function getIconType(name) {
  // console.log(name);
  let imgType = ['gif', 'jpeg', 'jpg', 'bmp', 'png'],
    videoType = ['avi', 'wmv', 'mkv', 'mp4', 'mov', '3gp', 'flv', 'mpg', 'rmvb'],
    textType = ['txt', 'pdf', 'text', 'doc', 'docx', 'ppt'],
    musicType = ['wav', 'acc', 'flac', 'ape', 'ogg', 'mp3'],
    archiveType = ['zip', 'tar', 'gz', 'rar', 'jar', '7z'];

  if (imgType.includes(name)) {
    return 'fas fa-file-image';
  } else if (videoType.includes(name)) {
    return 'fas fa-file-video';
  } else if (textType.includes(name)) {
    return 'fas fa-file-alt';
  } else if (musicType.includes(name)) {
    return 'fas fa-file-music ';
  } else if (name === 'folder') {
    return 'fas fa-folder';
  } else if (archiveType.includes(name)) {
    return 'fas fa-file-archive';
  } else {
    return 'fas fa-file';
  }
}

let listIcon = document.querySelectorAll('.list .list-item .icon-text .icon');
listIcon.forEach((item) => {
  const type = item.getAttribute('type');
  // console.log(type);
  const icon = item.getElementsByTagName('i')[0];
  // console.log(icon);
  icon.setAttribute('class', getIconType(type));
});

let isSpread = localStorage.getItem('spread');
console.log('spread', isSpread);

// for render img
(async function () {
  // console.log(isSpread);
  if (isSpread == 'true') {
    const imgType = ['gif', 'jpeg', 'jpg', 'bmp', 'png'];
    const listItem = document.querySelectorAll('.list .list-item');
    listItem.forEach((it) => {
      const type = it.getAttribute('type');
      if (!imgType.includes(type)) return;

      const a = it.querySelector('a');
      const href = a.getAttribute('href');
      const newLink = document.createElement('a');
      newLink.setAttribute('href', href);

      const img = document.createElement('img');
      const columns = document.createElement('div');
      columns.setAttribute('class', 'columns');
      columns.setAttribute('style', 'justify-content: center;');
      const column = document.createElement('div');
      column.setAttribute('class', 'column is-11 center');
      columns.appendChild(column);
      column.appendChild(img);
      img.setAttribute('style', 'width: 100%;');
      img.setAttribute('src', href);
      img.setAttribute('loading', 'lazy');

      newLink.appendChild(columns);
      it.appendChild(newLink);
    });
  }
})();

const spreadBtn = document.getElementById('spread-button');

const changeBtnBackground = (btn, spred) => {
  if (spred == null) {
    return;
  }
  if (spred == 'true') {
    spreadBtn.classList.remove('is-light');
    spreadBtn.classList.add('is-dark');
  } else {
    spreadBtn.classList.remove('is-dark');
    spreadBtn.classList.add('is-light');
  }
};
changeBtnBackground(spreadBtn, isSpread);

spreadBtn.addEventListener('click', function (e) {
  isSpread = isSpread == 'true' ? 'false' : 'true';
  changeBtnBackground(spreadBtn, isSpread);
  localStorage.setItem('spread', isSpread);
});

// for download div
(async function () {
  const downloadDiv = document.querySelectorAll('.column .download-div');
  downloadDiv.forEach((it) => {
    const url = it.getAttribute('url');
    const name = it.getAttribute('name');
    it.addEventListener('click', function (e) {
      e.stopPropagation();
      const link = document.createElement('a');
      link.download = name;
      link.style.display = 'none';
      link.href = url;
      link.click();
      document.body.appendChild(link);
      document.body.removeChild(link);
    });
  });
})();

//for fuse.js
async function fuseSearch(searchItem) {
  // console.log(searchItem);
  const listItem = document.querySelectorAll('.list .list-item');
  let data = [...listItem].map((it) => {
    it.style = 'display:block';
    const a = it.querySelector('.icon-text a');
    return {
      name: a.textContent,
    };
  });
  if (searchItem.trim() == '') return;
  const options = { keys: ['name'] };
  const index = Fuse.createIndex(options.keys, data);
  const fuse = new Fuse(data, options, index);
  let result = fuse.search(searchItem);
  // console.log(JSON.stringify(result));
  listItem.forEach((it) => {
    const textContent = it.querySelector('.icon-text a').textContent;
    if (!result.map((it) => it.item.name).includes(textContent)) {
      it.style.display = 'none';
    }
  });
}
function inputChange(e) {
  e.stopPropagation();
  fuseSearch(e.target.value);
}
