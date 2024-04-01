'use strict';

var App = (function () {
    var main = document.getElementById("main");
    var pageDoc = document.getElementById("doc");
    var pageTocContent = document.getElementById('toc-content').innerHTML;
    var pageTitle = document.getElementById("title");
    var docNav = document.getElementById("docs-nav");
    var content = document.getElementById("content");
    var rightToc = document.getElementById("right-toc");
    var header = document.getElementsByTagName("header")[0];
    var hasHistoy = false;
    var getBaseUrl = function () {
        var url = window.location.pathname;
        if ((url.lastIndexOf("/index.html") >= 0) || (url.lastIndexOf("/index.htm") >= 0)) {
            url = url.substring(0, url.lastIndexOf("/"));
        }
        if (url.substr(url.length - 1, url.length) === "/") {
            url = url.substr(0, url.length - 1)
        }
        return url;
    };

    var loading = {
        isDone: true,
        dom: document.getElementById("page-loading"),
        done: function () {
            loading.isDone = true;
            loading.dom.style.display = "none";
        },
        start: function () {
            loading.isDone = false;
            setTimeout(function () {
                if (!loading.isDone) {
                    loading.dom.style.display = "block";
                }
            }, 1200);
        }
    };

    var scrollIntoViewSmoothly = function (element) {
        element.scrollIntoView({behavior: "smooth", block: "start", inline: "nearest"});
    };

    var initHeader = function () {
        let header = document.querySelector('header');
        let logo = header.querySelector('a.logo');
        if (MdRestConfig.Logo) {
            if (MdRestConfig.Logo.Icon) {
                logo.querySelector('img').src = MdRestConfig.Logo.Icon;
            }
            if (MdRestConfig.Logo.Text) {
                logo.querySelector('span').innerHTML = MdRestConfig.Logo.Text;
            }
        }
        let navs = header.querySelector('.content ul');
        if (MdRestConfig.Headers.length > 0) {
            MdRestConfig.Headers.forEach(header => {
                let liElement = document.createElement('li');
                let aElement = document.createElement('a');
                aElement.textContent = header.Title;
                aElement.setAttribute('href', header.Href); // 设置a元素的href属性
                liElement.appendChild(aElement);
                navs.appendChild(liElement);
            });
        }
    }
    var baseActions = function () {
        HTMLDivElement.prototype.toggleInternalLink = function () {
            var anchors = this.querySelectorAll('a');
            for (var i = 0; i < anchors.length; i++) {
                anchors[i].onclick = function (e) {
                    var href = this.getAttribute("href");
                    if (!href) {
                        return
                    }
                    if (href.indexOf('#') === 0 && !MdRestConfig.SiglePage) {
                        e.preventDefault();
                        scrollIntoViewSmoothly(document.querySelector(href));
                    } else if (href.substr(-3) === ".md" && href.substring(0, 4) !== "http") {
                        e.preventDefault();
                        App.routes.goto("/page" + href.substring(0, href.length - 3));
                    } else if (!href.startsWith(MdRestConfig.GitPage) && href.substr(-3) === ".md") {
                        e.preventDefault();
                        href = href.substring(href.indexOf('/', 8), href.length - 3);
                        App.routes.goto("/page/" + href);
                    }
                };
            }
            anchors = this.querySelectorAll('div[data-link], li[data-link], a[data-link]');
            for (var i = 0; i < anchors.length; i++) {
                anchors[i].onclick = function (e) {
                    var href = this.getAttribute("href");
                    e.preventDefault();
                    if (href.indexOf('://') > 0 || href.indexOf('#') > 0 || this.getAttribute('target') === "_blank") {
                        link.href = href;
                        link.click();
                        return true;
                    }
                    App.routes.goto(href);
                };
            }
        };
        window.onhashchange = function () {
            if (!window.location.hash.indexOf("#/")) {
                App.routes.goto(decodeURI(window.location.hash).substr(1, window.location.hash.length));
            } else {
                App.routes.goto("/");
            }
        }

    };

    //get docs form documents by h2,h3
    var getPageToc = function (sel) {
        var documentRef = document.querySelector(sel);
        var toc = '';
        var level = 2;
        var headings = [].slice.call(documentRef.querySelectorAll('h2, h3'));
        if (headings.length < 2) {
            return ""
        }
        headings.forEach(function (heading, index) {
            var hIndex = parseInt(heading.nodeName.substring(1));
            var append = '<li><a href="#' + heading.id + '">' + heading.innerText + '</a>';

            if (hIndex > level) {
                append = '<ul>' + append;
            } else if (hIndex < level) {
                append = '</li></ul></li>' + append;
            }
            toc += append;
            level = hIndex;
        });
        var appendEnd = '</li>';
        if (level > 2) {
            appendEnd = '</li></ul></li>';
        }
        toc += appendEnd;
        return toc;
    };

    var activeTocIndex = -1;
    var toggleToc = function () {
        require.require("/assets/scripts/prism.js", "script", function (data, e) {
            if (!e) {
                Prism.highlightAll(pageDoc, false)
            }
        });
        rightToc.classList.add("hidden");
        var docs = document.querySelector(".docs");
        if (docs) {
            var pageTocHTML = getPageToc(".docs");
            if (pageTocHTML === "") {
                return
            }
            var tocContainer = document.getElementById('right-toc')
            tocContainer.innerHTML = pageTocContent + pageTocHTML;
            var navElements = document.querySelectorAll(".docs-toc a");
            for (var i = 0, l = navElements.length; i < l; i++) {
                navElements[i].addEventListener('click', function (e) {
                    var curActives = tocContainer.querySelectorAll("a.active");
                    if (curActives) {
                        curActives.forEach(function (e) {
                            e.classList.remove("active")
                        })
                    }
                    var id = this.getAttribute("href").substr(1);
                    e.preventDefault();
                    scrollIntoViewSmoothly(document.getElementById(id));

                    this.classList.add("active");
                });
            }
            var hElementsTops = undefined;
            window.addEventListener('scroll', function (e) {
                var docsTop = docs.getBoundingClientRect().top;
                if (docsTop >= 0) {
                    if (activeTocIndex !== -1) {
                        activeTocIndex = -1;
                        var curActive = tocContainer.querySelector("a.active");
                        if (curActive) {
                            curActive.classList.remove("active")
                        }
                    }
                    return
                }
                if (!hElementsTops) {
                    hElementsTops = [];
                    var hElements = document.querySelectorAll(".docs h2, .docs h3");
                    for (var i = 0, l = hElements.length; i < l; i++) {
                        var id = hElements[i].getAttribute("id");
                        if (id) {
                            hElementsTops.push({
                                id: id,
                                top: hElements[i].offsetTop - hElements[i].parentNode.offsetTop - 16,
                            });
                        }
                    }
                }
                var currentIndex = -1;
                for (var i = 0, l = hElementsTops.length; i < l; i++) {
                    if ((0 - docsTop) >= hElementsTops[i].top) {
                        currentIndex = i;
                        continue;
                    }
                }
                if (currentIndex > -1 && currentIndex !== activeTocIndex) {
                    activeTocIndex = currentIndex;
                    var curTocActive = tocContainer.querySelector("a.active");
                    if (curTocActive) {
                        curTocActive.classList.remove("active")
                    }
                    var tocDom = tocContainer.querySelector('.docs-toc a[href="#' + hElementsTops[currentIndex].id + '"]');
                    if (tocDom) {
                        tocDom.classList.add("active");
                    }
                }
            });
            rightToc.classList.remove("hidden");
        }
    };

    var require = {
        data: {},
        headEl: document.getElementsByTagName('head')[0],
        sync: true,
        reset: function (url) {
            require.data = {}
        },
        put: function (key, value) {
            require.data[key] = value;
        },
        require: function (url, type, callback) {
            if ("link" === type || "script" === type) {
                var status = require.data[url];
                if (undefined != status) {
                    if (200 === status) {
                        return callback();
                    } else {
                        return callback({}, {status: status});
                    }
                }
                var el = document.createElement(type), sync = false,
                    attrName, attributes;

                if ("link" === type) {
                    sync = true;
                    attributes = {rel: 'stylesheet', href: url, type: 'text/css'}
                } else {
                    attributes = {src: url}
                }
                for (attrName in attributes) {
                    el.setAttribute(attrName, attributes[attrName]);
                }
                if (callback) {
                    require.data[url] = 100;
                    el.addEventListener('load', function (e) {
                        require.data[url] = 200;
                        callback(e);
                    }, false);
                    setTimeout(function () {
                        if (200 !== require.data[url]) {
                            require.data[url] = 408;
                            callback({}, {status: 408});
                        }
                    }, 3000)
                } else {
                    require.data[url] = 200;
                }

                if (sync) {
                    require.headEl.appendChild(el);
                } else {
                    var s = document.getElementsByTagName(type)[0];
                    s.parentNode.insertBefore(el, s);
                }
                return;
            }
            var resp = require.data[url];
            if (resp) {
                callback(resp);
                return
            }
            var req = new XMLHttpRequest();
            req.open("GET", url, true);
            req.onreadystatechange = function () {
                if (req.readyState === 4) {
                    var data = null;
                    if (req.status === 200) {
                        if (type === "json") {
                            try {
                                resp = JSON.parse(req.responseText);
                            } catch (e) {
                                console.log("parse data error", url)
                                callback({}, {status: 415});
                            }
                        } else {
                            resp = req.responseText;
                        }
                        require.data[url] = resp;
                        callback(resp);

                    } else {
                        callback(resp, {status: req.status});
                    }
                }
            };
            req.send();
        }
    };

    var routes = {
        data: {
            "_error": '<div class="error"><i class="material-icons">info</i><h3>{{status}}</h3><p>{{msg}}</p></div>'
        },
        remove: function (url) {
            delete this.data[url];
        },
        exist: function (url) {
            return this.data.hasOwnProperty(url) && routes.data[url] !== null;
        },
        get: function (url) {
            return this.data[url];
        },
        add: function (url, title, templateUrl, mainCss, dataFunc) {
            var tmpData = {
                title: title,
                template: templateUrl,
                dataFunc: dataFunc
            };
            if (mainCss) {
                tmpData.mainCss = mainCss.split(" ")
            }
            this.data[url] = tmpData;

        },
        goto: function (url) {
            if (url === undefined || this.data.currentUrl === url) {
                return
            }
            loading.start();
            var curActive = docNav.querySelector("a.active");
            if (curActive) {
                curActive.classList.remove("active")
            }
            var headerActive = url;
            if ("/docs" === url.substring(0, 5)) {
                headerActive = "/docs"
            } else if ("/pages" === url.substring(0, 6)) {
                headerActive = "/"
            }
            curActive = docNav.querySelector('a[href="' + headerActive + '"]');
            if (curActive) {
                curActive.classList.add("active")
            }

            var urlBase = "/" + url.split("/", 2)[1];
            if (!this.exist(urlBase)) {
                content.innerHTML = Mustache.render(routes.data["_error"], {status: 404, msg: "页面 " + url + " 不存在"});
                return
            }
            var route = this.get(urlBase);
            var state = "#" + url;
            if (MdRestConfig.SiglePage) {
                state = url
            }
            if (url === "/") {
                state = "";
            }
            document.title = route.title;
            header.scrollIntoView();
            if (state !== window.location.hash) {
                var stateURL = window.location.pathname + state
                if (MdRestConfig.SiglePage) {
                    stateURL = window.location.origin + state
                }
                window.history.pushState(state, route.title, stateURL);
                hasHistoy = true
            }
            //clear header
            require.require(route.template, "html", function (tpl, err) {
                if (err) {
                    loading.done();
                    pageDoc.innerHTML = `<h2>错误 ${err.status}</h2>`
                    pageTitle.innerHTML = err.msg;
                    return
                }
                routes.data.currentUrl = url;
                if (route.mainCss) {
                    for (var i = 0, l = route.mainCss.length; i < l; i++) {
                        var mainCss = route.mainCss[i];
                        if (mainCss.substring(0, 1) === "-") {
                            main.classList.remove(mainCss.substring(1, mainCss.length))
                        } else {
                            main.classList.add(mainCss)
                        }
                    }
                }
                App.url = url;
                App.id = url;
                if (route.dataFunc) {
                    route.dataFunc(url, function (data, err) {
                        loading.done();
                        if (err) {
                            pageDoc.innerHTML = `<h2>错误 ${err.status}</h2>`
                            pageTitle.innerHTML = err.msg;
                            document.getElementsByClassName('docs-tools')[0].classList.add('hidden')
                        } else {
                            if (data.location) {
                                App.id = data.location;
                            }
                            if (data.title) {
                                document.title = data.title
                                pageTitle.innerHTML = data.title
                            }
                            if (data.siteMap) {
                                var nav = '<ul>';
                                var level = 0
                                var activeClass = ""
                                data.siteMap.forEach(function (page, index) {
                                    var item = '<li>';
                                    if (data.location) {
                                        if (data.location === page.location) {
                                            activeClass = `class="active"`
                                        }
                                    } else {
                                        if (index === 0) {
                                            activeClass = `class="active"`
                                        } else {
                                            activeClass = ''
                                        }
                                    }
                                    if (page.location.endsWith('/')) {
                                        item += `<a>${page.title}</a>`
                                        if (page.children) {
                                            item += `<ul>`
                                            page.children.forEach(function (childPage) {
                                                if (data.location === childPage.location) {
                                                    activeClass = `class="active"`
                                                } else {
                                                    activeClass = ''
                                                }
                                                item += `<li><a href="/page/${childPage.location}" ${activeClass} data-link>${childPage.title}</a></li>`
                                            })
                                            item += `</ul>`
                                        }
                                    } else {
                                        item += `<ul><li><a href="/page/${page.location}" ${activeClass} data-link>${page.title}</a></li></ul>`
                                    }
                                    item += '</li>'
                                    nav += item
                                });
                                var appendEnd = '</li></ul>';
                                if (level > 2) {
                                    appendEnd = '</li></ul></li></ul>';
                                }
                                nav += `</ul>`;
                                docNav.innerHTML = nav
                                main.toggleInternalLink()
                            }
                            if (data.page) {
                                pageDoc.innerHTML = data.page
                                pageDoc.toggleInternalLink()
                                document.getElementById('docs-git').href = `${MdRestConfig.GitPage}/blob/${MdRestConfig.GitBranch}/${data.location}.md`
                                toggleToc()
                            }
                        }
                    })
                } else {
                    loading.done();
                }
            })
        }
    };

    var toggleRoutes = function () {
        App.routes.add("/", MdRestConfig.Title, "", "", Blog.getDoc);
        App.routes.add("/page", MdRestConfig.Title, "", "", Blog.getDoc);
        if (MdRestConfig.SiglePage){
            if (window.location.pathname === '/') {
                App.routes.goto("/")
                return;
            }
            App.routes.goto(decodeURI(window.location.pathname).substr(0, window.location.pathname.length))
            return;
        }
        if ("#/ncr" === window.location) {
            return
        }
        if ("#/" === window.location.hash.substr(0, 2)) {
            App.routes.goto(decodeURI(window.location.hash).substr(1, window.location.hash.length));
        } else if ("#!/" === window.location.hash.substr(0, 3)) {
            App.routes.goto(decodeURI(window.location.hash).substr(2, window.location.hash.length));
        } else {
            App.routes.goto("/")
        }
        main.toggleInternalLink();
    };

    function init() {
        initHeader();
        baseActions();
        toggleRoutes();
    }

    return {
        init: init,
        require: require,
        routes: routes,
        baseUrl: getBaseUrl()
    };
})();


document.addEventListener("DOMContentLoaded", function () {
    App.init();
});


