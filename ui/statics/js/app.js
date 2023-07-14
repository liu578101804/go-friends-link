


// 默认数据
var k_data = {
    api_url: '',
    error_img: 'https://friend-links.henjinet.com/ui/statics/imgs/default.jpg'
}

// 可通过 var kun_user 替换默认值
if(typeof(kun_user) !=="undefined"){
    for(let key in kun_user) {
        if(kun_user[key]){
            k_data[key] = kun_user[key];
        }
    }
}

// 主容器
var container = document.getElementById('k-container');
// 加载的页面
var page = 1;
// 坐标缓存
var articleIndex = 1

// 打印文章内容 k-article
function loadArticleItem(datalist){
    let articleItem = ""
    for (let i = 0;i<datalist.length;i++){

        let item = datalist[i];
        articleItem +=`
      <div class="k-article">
        <a class="k-article-title" href="${item.link}" target="_blank" rel="noopener nofollow" data-title="${item.title}">${item.title}</a>
        <span class="k-article-floor">${articleIndex}</span>
        <div class="k-article-avatar no-lightbox flink-item-icon">
          <img class="k-img-avatar avatar" src="${item.avatar}" alt="avatar" onerror="this.src='${k_data.error_img}'; this.onerror = null;">
          <a onclick="openMeShow(event)" data-link="${item.link}" class="" target="_blank" rel="noopener nofollow" href="javascript:;"><span class="k-article-author">${item.author}</span></a>
          <span class="k-article-time">
            <span class="k-time-created"><i class="far fa-calendar-alt">发布于</i>${item.created}</span>
          </span>
        </div>
      </div>
      `;

        articleIndex = articleIndex+1
    }
    container.insertAdjacentHTML('beforeend', articleItem);
}

// 输出基本结构
function loadStatistical(sdata){
    let messageBoard =`
  <div id="k-state" class="k-new-add">
    <div class="k-state-data">
      <div class="k-data-friends">
        <span class="k-label">订阅</span>
        <span class="k-message">${sdata.friends_num}</span>
      </div>
      <div class="k-data-article">
        <span class="k-label">日志</span>
        <span class="k-message">${sdata.article_num}</span>
      </div>
    </div>
    <div id="k-change">
        <span>更新于：${sdata.last_updated_time}</span>
    </div>
  </div>
  `;
    let loadMoreBtn = `
    <div id="k-more" class="k-new-add" onclick="loadNextArticle()">加载更多</div>
    <div id="k-footer" class="k-new-add">
     <span id="k-version-up"></span>
      Powered by <a target="_blank" href="https://github.com/liu578101804/go-friends-link" target="_blank">friends-links</a>
    </div>
    <div id="k-overlay" class="k-new-add" onclick="closeShow()"></div>
    <div id="k-overshow" class="k-new-add"></div>
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
        loadArticleItem(articleData);
    })

}

// 获取文章列表
function fetchArticle(call_back){
    let fetchUrl = k_data.api_url+"/articles?page="+page;
    fetch(fetchUrl)
        .then(res => res.json())
        .then(json =>{
            // 判断是否最后一页
            let articleData = eval(json.article_data);
            // 回调
            call_back(json)
            // 判断是否最后一页
            if(articleData.length<10){
                document.getElementById('k-more').outerHTML = `<div id="k-more" class="k-new-add"><small>没有更多了！</small></div>`
            }
        })
}

function init(){
    // 检查配置文件信息
    if (k_data.api_url==""){
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

// 启动
init()




