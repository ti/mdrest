var AppSearch = (function () {
    'use strict';
    var datas = [];
    var defaultResult = [];
    var input, resultContent,orgHtml,hasMatch;
    var isSearching = false;
    var isInit = false;

    function init(indexPath, defaultResultPath, searchId, contentId) {
        if (isInit) {
            return
        }
        input = document.getElementById(searchId);
        resultContent = resultContent = document.getElementById(contentId);
        orgHtml = "<ul><li>没有查询到相关文档</li></ul>";
        if(defaultResultPath) {
            var reqDefaultResult = new XMLHttpRequest();
            reqDefaultResult.open("GET", defaultResultPath, true);
            reqDefaultResult.setRequestHeader("Accept", "application/json");
            reqDefaultResult.onreadystatechange = function() {
                if (!(reqDefaultResult.readyState == 4 && reqDefaultResult.status == 200)) {
                    return
                }
                defaultResult = JSON.parse(reqDefaultResult.responseText);
                if (!isSearching) {
                    search();
                }
            };
            reqDefaultResult.send();
        }
        var req = new XMLHttpRequest();
        req.open("GET", indexPath, true);
        req.setRequestHeader("Accept", "application/json");
        req.onreadystatechange = function() {
            if (!(req.readyState == 4 && req.status == 200)) {
                return
            }
            isInit = true;
            datas = JSON.parse(req.responseText);
            if (!isSearching) {
                search();
            }
        };
        req.send();
        input.addEventListener('keyup',search)
    }


    function search() {
        var keywords = input.value.trim().toLowerCase().split(/[\s\-]+/);
        if(undefined === keywords || keywords.length === 0) {
            return
        }
        if (datas.length === 0) {
            resultContent.innerHTML = "<ul>...</ul>";
            return;
        } else {
            resultContent.innerHTML = "<ul></ul>";
        }
        var resultLen = 0;
        isSearching = true;
        var resultUl = resultContent.firstElementChild;
        if (input.value.trim().length <= 0) {
            defaultResult.forEach(function (data) {
                var resultLi = "<li data-link href='/page/" + data.location + "'><div>" + data.title + "</div>" + "<small>" + data.text + "...</small></li>";
                resultUl.insertAdjacentHTML('beforeend', resultLi);
            });
            resultContent.toggleInternalLink();
            return;
        }
        hasMatch = false;
        datas.forEach(function (data) {
            if ("readme" === data.location || "README" === data.location || "about" === data.catalog) {
                return
            }
            var isMatch = true;
            var data_title = data.title.trim().toLowerCase();
            var content = data.text.trim().replace(/<[^>]+>/g, "");
            var data_content = content.toLowerCase();
            var tagsContent = "";
            var data_tags = data.tags;
            var data_url = data.location;
            var index_title = -1;
            var index_content = -1;
            var index_tags = -1;
            var first_occur = -1;
            // only match artiles with not empty titles and contents
            if (data_title != '' && data_content != '') {
                keywords.forEach(function (keyword, i) {
                    index_title = data_title.indexOf(keyword);
                    if(data_tags){
                        index_tags =  data_tags.indexOf(keyword);
                    }
                    if(index_tags < 0) {
                        index_content = data_content.indexOf(keyword);
                    } else {
                        tagsContent = "[" + data.tags.join("] [") + "] ";
                        index_content = tagsContent.indexOf(keyword);
                    }
                    if (index_title < 0 && index_content < 0) {
                        isMatch = false;
                    } else {
                        if (index_content < 0) {
                            index_content = 0;
                        }
                        if (i == 0) {
                            first_occur = index_content;
                        }
                    }
                });
            }
            // show search results
            if (isMatch) {
                if (resultLen > 10) {
                    return
                }
                resultLen ++;
                var resultLi = "<li data-link href='/page/" + data_url + "'><div>" + data.title + "</div>";
                if (first_occur >= 0) {
                    // cut out characters
                    var start = first_occur - 6;
                    var end = 100;
                    if (start < 0) {
                        start = 0;
                    }
                    if (end > (content.length - start)) {
                        end = (content.length - start);
                    }
                    var match_content = content.substr(start, end);
                    if(index_tags>=0) {
                        match_content = tagsContent + match_content;
                    }
                    keywords.forEach(function (keyword) {
                        var regS = new RegExp(keyword, "gi");
                        match_content = match_content.replace(regS, "<em>" + keyword + "</em>");
                    });
                    resultLi += "<small>" + match_content + "...</small></li>";
                }
                hasMatch = true;
                resultUl.insertAdjacentHTML('beforeend', resultLi);
            }
        });

        if (!hasMatch) {
            resultContent.innerHTML = orgHtml;
        } else {
            resultContent.toggleInternalLink();
        }
    }
    return {
        init: init,
        search:search
    };
})();

document.addEventListener("DOMContentLoaded", function() {
    var box = document.querySelector(".right-nav");
    box.onclick = function () {
        this.classList.add("is-focused");
        this.querySelector("input").focus();
        this.querySelector("input").addEventListener('keyup',function (){box.classList.add("is-focused")})
        AppSearch.init(MdRestConfig.BasePath + "mdrest_search_index.json",undefined,"search-field","search-result");
    };
    document.onclick = function (e) {
        if(e.target.id !== "search-field"){
            box.classList.remove("is-focused")
        }
    }



});