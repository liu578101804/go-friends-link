

//默认数据
const fdata = {
    api_url: '',
    error_img: 'https://sdn.geekzu.org/avatar/57d8260dfb55501c37dde588e7c3852c'
}

//可通过 var fdataUser 替换默认值
if(typeof(fdataUser) !=="undefined"){
    for(let key in fdataUser) {
        if(fdataUser[key]){
            fdata[key] = fdataUser[key];
        }
    }
}

let container = document.getElementById('cf-container') ||  document.getElementById('fcircleContainer');

// 打印文章内容 cf-article
function loadArticleItem(datalist){
    let articleItem = ""
    for (let i = 0;i<datalist.length;i++){
        let item = datalist[i];
        articleItem +=`
      <div class="cf-article">
        <a class="cf-article-title" href="${item.link}" target="_blank" rel="noopener nofollow" data-title="${item.title}">${item.title}</a>
        <span class="cf-article-floor"></span>
        <div class="cf-article-avatar no-lightbox flink-item-icon">
          <img class="cf-img-avatar avatar" src="${item.avatar}" alt="avatar" onerror="this.src='${fdata.error_img}'; this.onerror = null;">
          <a onclick="openMeShow(event)" data-link="${item.link}" class="" target="_blank" rel="noopener nofollow" href="javascript:;"><span class="cf-article-author">${item.author}</span></a>
          <span class="cf-article-time">
            <span class="cf-time-created"><i class="far fa-calendar-alt">更新于</i>${item.updated}</span>
          </span>
        </div>
      </div>
      `;
    }
    container.insertAdjacentHTML('beforeend', articleItem);
}


// 输出基本结构
function loadStatistical(sdata){
    let messageBoard =`
  <div id="cf-state" class="cf-new-add">
    <div class="cf-state-data">
      <div class="cf-data-friends" onclick="openToShow()">
        <span class="cf-label">订阅</span>
        <span class="cf-message">${sdata.friends_num}</span>
      </div>
      <div class="cf-data-active" onclick="changeEgg()">
        <span class="cf-label">活跃</span>
        <span class="cf-message">${sdata.active_num}</span>
      </div>
      <div class="cf-data-article" onclick="clearLocal()">
        <span class="cf-label">日志</span>
        <span class="cf-message">${sdata.article_num}</span>
      </div>
    </div>
    <div id="cf-change">
        <span>更新于：${sdata.last_updated_time}</span>
    </div>
  </div>
  `;
    let loadMoreBtn = `
    <div id="cf-more" class="cf-new-add" onclick="loadNextArticle()">加载更多</div>
    <div id="cf-footer" class="cf-new-add">
     <span id="cf-version-up" onclick="checkVersion()"></span>
      Powered by <a target="_blank" href="https://github.com/Rock-Candy-Tea/hexo-circle-of-friends" target="_blank">friends-links</a>
    </div>
    <div id="cf-overlay" class="cf-new-add" onclick="closeShow()"></div>
    <div id="cf-overshow" class="cf-new-add"></div>
  `;
    if(container){
        container.insertAdjacentHTML('beforebegin', messageBoard);
        container.insertAdjacentHTML('afterend', loadMoreBtn);
    }
}


/*加载下一页*/
function loadNextArticle(){
    // 页码+1
    page = page+1
    fetchArticle(function (json){
        let articleData = eval(json.article_data);
        let articleItem = ""
        for (let i = 0;i<articleData.length;i++){
            let item = articleData[i];
            articleItem +=`
      <div class="cf-article">
        <a class="cf-article-title" href="${item.link}" target="_blank" rel="noopener nofollow" data-title="${item.title}">${item.title}</a>
        <span class="cf-article-floor"></span>
        <div class="cf-article-avatar no-lightbox flink-item-icon">
          <img class="cf-img-avatar avatar" src="${item.avatar}" alt="avatar" onerror="this.src='${fdata.error_img}'; this.onerror = null;">
          <a onclick="openMeShow(event)" data-link="${item.link}" class="" target="_blank" rel="noopener nofollow" href="javascript:;"><span class="cf-article-author">${item.author}</span></a>
          <span class="cf-article-time">
            <span class="cf-time-created"><i class="far fa-calendar-alt">更新于</i>${item.updated}</span>
          </span>
        </div>
      </div>
      `;
        }
        container.insertAdjacentHTML('beforeend', articleItem);
    })

}

// 获取文章列表
let page = 1;
function fetchArticle(call_back){
    let fetchUrl = fdata.api_url+"all?page="+page;
    fetch(fetchUrl)
        .then(res => res.json())
        .then(json =>{
            // 判断是否最后一页
            let articleData = eval(json.article_data);
            // 回调
            call_back(json)
            // 判断是否最后一页
            if(articleData.length<10){
                document.getElementById('cf-more').outerHTML = `<div id="cf-more" class="cf-new-add"><small>没有更多了！</small></div>`
            }
        })
}

function init(){

    // 检查配置文件信息
    if (fdata.api_url==""){
        let errMsg = "请配置主机地址"
        container.innerHTML = `<span style="color: #ff0000;">${errMsg}</span>`;
        return
    }

    // 加载第一页
    fetchArticle(function (json){
        container.innerHTML = "";
        let statisticalData = json.statistical_data;
        let articleData = eval(json.article_data);
        // 加载基本信息
        loadStatistical(statisticalData);
        // 加载文章信息
        loadArticleItem(articleData)
    })
}

init()